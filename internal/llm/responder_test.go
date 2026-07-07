package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tiroq/praxis/internal/storage/userfacts"
)

func TestConversationResponder_PersistsFactsWithoutChangingReply(t *testing.T) {
	ctx := context.Background()
	store, err := userfacts.OpenStore(ctx, ":memory:")
	if err != nil {
		t.Fatalf("OpenStore(:memory:) failed: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/extract-facts":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"facts": []map[string]any{{
					"type":       "location",
					"value":      "Bangkok",
					"confidence": 0.94,
				}},
			})
		case "/v1/reply":
			_ = json.NewEncoder(w).Encode(map[string]any{"reply_text": "That's nice."})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	responder := NewConversationResponder(NewClient(server.URL+"/v1/reply", server.Client()), nil, store)
	resp, err := responder.Respond(ctx, RespondRequest{
		InputEventID:  "evt_1",
		CorrelationID: "corr_1",
		UserMessage:   "I live in Bangkok.",
		Metadata:      map[string]string{"chat_id": "42"},
	})
	if err != nil {
		t.Fatalf("Respond() error = %v", err)
	}
	if resp.ReplyText != "That's nice." {
		t.Fatalf("unexpected reply: %q", resp.ReplyText)
	}

	stored, err := store.ListBySourceEvent(ctx, "evt_1")
	if err != nil {
		t.Fatalf("ListBySourceEvent() failed: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("expected 1 stored fact, got %d", len(stored))
	}
	if stored[0].Type != "location" || stored[0].Value != "Bangkok" || stored[0].CorrelationID != "corr_1" {
		t.Fatalf("unexpected stored fact: %#v", stored[0])
	}
}

func TestConversationResponder_ExtractionFailureDoesNotFailReply(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/extract-facts":
			http.Error(w, "extract unavailable", http.StatusBadGateway)
		case "/v1/reply":
			_ = json.NewEncoder(w).Encode(map[string]any{"reply_text": "Reply still works."})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	responder := NewConversationResponder(NewClient(server.URL+"/v1/reply", server.Client()), nil, nil)
	resp, err := responder.Respond(context.Background(), RespondRequest{
		InputEventID: "evt_1",
		UserMessage:  "hello",
	})
	if err != nil {
		t.Fatalf("Respond() error = %v", err)
	}
	if resp.ReplyText != "Reply still works." {
		t.Fatalf("unexpected reply: %q", resp.ReplyText)
	}
}

func TestConversationResponder_ReplyFailureDoesNotSkipFactPersistence(t *testing.T) {
	ctx := context.Background()
	store, err := userfacts.OpenStore(ctx, ":memory:")
	if err != nil {
		t.Fatalf("OpenStore(:memory:) failed: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/extract-facts":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"facts": []map[string]any{{
					"type":       "preference",
					"value":      "Go",
					"confidence": 0.8,
				}},
			})
		case "/v1/reply":
			http.Error(w, "reply unavailable", http.StatusBadGateway)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	responder := NewConversationResponder(NewClient(server.URL+"/v1/reply", server.Client()), nil, store)
	resp, err := responder.Respond(ctx, RespondRequest{
		InputEventID:  "evt_reply_fail",
		CorrelationID: "corr_reply_fail",
		UserMessage:   "I prefer Go.",
	})
	if err == nil {
		t.Fatal("expected reply failure")
	}
	if len(resp.Facts) != 1 {
		t.Fatalf("expected returned facts despite reply failure, got %#v", resp.Facts)
	}

	stored, err := store.ListBySourceEvent(ctx, "evt_reply_fail")
	if err != nil {
		t.Fatalf("ListBySourceEvent() failed: %v", err)
	}
	if len(stored) != 1 || stored[0].Value != "Go" {
		t.Fatalf("expected stored fact despite reply failure, got %#v", stored)
	}
}
