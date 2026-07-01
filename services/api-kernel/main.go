// api-kernel is a minimal HTTP wrapper around the Core Kernel pipeline.
// It exposes a single endpoint for local testing without any external
// dependencies (no NATS, no storage, no LLM).
//
// Usage:
//
//	go run ./services/api-kernel
//
// Endpoint: POST /v1/kernel/run
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tiroq/praxis/internal/core/kernel"
)

// runRequest is the inbound payload for POST /v1/kernel/run.
type runRequest struct {
	Text   string `json:"text"`
	Source string `json:"source"`
}

// errorResponse is returned on 4xx/5xx.
type errorResponse struct {
	Error string `json:"error"`
}

// defaultKeywords covers common Russian and English trigger words for the
// default KeywordReviewer used in local testing.
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

// buildKernel wires the default pipeline components.
func buildKernel() *kernel.Kernel {
	reviewer := kernel.NewKeywordReviewer(defaultKeywords, "keyword-reviewer-v1")
	dm := kernel.NewRuleBasedDecisionMaker(kernel.DefaultConfidenceThreshold, "rule-based-policy-v1", "api-kernel")
	planner := kernel.NewSimpleActionPlanner(nil)
	return kernel.New(reviewer, dm, planner)
}

// kernelRunHandler returns an http.HandlerFunc for POST /v1/kernel/run.
func kernelRunHandler(k *kernel.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
			return
		}

		var req runRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON: " + err.Error()})
			return
		}
		if strings.TrimSpace(req.Text) == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "text must not be empty"})
			return
		}

		source := strings.TrimSpace(req.Source)
		if source == "" {
			source = "api"
		}

		now := time.Now().UTC()
		event := kernel.Event{
			ID:               fmt.Sprintf("api-%d", now.UnixMicro()),
			Text:             req.Text,
			Type:             "api.input",
			Source:           source,
			Actor:            "api-kernel",
			OccurredAt:       now,
			ObservedAt:       now,
			Confidence:       0.75,
			TrustLevel:       kernel.TrustLevelMedium,
			ValidationStatus: kernel.ValidationStatusPending,
			Origin:           kernel.EventOriginUser,
			SchemaVersion:    "1",
		}

		result, err := k.Run(context.Background(), event)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "kernel error: " + err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, result)
	}
}

// writeJSON serialises v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	k := buildKernel()

	mux := http.NewServeMux()
	mux.Handle("/v1/kernel/run", kernelRunHandler(k))

	addr := ":8080"
	fmt.Printf("api-kernel listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
	}
}
