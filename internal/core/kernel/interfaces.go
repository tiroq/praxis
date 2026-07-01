package kernel

import "context"

// Reviewer evaluates an Event and returns a Review.
// Implementations must be deterministic and side-effect-free.
// A Reviewer must never commit a Decision or execute an Action.
// Per RFC-020: reviews are advisory input for subsequent Decision processes.
type Reviewer interface {
	Review(ctx context.Context, event Event) (Review, error)
}

// DecisionMaker consumes a Review and produces a Decision.
// Implementations must be deterministic and policy-driven.
// A DecisionMaker must never execute Actions directly.
// Per RFC-021: decisions are explicit, immutable commitments.
type DecisionMaker interface {
	Decide(ctx context.Context, event Event, review Review) (Decision, error)
}

// ActionPlanner derives a set of planned Actions from a Decision.
// Implementations must be deterministic and side-effect-free.
// The returned Actions are plans only; execution is the responsibility
// of the calling service layer.
// Per RFC-023: planning, requesting, and execution are distinct.
type ActionPlanner interface {
	Plan(ctx context.Context, event Event, decision Decision) ([]Action, error)
}
