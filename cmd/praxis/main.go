package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/tiroq/praxis/internal/cli/praxiscli"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
	"github.com/tiroq/praxis/internal/transport/natscli"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := natstransport.ConfigFromEnv()
	app := praxiscli.NewApp(cfg, logger, natscli.NewTransport)

	rootCmd := &cobra.Command{
		Use:   "praxis",
		Short: "Praxis transport CLI over NATS",
	}

	rootCmd.AddCommand(newPublishCommand(app))
	rootCmd.AddCommand(newWatchCommand(app))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newPublishCommand(app *praxiscli.App) *cobra.Command {
	var text string
	var source string
	var messageID string
	var correlationID string

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish an input message to Praxis via NATS",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			published, err := app.Publish(ctx, praxiscli.PublishRequest{
				Text:          text,
				Source:        source,
				MessageID:     messageID,
				CorrelationID: correlationID,
			})
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "published message_id=%s correlation_id=%s\n", published.MessageID, published.CorrelationID)
			return err
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Input text to publish")
	cmd.Flags().StringVar(&source, "source", "praxis-cli", "Input source identifier")
	cmd.Flags().StringVar(&messageID, "id", "", "Optional explicit message id")
	cmd.Flags().StringVar(&correlationID, "correlation-id", "", "Optional explicit correlation id")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}

func newWatchCommand(app *praxiscli.App) *cobra.Command {
	var maxMessages int
	var pollTimeout time.Duration

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch output messages from Praxis over NATS",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			return app.Watch(ctx, praxiscli.WatchRequest{
				Writer:      cmd.OutOrStdout(),
				MaxMessages: maxMessages,
				PollTimeout: pollTimeout,
			})
		},
	}

	cmd.Flags().IntVar(&maxMessages, "max-messages", 0, "Optional number of messages to process before exit (0 = unlimited)")
	cmd.Flags().DurationVar(&pollTimeout, "poll-timeout", time.Second, "Polling timeout for next output message")
	return cmd
}
