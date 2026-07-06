package llm

import (
	"context"
	"fmt"

	"github.com/tiroq/praxis/internal/storage/conversationstore"
)

// RespondRequest is the semantic input to conversation response generation.
// It contains only what the caller (e.g., worker) needs to provide.
// This is the stable public API boundary for the LLM layer.
type RespondRequest struct {
	// CorrelationID groups messages into a logical conversation.
	CorrelationID string
	// UserMessage is the latest user input.
	UserMessage string
	// Source identifies the originating system (e.g., "telegram").
	Source string
	// Metadata carries optional enrichment fields.
	Metadata map[string]string
}

// RespondResponse is the semantic output from conversation response generation.
// This is the stable public API boundary for the LLM layer.
type RespondResponse struct {
	ReplyText string
}

// ConversationResponder orchestrates reply generation with bounded conversation context.
//
// It owns the orchestration logic only. Context loading, prompt construction,
// and provider interaction are delegated to lower layers:
//   - ConversationContext (context loading from projection)
//   - Client (HTTP transport to llm-router)
//   - llm-router (prompt construction and provider routing)
//
// The responder does not know about storage implementation (SQLite, etc.)
// or context policy (bounded history limits, etc.). Those are private details
// of the context loading layer.
//
// If the conversation store is nil, Respond will work but without conversation
// history (graceful degradation).
type ConversationResponder struct {
	client            *Client
	conversationStore *conversationstore.SQLiteStore
}

// NewConversationResponder constructs a ConversationResponder.
// The conversation store is optional; nil means no conversation history is available.
func NewConversationResponder(
	client *Client,
	conversationStore *conversationstore.SQLiteStore,
) *ConversationResponder {
	return &ConversationResponder{
		client:            client,
		conversationStore: conversationStore,
	}
}

// Respond takes a semantic request and returns a response with bounded-context reply.
//
// The responder delegates:
//   - Context loading to loadConversationContext (private to llm package)
//   - Prompt construction to llm-router service
//   - Provider interaction to Client
//
// If context loading fails, the response is generated from the user message alone
// (graceful degradation). Errors are returned only for critical failures
// (e.g., nil client, empty user message, network errors).
func (r *ConversationResponder) Respond(ctx context.Context, req RespondRequest) (RespondResponse, error) {
	if r == nil || r.client == nil {
		return RespondResponse{}, fmt.Errorf("conversation responder not properly configured")
	}

	if req.UserMessage == "" {
		return RespondResponse{}, fmt.Errorf("user_message is required")
	}

	// Load bounded conversation context (handles graceful degradation internally)
	conversationCtx := loadConversationContext(ctx, req.CorrelationID, r.conversationStore)

	// Build the request to llm-router with optional conversation context
	genReq := GenerateRequest{
		CorrelationID: req.CorrelationID,
		Source:        req.Source,
		UserMessage:   req.UserMessage,
		Metadata:      req.Metadata,
		Conversation:  conversationCtx.Messages(),
	}

	// Call the client to generate the reply
	reply, err := r.client.Generate(ctx, genReq)
	if err != nil {
		return RespondResponse{}, err
	}

	return RespondResponse{ReplyText: reply.ReplyText}, nil
}
