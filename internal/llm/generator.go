package llm

import (
	"context"

	"github.com/tiroq/praxis/internal/storage/conversationstore"
)

// ReplyRequest is a backward-compatibility alias for RespondRequest.
// Deprecated: Use RespondRequest instead. Existing code continues to work.
type ReplyRequest = RespondRequest

// ReplyResponse is a backward-compatibility alias for RespondResponse.
// Deprecated: Use RespondResponse instead. Existing code continues to work.
type ReplyResponse = RespondResponse

// ReplyGenerator is a backward-compatibility alias for ConversationResponder.
// Deprecated: Use ConversationResponder instead. Existing code continues to work.
type ReplyGenerator = ConversationResponder

// NewReplyGenerator is a backward-compatibility function that creates a ConversationResponder.
// Deprecated: Use NewConversationResponder instead. This exists only for compatibility.
func NewReplyGenerator(
	client *Client,
	conversationStore interface{}, // Accept any type for compatibility
) *ConversationResponder {
	// Handle both *SQLiteStore and nil
	var store *conversationstore.SQLiteStore
	if cs, ok := conversationStore.(*conversationstore.SQLiteStore); ok {
		store = cs
	}
	return NewConversationResponder(client, store)
}

// Compatibility method: GenerateReply wraps Respond for backward compatibility.
// Deprecated: Use Respond method on ConversationResponder instead.
func (r *ConversationResponder) GenerateReply(ctx context.Context, req ReplyRequest) (ReplyResponse, error) {
	return r.Respond(ctx, req)
}

// Compatibility method: WithContextLimit is deprecated and does nothing.
// Context limits are now private to the context loading layer.
// Deprecated: Context policy is now internal to the llm package.
func (r *ConversationResponder) WithContextLimit(limit int) *ConversationResponder {
	// This method no longer has an effect.
	// Context limits are managed internally by loadConversationContext.
	return r
}
