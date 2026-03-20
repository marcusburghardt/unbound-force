package impediment

import (
	"testing"
	"time"
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
