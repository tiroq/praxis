// worker is the NATS JetStream transport adapter for the Core Kernel.
//
// It subscribes to an input subject, runs each message through the Praxis
// kernel pipeline (Event → Review → Decision → Action), and publishes the
// result to an output subject.  The kernel itself is never modified; all
// NATS concerns are isolated to internal/transport/nats.
//
// Configuration is read from environment variables; see ConfigFromEnv for
// the full list and defaults.
//
// Usage:
//
//	go run ./services/worker
//	NATS_URL=nats://remote:4222 go run ./services/worker
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tiroq/praxis/internal/core/kernel"
	"github.com/tiroq/praxis/internal/storage"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
	"github.com/tiroq/praxis/internal/transport/natsworker"
)

// defaultKeywords mirrors the keyword set used by the api-kernel service so
// that the worker applies the same deterministic review logic.
var defaultKeywords = map[string]string{
	"билет":   "travel",
	"билеты":  "travel",
	"купить":  "purchase",
	"buy":     "purchase",
	"urgent":  "time-sensitive",
	"срочно":  "time-sensitive",
	"review":  "needs-review",
	"approve": "approval",
	"blocked": "blocker",
}

// buildKernel wires the default, deterministic pipeline components.
// No LLM calls, no external dependencies.
// Accepts optional kernel.Option parameters for configuration (e.g., WithEventRecorder).
func buildKernel(opts ...kernel.Option) *kernel.Kernel {
	reviewer := kernel.NewKeywordReviewer(defaultKeywords, "keyword-reviewer-v1")
	dm := kernel.NewRuleBasedDecisionMaker(kernel.DefaultConfidenceThreshold, "rule-based-policy-v1", "worker")
	planner := kernel.NewSimpleActionPlanner(nil)
	return kernel.New(reviewer, dm, planner, opts...)
}

// tryOpenStorage attempts to open storage using environment configuration.
// If storage cannot be opened, logs the error and returns nil (non-fatal).
// The worker will continue without event recording if storage is unavailable.
func tryOpenStorage(ctx context.Context, logger *slog.Logger) *storage.Storage {
	cfg := storage.ConfigFromEnv()
	logger.Info("attempting to open storage",
		"backend", cfg.Backend,
		"sqlite_path", cfg.SQLitePath,
	)

	store, err := storage.Open(ctx, cfg)
	if err != nil {
		logger.Warn("failed to open storage; worker will continue without event recording",
			"err", err,
			"backend", cfg.Backend,
		)
		return nil
	}

	logger.Info("storage opened successfully",
		"backend", cfg.Backend,
	)
	return store
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := natstransport.ConfigFromEnv()
	logger.Info("worker starting",
		"nats_url", cfg.URL,
		"stream", cfg.StreamName,
		"input_subject", cfg.InputSubject,
		"output_subject", cfg.OutputSubject,
		"durable", cfg.Durable,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Attempt to open storage; non-fatal if it fails
	store := tryOpenStorage(ctx, logger)
	if store != nil {
		defer func() {
			if err := store.Close(); err != nil {
				logger.Error("failed to close storage", "err", err)
			} else {
				logger.Info("storage closed")
			}
		}()
	}

	client, err := natstransport.NewClient(cfg)
	if err != nil {
		logger.Error("failed to connect to NATS", "err", err)
		os.Exit(1)
	}
	defer client.Close()

	js := client.JetStream()
	pub := natstransport.NewPublisher(js, cfg.OutputSubject)

	// Build kernel with optional event recording
	var k *kernel.Kernel
	if store != nil {
		recorder := newEventRecorderAdapter(store.Events)
		k = buildKernel(kernel.WithEventRecorder(recorder))
		logger.Info("kernel built with event recording enabled")
	} else {
		k = buildKernel()
		logger.Info("kernel built without event recording")
	}

	sub := natsworker.NewSubscriber(js, cfg, k, pub, logger)

	if err := sub.Run(ctx); err != nil {
		logger.Error("subscriber exited with error", "err", err)
		os.Exit(1)
	}

	logger.Info("worker shut down cleanly")
}
