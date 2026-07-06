package llm

import (
	"context"
	"fmt"

	"github.com/tiroq/praxis/internal/storage/conversationstore"
)

// ReplyRequest is the semantic input to reply generation.
// It contains only what the caller (e.g., worker) needs to provide.
type ReplyRequest struct {
	// CorrelationID groups messages into a logical conversation.
	CorrelationID string
	// UserMessage is the latest user input.
	UserMessage string
	// Source identifies the originating system (e.g., "telegram").
	Source string
	// Metadata carries optional enrichment fields.
	Metadata map[string]string
}

// ReplyResponse is the semantic output from reply generation.
type ReplyResponse struct {
	ReplyText string
}

// ReplyGenerator assembles bounded conversation context and generates an LLM reply.
// It owns the logic for loading conversation history from storage and shaping
// it for the llm-router, keeping this complexity out of the worker layer.
//
// If the conversation store is nil, GenerateReply will work but without
// conversation history (graceful degradation).
type ReplyGenerator struct {
	client                  *Client
	conversationStore       *conversationstore.SQLiteStore
	contextLimit            int // max messages to include in context
	inputEventIDPlaceholder string
}

const defaultContextLimit = 10

// NewReplyGenerator constructs a ReplyGenerator with required and optional dependencies.
func NewReplyGenerator(
	client *Client,
	conversationStore *conversationstore.SQLiteStore,
) *ReplyGenerator {
	return &ReplyGenerator{
		client:            client,
		conversationStore: conversationStore,
		contextLimit:      defaultContextLimit,
	}
}

// WithContextLimit sets the maximum number of recent messages to include in context.
// Returns the ReplyGenerator for method chaining.
func (g *ReplyGenerator) WithContextLimit(limit int) *ReplyGenerator {
	if limit > 0 {
		g.contextLimit = limit
	}
	return g
}

// GenerateReply takes semantic input and returns a bounded-context reply.
//
// If the conversation store is available and the correlation_id is non-empty,
// it loads recent messages and includes them in the request to llm-router.
// If the store is nil or conversation not found, the reply is generated
// from the user message alone (graceful degradation).
func (g *ReplyGenerator) GenerateReply(ctx context.Context, req ReplyRequest) (ReplyResponse, error) {
	if g == nil || g.client == nil {
		return ReplyResponse{}, fmt.Errorf("reply generator not properly configured")
	}

	if req.UserMessage == "" {
		return ReplyResponse{}, fmt.Errorf("user_message is required")
	}

	// Build the request to llm-router
	genReq := GenerateRequest{
		CorrelationID: req.CorrelationID,
		Source:        req.Source,
		UserMessage:   req.UserMessage,
		Metadata:      req.Metadata,
	}

	// Optionally load and populate bounded conversation context
	if g.conversationStore != nil && req.CorrelationID != "" {
		history, err := g.loadBoundedContext(ctx, req.CorrelationID)
		if err != nil {
			// Non-fatal: log would happen in caller; continue without history
			// (similar to current worker behavior)
		} else if len(history) > 0 {
			genReq.Conversation = history
		}
	}

	// Call the client to generate the reply
	reply, err := g.client.Generate(ctx, genReq)
	if err != nil {
		return ReplyResponse{}, err
	}

	return ReplyResponse{ReplyText: reply.ReplyText}, nil
}

// loadBoundedContext fetches recent messages from the conversation projection
// and converts them to ConversationMessage format.
// Returns an empty slice if conversation not found or fetch fails (non-fatal).
func (g *ReplyGenerator) loadBoundedContext(ctx context.Context, correlationID string) ([]ConversationMessage, error) {
	// Fetch the conversation by correlation_id
	conv, err := g.conversationStore.GetConversationByCorrelationID(ctx, correlationID)
	if err != nil || conv == nil {
		return nil, nil // Non-fatal; conversation not found
	}

	// Fetch recent messages (bounded by context limit)
	msgs, err := g.conversationStore.ListMessages(
		ctx,
		conv.ID,
		conversationstore.ListFilter{Limit: g.contextLimit},
	)
	if err != nil || len(msgs) == 0 {
		return nil, nil // Non-fatal; no messages found
	}

	// Convert to ConversationMessage format
	history := make([]ConversationMessage, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Role == "user" || msg.Role == "assistant" {
			history = append(history, ConversationMessage{
				Role: msg.Role,
				Text: msg.Content,
			})
		}
	}

	return history, nil
}
