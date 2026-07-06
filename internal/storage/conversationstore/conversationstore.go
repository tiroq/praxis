package conversationstore

import (
	"context"
)

// ListFilter specifies criteria for filtering messages in a conversation.
type ListFilter struct {
	Limit  int // maximum number of messages to return (0 = default)
	Offset int // number of messages to skip
}

// ConversationStore defines the interface for persistent conversation storage.
// Implements RFC-033 Canonical Store (conversations) + Projection Store (messages).
type ConversationStore interface {
	// CreateConversation persists a new Conversation object.
	// Returns error if conversation with same id already exists.
	CreateConversation(ctx context.Context, conv *Conversation) error

	// GetConversation retrieves a Conversation by ID.
	// Returns ErrConversationNotFound if not found.
	GetConversation(ctx context.Context, id string) (*Conversation, error)

	// GetConversationByCorrelationID retrieves or creates a Conversation by correlation_id.
	// If a conversation with this correlation_id exists, returns it.
	// Otherwise creates a new conversation with the given correlation_id.
	// Used during message ingestion to ensure a conversation exists.
	GetConversationByCorrelationID(ctx context.Context, correlationID string) (*Conversation, error)

	// AppendMessage persists a new Message in append-only fashion.
	// Returns error if message with same id already exists.
	// Returns error if referenced conversation does not exist.
	AppendMessage(ctx context.Context, msg *Message) error

	// ListMessages retrieves all messages for a conversation, ordered by timestamp ascending.
	// Returns empty list if conversation has no messages.
	ListMessages(ctx context.Context, conversationID string, filter ListFilter) ([]*Message, error)

	// UpdateConversationMetadata updates last_message_at and updated_at for a conversation.
	// Used after appending a message to maintain conversation state.
	// Does NOT allow updating id, correlation_id, or created_at.
	UpdateConversationMetadata(ctx context.Context, id string) error

	// ArchiveConversation marks a conversation as archived (lifecycle = "archived").
	ArchiveConversation(ctx context.Context, id string) error

	// Close releases any resources held by the store.
	Close() error
}
