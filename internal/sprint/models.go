package sprint

import "time"

// SprintState tracks a sprint lifecycle.
type SprintState struct {
	SprintName     string   `json:"sprint_name"`
	Goal           string   `json:"goal"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date"`
	PlannedItems   []string `json:"planned_items"`
	CompletedItems []string `json:"completed_items"`
	Velocity       float64  `json:"velocity"`
	Status         string   `json:"status"` // planning, active, complete
}

// ComputeVelocity sets velocity based on completed items count.
func (s *SprintState) ComputeVelocity() {
	s.Velocity = float64(len(s.CompletedItems))
}

// DurationDays returns the sprint duration in days.
func (s *SprintState) DurationDays() int {
	start, err1 := time.Parse("2006-01-02", s.StartDate)
	end, err2 := time.Parse("2006-01-02", s.EndDate)
	if err1 != nil || err2 != nil {
		return 14 // default 2-week sprint
	}
	return int(end.Sub(start).Hours() / 24)
}
