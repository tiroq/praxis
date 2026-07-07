package llm

import (
	"context"

	"github.com/tiroq/praxis/internal/storage/conversationstore"
)

// ConversationContext is the context envelope passed between LLM capability services.
// Today it contains only conversation messages and can grow in future sprints.
type ConversationContext struct {
	Messages []ConversationMessage
}

type ConversationContextLoader struct {
	store *conversationstore.SQLiteStore
}

func NewConversationContextLoader(store *conversationstore.SQLiteStore) *ConversationContextLoader {
	return &ConversationContextLoader{store: store}
}

// Load fetches and assembles bounded context from the conversation projection.
func (l *ConversationContextLoader) Load(
	ctx context.Context,
	correlationID string,
) ConversationContext {
	if l == nil || l.store == nil || correlationID == "" {
		return ConversationContext{
			Messages: []ConversationMessage{},
		}
	}

	conv, err := l.store.GetConversationByCorrelationID(ctx, correlationID)
	if err != nil || conv == nil {
		return ConversationContext{
			Messages: []ConversationMessage{},
		}
	}

	const contextLimit = 10
	msgs, err := l.store.ListMessages(
		ctx,
		conv.ID,
		conversationstore.ListFilter{Limit: contextLimit},
	)
	if err != nil || len(msgs) == 0 {
		return ConversationContext{
			Messages: []ConversationMessage{},
		}
	}

	messages := make([]ConversationMessage, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Role == "user" || msg.Role == "assistant" {
			messages = append(messages, ConversationMessage{
				Role: msg.Role,
				Text: msg.Content,
			})
		}
	}

	return ConversationContext{
		Messages: messages,
	}
}
