package llm

import (
	"context"
	"strings"

	"github.com/tiroq/praxis/internal/storage/userfacts"
)

// FactExtractor owns extraction-specific orchestration and persistence policy.
type FactExtractor struct {
	client    *Client
	factStore *userfacts.SQLiteStore
}

func NewFactExtractor(client *Client, factStore *userfacts.SQLiteStore) *FactExtractor {
	return &FactExtractor{
		client:    client,
		factStore: factStore,
	}
}

type FactExtractionInput struct {
	SourceEventID       string
	SourceMessageID     string
	CorrelationID       string
	Source              string
	UserMessage         string
	Metadata            map[string]string
	ConversationContext ConversationContext
}

func (e *FactExtractor) Extract(ctx context.Context, input FactExtractionInput) ([]userfacts.CandidateFact, error) {
	extracted, err := e.client.ExtractFacts(ctx, ExtractFactsRequest{
		CorrelationID:       input.CorrelationID,
		Source:              input.Source,
		LatestUserMessage:   input.UserMessage,
		Metadata:            input.Metadata,
		ConversationContext: input.ConversationContext,
	})
	if err != nil {
		return nil, err
	}

	if len(extracted.Facts) == 0 {
		return []userfacts.CandidateFact{}, nil
	}

	userID := deriveUserID(input.Metadata, input.CorrelationID)
	facts := make([]userfacts.CandidateFact, 0, len(extracted.Facts))
	for _, extractedFact := range extracted.Facts {
		fact := userfacts.NewCandidateFact(
			userID,
			input.CorrelationID,
			extractedFact.Type,
			extractedFact.Value,
			extractedFact.Confidence,
			input.SourceEventID,
			input.SourceMessageID,
		)
		if e.factStore != nil {
			if err := e.factStore.Append(ctx, fact); err != nil {
				return facts, err
			}
		}
		facts = append(facts, *fact)
	}

	return facts, nil
}

func deriveUserID(metadata map[string]string, correlationID string) string {
	for _, key := range []string{"user_id", "telegram_user_id", "username", "chat_id"} {
		if value := strings.TrimSpace(metadata[key]); value != "" {
			return value
		}
	}
	return correlationID
}
