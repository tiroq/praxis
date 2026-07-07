package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/tiroq/praxis/internal/storage/conversationstore"
	"github.com/tiroq/praxis/internal/storage/userfacts"
)

// RespondRequest is the semantic input to conversation response generation.
type RespondRequest struct {
	// InputEventID identifies the source event that caused this response.
	InputEventID string
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
type RespondResponse struct {
	ReplyText string
	Facts     []userfacts.CandidateFact
}

// ConversationResponder orchestrates LLM capabilities for one user message.
type ConversationResponder struct {
	contextLoader *ConversationContextLoader
	replyService  *ReplyService
	factExtractor *FactExtractor
}

// NewConversationResponder constructs a concrete orchestration boundary for LLM capabilities.
func NewConversationResponder(
	client *Client,
	conversationStore *conversationstore.SQLiteStore,
	factStore *userfacts.SQLiteStore,
) *ConversationResponder {
	return &ConversationResponder{
		contextLoader: NewConversationContextLoader(conversationStore),
		replyService:  NewReplyService(client),
		factExtractor: NewFactExtractor(client, factStore),
	}
}

// Respond runs reply generation and fact extraction independently.
func (r *ConversationResponder) Respond(ctx context.Context, req RespondRequest) (RespondResponse, error) {
	if r == nil || r.replyService == nil || r.factExtractor == nil || r.contextLoader == nil {
		return RespondResponse{}, fmt.Errorf("conversation responder not properly configured")
	}

	req.UserMessage = strings.TrimSpace(req.UserMessage)
	if req.UserMessage == "" {
		return RespondResponse{}, fmt.Errorf("user_message is required")
	}

	correlationID := normalizeCorrelationID(req.CorrelationID, req.InputEventID)
	sourceEventID := normalizeSourceEventID(req.InputEventID, correlationID)
	sourceMessageID := resolveSourceMessageID(req.Metadata, sourceEventID)
	conversationCtx := r.contextLoader.Load(ctx, correlationID)

	replyResult := make(chan replyOutcome, 1)
	factResult := make(chan factOutcome, 1)

	go func() {
		reply, err := r.replyService.Generate(ctx, ReplyRequest{
			InputEventID:        sourceEventID,
			CorrelationID:       correlationID,
			Source:              req.Source,
			UserMessage:         req.UserMessage,
			Metadata:            req.Metadata,
			ConversationContext: conversationCtx,
		})
		replyResult <- replyOutcome{reply: reply, err: err}
	}()

	go func() {
		facts, err := r.factExtractor.Extract(ctx, FactExtractionInput{
			SourceEventID:       sourceEventID,
			SourceMessageID:     sourceMessageID,
			CorrelationID:       correlationID,
			Source:              req.Source,
			UserMessage:         req.UserMessage,
			Metadata:            req.Metadata,
			ConversationContext: conversationCtx,
		})
		factResult <- factOutcome{facts: facts, err: err}
	}()

	reply := <-replyResult
	facts := <-factResult

	if reply.err != nil {
		return RespondResponse{
			Facts: facts.facts,
		}, reply.err
	}

	// Fact extraction remains best effort and must not fail user replies.
	if facts.err != nil {
		return RespondResponse{
			ReplyText: reply.reply,
		}, nil
	}

	return RespondResponse{
		ReplyText: reply.reply,
		Facts:     facts.facts,
	}, nil
}

type replyOutcome struct {
	reply string
	err   error
}

type factOutcome struct {
	facts []userfacts.CandidateFact
	err   error
}

func normalizeCorrelationID(correlationID string, inputEventID string) string {
	normalized := strings.TrimSpace(correlationID)
	if normalized != "" {
		return normalized
	}
	return strings.TrimSpace(inputEventID)
}

func normalizeSourceEventID(inputEventID string, correlationID string) string {
	normalized := strings.TrimSpace(inputEventID)
	if normalized != "" {
		return normalized
	}
	return strings.TrimSpace(correlationID)
}

func resolveSourceMessageID(metadata map[string]string, sourceEventID string) string {
	if metadata != nil {
		if messageID := strings.TrimSpace(metadata["message_id"]); messageID != "" {
			return messageID
		}
	}
	return sourceEventID
}
