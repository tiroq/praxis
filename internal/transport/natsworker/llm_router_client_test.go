package natsworker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tiroq/praxis/internal/core/kernel"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

func TestLLMRouterClientGenerateReply(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			var req map[string]any
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			if req["input_event_id"] != "evt_1" {
				t.Fatalf("expected input_event_id=evt_1, got %#v", req["input_event_id"])
			}
			if req["correlation_id"] != "corr_1" {
				t.Fatalf("expected correlation_id=corr_1, got %#v", req["correlation_id"])
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"assistant_reply": "generated reply"})
		}))
		defer server.Close()

		client := NewLLMRouterClient(server.URL, server.Client())
		reply, err := client.GenerateReply(context.Background(), natstransport.InputMessage{
			ID:            "evt_1",
			CorrelationID: "corr_1",
			Source:        "telegram",
			Text:          "hello",
		}, kernel.PipelineResult{EventID: "evt_1"})
		if err != nil {
			t.Fatalf("GenerateReply() error = %v", err)
		}
		if reply != "generated reply" {
			t.Fatalf("expected generated reply, got %q", reply)
		}
	})

	t.Run("timeout via context", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(120 * time.Millisecond)
			_ = json.NewEncoder(w).Encode(map[string]any{"assistant_reply": "late reply"})
		}))
		defer server.Close()

		client := NewLLMRouterClient(server.URL, server.Client())
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
		defer cancel()

		_, err := client.GenerateReply(ctx, natstransport.InputMessage{
			ID:            "evt_timeout",
			CorrelationID: "corr_timeout",
			Text:          "hello",
		}, kernel.PipelineResult{EventID: "evt_timeout"})
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
	})
}
