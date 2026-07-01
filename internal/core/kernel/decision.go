package kernel

import "time"

// DecisionOutcome is the explicit, authoritative commitment produced by a Decision.
// Per RFC-021 §11. Outcomes are commitments, not evaluations.
type DecisionOutcome string

const (
	DecisionOutcomeApprove       DecisionOutcome = "approve"
	DecisionOutcomeReject        DecisionOutcome = "reject"
	DecisionOutcomeNeedsRevision DecisionOutcome = "needs_revision"
	DecisionOutcomeEscalate      DecisionOutcome = "escalate"
	DecisionOutcomeDefer         DecisionOutcome = "defer"
	DecisionOutcomeSplit         DecisionOutcome = "split"
	DecisionOutcomeMerge         DecisionOutcome = "merge"
	DecisionOutcomeArchive       DecisionOutcome = "archive"
	DecisionOutcomeNoAction      DecisionOutcome = "no_action"
)

// Decision is an immutable, explicit commitment made based on one or more Reviews.
// Decisions never evaluate; they commit. They are always auditable and traceable.
// Per RFC-021: decisions are append-only and must never be modified after creation.
type Decision struct {
	ID        string
	EventID   string   // the original event this decision concerns
	ReviewIDs []string // IDs of reviews that informed this decision

	Timestamp  time.Time
	Outcome    DecisionOutcome
	Confidence float64 // 0.0–1.0

	Reasoning       string // human-readable rationale
	EvidenceSummary string // summary of the evidence that drove the decision
	Policy          string // name/version of the policy applied
	Actor           string // who or what made this decision
}
