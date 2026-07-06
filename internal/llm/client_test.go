package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientGenerateReply(t *testing.T) {
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
			if req["user_message"] != "hello" {
				t.Fatalf("expected user_message=hello, got %#v", req["user_message"])
			}

			_ = json.NewEncoder(w).Encode(map[string]any{"reply_text": "generated reply"})
		}))
		defer server.Close()

		client := NewClient(server.URL, server.Client())
		reply, err := client.GenerateReply(context.Background(), GenerateReplyInput{
			InputEventID:  "evt_1",
			CorrelationID: "corr_1",
			Source:        "telegram",
			UserMessage:   "hello",
		})
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
			_ = json.NewEncoder(w).Encode(map[string]any{"reply_text": "late reply"})
		}))
		defer server.Close()

		client := NewClient(server.URL, server.Client())
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
		defer cancel()

		_, err := client.GenerateReply(ctx, GenerateReplyInput{
			InputEventID:  "evt_timeout",
			CorrelationID: "corr_timeout",
			UserMessage:   "hello",
		})
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
	})

	t.Run("legacy response compatibility", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{"assistant_reply": "legacy reply"})
		}))
		defer server.Close()

		client := NewClient(server.URL, server.Client())
		reply, err := client.GenerateReply(context.Background(), GenerateReplyInput{UserMessage: "hello"})
		if err != nil {
			t.Fatalf("GenerateReply() error = %v", err)
		}
		if reply != "legacy reply" {
			t.Fatalf("expected legacy reply, got %q", reply)
		}
	})
}
