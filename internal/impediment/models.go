package impediment

import "time"

// Impediment represents a tracked blocker affecting team flow.
type Impediment struct {
	ID          string     `yaml:"id"`
	Title       string     `yaml:"title"`
	Severity    string     `yaml:"severity"`
	Owner       string     `yaml:"owner,omitempty"`
	Status      string     `yaml:"status"`
	CreatedAt   time.Time  `yaml:"created_at"`
	ResolvedAt  *time.Time `yaml:"resolved_at,omitempty"`
	Resolution  string     `yaml:"resolution,omitempty"`
	Source      string     `yaml:"source"` // manual or detected
	Description string     `yaml:"-"`      // Stored as Markdown body
}

// AgeDays returns the number of days since the impediment was created.
func (imp *Impediment) AgeDays() int {
	end := time.Now()
	if imp.ResolvedAt != nil {
		end = *imp.ResolvedAt
	}
	return int(end.Sub(imp.CreatedAt).Hours() / 24)
}

// IsStale returns true if the impediment is open and older than 14 days.
func (imp *Impediment) IsStale() bool {
	return imp.Status == "open" && imp.AgeDays() > 14
}
