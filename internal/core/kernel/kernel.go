package kernel

import (
	"context"
	"fmt"
)

// Kernel is the immutable Praxis execution spine:
//
//	Event → Review → Decision → Action
//
// It is transport-agnostic, deterministic, and stateless. It does not read
// environment variables, use global state, call LLMs, publish events, or
// persist anything. All side effects are the responsibility of the caller.
//
// Kernel is safe for concurrent use: each Run call is fully isolated.
type Kernel struct {
	reviewer      Reviewer
	decisionMaker DecisionMaker
	planner       ActionPlanner
}

// New creates a Kernel wired with the provided pipeline components.
// All three components are required; passing nil causes a panic so that
// misconfiguration is caught at startup rather than at runtime.
func New(reviewer Reviewer, decisionMaker DecisionMaker, planner ActionPlanner) *Kernel {
	if reviewer == nil {
		panic("kernel: reviewer must not be nil")
	}
	if decisionMaker == nil {
		panic("kernel: decisionMaker must not be nil")
	}
	if planner == nil {
		panic("kernel: planner must not be nil")
	}
	return &Kernel{
		reviewer:      reviewer,
		decisionMaker: decisionMaker,
		planner:       planner,
	}
}

// Run executes the full Event → Review → Decision → Action pipeline.
//
// The pipeline is strictly sequential and deterministic:
//  1. Validate the incoming Event.
//  2. Pass the Event to the Reviewer; collect the Review.
//  3. Pass the Event + Review to the DecisionMaker; collect the Decision.
//  4. Pass the Event + Decision to the ActionPlanner; collect the Actions.
//
// Any error at any stage aborts the pipeline immediately.
// The returned PipelineResult is the complete, auditable trace of the run.
func (k *Kernel) Run(ctx context.Context, event Event) (PipelineResult, error) {
	if err := validateEvent(event); err != nil {
		return PipelineResult{}, err
	}

	review, err := k.reviewer.Review(ctx, event)
	if err != nil {
		return PipelineResult{}, fmt.Errorf("%w: %w", ErrReviewFailed, err)
	}

	decision, err := k.decisionMaker.Decide(ctx, event, review)
	if err != nil {
		return PipelineResult{}, fmt.Errorf("%w: %w", ErrDecisionFailed, err)
	}

	if err := validateDecision(decision); err != nil {
		return PipelineResult{}, err
	}

	actions, err := k.planner.Plan(ctx, event, decision)
	if err != nil {
		return PipelineResult{}, fmt.Errorf("%w: %w", ErrPlanFailed, err)
	}

	return PipelineResult{
		EventID:  event.ID,
		Review:   review,
		Decision: decision,
		Actions:  actions,
	}, nil
}
