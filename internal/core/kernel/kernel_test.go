package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// ---- test helpers / stubs ----

type stubReviewer struct {
	review Review
	err    error
}

func (s *stubReviewer) Review(_ context.Context, _ Event) (Review, error) {
	return s.review, s.err
}

type stubDecisionMaker struct {
	decision Decision
	err      error
}

func (s *stubDecisionMaker) Decide(_ context.Context, _ Event, _ Review) (Decision, error) {
	return s.decision, s.err
}

type stubPlanner struct {
	actions []Action
	err     error
}

func (s *stubPlanner) Plan(_ context.Context, _ Event, _ Decision) ([]Action, error) {
	return s.actions, s.err
}

// makeValidEvent returns a minimal valid Event for tests.
func makeValidEvent() Event {
	return Event{
		ID:         "evt-001",
		Text:       "please review this proposal",
		Confidence: 0.8,
	}
}

// makeApprovalReview returns a Review that drives an Approve decision.
func makeApprovalReview(eventID string) Review {
	return Review{
		ID:             "rev-evt-001",
		EventID:        eventID,
		Reviewer:       "test-reviewer",
		Recommendation: ReviewRecommendationApprove,
		Assessment: ReviewAssessment{
			Score:      0.85,
			Confidence: 0.8,
		},
	}
}

// makeApproveDecision returns a well-formed Decision for tests.
func makeApproveDecision(eventID, reviewID string) Decision {
	return Decision{
		ID:        "dec-evt-001",
		EventID:   eventID,
		ReviewIDs: []string{reviewID},
		Outcome:   DecisionOutcomeApprove,
		Actor:     "test-actor",
	}
}

// ---- Kernel.Run — happy path ----

func TestKernel_Run_HappyPath(t *testing.T) {
	evt := makeValidEvent()
	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)
	acts := []Action{{ID: "act-1", Type: ActionTypeNotify}}

	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{decision: dec}, &stubPlanner{actions: acts})

	result, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EventID != evt.ID {
		t.Errorf("result.EventID = %q; want %q", result.EventID, evt.ID)
	}
	if result.Review.ID != rev.ID {
		t.Errorf("result.Review.ID = %q; want %q", result.Review.ID, rev.ID)
	}
	if result.Decision.Outcome != DecisionOutcomeApprove {
		t.Errorf("result.Decision.Outcome = %q; want %q", result.Decision.Outcome, DecisionOutcomeApprove)
	}
	if len(result.Actions) != 1 {
		t.Errorf("len(result.Actions) = %d; want 1", len(result.Actions))
	}
}

// ---- Validation: invalid event ----

func TestKernel_Run_EmptyEventID(t *testing.T) {
	evt := makeValidEvent()
	evt.ID = ""

	k := New(&stubReviewer{}, &stubDecisionMaker{}, &stubPlanner{})
	_, err := k.Run(context.Background(), evt)
	if !errors.Is(err, ErrEmptyEventID) {
		t.Errorf("err = %v; want ErrEmptyEventID", err)
	}
}

func TestKernel_Run_EmptyEventText(t *testing.T) {
	evt := makeValidEvent()
	evt.Text = ""

	k := New(&stubReviewer{}, &stubDecisionMaker{}, &stubPlanner{})
	_, err := k.Run(context.Background(), evt)
	if !errors.Is(err, ErrEmptyEventText) {
		t.Errorf("err = %v; want ErrEmptyEventText", err)
	}
}

func TestKernel_Run_WhitespaceOnlyEventText(t *testing.T) {
	evt := makeValidEvent()
	evt.Text = "   "

	k := New(&stubReviewer{}, &stubDecisionMaker{}, &stubPlanner{})
	_, err := k.Run(context.Background(), evt)
	if !errors.Is(err, ErrEmptyEventText) {
		t.Errorf("err = %v; want ErrEmptyEventText", err)
	}
}

func TestKernel_Run_ConfidenceBelowZero(t *testing.T) {
	evt := makeValidEvent()
	evt.Confidence = -0.1

	k := New(&stubReviewer{}, &stubDecisionMaker{}, &stubPlanner{})
	_, err := k.Run(context.Background(), evt)
	if !errors.Is(err, ErrInvalidConfidence) {
		t.Errorf("err = %v; want ErrInvalidConfidence", err)
	}
}

func TestKernel_Run_ConfidenceAboveOne(t *testing.T) {
	evt := makeValidEvent()
	evt.Confidence = 1.01

	k := New(&stubReviewer{}, &stubDecisionMaker{}, &stubPlanner{})
	_, err := k.Run(context.Background(), evt)
	if !errors.Is(err, ErrInvalidConfidence) {
		t.Errorf("err = %v; want ErrInvalidConfidence", err)
	}
}

func TestKernel_Run_ConfidenceBoundaryZero(t *testing.T) {
	evt := makeValidEvent()
	evt.Confidence = 0.0

	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)

	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{decision: dec}, &stubPlanner{actions: []Action{{ID: "act-1"}}})
	_, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Errorf("unexpected error for confidence=0.0: %v", err)
	}
}

func TestKernel_Run_ConfidenceBoundaryOne(t *testing.T) {
	evt := makeValidEvent()
	evt.Confidence = 1.0

	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)

	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{decision: dec}, &stubPlanner{actions: []Action{{ID: "act-1"}}})
	_, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Errorf("unexpected error for confidence=1.0: %v", err)
	}
}

// ---- Missing information / low confidence ----

func TestKernel_Run_LowConfidence(t *testing.T) {
	evt := makeValidEvent()
	evt.Confidence = 0.1

	// Reviewer returns a review with low confidence.
	rev := Review{
		ID:             "rev-001",
		EventID:        evt.ID,
		Reviewer:       "test",
		Recommendation: ReviewRecommendationInformational,
		Assessment: ReviewAssessment{
			Score:      0.2,
			Confidence: 0.1,
		},
	}
	dec := Decision{
		ID:      "dec-001",
		EventID: evt.ID,
		Outcome: DecisionOutcomeNeedsRevision,
		Actor:   "test",
	}

	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{decision: dec}, &stubPlanner{actions: []Action{{ID: "act-1"}}})
	result, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Decision.Outcome != DecisionOutcomeNeedsRevision {
		t.Errorf("outcome = %q; want NeedsRevision", result.Decision.Outcome)
	}
}

// ---- Multiple actions ----

func TestKernel_Run_MultipleActions(t *testing.T) {
	evt := makeValidEvent()
	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)
	acts := []Action{
		{ID: "act-1", Type: ActionTypeNotify},
		{ID: "act-2", Type: ActionTypeSchedule},
		{ID: "act-3", Type: ActionTypeCreate},
	}

	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{decision: dec}, &stubPlanner{actions: acts})
	result, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Actions) != 3 {
		t.Errorf("len(result.Actions) = %d; want 3", len(result.Actions))
	}
}

// ---- Error propagation ----

var errSentinel = errors.New("test error")

func TestKernel_Run_ReviewerError(t *testing.T) {
	k := New(&stubReviewer{err: errSentinel}, &stubDecisionMaker{}, &stubPlanner{})
	_, err := k.Run(context.Background(), makeValidEvent())

	if !errors.Is(err, ErrReviewFailed) {
		t.Errorf("err = %v; want ErrReviewFailed wrapper", err)
	}
	if !errors.Is(err, errSentinel) {
		t.Errorf("err = %v; want original errSentinel wrapped inside", err)
	}
}

func TestKernel_Run_DecisionMakerError(t *testing.T) {
	rev := makeApprovalReview("evt-001")
	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{err: errSentinel}, &stubPlanner{})
	_, err := k.Run(context.Background(), makeValidEvent())

	if !errors.Is(err, ErrDecisionFailed) {
		t.Errorf("err = %v; want ErrDecisionFailed wrapper", err)
	}
	if !errors.Is(err, errSentinel) {
		t.Errorf("err = %v; want original errSentinel wrapped inside", err)
	}
}

func TestKernel_Run_PlannerError(t *testing.T) {
	rev := makeApprovalReview("evt-001")
	dec := makeApproveDecision("evt-001", rev.ID)
	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{decision: dec}, &stubPlanner{err: errSentinel})
	_, err := k.Run(context.Background(), makeValidEvent())

	if !errors.Is(err, ErrPlanFailed) {
		t.Errorf("err = %v; want ErrPlanFailed wrapper", err)
	}
	if !errors.Is(err, errSentinel) {
		t.Errorf("err = %v; want original errSentinel wrapped inside", err)
	}
}

// ---- Decision validation: empty outcome ----

func TestKernel_Run_EmptyDecisionOutcome(t *testing.T) {
	rev := makeApprovalReview("evt-001")
	dec := Decision{
		ID:      "dec-001",
		EventID: "evt-001",
		Outcome: "", // intentionally empty
		Actor:   "test",
	}

	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{decision: dec}, &stubPlanner{})
	_, err := k.Run(context.Background(), makeValidEvent())
	if !errors.Is(err, ErrEmptyDecisionOutcome) {
		t.Errorf("err = %v; want ErrEmptyDecisionOutcome", err)
	}
}

// ---- Default implementations ----

func TestKeywordReviewer_NoKeywordsMatched(t *testing.T) {
	r := NewKeywordReviewer(map[string]string{"urgent": "time-sensitive"}, "")
	evt := makeValidEvent()
	evt.Text = "nothing relevant here"

	rev, err := r.Review(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rev.Recommendation != ReviewRecommendationInformational {
		t.Errorf("recommendation = %q; want informational", rev.Recommendation)
	}
	if len(rev.Assessment.Findings) != 0 {
		t.Errorf("expected no findings, got %v", rev.Assessment.Findings)
	}
}

func TestKeywordReviewer_KeywordsMatched(t *testing.T) {
	r := NewKeywordReviewer(map[string]string{"urgent": "time-sensitive", "blocked": "blocker"}, "")
	evt := makeValidEvent()
	evt.Text = "this is urgent and blocked"
	evt.Confidence = 0.9

	rev, err := r.Review(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rev.Recommendation != ReviewRecommendationApprove {
		t.Errorf("recommendation = %q; want approve", rev.Recommendation)
	}
	if len(rev.Assessment.Findings) == 0 {
		t.Error("expected findings, got none")
	}
}

func TestKeywordReviewer_CaseInsensitive(t *testing.T) {
	r := NewKeywordReviewer(map[string]string{"urgent": "time-sensitive"}, "")
	evt := makeValidEvent()
	evt.Text = "URGENT task needs attention"
	evt.Confidence = 0.7

	rev, err := r.Review(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rev.Recommendation != ReviewRecommendationApprove {
		t.Errorf("recommendation = %q; want approve (case-insensitive match)", rev.Recommendation)
	}
}

func TestRuleBasedDecisionMaker_Approve(t *testing.T) {
	dm := NewRuleBasedDecisionMaker(DefaultConfidenceThreshold, "", "")
	evt := makeValidEvent()
	rev := Review{
		ID:             "rev-1",
		EventID:        evt.ID,
		Recommendation: ReviewRecommendationApprove,
		Assessment:     ReviewAssessment{Confidence: 0.8},
	}

	dec, err := dm.Decide(context.Background(), evt, rev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Outcome != DecisionOutcomeApprove {
		t.Errorf("outcome = %q; want approve", dec.Outcome)
	}
}

func TestRuleBasedDecisionMaker_Reject(t *testing.T) {
	dm := NewRuleBasedDecisionMaker(DefaultConfidenceThreshold, "", "")
	evt := makeValidEvent()
	rev := Review{
		ID:             "rev-1",
		EventID:        evt.ID,
		Recommendation: ReviewRecommendationReject,
		Assessment:     ReviewAssessment{Confidence: 0.9},
	}

	dec, err := dm.Decide(context.Background(), evt, rev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Outcome != DecisionOutcomeReject {
		t.Errorf("outcome = %q; want reject", dec.Outcome)
	}
}

func TestRuleBasedDecisionMaker_Escalate(t *testing.T) {
	dm := NewRuleBasedDecisionMaker(DefaultConfidenceThreshold, "", "")
	evt := makeValidEvent()
	rev := Review{
		ID:             "rev-1",
		EventID:        evt.ID,
		Recommendation: ReviewRecommendationEscalate,
		Assessment:     ReviewAssessment{Confidence: 0.7},
	}

	dec, err := dm.Decide(context.Background(), evt, rev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Outcome != DecisionOutcomeEscalate {
		t.Errorf("outcome = %q; want escalate", dec.Outcome)
	}
}

func TestRuleBasedDecisionMaker_Defer(t *testing.T) {
	dm := NewRuleBasedDecisionMaker(DefaultConfidenceThreshold, "", "")
	evt := makeValidEvent()
	rev := Review{
		ID:             "rev-1",
		EventID:        evt.ID,
		Recommendation: ReviewRecommendationInformational,
		Assessment:     ReviewAssessment{Confidence: 0.4}, // between Defer and Approve
	}

	dec, err := dm.Decide(context.Background(), evt, rev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Outcome != DecisionOutcomeDefer {
		t.Errorf("outcome = %q; want defer", dec.Outcome)
	}
}

func TestRuleBasedDecisionMaker_NeedsRevision(t *testing.T) {
	dm := NewRuleBasedDecisionMaker(DefaultConfidenceThreshold, "", "")
	evt := makeValidEvent()
	rev := Review{
		ID:             "rev-1",
		EventID:        evt.ID,
		Recommendation: ReviewRecommendationInformational,
		Assessment:     ReviewAssessment{Confidence: 0.1}, // below Defer threshold
	}

	dec, err := dm.Decide(context.Background(), evt, rev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Outcome != DecisionOutcomeNeedsRevision {
		t.Errorf("outcome = %q; want needs_revision", dec.Outcome)
	}
}

func TestSimpleActionPlanner_KnownOutcome(t *testing.T) {
	p := NewSimpleActionPlanner(nil) // uses defaultActionRules
	evt := makeValidEvent()
	dec := makeApproveDecision(evt.ID, "rev-1")

	acts, err := p.Plan(context.Background(), evt, dec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(acts) == 0 {
		t.Error("expected at least one action for approve outcome")
	}
	if acts[0].Type != ActionTypeNotify {
		t.Errorf("action type = %q; want notify", acts[0].Type)
	}
}

func TestSimpleActionPlanner_UnknownOutcome(t *testing.T) {
	p := NewSimpleActionPlanner(nil)
	evt := makeValidEvent()
	dec := Decision{
		ID:      "dec-x",
		EventID: evt.ID,
		Outcome: DecisionOutcome("unknown_outcome"),
		Actor:   "test",
	}

	acts, err := p.Plan(context.Background(), evt, dec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(acts) != 1 {
		t.Fatalf("expected 1 no-op action, got %d", len(acts))
	}
	if acts[0].Type != ActionTypeNoOp {
		t.Errorf("action type = %q; want noop", acts[0].Type)
	}
}

// ---- End-to-end with default implementations ----

func TestKernel_EndToEnd_DefaultImpls(t *testing.T) {
	reviewer := NewKeywordReviewer(map[string]string{"review": "needs-review"}, "")
	dm := NewRuleBasedDecisionMaker(DefaultConfidenceThreshold, "", "")
	planner := NewSimpleActionPlanner(nil)

	k := New(reviewer, dm, planner)

	evt := Event{
		ID:         "e2e-001",
		Text:       "please review this proposal",
		Confidence: 0.75,
	}

	result, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EventID != evt.ID {
		t.Errorf("result.EventID = %q; want %q", result.EventID, evt.ID)
	}
	if result.Decision.Outcome == "" {
		t.Error("decision outcome must not be empty")
	}
	if len(result.Actions) == 0 {
		t.Error("expected at least one action")
	}
}

func TestKernel_EndToEnd_LowConfidence(t *testing.T) {
	reviewer := NewKeywordReviewer(nil, "")
	dm := NewRuleBasedDecisionMaker(DefaultConfidenceThreshold, "", "")
	planner := NewSimpleActionPlanner(nil)

	k := New(reviewer, dm, planner)

	evt := Event{
		ID:         "e2e-low-001",
		Text:       "nothing here",
		Confidence: 0.05, // below defer threshold
	}

	result, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Decision.Outcome != DecisionOutcomeNeedsRevision {
		t.Errorf("outcome = %q; want needs_revision", result.Decision.Outcome)
	}
}

func TestNew_PanicsOnNilReviewer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil reviewer")
		}
	}()
	New(nil, &stubDecisionMaker{}, &stubPlanner{})
}

func TestNew_PanicsOnNilDecisionMaker(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil decisionMaker")
		}
	}()
	New(&stubReviewer{}, nil, &stubPlanner{})
}

func TestNew_PanicsOnNilPlanner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil planner")
		}
	}()
	New(&stubReviewer{}, &stubDecisionMaker{}, nil)
}

// ---- Event Recording / DI tests ----

type stubEventRecorder struct {
	records []EventRecord
	err     error
}

func (s *stubEventRecorder) Append(_ context.Context, event EventRecord) error {
	if s.err != nil {
		return s.err
	}
	s.records = append(s.records, event)
	return nil
}

func TestKernel_WithoutEventRecorder_BehavesAsUsual(t *testing.T) {
	evt := makeValidEvent()
	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)
	acts := []Action{{ID: "act-1", Type: ActionTypeNotify}}

	// Create kernel WITHOUT event recorder option
	k := New(&stubReviewer{review: rev}, &stubDecisionMaker{decision: dec}, &stubPlanner{actions: acts})

	result, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EventID != evt.ID {
		t.Errorf("result.EventID = %q; want %q", result.EventID, evt.ID)
	}
}

func TestKernel_WithEventRecorder_RecordsEvent(t *testing.T) {
	evt := makeValidEvent()
	evt.CorrelationID = "corr-123"
	evt.TraceID = "trace-456"

	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)
	dec.Policy = "test-policy"
	acts := []Action{{ID: "act-1", Type: ActionTypeNotify}, {ID: "act-2", Type: ActionTypeCreate}}

	recorder := &stubEventRecorder{}
	k := New(
		&stubReviewer{review: rev},
		&stubDecisionMaker{decision: dec},
		&stubPlanner{actions: acts},
		WithEventRecorder(recorder),
	)

	result, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify event was recorded
	if len(recorder.records) != 1 {
		t.Fatalf("expected 1 recorded event, got %d", len(recorder.records))
	}

	record := recorder.records[0]

	// Verify event type
	if record.Type != "kernel.pipeline.completed" {
		t.Errorf("record.Type = %q; want %q", record.Type, "kernel.pipeline.completed")
	}

	// Verify source
	if record.Source != "kernel" {
		t.Errorf("record.Source = %q; want %q", record.Source, "kernel")
	}

	// Verify subject is decision ID
	if record.SubjectID != dec.ID {
		t.Errorf("record.SubjectID = %q; want %q", record.SubjectID, dec.ID)
	}

	// Verify correlation/causation/trace mapping
	if record.CorrelationID != evt.CorrelationID {
		t.Errorf("record.CorrelationID = %q; want %q", record.CorrelationID, evt.CorrelationID)
	}
	if record.CausationID != evt.ID {
		t.Errorf("record.CausationID = %q; want %q (input event ID)", record.CausationID, evt.ID)
	}
	if record.TraceID != evt.TraceID {
		t.Errorf("record.TraceID = %q; want %q", record.TraceID, evt.TraceID)
	}

	// Verify metadata
	if record.Metadata["decision_outcome"] != string(dec.Outcome) {
		t.Errorf("metadata[decision_outcome] = %q; want %q", record.Metadata["decision_outcome"], dec.Outcome)
	}
	if record.Metadata["action_count"] != "2" {
		t.Errorf("metadata[action_count] = %q; want %q", record.Metadata["action_count"], "2")
	}
	if record.Metadata["policy"] != "test-policy" {
		t.Errorf("metadata[policy] = %q; want %q", record.Metadata["policy"], "test-policy")
	}

	// Verify ID is generated
	if record.ID == "" {
		t.Error("record.ID must not be empty")
	}

	// Verify timestamps are set
	if record.OccurredAt.IsZero() {
		t.Error("record.OccurredAt must be set")
	}
	if record.CreatedAt.IsZero() {
		t.Error("record.CreatedAt must be set")
	}

	// Verify result is still returned correctly
	if result.EventID != evt.ID {
		t.Errorf("result.EventID = %q; want %q", result.EventID, evt.ID)
	}
}

func TestKernel_EventRecord_HasValidJSON(t *testing.T) {
	evt := makeValidEvent()
	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)
	acts := []Action{{ID: "act-1", Type: ActionTypeNotify}}

	recorder := &stubEventRecorder{}
	k := New(
		&stubReviewer{review: rev},
		&stubDecisionMaker{decision: dec},
		&stubPlanner{actions: acts},
		WithEventRecorder(recorder),
	)

	_, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recorder.records) != 1 {
		t.Fatalf("expected 1 recorded event, got %d", len(recorder.records))
	}

	record := recorder.records[0]

	// Verify payload is valid JSON
	var payload map[string]interface{}
	if err := json.Unmarshal(record.Payload, &payload); err != nil {
		t.Fatalf("payload is not valid JSON: %v", err)
	}

	// Verify payload structure
	if payload["event_id"] != evt.ID {
		t.Errorf("payload[event_id] = %v; want %q", payload["event_id"], evt.ID)
	}
	if payload["action_count"] != float64(1) {
		t.Errorf("payload[action_count] = %v; want 1", payload["action_count"])
	}

	decision, ok := payload["decision"].(map[string]interface{})
	if !ok {
		t.Fatal("payload[decision] is not a map")
	}
	if decision["id"] != dec.ID {
		t.Errorf("payload[decision][id] = %v; want %q", decision["id"], dec.ID)
	}
	if decision["outcome"] != string(dec.Outcome) {
		t.Errorf("payload[decision][outcome] = %v; want %q", decision["outcome"], dec.Outcome)
	}
}

func TestKernel_EventRecorderError_ReturnsError(t *testing.T) {
	evt := makeValidEvent()
	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)
	acts := []Action{{ID: "act-1", Type: ActionTypeNotify}}

	recorderErr := errors.New("storage unavailable")
	recorder := &stubEventRecorder{err: recorderErr}
	k := New(
		&stubReviewer{review: rev},
		&stubDecisionMaker{decision: dec},
		&stubPlanner{actions: acts},
		WithEventRecorder(recorder),
	)

	result, err := k.Run(context.Background(), evt)

	// Pipeline execution should complete successfully
	if result.EventID != evt.ID {
		t.Errorf("result.EventID = %q; want %q", result.EventID, evt.ID)
	}

	// But error should be returned for event recording failure
	if err == nil {
		t.Fatal("expected error for event recording failure, got nil")
	}
	if !errors.Is(err, recorderErr) {
		t.Errorf("expected error to wrap recorderErr, got %v", err)
	}
}

func TestKernel_EventRecorder_EmptyCorrelationID(t *testing.T) {
	evt := makeValidEvent()
	evt.CorrelationID = "" // empty
	evt.TraceID = ""       // empty

	rev := makeApprovalReview(evt.ID)
	dec := makeApproveDecision(evt.ID, rev.ID)
	acts := []Action{{ID: "act-1", Type: ActionTypeNotify}}

	recorder := &stubEventRecorder{}
	k := New(
		&stubReviewer{review: rev},
		&stubDecisionMaker{decision: dec},
		&stubPlanner{actions: acts},
		WithEventRecorder(recorder),
	)

	_, err := k.Run(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recorder.records) != 1 {
		t.Fatalf("expected 1 recorded event, got %d", len(recorder.records))
	}

	record := recorder.records[0]

	// Empty IDs should be passed through as-is
	if record.CorrelationID != "" {
		t.Errorf("record.CorrelationID = %q; want empty", record.CorrelationID)
	}
	if record.TraceID != "" {
		t.Errorf("record.TraceID = %q; want empty", record.TraceID)
	}
	// Causation should still be set to event ID
	if record.CausationID != evt.ID {
		t.Errorf("record.CausationID = %q; want %q", record.CausationID, evt.ID)
	}
}
