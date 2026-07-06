package llm

import (
	"context"
	"testing"
)

// TestNewReplyGenerator verifies that NewReplyGenerator creates a generator with expected defaults.
func TestNewReplyGenerator(t *testing.T) {
	client := NewClient("http://localhost:8081/v1/reply", nil)
	gen := NewReplyGenerator(client, nil)

	if gen == nil {
		t.Fatal("NewReplyGenerator returned nil")
	}
	if gen.client != client {
		t.Error("client not set correctly")
	}
	if gen.conversationStore != nil {
		t.Error("conversation store should be nil when not provided")
	}
	if gen.contextLimit != defaultContextLimit {
		t.Errorf("expected context limit %d, got %d", defaultContextLimit, gen.contextLimit)
	}
}

// TestReplyGenerator_WithContextLimit verifies that WithContextLimit sets the limit correctly.
func TestReplyGenerator_WithContextLimit(t *testing.T) {
	client := NewClient("http://localhost:8081/v1/reply", nil)
	gen := NewReplyGenerator(client, nil).WithContextLimit(5)

	if gen.contextLimit != 5 {
		t.Errorf("expected context limit 5, got %d", gen.contextLimit)
	}

	// Verify negative/zero limit is ignored
	gen.WithContextLimit(-1)
	if gen.contextLimit != 5 {
		t.Errorf("expected context limit to remain 5 after negative input, got %d", gen.contextLimit)
	}

	gen.WithContextLimit(0)
	if gen.contextLimit != 5 {
		t.Errorf("expected context limit to remain 5 after zero input, got %d", gen.contextLimit)
	}
}

// TestReplyGenerator_GenerateReply_EmptyUserMessage verifies that GenerateReply
// rejects empty user messages.
func TestReplyGenerator_GenerateReply_EmptyUserMessage(t *testing.T) {
	client := NewClient("http://localhost:8081/v1/reply", nil)
	gen := NewReplyGenerator(client, nil)

	ctx := context.Background()
	_, err := gen.GenerateReply(ctx, ReplyRequest{
		CorrelationID: "conv-123",
		UserMessage:   "",
		Source:        "test",
	})

	if err == nil {
		t.Fatal("expected error for empty user message")
	}
}

// TestReplyGenerator_GenerateReply_NilGenerator verifies that GenerateReply
// handles nil generator gracefully.
func TestReplyGenerator_GenerateReply_NilGenerator(t *testing.T) {
	var gen *ReplyGenerator

	_, err := gen.GenerateReply(context.Background(), ReplyRequest{
		UserMessage: "Hello",
	})

	if err == nil {
		t.Fatal("expected error for nil generator")
	}
}

// TestReplyGenerator_GenerateReply_NilClient verifies that GenerateReply
// handles nil client gracefully.
func TestReplyGenerator_GenerateReply_NilClient(t *testing.T) {
	gen := NewReplyGenerator(nil, nil)

	_, err := gen.GenerateReply(context.Background(), ReplyRequest{
		UserMessage: "Hello",
	})

	if err == nil {
		t.Fatal("expected error for nil client")
	}
}
