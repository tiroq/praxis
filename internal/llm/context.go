package llm

import (
	"context"

	"github.com/tiroq/praxis/internal/storage/conversationstore"
)

// ConversationContext represents bounded conversation history for LLM context.
// It is owned internally by the LLM package and not exposed outside.
type ConversationContext struct {
	messages []ConversationMessage
}

// loadConversationContext fetches and assembles bounded context from the conversation projection.
// This is an internal function that owns context policy and projection interaction.
//
// Context policy (last 10 messages, user/assistant only) is private to this function.
// If projection unavailable, returns empty context (graceful degradation).
func loadConversationContext(
	ctx context.Context,
	correlationID string,
	store *conversationstore.SQLiteStore,
) ConversationContext {
	if store == nil || correlationID == "" {
		return ConversationContext{} // No context available
	}

	// Fetch the conversation by correlation_id
	conv, err := store.GetConversationByCorrelationID(ctx, correlationID)
	if err != nil || conv == nil {
		return ConversationContext{} // Conversation not found (non-fatal)
	}

	// Context policy: fetch last 10 messages (private constant, not exposed)
	const contextLimit = 10

	// Fetch recent messages (bounded by context policy)
	msgs, err := store.ListMessages(
		ctx,
		conv.ID,
		conversationstore.ListFilter{Limit: contextLimit},
	)
	if err != nil || len(msgs) == 0 {
		return ConversationContext{} // No messages found (non-fatal)
	}

	// Convert to ConversationMessage format, filtering by role
	messages := make([]ConversationMessage, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Role == "user" || msg.Role == "assistant" {
			messages = append(messages, ConversationMessage{
				Role: msg.Role,
				Text: msg.Content,
			})
		}
	}

	return ConversationContext{messages: messages}
}

// Messages returns the bounded conversation history.
func (cc ConversationContext) Messages() []ConversationMessage {
	if cc.messages == nil {
		return []ConversationMessage{} // Empty slice, not nil
	}
	return cc.messages
}
