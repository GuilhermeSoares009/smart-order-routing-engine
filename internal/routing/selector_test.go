package routing

import "testing"

func TestSelectTargetFallbackIsDeterministic(t *testing.T) {
	targets := []Target{
		{ID: "b", LatencyMs: 20, Availability: 0.1, Priority: 2},
		{ID: "a", LatencyMs: 20, Availability: 0.1, Priority: 1},
		{ID: "c", LatencyMs: 30, Availability: 0.2, Priority: 3},
	}

	first, err := SelectTarget(targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !first.Fallback {
		t.Fatalf("expected fallback decision")
	}

	for i := 0; i < 25; i++ {
		next, callErr := SelectTarget(targets)
		if callErr != nil {
			t.Fatalf("unexpected error in iteration %d: %v", i, callErr)
		}
		if next.Target.ID != first.Target.ID {
			t.Fatalf("non-deterministic target: got %s want %s", next.Target.ID, first.Target.ID)
		}
		if next.Reason != first.Reason {
			t.Fatalf("non-deterministic reason: got %s want %s", next.Reason, first.Reason)
		}
	}
}

func TestSelectTargetRejectsEmptyTargets(t *testing.T) {
	_, err := SelectTarget(nil)
	if err != ErrNoTargets {
		t.Fatalf("expected ErrNoTargets, got %v", err)
	}
}
