package kernel

import (
	"context"
	"fmt"
	"time"
)

// ConfidenceThreshold groups the thresholds that drive RuleBasedDecisionMaker.
type ConfidenceThreshold struct {
	// Approve is the minimum Review assessment confidence required for an
	// Approve outcome. Below this, the decision escalates or defers.
	Approve float64

	// Defer is the minimum confidence for a Defer outcome.
	// Confidence below Defer yields a NeedsRevision outcome.
	Defer float64
}

// DefaultConfidenceThreshold provides sensible defaults.
var DefaultConfidenceThreshold = ConfidenceThreshold{
	Approve: 0.6,
	Defer:   0.3,
}

// RuleBasedDecisionMaker is a deterministic DecisionMaker.
// It applies the following policy in order:
//
//  1. If the Review recommends Reject → outcome is Reject.
//  2. If the Review recommends Escalate → outcome is Escalate.
//  3. If confidence >= Approve threshold → outcome is Approve.
//  4. If confidence >= Defer threshold  → outcome is Defer.
//  5. Otherwise                          → outcome is NeedsRevision.
//
// No LLM calls are made.
type RuleBasedDecisionMaker struct {
	Thresholds ConfidenceThreshold
	PolicyName string
	ActorName  string
}

// NewRuleBasedDecisionMaker creates a RuleBasedDecisionMaker with the given
// thresholds. Zero-value thresholds are replaced by DefaultConfidenceThreshold.
// PolicyName and ActorName default to "rule-based-policy" and "rule-based-decision-maker".
func NewRuleBasedDecisionMaker(thresholds ConfidenceThreshold, policyName, actorName string) *RuleBasedDecisionMaker {
	if thresholds.Approve == 0 && thresholds.Defer == 0 {
		thresholds = DefaultConfidenceThreshold
	}
	if policyName == "" {
		policyName = "rule-based-policy"
	}
	if actorName == "" {
		actorName = "rule-based-decision-maker"
	}
	return &RuleBasedDecisionMaker{
		Thresholds: thresholds,
		PolicyName: policyName,
		ActorName:  actorName,
	}
}

// Decide applies rule-based policy to derive an explicit, immutable Decision.
func (d *RuleBasedDecisionMaker) Decide(ctx context.Context, event Event, review Review) (Decision, error) {
	confidence := review.Assessment.Confidence

	var outcome DecisionOutcome
	var reasoning string

	switch {
	case review.Recommendation == ReviewRecommendationReject:
		outcome = DecisionOutcomeReject
		reasoning = "reviewer recommended rejection"

	case review.Recommendation == ReviewRecommendationEscalate:
		outcome = DecisionOutcomeEscalate
		reasoning = "reviewer recommended escalation"

	case confidence >= d.Thresholds.Approve:
		outcome = DecisionOutcomeApprove
		reasoning = fmt.Sprintf("confidence %.2f meets approval threshold %.2f", confidence, d.Thresholds.Approve)

	case confidence >= d.Thresholds.Defer:
		outcome = DecisionOutcomeDefer
		reasoning = fmt.Sprintf("confidence %.2f is below approval threshold %.2f; deferring", confidence, d.Thresholds.Approve)

	default:
		outcome = DecisionOutcomeNeedsRevision
		reasoning = fmt.Sprintf("confidence %.2f is too low; revision required", confidence)
	}

	evidenceSummary := fmt.Sprintf("review %s (recommendation: %s, score: %.2f, confidence: %.2f)",
		review.ID, review.Recommendation, review.Assessment.Score, confidence)

	return Decision{
		ID:              fmt.Sprintf("dec-%s", event.ID),
		EventID:         event.ID,
		ReviewIDs:       []string{review.ID},
		Timestamp:       time.Now().UTC(),
		Outcome:         outcome,
		Confidence:      confidence,
		Reasoning:       reasoning,
		EvidenceSummary: evidenceSummary,
		Policy:          d.PolicyName,
		Actor:           d.ActorName,
	}, nil
}
