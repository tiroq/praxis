package kernel

import "errors"

// Sentinel errors for the kernel pipeline.
// Callers should use errors.Is for matching.
var (
	// ErrEmptyEventText is returned when Event.Text is empty.
	ErrEmptyEventText = errors.New("kernel: event text must not be empty")

	// ErrEmptyEventID is returned when Event.ID is empty.
	ErrEmptyEventID = errors.New("kernel: event ID must not be empty")

	// ErrInvalidConfidence is returned when Event.Confidence is outside [0, 1].
	ErrInvalidConfidence = errors.New("kernel: confidence must be between 0.0 and 1.0")

	// ErrEmptyDecisionOutcome is returned when a Decision has no explicit outcome.
	ErrEmptyDecisionOutcome = errors.New("kernel: decision outcome must be explicit and non-empty")

	// ErrReviewFailed is a wrapper sentinel for errors originating in the Reviewer.
	ErrReviewFailed = errors.New("kernel: reviewer failed")

	// ErrDecisionFailed is a wrapper sentinel for errors originating in the DecisionMaker.
	ErrDecisionFailed = errors.New("kernel: decision maker failed")

	// ErrPlanFailed is a wrapper sentinel for errors originating in the ActionPlanner.
	ErrPlanFailed = errors.New("kernel: action planner failed")
)
