package userfacts

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ValidationState string

const (
	ValidationStateObserved      ValidationState = "observed"
	ValidationStateExtracted     ValidationState = "extracted"
	ValidationStateCorrelated    ValidationState = "correlated"
	ValidationStateReviewed      ValidationState = "reviewed"
	ValidationStateHumanApproved ValidationState = "human_approved"
	ValidationStateCanonical     ValidationState = "canonical"
)

var validationStateOrder = map[ValidationState]int{
	ValidationStateObserved:      1,
	ValidationStateExtracted:     2,
	ValidationStateCorrelated:    3,
	ValidationStateReviewed:      4,
	ValidationStateHumanApproved: 5,
	ValidationStateCanonical:     6,
}

func IsValidValidationState(state ValidationState) bool {
	_, ok := validationStateOrder[state]
	return ok
}

func CanTransitionValidationState(from ValidationState, to ValidationState) bool {
	fromOrder, okFrom := validationStateOrder[from]
	toOrder, okTo := validationStateOrder[to]
	if !okFrom || !okTo {
		return false
	}
	return toOrder == fromOrder+1
}

// CandidateFact is an extracted, unvalidated user fact candidate.
// It is not memory retrieval input and is not canonical truth.
type CandidateFact struct {
	ID                  string
	UserID              string
	CorrelationID       string
	Type                string
	Value               string
	Confidence          float64
	SourceEventID       string
	SourceMessageID     string
	ValidationState     ValidationState
	ValidationUpdatedAt string
	CreatedAt           string
}

// NewCandidateFact creates a candidate fact with a stable local identifier.
func NewCandidateFact(
	userID, correlationID, factType, value string,
	confidence float64,
	sourceEventID, sourceMessageID string,
) *CandidateFact {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return &CandidateFact{
		ID:                  generateID("fact"),
		UserID:              userID,
		CorrelationID:       correlationID,
		Type:                factType,
		Value:               value,
		Confidence:          confidence,
		SourceEventID:       sourceEventID,
		SourceMessageID:     sourceMessageID,
		ValidationState:     ValidationStateExtracted,
		ValidationUpdatedAt: now,
		CreatedAt:           now,
	}
}

func generateID(prefix string) string {
	id := uuid.New().String()
	hexID := id[:8] + id[9:13] + id[14:18]
	return fmt.Sprintf("%s_%s", prefix, hexID)
}
