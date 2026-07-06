package conversationstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Conversation represents a persistent dialogue between participants.
// It is a Canonical Object per RFC-014, owned by and stored in the Canonical Store per RFC-033.
//
// Fields:
//
//	ID: Globally unique, immutable conversation identifier (conv_*)
//	CorrelationID: External grouping identifier (e.g. telegram-chat-{chat_id})
//	Lifecycle: Conversation state (created, active, archived)
//	CreatedAt: ISO8601 timestamp of conversation creation
//	UpdatedAt: ISO8601 timestamp of last metadata update
//	LastMessageAt: ISO8601 timestamp of last message in conversation (or empty)
type Conversation struct {
	ID            string
	CorrelationID string
	Lifecycle     string // "created", "active", or "archived"
	CreatedAt     string // ISO8601
	UpdatedAt     string // ISO8601
	LastMessageAt string // ISO8601 or empty
}

// Message represents a single message in a conversation.
// It is a Projection per RFC-033, referencing immutable Events per RFC-013.
// Messages are append-only; once created, they are never modified.
//
// Fields:
//
//	ID: Globally unique, immutable message identifier (msg_*)
//	ConversationID: Foreign key to Conversation
//	EventID: Reference to immutable event in Event Store (never changes)
//	Role: "user" or "assistant" (role of message author)
//	Content: Message text content
//	Timestamp: ISO8601 timestamp of message
//	Metadata: Optional enrichment fields (e.g. username, message_id, chat_id)
//	CreatedAt: ISO8601 timestamp of when message was stored
type Message struct {
	ID             string
	ConversationID string
	EventID        string
	Role           string // "user" or "assistant"
	Content        string
	Timestamp      string            // ISO8601
	Metadata       map[string]string // optional enrichment
	CreatedAt      string            // ISO8601
}

// NewConversation creates a new Conversation with a given correlation_id.
// Assigns a unique ID and sets timestamps to now.
func NewConversation(correlationID string) *Conversation {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return &Conversation{
		ID:            generateID("conv"),
		CorrelationID: correlationID,
		Lifecycle:     "created",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewMessage creates a new Message in a conversation.
// Assigns a unique ID and sets CreatedAt to now.
func NewMessage(conversationID, eventID, role, content string, timestamp string, metadata map[string]string) *Message {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return &Message{
		ID:             generateID("msg"),
		ConversationID: conversationID,
		EventID:        eventID,
		Role:           role,
		Content:        content,
		Timestamp:      timestamp,
		Metadata:       metadata,
		CreatedAt:      now,
	}
}

// generateID creates a unique identifier with the given prefix using UUID.
// Format: prefix_randomhex (12 hex chars from UUID)
func generateID(prefix string) string {
	// Use UUID v4 and truncate to 12 hex chars for compact IDs
	id := uuid.New().String()
	// Remove hyphens and take first 12 chars
	hexID := id[:8] + id[9:13] + id[14:18]
	return fmt.Sprintf("%s_%s", prefix, hexID)
}
