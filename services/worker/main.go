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
	"strconv"
	"syscall"
	"time"

	"github.com/tiroq/praxis/internal/core/kernel"
	"github.com/tiroq/praxis/internal/llm"
	"github.com/tiroq/praxis/internal/storage"
	"github.com/tiroq/praxis/internal/storage/conversationstore"
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
// LLM reply generation is orchestrated by the worker subscriber.
// Accepts optional kernel.Option parameters for configuration (e.g., WithEventRecorder).
func buildKernel(opts ...kernel.Option) *kernel.Kernel {
	reviewer := kernel.NewKeywordReviewer(defaultKeywords, "keyword-reviewer-v1")
	dm := kernel.NewRuleBasedDecisionMaker(kernel.DefaultConfidenceThreshold, "rule-based-policy-v1", "worker")
	planner := kernel.NewSimpleActionPlanner(nil)
	return kernel.New(reviewer, dm, planner, opts...)
}

type llmRuntimeConfig struct {
	Enabled  bool
	Endpoint string
	Timeout  time.Duration
}

func llmRuntimeConfigFromEnv() llmRuntimeConfig {
	timeout := 5 * time.Second
	if raw := os.Getenv("PRAXIS_LLM_TIMEOUT"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			timeout = parsed
		}
	}

	enabled := true
	if raw := os.Getenv("PRAXIS_LLM_ENABLED"); raw != "" {
		if parsed, err := strconv.ParseBool(raw); err == nil {
			enabled = parsed
		}
	}

	endpoint := os.Getenv("PRAXIS_LLM_ROUTER_URL")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:8081/v1/reply"
	}

	return llmRuntimeConfig{
		Enabled:  enabled,
		Endpoint: endpoint,
		Timeout:  timeout,
	}
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

// tryOpenConversationStore attempts to open the SQLite-backed conversation projection store.
// If store cannot be opened, logs the error and returns nil (non-fatal).
// The worker will continue without conversation persistence if store is unavailable.
func tryOpenConversationStore(ctx context.Context, logger *slog.Logger) *conversationstore.SQLiteStore {
	// For now, conversation store uses SQLite at the same path as event store.
	// Environment: PRAXIS_STORAGE_PATH (defaults to build/praxis.db)
	cfg := storage.ConfigFromEnv()
	if cfg.Backend != storage.BackendSQLite {
		logger.Warn("conversation store only supports SQLite backend; skipping initialization",
			"backend", cfg.Backend,
		)
		return nil
	}

	logger.Info("attempting to open conversation store",
		"backend", cfg.Backend,
		"sqlite_path", cfg.SQLitePath,
	)

	store, err := conversationstore.OpenStore(ctx, cfg.SQLitePath)
	if err != nil {
		logger.Warn("failed to open conversation store; worker will continue without conversation persistence",
			"err", err,
		)
		return nil
	}

	logger.Info("conversation store opened successfully")
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

	// Attempt to open conversation store; non-fatal if it fails
	convStore := tryOpenConversationStore(ctx, logger)
	if convStore != nil {
		defer func() {
			if err := convStore.Close(); err != nil {
				logger.Error("failed to close conversation store", "err", err)
			} else {
				logger.Info("conversation store closed")
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
	if convStore != nil {
		sub.WithConversationStore(convStore)
		logger.Info("subscriber configured with conversation store")
	}

	llmCfg := llmRuntimeConfigFromEnv()
	if llmCfg.Enabled {
		llmClient := llm.NewClient(llmCfg.Endpoint, nil)
		sub.WithReplyGenerator(func(ctx context.Context, input natstransport.InputMessage, _ kernel.PipelineResult) (string, error) {
			// Build the LLM request with bounded conversation context
			genReq := llm.GenerateRequest{
				InputEventID:  input.ID,
				CorrelationID: input.CorrelationID,
				Source:        input.Source,
				UserMessage:   input.Text,
				Metadata:      input.Metadata,
			}

			// If conversation store is available, fetch and populate recent conversation history
			if convStore != nil && input.CorrelationID != "" {
				conv, err := convStore.GetConversationByCorrelationID(ctx, input.CorrelationID)
				if err == nil && conv != nil {
					// Fetch recent messages (bounded context: last 10 messages)
					msgs, err := convStore.ListMessages(ctx, conv.ID, conversationstore.ListFilter{Limit: 10})
					if err == nil && len(msgs) > 0 {
						// Convert to ConversationMessage structs for llm-router
						conversationHistory := make([]llm.ConversationMessage, 0, len(msgs))
						for _, msg := range msgs {
							if msg.Role == "user" || msg.Role == "assistant" {
								conversationHistory = append(conversationHistory, llm.ConversationMessage{
									Role: msg.Role,
									Text: msg.Content,
								})
							}
						}
						genReq.Conversation = conversationHistory
					}
				}
			}

			resp, err := llmClient.Generate(ctx, genReq)
			if err != nil {
				return "", err
			}
			return resp.ReplyText, nil
		}, llmCfg.Timeout)
		logger.Info("subscriber configured with llm reply generator",
			"llm_router_url", llmCfg.Endpoint,
			"llm_timeout", llmCfg.Timeout,
		)
	} else {
		logger.Info("llm reply generator disabled")
	}

	if err := sub.Run(ctx); err != nil {
		logger.Error("subscriber exited with error", "err", err)
		os.Exit(1)
	}

	logger.Info("worker shut down cleanly")
}
