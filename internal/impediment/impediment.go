package impediment

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Repository manages impediment records on the filesystem.
type Repository struct {
	Dir string
}

// NewRepository creates a new impediment repository.
func NewRepository(dir string) *Repository {
	return &Repository{Dir: dir}
}

func (r *Repository) ensureDir() error {
	return os.MkdirAll(r.Dir, 0755)
}

// Add creates a new impediment with an auto-assigned ID.
func (r *Repository) Add(title, severity, owner, description string, now time.Time) (*Impediment, error) {
	if err := r.ensureDir(); err != nil {
		return nil, fmt.Errorf("create impediment dir: %w", err)
	}

	id, err := r.nextID()
	if err != nil {
		return nil, err
	}

	if owner == "" {
		owner = "unassigned"
	}

	imp := &Impediment{
		ID:          id,
		Title:       title,
		Severity:    severity,
		Owner:       owner,
		Status:      "open",
		CreatedAt:   now,
		Source:      "manual",
		Description: description,
	}

	if err := r.save(imp); err != nil {
		return nil, err
	}
	return imp, nil
}

// List returns impediments filtered by status, sorted by severity.
func (r *Repository) List(statusFilter string) ([]Impediment, error) {
	entries, err := os.ReadDir(r.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read impediments dir: %w", err)
	}

	var results []Impediment
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		imp, err := r.load(filepath.Join(r.Dir, e.Name()))
		if err != nil {
			continue
		}
		if statusFilter == "all" || imp.Status == statusFilter {
			results = append(results, *imp)
		}
	}

	severityOrder := map[string]int{
		"critical": 0, "high": 1, "medium": 2, "low": 3,
	}
	sort.Slice(results, func(i, j int) bool {
		si := severityOrder[results[i].Severity]
		sj := severityOrder[results[j].Severity]
		if si != sj {
			return si < sj
		}
		return results[i].ID < results[j].ID
	})

	return results, nil
}

// Resolve marks an impediment as resolved.
func (r *Repository) Resolve(id, resolution string, now time.Time) error {
	imp, err := r.Get(id)
	if err != nil {
		return err
	}

	imp.Status = "resolved"
	imp.Resolution = resolution
	imp.ResolvedAt = &now

	return r.save(imp)
}

// Get retrieves a single impediment by ID.
func (r *Repository) Get(id string) (*Impediment, error) {
	path := filepath.Join(r.Dir, id+".md")
	return r.load(path)
}

func (r *Repository) save(imp *Impediment) error {
	data, err := yaml.Marshal(imp)
	if err != nil {
		return fmt.Errorf("marshal impediment: %w", err)
	}

	content := fmt.Sprintf("---\n%s---\n\n%s\n", string(data), imp.Description)
	path := filepath.Join(r.Dir, imp.ID+".md")
	return os.WriteFile(path, []byte(content), 0644)
}

func (r *Repository) load(path string) (*Impediment, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read impediment %q: %w", path, err)
	}

	content := string(data)
	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid impediment format in %q", path)
	}

	var imp Impediment
	if err := yaml.Unmarshal([]byte(parts[1]), &imp); err != nil {
		return nil, fmt.Errorf("parse impediment frontmatter: %w", err)
	}
	imp.Description = strings.TrimSpace(parts[2])

	return &imp, nil
}

var impIDRegex = regexp.MustCompile(`^IMP-(\d+)$`)

func (r *Repository) nextID() (string, error) {
	entries, err := os.ReadDir(r.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "IMP-001", nil
		}
		return "", err
	}

	maxNum := 0
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), ".md")
		matches := impIDRegex.FindStringSubmatch(name)
		if len(matches) == 2 {
			n, err := strconv.Atoi(matches[1])
			if err == nil && n > maxNum {
				maxNum = n
			}
		}
	}

	return fmt.Sprintf("IMP-%03d", maxNum+1), nil
}
