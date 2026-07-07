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
		reply, err := client.GenerateReply(context.Background(), GenerateReplyRequest{
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

		_, err := client.GenerateReply(ctx, GenerateReplyRequest{
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
		reply, err := client.GenerateReply(context.Background(), GenerateReplyRequest{UserMessage: "hello"})
		if err != nil {
			t.Fatalf("GenerateReply() error = %v", err)
		}
		if reply != "legacy reply" {
			t.Fatalf("expected legacy reply, got %q", reply)
		}
	})
}

func TestClientExtractFacts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/v1/extract-facts" {
				t.Fatalf("expected /v1/extract-facts, got %s", r.URL.Path)
			}

			var req map[string]any
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			if req["latest_user_message"] != "I live in Bangkok." {
				t.Fatalf("expected latest_user_message, got %#v", req["latest_user_message"])
			}

			_ = json.NewEncoder(w).Encode(map[string]any{
				"facts": []map[string]any{{
					"type":       "location",
					"value":      "Bangkok",
					"confidence": 0.94,
				}},
			})
		}))
		defer server.Close()

		client := NewClient(server.URL+"/v1/reply", server.Client())
		resp, err := client.ExtractFacts(context.Background(), ExtractFactsRequest{
			LatestUserMessage: "I live in Bangkok.",
		})
		if err != nil {
			t.Fatalf("ExtractFacts() error = %v", err)
		}
		if len(resp.Facts) != 1 {
			t.Fatalf("expected 1 fact, got %d", len(resp.Facts))
		}
		if resp.Facts[0].Type != "location" || resp.Facts[0].Value != "Bangkok" {
			t.Fatalf("unexpected fact: %#v", resp.Facts[0])
		}
	})

	t.Run("invalid confidence fails", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"facts": []map[string]any{{
					"type":       "location",
					"value":      "Bangkok",
					"confidence": 1.5,
				}},
			})
		}))
		defer server.Close()

		client := NewClient(server.URL+"/v1/reply", server.Client())
		_, err := client.ExtractFacts(context.Background(), ExtractFactsRequest{LatestUserMessage: "I live in Bangkok."})
		if err == nil {
			t.Fatal("expected invalid confidence error")
		}
	})
}
