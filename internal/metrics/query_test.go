package metrics

import (
	"testing"
	"time"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return NewStore(t.TempDir())
}

func writeFixtureSnapshots(t *testing.T, store *Store, snaps []MetricsSnapshot) {
	t.Helper()
	for _, s := range snaps {
		if err := store.WriteSnapshot(s); err != nil {
			t.Fatalf("WriteSnapshot: %v", err)
		}
	}
}

func TestQuery_Summary_ReturnsLatest(t *testing.T) {
	store := newTestStore(t)

	snaps := []MetricsSnapshot{
		{
			Timestamp: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			Velocity:  5.0,
		},
		{
			Timestamp: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			Velocity:  9.0,
		},
	}
	writeFixtureSnapshots(t, store, snaps)

	q := NewQuery(store)
	summary, err := q.Summary(365 * 24 * time.Hour)
	if err != nil {
		t.Fatalf("Summary: %v", err)
	}
	if summary.Velocity != 9.0 {
		t.Errorf("Velocity = %f, want 9.0 (most recent)", summary.Velocity)
	}
}

func TestQuery_Summary_IncludesHealth(t *testing.T) {
	store := newTestStore(t)

	snaps := []MetricsSnapshot{
		{
			Timestamp:        time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			Velocity:         10.0,
			DefectRate:       0.05,
			ReviewIterations: 1.5,
			BacklogHealth:    BacklogHealth{Total: 100, Ready: 85},
			FlowEfficiency:   75.0,
		},
	}
	writeFixtureSnapshots(t, store, snaps)

	q := NewQuery(store)
	summary, err := q.Summary(365 * 24 * time.Hour)
	if err != nil {
		t.Fatalf("Summary: %v", err)
	}
	if len(summary.HealthIndicators) == 0 {
		t.Error("HealthIndicators is empty, want populated")
	}
	if len(summary.HealthIndicators) != 5 {
		t.Errorf("got %d health indicators, want 5", len(summary.HealthIndicators))
	}
}

func TestQuery_Summary_NoData(t *testing.T) {
	store := newTestStore(t)
	q := NewQuery(store)

	_, err := q.Summary(24 * time.Hour)
	if err == nil {
		t.Fatal("expected error for empty store, got nil")
	}
}

func TestQuery_Velocity_AllSprints(t *testing.T) {
	store := newTestStore(t)

	snaps := []MetricsSnapshot{
		{Timestamp: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Velocity: 5.0},
		{Timestamp: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), Velocity: 8.0},
		{Timestamp: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Velocity: 12.0},
	}
	writeFixtureSnapshots(t, store, snaps)

	q := NewQuery(store)
	points, err := q.Velocity(0) // 0 means all sprints
	if err != nil {
		t.Fatalf("Velocity: %v", err)
	}
	if len(points) != 3 {
		t.Fatalf("got %d points, want 3", len(points))
	}
	if points[0].Velocity != 5.0 {
		t.Errorf("points[0].Velocity = %f, want 5.0", points[0].Velocity)
	}
	if points[2].Velocity != 12.0 {
		t.Errorf("points[2].Velocity = %f, want 12.0", points[2].Velocity)
	}
}

func TestQuery_Velocity_LimitSprints(t *testing.T) {
	store := newTestStore(t)

	snaps := []MetricsSnapshot{
		{Timestamp: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Velocity: 5.0},
		{Timestamp: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), Velocity: 8.0},
		{Timestamp: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Velocity: 12.0},
	}
	writeFixtureSnapshots(t, store, snaps)

	q := NewQuery(store)
	points, err := q.Velocity(2) // last 2 only
	if err != nil {
		t.Fatalf("Velocity: %v", err)
	}
	if len(points) != 2 {
		t.Fatalf("got %d points, want 2", len(points))
	}
	// Should be the last two: Sprint 2 (8.0) and Sprint 3 (12.0)
	if points[0].Velocity != 8.0 {
		t.Errorf("points[0].Velocity = %f, want 8.0", points[0].Velocity)
	}
	if points[1].Velocity != 12.0 {
		t.Errorf("points[1].Velocity = %f, want 12.0", points[1].Velocity)
	}
}

func TestQuery_CycleTime_ComputesStats(t *testing.T) {
	store := newTestStore(t)

	snaps := []MetricsSnapshot{
		{
			Timestamp: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			CycleTime: CycleTimeStats{Avg: 10.0},
		},
		{
			Timestamp: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			CycleTime: CycleTimeStats{Avg: 20.0},
		},
		{
			Timestamp: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			CycleTime: CycleTimeStats{Avg: 30.0},
		},
	}
	writeFixtureSnapshots(t, store, snaps)

	q := NewQuery(store)
	stats, err := q.CycleTime(365 * 24 * time.Hour)
	if err != nil {
		t.Fatalf("CycleTime: %v", err)
	}
	// Values fed to ComputeCycleTimeFromValues: [10, 20, 30]
	// Avg = 20, Median = 20 (p50 of [10,20,30] → index 1.0 → 20)
	if stats.Avg != 20.0 {
		t.Errorf("Avg = %f, want 20.0", stats.Avg)
	}
	if stats.Median != 20.0 {
		t.Errorf("Median = %f, want 20.0", stats.Median)
	}
	if stats.P90 <= 0 {
		t.Errorf("P90 = %f, want > 0", stats.P90)
	}
}

func TestQuery_Bottlenecks_SortedDescending(t *testing.T) {
	store := newTestStore(t)

	snaps := []MetricsSnapshot{
		{
			Timestamp:        time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			Velocity:         8.0,
			CycleTime:        CycleTimeStats{Avg: 48.0}, // 48 hours = 2 days
			LeadTime:         120.0,                     // 5 days
			ReviewIterations: 3.0,
		},
	}
	writeFixtureSnapshots(t, store, snaps)

	q := NewQuery(store)
	results, err := q.Bottlenecks()
	if err != nil {
		t.Fatalf("Bottlenecks: %v", err)
	}
	if len(results) < 2 {
		t.Fatalf("got %d bottlenecks, want at least 2", len(results))
	}

	// Verify sorted descending by AvgWaitDays
	for i := 1; i < len(results); i++ {
		if results[i].AvgWaitDays > results[i-1].AvgWaitDays {
			t.Errorf("bottlenecks not sorted descending: [%d]=%f > [%d]=%f",
				i, results[i].AvgWaitDays, i-1, results[i-1].AvgWaitDays)
		}
	}

	// Verify each result has a stage name
	for _, r := range results {
		if r.Stage == "" {
			t.Error("bottleneck stage is empty")
		}
	}
}

func TestQuery_Health_ProducesIndicators(t *testing.T) {
	store := newTestStore(t)

	snaps := []MetricsSnapshot{
		{
			Timestamp:        time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			Velocity:         10.0,
			DefectRate:       0.05,
			ReviewIterations: 1.5,
			BacklogHealth:    BacklogHealth{Total: 50, Ready: 40, Stale: 1},
			FlowEfficiency:   72.0,
		},
	}
	writeFixtureSnapshots(t, store, snaps)

	q := NewQuery(store)
	indicators, err := q.Health()
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	if len(indicators) == 0 {
		t.Fatal("expected non-empty indicators")
	}
	if len(indicators) != 5 {
		t.Errorf("got %d indicators, want 5", len(indicators))
	}

	// Verify each indicator has a dimension and valid status
	validStatuses := map[string]bool{"green": true, "yellow": true, "red": true}
	for _, ind := range indicators {
		if ind.Dimension == "" {
			t.Error("indicator dimension is empty")
		}
		if !validStatuses[ind.Status] {
			t.Errorf("indicator %q has invalid status %q", ind.Dimension, ind.Status)
		}
	}
}
