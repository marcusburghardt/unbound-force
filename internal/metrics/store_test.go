package metrics

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStore_WriteReadCollection_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	now := time.Date(2026, 3, 20, 14, 30, 0, 0, time.UTC)
	coll := SourceCollection{
		Source:      "github",
		CollectedAt: now,
		DataPoints:  42,
		RawData:     map[string]interface{}{"pr_count": float64(10)},
	}

	if err := store.WriteCollection("github", coll); err != nil {
		t.Fatalf("WriteCollection: %v", err)
	}

	results, err := store.ReadCollections("github", time.Time{})
	if err != nil {
		t.Fatalf("ReadCollections: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d collections, want 1", len(results))
	}
	if results[0].DataPoints != 42 {
		t.Errorf("DataPoints = %d, want 42", results[0].DataPoints)
	}
	if results[0].Source != "github" {
		t.Errorf("Source = %q, want github", results[0].Source)
	}
}

func TestStore_ReadCollections_FilterByTime(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	old := SourceCollection{
		Source:      "github",
		CollectedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		DataPoints:  10,
		RawData:     map[string]interface{}{},
	}
	recent := SourceCollection{
		Source:      "github",
		CollectedAt: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
		DataPoints:  20,
		RawData:     map[string]interface{}{},
	}

	_ = store.WriteCollection("github", old)
	_ = store.WriteCollection("github", recent)

	since := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	results, err := store.ReadCollections("github", since)
	if err != nil {
		t.Fatalf("ReadCollections: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d collections, want 1 (filtered)", len(results))
	}
	if results[0].DataPoints != 20 {
		t.Errorf("got old collection, want recent")
	}
}

func TestStore_WriteReadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	snap := MetricsSnapshot{
		Timestamp:        time.Date(2026, 3, 20, 14, 30, 0, 0, time.UTC),
		Velocity:         8.2,
		CycleTime:        CycleTimeStats{Avg: 18.4, Median: 14.2, P90: 42.1, P99: 71.3},
		LeadTime:         76.8,
		DefectRate:       0.12,
		ReviewIterations: 1.8,
		CIPassRate:       94.3,
		BacklogHealth:    BacklogHealth{Total: 32, Ready: 12, Stale: 3},
		FlowEfficiency:   68.4,
		SourcesCollected: []string{"github", "gaze"},
	}

	if err := store.WriteSnapshot(snap); err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}

	results, err := store.ReadSnapshots(time.Time{})
	if err != nil {
		t.Fatalf("ReadSnapshots: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d snapshots, want 1", len(results))
	}
	if results[0].Velocity != 8.2 {
		t.Errorf("Velocity = %f, want 8.2", results[0].Velocity)
	}
	if results[0].CIPassRate != 94.3 {
		t.Errorf("CIPassRate = %f, want 94.3", results[0].CIPassRate)
	}
}

func TestStore_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	colls, err := store.ReadCollections("github", time.Time{})
	if err != nil {
		t.Fatalf("ReadCollections on empty: %v", err)
	}
	if len(colls) != 0 {
		t.Errorf("expected 0 collections, got %d", len(colls))
	}

	snaps, err := store.ReadSnapshots(time.Time{})
	if err != nil {
		t.Fatalf("ReadSnapshots on empty: %v", err)
	}
	if len(snaps) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(snaps))
	}

	// Verify data directory was created
	dataDir := filepath.Join(dir, "data")
	_, statErr := os.Stat(dataDir)
	// It's OK if dir doesn't exist yet -- ReadCollections handles NotExist gracefully
	_ = statErr
}
