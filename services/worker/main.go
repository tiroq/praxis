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
func buildKernel() *kernel.Kernel {
	reviewer := kernel.NewKeywordReviewer(defaultKeywords, "keyword-reviewer-v1")
	dm := kernel.NewRuleBasedDecisionMaker(kernel.DefaultConfidenceThreshold, "rule-based-policy-v1", "worker")
	planner := kernel.NewSimpleActionPlanner(nil)
	return kernel.New(reviewer, dm, planner)
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

	client, err := natstransport.NewClient(cfg)
	if err != nil {
		logger.Error("failed to connect to NATS", "err", err)
		os.Exit(1)
	}
	defer client.Close()

	js := client.JetStream()
	pub := natstransport.NewPublisher(js, cfg.OutputSubject)
	k := buildKernel()
	sub := natsworker.NewSubscriber(js, cfg, k, pub, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := sub.Run(ctx); err != nil {
		logger.Error("subscriber exited with error", "err", err)
		os.Exit(1)
	}

	logger.Info("worker shut down cleanly")
}
