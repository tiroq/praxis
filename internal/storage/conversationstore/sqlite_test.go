package conversationstore

import (
	"context"
	"errors"
	"testing"
)

func openTestStore(t *testing.T) *SQLiteStore {
	t.Helper()

	store, err := OpenStore(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("OpenStore(:memory:) failed: %v", err)
	}
	return store
}

func TestConversationProjection_AppendUserMessage(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	conv := NewConversation("telegram-chat-100")
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation() failed: %v", err)
	}

	msg := NewMessage(conv.ID, "evt_user_1", "user", "hello", "2026-07-06T10:00:00Z", map[string]string{"chat_id": "100"})
	if err := store.AppendMessage(ctx, msg); err != nil {
		t.Fatalf("AppendMessage() failed: %v", err)
	}

	msgs, err := store.ListMessages(ctx, conv.ID, ListFilter{})
	if err != nil {
		t.Fatalf("ListMessages() failed: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].Role != "user" {
		t.Fatalf("expected role=user, got %q", msgs[0].Role)
	}
	if msgs[0].EventID != "evt_user_1" {
		t.Fatalf("expected event_id=evt_user_1, got %q", msgs[0].EventID)
	}
}

func TestConversationProjection_AppendAssistantMessage(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	conv := NewConversation("telegram-chat-101")
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation() failed: %v", err)
	}

	user := NewMessage(conv.ID, "evt_user_2", "user", "question", "2026-07-06T10:00:00Z", nil)
	assistant := NewMessage(conv.ID, "evt_user_2", "assistant", "answer", "2026-07-06T10:00:01Z", nil)

	if err := store.AppendMessage(ctx, user); err != nil {
		t.Fatalf("AppendMessage(user) failed: %v", err)
	}
	if err := store.AppendMessage(ctx, assistant); err != nil {
		t.Fatalf("AppendMessage(assistant) failed: %v", err)
	}

	msgs, err := store.ListMessages(ctx, conv.ID, ListFilter{})
	if err != nil {
		t.Fatalf("ListMessages() failed: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[1].Role != "assistant" {
		t.Fatalf("expected second role=assistant, got %q", msgs[1].Role)
	}
}

func TestConversationProjection_ListOrderedHistory(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	conv := NewConversation("telegram-chat-102")
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation() failed: %v", err)
	}

	later := NewMessage(conv.ID, "evt_order_2", "assistant", "later", "2026-07-06T10:00:02Z", nil)
	earlier := NewMessage(conv.ID, "evt_order_1", "user", "earlier", "2026-07-06T10:00:01Z", nil)

	if err := store.AppendMessage(ctx, later); err != nil {
		t.Fatalf("AppendMessage(later) failed: %v", err)
	}
	if err := store.AppendMessage(ctx, earlier); err != nil {
		t.Fatalf("AppendMessage(earlier) failed: %v", err)
	}

	msgs, err := store.ListMessages(ctx, conv.ID, ListFilter{})
	if err != nil {
		t.Fatalf("ListMessages() failed: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].EventID != "evt_order_1" || msgs[1].EventID != "evt_order_2" {
		t.Fatalf("messages not ordered by timestamp asc: %#v %#v", msgs[0], msgs[1])
	}
}

func TestConversationProjection_DuplicateMessageHandling(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	conv := NewConversation("telegram-chat-103")
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation() failed: %v", err)
	}

	msg := NewMessage(conv.ID, "evt_dup_1", "user", "hello", "2026-07-06T10:00:00Z", nil)
	msg.ID = "msg_duplicate"

	if err := store.AppendMessage(ctx, msg); err != nil {
		t.Fatalf("first AppendMessage() failed: %v", err)
	}

	err := store.AppendMessage(ctx, msg)
	if err == nil {
		t.Fatal("expected duplicate append error, got nil")
	}

	var dupErr ErrDuplicateMessage
	if !errors.As(err, &dupErr) {
		t.Fatalf("expected ErrDuplicateMessage, got %T: %v", err, err)
	}
}

func TestConversationProjection_MissingConversationHandling(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	msg := NewMessage("conv_missing", "evt_missing", "user", "hello", "2026-07-06T10:00:00Z", nil)
	err := store.AppendMessage(ctx, msg)
	if err == nil {
		t.Fatal("expected missing conversation error, got nil")
	}

	var notFound ErrConversationNotFound
	if !errors.As(err, &notFound) {
		t.Fatalf("expected ErrConversationNotFound in chain, got %T: %v", err, err)
	}
}

func TestConversationProjection_CorrelationIDPreservation(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	conv1, err := store.GetConversationByCorrelationID(ctx, "telegram-chat-104")
	if err != nil {
		t.Fatalf("GetConversationByCorrelationID(first) failed: %v", err)
	}
	if conv1.CorrelationID != "telegram-chat-104" {
		t.Fatalf("expected correlation_id=telegram-chat-104, got %q", conv1.CorrelationID)
	}

	conv2, err := store.GetConversationByCorrelationID(ctx, "telegram-chat-104")
	if err != nil {
		t.Fatalf("GetConversationByCorrelationID(second) failed: %v", err)
	}
	if conv2.ID != conv1.ID {
		t.Fatalf("expected same conversation id for same correlation id, got %q != %q", conv2.ID, conv1.ID)
	}
}
