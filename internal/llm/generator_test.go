package llm

import (
	"context"
	"testing"
)

// TestNewConversationResponder verifies construction.
func TestNewConversationResponder(t *testing.T) {
	client := NewClient("http://localhost:8081/v1/reply", nil)
	responder := NewConversationResponder(client, nil)

	if responder == nil {
		t.Fatal("NewConversationResponder returned nil")
	}
	if responder.client != client {
		t.Error("client not set correctly")
	}
	if responder.conversationStore != nil {
		t.Error("conversation store should be nil when not provided")
	}
}

// TestConversationResponder_Respond_EmptyUserMessage verifies rejection of empty messages.
func TestConversationResponder_Respond_EmptyUserMessage(t *testing.T) {
	client := NewClient("http://localhost:8081/v1/reply", nil)
	responder := NewConversationResponder(client, nil)

	_, err := responder.Respond(context.Background(), RespondRequest{
		CorrelationID: "conv-123",
		UserMessage:   "",
		Source:        "test",
	})

	if err == nil {
		t.Fatal("expected error for empty user message")
	}
}

// TestConversationResponder_Respond_NilResponder verifies error handling.
func TestConversationResponder_Respond_NilResponder(t *testing.T) {
	var responder *ConversationResponder

	_, err := responder.Respond(context.Background(), RespondRequest{
		UserMessage: "Hello",
	})

	if err == nil {
		t.Fatal("expected error for nil responder")
	}
}

// TestConversationResponder_Respond_NilClient verifies error handling.
func TestConversationResponder_Respond_NilClient(t *testing.T) {
	responder := NewConversationResponder(nil, nil)

	_, err := responder.Respond(context.Background(), RespondRequest{
		UserMessage: "Hello",
	})

	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

// TestBackwardCompatibility_ReplyGenerator verifies old API still works.
func TestBackwardCompatibility_ReplyGenerator(t *testing.T) {
	client := NewClient("http://localhost:8081/v1/reply", nil)

	// Old API: NewReplyGenerator (should work via alias)
	gen := NewReplyGenerator(client, nil)
	if gen == nil {
		t.Fatal("NewReplyGenerator returned nil")
	}

	// Verify type
	if _, ok := interface{}(gen).(*ConversationResponder); !ok {
		t.Error("ReplyGenerator alias not pointing to ConversationResponder")
	}
}

// TestBackwardCompatibility_WithContextLimit verifies deprecated method exists.
func TestBackwardCompatibility_WithContextLimit(t *testing.T) {
	client := NewClient("http://localhost:8081/v1/reply", nil)
	responder := NewConversationResponder(client, nil)

	// Old API: WithContextLimit (should be no-op)
	result := responder.WithContextLimit(5)
	if result != responder {
		t.Error("WithContextLimit should return same responder")
	}
}

// TestLoadConversationContext_NoStore verifies graceful degradation.
func TestLoadConversationContext_NoStore(t *testing.T) {
	ctx := context.Background()

	// With nil store, should return empty context
	conversationCtx := loadConversationContext(ctx, "any-id", nil)
	msgs := conversationCtx.Messages()

	if len(msgs) != 0 {
		t.Errorf("expected empty messages, got %d", len(msgs))
	}
}
