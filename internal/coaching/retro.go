package coaching

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// RetroStore manages retrospective records on the filesystem.
type RetroStore struct {
	Dir string
}

// NewRetroStore creates a new retrospective store.
func NewRetroStore(dir string) *RetroStore {
	return &RetroStore{Dir: dir}
}

// StartRetro creates a new retrospective record with pre-populated metrics.
func (s *RetroStore) StartRetro(date string, metricsData map[string]interface{}) (*RetroRecord, error) {
	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return nil, fmt.Errorf("create retros dir: %w", err)
	}

	record := &RetroRecord{
		Date:          date,
		DataPresented: metricsData,
	}

	return record, nil
}

// SaveRetro writes a retrospective record to disk.
func (s *RetroStore) SaveRetro(record *RetroRecord) error {
	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return fmt.Errorf("create retros dir: %w", err)
	}

	data, err := yaml.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal retro: %w", err)
	}

	content := fmt.Sprintf("---\n%s---\n\n%s\n", string(data), record.Notes)
	path := filepath.Join(s.Dir, record.Date+"-retro.md")
	return os.WriteFile(path, []byte(content), 0644)
}

// LoadRetro reads a retrospective record from disk.
func (s *RetroStore) LoadRetro(date string) (*RetroRecord, error) {
	path := filepath.Join(s.Dir, date+"-retro.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read retro %q: %w", path, err)
	}

	content := string(data)
	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid retro format in %q", path)
	}

	var record RetroRecord
	if err := yaml.Unmarshal([]byte(parts[1]), &record); err != nil {
		return nil, fmt.Errorf("parse retro frontmatter: %w", err)
	}
	record.Notes = strings.TrimSpace(parts[2])
	return &record, nil
}

// ListRetros returns all retrospective records, sorted by date descending.
func (s *RetroStore) ListRetros() ([]RetroRecord, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var records []RetroRecord
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "-retro.md") {
			continue
		}
		date := strings.TrimSuffix(e.Name(), "-retro.md")
		record, err := s.LoadRetro(date)
		if err != nil {
			continue
		}
		records = append(records, *record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Date > records[j].Date
	})
	return records, nil
}

// NextActionID returns the next auto-incrementing AI-NNN ID.
func NextActionID(retros []RetroRecord) string {
	maxNum := 0
	re := regexp.MustCompile(`^AI-(\d+)$`)
	for _, r := range retros {
		for _, ai := range r.ActionItems {
			matches := re.FindStringSubmatch(ai.ID)
			if len(matches) == 2 {
				n, err := strconv.Atoi(matches[1])
				if err == nil && n > maxNum {
					maxNum = n
				}
			}
		}
	}
	return fmt.Sprintf("AI-%03d", maxNum+1)
}

// ReviewPreviousActions finds all non-completed action items and marks stale ones.
func ReviewPreviousActions(retros []RetroRecord) []ActionItem {
	var items []ActionItem
	for _, r := range retros {
		for _, ai := range r.ActionItems {
			if ai.Status != "completed" {
				if ai.IsStale() {
					ai.Status = "stale"
				}
				items = append(items, ai)
			}
		}
	}
	return items
}
