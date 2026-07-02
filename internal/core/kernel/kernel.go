package kernel

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventRecorder defines the minimal interface for recording pipeline execution events.
// This interface allows the kernel to optionally persist its execution trace
// without coupling to any specific storage backend.
type EventRecorder interface {
	// Append stores a new event record.
	// Implementations must be append-only and never modify or delete events.
	Append(ctx context.Context, event EventRecord) error
}

// EventRecord is the stored representation of a kernel execution event.
// This is a minimal subset of the full eventstore.EventRecord to avoid
// importing storage packages into the kernel.
type EventRecord struct {
	ID            string
	CorrelationID string
	CausationID   string
	TraceID       string
	Type          string
	Source        string
	SubjectID     string
	OccurredAt    time.Time
	Payload       json.RawMessage
	Metadata      map[string]string
	CreatedAt     time.Time
}

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
	eventRecorder EventRecorder // optional; nil means no event recording
}

// Option is a functional option for configuring a Kernel.
type Option func(*Kernel)

// WithEventRecorder configures the Kernel to record pipeline execution events.
// When provided, the Kernel will append an event record for each successful
// pipeline completion. Event recording failures do not prevent pipeline execution
// but are returned as errors from Run().
func WithEventRecorder(recorder EventRecorder) Option {
	return func(k *Kernel) {
		k.eventRecorder = recorder
	}
}

// New creates a Kernel wired with the provided pipeline components.
// All three components are required; passing nil causes a panic so that
// misconfiguration is caught at startup rather than at runtime.
//
// Optional configuration can be provided via Option functions:
//   - WithEventRecorder: enables event recording for pipeline executions
func New(reviewer Reviewer, decisionMaker DecisionMaker, planner ActionPlanner, opts ...Option) *Kernel {
	if reviewer == nil {
		panic("kernel: reviewer must not be nil")
	}
	if decisionMaker == nil {
		panic("kernel: decisionMaker must not be nil")
	}
	if planner == nil {
		panic("kernel: planner must not be nil")
	}

	k := &Kernel{
		reviewer:      reviewer,
		decisionMaker: decisionMaker,
		planner:       planner,
	}

	for _, opt := range opts {
		opt(k)
	}

	return k
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
//
// If an EventRecorder was configured via WithEventRecorder, the kernel will
// record a pipeline completion event after successful execution. Event recording
// failures are returned as errors but do not invalidate the pipeline result.
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

	result := PipelineResult{
		EventID:  event.ID,
		Review:   review,
		Decision: decision,
		Actions:  actions,
	}

	// Record pipeline execution event if recorder is configured
	if k.eventRecorder != nil {
		if err := k.recordPipelineCompletion(ctx, event, result); err != nil {
			return result, fmt.Errorf("kernel: failed to record pipeline event: %w", err)
		}
	}

	return result, nil
}

// recordPipelineCompletion creates and stores an event record for a completed pipeline execution.
func (k *Kernel) recordPipelineCompletion(ctx context.Context, event Event, result PipelineResult) error {
	// Create a minimal payload with the pipeline result
	payload := map[string]any{
		"event_id": result.EventID,
		"decision": map[string]any{
			"id":         result.Decision.ID,
			"outcome":    result.Decision.Outcome,
			"confidence": result.Decision.Confidence,
			"policy":     result.Decision.Policy,
		},
		"action_count": len(result.Actions),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal pipeline result payload: %w", err)
	}

	// Create metadata with additional context
	metadata := map[string]string{
		"decision_outcome": string(result.Decision.Outcome),
		"action_count":     fmt.Sprintf("%d", len(result.Actions)),
	}
	if result.Decision.Policy != "" {
		metadata["policy"] = result.Decision.Policy
	}

	record := EventRecord{
		ID:            uuid.New().String(),
		Type:          "kernel.pipeline.completed",
		Source:        "kernel",
		SubjectID:     result.Decision.ID,
		CorrelationID: event.CorrelationID,
		CausationID:   event.ID,
		TraceID:       event.TraceID,
		OccurredAt:    time.Now().UTC(),
		Payload:       payloadJSON,
		Metadata:      metadata,
		CreatedAt:     time.Now().UTC(),
	}

	return k.eventRecorder.Append(ctx, record)
}
