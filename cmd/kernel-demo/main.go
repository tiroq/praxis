package main
// kernel-demo is a minimal CLI that exercises the Core Kernel pipeline
// end-to-end without any external dependencies.
//
// Usage:
//
//	go run ./cmd/kernel-demo "text to evaluate"
//
// The CLI constructs a single Event from the supplied text, runs it through
// the default KeywordReviewer → RuleBasedDecisionMaker → SimpleActionPlanner
// pipeline, and prints the PipelineResult as pretty-printed JSON.
// It exits non-zero on validation or runtime errors.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tiroq/praxis/internal/core/kernel"
)

func main() {
	if len(os.Args) < 2 || strings.TrimSpace(strings.Join(os.Args[1:], " ")) == "" {
		fmt.Fprintln(os.Stderr, "usage: kernel-demo <text>")
		os.Exit(1)
	}

	text := strings.Join(os.Args[1:], " ")
	now := time.Now().UTC()

	// Build a minimal Event. ID uses a timestamp-based key so the demo
	// is reproducible without requiring a UUID library.
	event := kernel.Event{
		ID:               fmt.Sprintf("demo-%d", now.UnixMicro()),
		Text:             text,
		Type:             "manual.input",
		Source:           "kernel-demo-cli",
		Actor:            "user",
		OccurredAt:       now,
		ObservedAt:       now,
		Confidence:       0.75,
		TrustLevel:       kernel.TrustLevelMedium,
		ValidationStatus: kernel.ValidationStatusPending,
		Origin:           kernel.EventOriginUser,
		SchemaVersion:    "1",
	}

	// Default keyword set covers common Russian and English trigger words.
	keywords := map[string]string{
		"билет":    "travel",
		"билеты":   "travel",
		"купить":   "purchase",
		"buy":      "purchase",
		"urgent":   "time-sensitive",
		"срочно":   "time-sensitive",
		"review":   "needs-review",
		"approve":  "approval",
		"blocked":  "blocker",
	}

	reviewer := kernel.NewKeywordReviewer(keywords, "keyword-reviewer-v1")
	decisionMaker := kernel.NewRuleBasedDecisionMaker(kernel.DefaultConfidenceThreshold, "rule-based-policy-v1", "kernel-demo")
	planner := kernel.NewSimpleActionPlanner(nil)

	k := kernel.New(reviewer, decisionMaker, planner)

	result, err := k.Run(context.Background(), event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kernel error: %v\n", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "json encode error: %v\n", err)
		os.Exit(1)
	}
}
