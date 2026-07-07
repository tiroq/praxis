package userfacts

import (
	"context"
	"errors"
	"testing"
)

func openTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	store, err := OpenStore(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("OpenStore(:memory:) failed: %v", err)
	}
	return store
}

func TestUserFactsStore_AppendAndListBySourceEvent(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	fact := NewCandidateFact("telegram-chat-42", "corr_1", "location", "Bangkok", 0.94, "evt_1", "msg_1")
	if err := store.Append(ctx, fact); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	facts, err := store.ListBySourceEvent(ctx, "evt_1")
	if err != nil {
		t.Fatalf("ListBySourceEvent() failed: %v", err)
	}
	if len(facts) != 1 {
		t.Fatalf("expected 1 fact, got %d", len(facts))
	}
	if facts[0].Type != "location" || facts[0].Value != "Bangkok" || facts[0].UserID != "telegram-chat-42" {
		t.Fatalf("unexpected fact: %#v", facts[0])
	}
	if facts[0].SourceMessageID != "msg_1" {
		t.Fatalf("expected SourceMessageID=msg_1, got %q", facts[0].SourceMessageID)
	}
	if facts[0].ValidationState != ValidationStateExtracted {
		t.Fatalf("expected ValidationState=extracted, got %q", facts[0].ValidationState)
	}
	if facts[0].ValidationUpdatedAt == "" {
		t.Fatal("expected ValidationUpdatedAt to be set")
	}
}

func TestUserFactsStore_RejectsInvalidConfidence(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	fact := NewCandidateFact("user_1", "corr_1", "location", "Bangkok", 1.5, "evt_1", "msg_1")
	err := store.Append(ctx, fact)
	if err == nil {
		t.Fatal("expected invalid confidence error")
	}
	var confidenceErr ErrInvalidConfidence
	if !errors.As(err, &confidenceErr) {
		t.Fatalf("expected ErrInvalidConfidence, got %T: %v", err, err)
	}
}

func TestUserFactsStore_DuplicateFactHandling(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	fact := NewCandidateFact("user_1", "corr_1", "preference", "Go", 0.8, "evt_1", "msg_1")
	fact.ID = "fact_duplicate"
	if err := store.Append(ctx, fact); err != nil {
		t.Fatalf("first Append() failed: %v", err)
	}

	err := store.Append(ctx, fact)
	if err == nil {
		t.Fatal("expected duplicate fact error")
	}
	var dupErr ErrDuplicateFact
	if !errors.As(err, &dupErr) {
		t.Fatalf("expected ErrDuplicateFact, got %T: %v", err, err)
	}
}

func TestUserFactsStore_TransitionValidationState(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	fact := NewCandidateFact("user_1", "corr_1", "preference", "Go", 0.8, "evt_1", "msg_1")
	if err := store.Append(ctx, fact); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	if err := store.TransitionValidationState(ctx, fact.ID, ValidationStateCorrelated, "fact-validator", "linked to similar facts"); err != nil {
		t.Fatalf("TransitionValidationState() failed: %v", err)
	}

	facts, err := store.ListBySourceEvent(ctx, "evt_1")
	if err != nil {
		t.Fatalf("ListBySourceEvent() failed: %v", err)
	}
	if len(facts) != 1 {
		t.Fatalf("expected 1 fact, got %d", len(facts))
	}
	if facts[0].ValidationState != ValidationStateCorrelated {
		t.Fatalf("expected state correlated, got %q", facts[0].ValidationState)
	}
}

func TestUserFactsStore_RejectsInvalidValidationTransition(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("store.Close() failed: %v", err)
		}
	}()

	fact := NewCandidateFact("user_1", "corr_1", "preference", "Go", 0.8, "evt_1", "msg_1")
	if err := store.Append(ctx, fact); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	err := store.TransitionValidationState(ctx, fact.ID, ValidationStateReviewed, "fact-validator", "skip correlation")
	if err == nil {
		t.Fatal("expected invalid transition error")
	}
	var transitionErr ErrInvalidValidationTransition
	if !errors.As(err, &transitionErr) {
		t.Fatalf("expected ErrInvalidValidationTransition, got %T: %v", err, err)
	}
}
