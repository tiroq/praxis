package kernel

import "strings"

// validateEvent checks that an Event satisfies all kernel invariants before
// the pipeline begins. Validation is the only system boundary where input
// rejection is appropriate.
func validateEvent(e Event) error {
	if strings.TrimSpace(e.ID) == "" {
		return ErrEmptyEventID
	}
	if strings.TrimSpace(e.Text) == "" {
		return ErrEmptyEventText
	}
	if e.Confidence < 0.0 || e.Confidence > 1.0 {
		return ErrInvalidConfidence
	}
	return nil
}

// validateDecision ensures the Decision produced by a DecisionMaker is well-formed.
func validateDecision(d Decision) error {
	if d.Outcome == "" {
		return ErrEmptyDecisionOutcome
	}
	return nil
}
