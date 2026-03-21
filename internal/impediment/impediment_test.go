package impediment

import (
	"strings"
	"testing"
	"time"

	"github.com/unbound-force/unbound-force/internal/metrics"
)

func TestRepository_Add_AutoID(t *testing.T) {
	dir := t.TempDir()
	repo := NewRepository(dir)
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	imp, err := repo.Add("Test impediment", "high", "@dev", "Description", now)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if imp.ID != "IMP-001" {
		t.Errorf("ID = %q, want IMP-001", imp.ID)
	}

	imp2, err := repo.Add("Second impediment", "medium", "", "Desc 2", now)
	if err != nil {
		t.Fatalf("Add second: %v", err)
	}
	if imp2.ID != "IMP-002" {
		t.Errorf("ID = %q, want IMP-002", imp2.ID)
	}
}

func TestRepository_List_SortBySeverity(t *testing.T) {
	dir := t.TempDir()
	repo := NewRepository(dir)
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	_, _ = repo.Add("Low issue", "low", "@dev", "", now)
	_, _ = repo.Add("Critical issue", "critical", "@lead", "", now)
	_, _ = repo.Add("High issue", "high", "@dev", "", now)

	list, err := repo.List("all")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("got %d impediments, want 3", len(list))
	}
	if list[0].Severity != "critical" {
		t.Errorf("first = %q, want critical", list[0].Severity)
	}
	if list[1].Severity != "high" {
		t.Errorf("second = %q, want high", list[1].Severity)
	}
	if list[2].Severity != "low" {
		t.Errorf("third = %q, want low", list[2].Severity)
	}
}

func TestRepository_Resolve_UpdatesFile(t *testing.T) {
	dir := t.TempDir()
	repo := NewRepository(dir)
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	imp, _ := repo.Add("Flaky CI", "high", "@dev", "CI is flaky", now)

	resolveTime := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)
	if err := repo.Resolve(imp.ID, "Pinned base image", resolveTime); err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	resolved, err := repo.Get(imp.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if resolved.Status != "resolved" {
		t.Errorf("Status = %q, want resolved", resolved.Status)
	}
	if resolved.Resolution != "Pinned base image" {
		t.Errorf("Resolution = %q", resolved.Resolution)
	}
}

func TestRepository_List_StaleDetection(t *testing.T) {
	dir := t.TempDir()
	repo := NewRepository(dir)
	old := time.Now().Add(-15 * 24 * time.Hour) // 15 days ago

	imp, _ := repo.Add("Old issue", "medium", "@dev", "", old)

	got, err := repo.Get(imp.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !got.IsStale() {
		t.Error("expected impediment to be stale (>14 days)")
	}
}

func TestRepository_Add_NoOwner(t *testing.T) {
	dir := t.TempDir()
	repo := NewRepository(dir)
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	imp, err := repo.Add("Unowned", "low", "", "No owner", now)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if imp.Owner != "unassigned" {
		t.Errorf("Owner = %q, want unassigned", imp.Owner)
	}
}

// writeSnapshot is a test helper that writes a MetricsSnapshot to the store.
func writeSnapshot(t *testing.T, store *metrics.Store, ts time.Time, velocity, ciPassRate, reviewIter float64) {
	t.Helper()
	snap := metrics.MetricsSnapshot{
		Timestamp:        ts,
		Velocity:         velocity,
		CIPassRate:       ciPassRate,
		ReviewIterations: reviewIter,
		CycleTime:        metrics.CycleTimeStats{Avg: 24, Median: 18},
		FlowEfficiency:   65,
	}
	if err := store.WriteSnapshot(snap); err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}
}

func TestDetect_CISpike(t *testing.T) {
	tmp := t.TempDir()
	store := metrics.NewStore(tmp)
	repo := NewRepository(tmp + "/impediments")
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	// Two snapshots: CI drops from 95% to 70% (25% drop > 15% threshold)
	writeSnapshot(t, store, now.Add(-24*time.Hour), 10, 95, 2.0)
	writeSnapshot(t, store, now, 10, 70, 2.0)

	detected, err := Detect(store, repo, now)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if len(detected) == 0 {
		t.Fatal("expected at least one impediment for CI spike")
	}

	found := false
	for _, imp := range detected {
		if strings.Contains(strings.ToLower(imp.Title), "ci pass rate") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected impediment with 'CI pass rate' in title, got: %v", detected)
	}
}

func TestDetect_ReviewTurnaroundIncrease(t *testing.T) {
	tmp := t.TempDir()
	store := metrics.NewStore(tmp)
	repo := NewRepository(tmp + "/impediments")
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	// 3 snapshots: review iterations stable at 2.0, then jumps to 4.0
	// Rolling avg of first 2 = 2.0; 4.0 > 2.0*1.5=3.0 → triggers
	writeSnapshot(t, store, now.Add(-48*time.Hour), 10, 90, 2.0)
	writeSnapshot(t, store, now.Add(-24*time.Hour), 10, 90, 2.0)
	writeSnapshot(t, store, now, 10, 90, 4.0)

	detected, err := Detect(store, repo, now)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if len(detected) == 0 {
		t.Fatal("expected at least one impediment for review turnaround increase")
	}

	found := false
	for _, imp := range detected {
		if strings.Contains(strings.ToLower(imp.Title), "review iterations") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected impediment with 'review iterations' in title, got titles: %v",
			func() []string {
				var titles []string
				for _, imp := range detected {
					titles = append(titles, imp.Title)
				}
				return titles
			}())
	}
}

func TestDetect_VelocityDrop(t *testing.T) {
	tmp := t.TempDir()
	store := metrics.NewStore(tmp)
	repo := NewRepository(tmp + "/impediments")
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	// 3 snapshots: velocity stable at 20, then drops to 10
	// Rolling avg of first 2 = 20; 10 < 20*0.75=15 → triggers
	writeSnapshot(t, store, now.Add(-48*time.Hour), 20, 90, 2.0)
	writeSnapshot(t, store, now.Add(-24*time.Hour), 20, 90, 2.0)
	writeSnapshot(t, store, now, 10, 90, 2.0)

	detected, err := Detect(store, repo, now)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if len(detected) == 0 {
		t.Fatal("expected at least one impediment for velocity drop")
	}

	found := false
	for _, imp := range detected {
		if strings.Contains(strings.ToLower(imp.Title), "velocity") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected impediment with 'velocity' in title, got titles: %v",
			func() []string {
				var titles []string
				for _, imp := range detected {
					titles = append(titles, imp.Title)
				}
				return titles
			}())
	}
}

func TestDetect_NoAnomalies(t *testing.T) {
	tmp := t.TempDir()
	store := metrics.NewStore(tmp)
	repo := NewRepository(tmp + "/impediments")
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	// 2 stable snapshots: no spikes, no drops
	writeSnapshot(t, store, now.Add(-24*time.Hour), 15, 90, 2.0)
	writeSnapshot(t, store, now, 15, 90, 2.0)

	detected, err := Detect(store, repo, now)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if len(detected) != 0 {
		t.Errorf("expected no impediments for stable metrics, got %d: %v", len(detected), detected)
	}
}

func TestDetect_InsufficientData(t *testing.T) {
	tmp := t.TempDir()
	store := metrics.NewStore(tmp)
	repo := NewRepository(tmp + "/impediments")
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	// Only 1 snapshot — should return an error
	writeSnapshot(t, store, now, 15, 90, 2.0)

	_, err := Detect(store, repo, now)
	if err == nil {
		t.Fatal("expected error for insufficient data")
	}
	if !strings.Contains(err.Error(), "insufficient") {
		t.Errorf("error should mention 'insufficient', got: %v", err)
	}
}
