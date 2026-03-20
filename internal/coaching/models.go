package coaching

import "time"

// RetroRecord represents a structured retrospective session.
type RetroRecord struct {
	Date                 string                 `yaml:"date"`
	Participants         []string               `yaml:"participants,omitempty"`
	DataPresented        map[string]interface{} `yaml:"data_presented,omitempty"`
	PatternsIdentified   []string               `yaml:"patterns_identified,omitempty"`
	RootCauses           []string               `yaml:"root_causes,omitempty"`
	ImprovementProposals []string               `yaml:"improvement_proposals,omitempty"`
	ActionItems          []ActionItem           `yaml:"action_items,omitempty"`
	Notes                string                 `yaml:"-"` // Markdown body
}

// ActionItem represents a tracked improvement commitment.
type ActionItem struct {
	ID              string `yaml:"id"`
	Description     string `yaml:"description"`
	Owner           string `yaml:"owner"`
	Deadline        string `yaml:"deadline"`
	Status          string `yaml:"status"` // pending, in-progress, completed, stale
	RetrospectiveID string `yaml:"retrospective_id"`
}

// CoachingInteraction records a coaching session.
type CoachingInteraction struct {
	Topic            string    `json:"topic"`
	QuestionsAsked   []string  `json:"questions_asked"`
	InsightsSurfaced []string  `json:"insights_surfaced"`
	Outcome          string    `json:"outcome"` // action_item, escalation, resolved, deferred
	Timestamp        time.Time `json:"timestamp"`
}

// IsStale returns true if the action item is past deadline and not completed.
func (ai *ActionItem) IsStale() bool {
	if ai.Status == "completed" {
		return false
	}
	deadline, err := time.Parse("2006-01-02", ai.Deadline)
	if err != nil {
		return false
	}
	return time.Now().After(deadline)
}
