package sprint

import (
	"fmt"
	"testing"
)

func TestSprintStore_PlanAndReview(t *testing.T) {
	dir := t.TempDir()
	store := NewSprintStore(dir)

	items := []string{"BI-001", "BI-002", "BI-003"}
	state, err := store.Plan("Test sprint", 10.0, items)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if state.Status != "active" {
		t.Errorf("Status = %q, want active", state.Status)
	}
	if len(state.PlannedItems) != 3 {
		t.Errorf("PlannedItems = %d, want 3", len(state.PlannedItems))
	}

	// Simulate completing items
	state.CompletedItems = []string{"BI-001", "BI-002"}
	_ = store.Save(state)

	reviewed, err := store.Review(state.SprintName)
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	if reviewed.Status != "complete" {
		t.Errorf("Status = %q, want complete", reviewed.Status)
	}
	if reviewed.Velocity != 2.0 {
		t.Errorf("Velocity = %f, want 2.0", reviewed.Velocity)
	}
}

func TestSprintPlan_CapacityCalculation(t *testing.T) {
	dir := t.TempDir()
	store := NewSprintStore(dir)

	// More items than velocity allows
	items := make([]string, 20)
	for i := range items {
		items[i] = fmt.Sprintf("BI-%03d", i+1)
	}

	state, err := store.Plan("Capacity test", 10.0, items)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if len(state.PlannedItems) != 10 {
		t.Errorf("PlannedItems = %d, want 10 (capped by velocity)", len(state.PlannedItems))
	}
}
