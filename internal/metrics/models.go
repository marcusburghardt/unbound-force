package metrics

import "time"

// SourceCollection represents raw data from a single collection run.
type SourceCollection struct {
	Source      string                 `json:"source"`
	CollectedAt time.Time              `json:"collected_at"`
	DataPoints  int                    `json:"data_points"`
	RawData     map[string]interface{} `json:"raw_data"`
}

// MetricsSnapshot is a point-in-time collection of all computed metrics.
type MetricsSnapshot struct {
	Timestamp        time.Time         `json:"timestamp"`
	Velocity         float64           `json:"velocity"`
	CycleTime        CycleTimeStats    `json:"cycle_time"`
	LeadTime         float64           `json:"lead_time"`
	DefectRate       float64           `json:"defect_rate"`
	ReviewIterations float64           `json:"review_iterations"`
	CIPassRate       float64           `json:"ci_pass_rate"`
	BacklogHealth    BacklogHealth     `json:"backlog_health"`
	FlowEfficiency   float64           `json:"flow_efficiency"`
	SourcesCollected []string          `json:"sources_collected"`
	HealthIndicators []HealthIndicator `json:"health_indicators,omitempty"`
}

// CycleTimeStats provides statistical breakdown of cycle time.
type CycleTimeStats struct {
	Avg    float64 `json:"avg"`
	Median float64 `json:"median"`
	P90    float64 `json:"p90"`
	P99    float64 `json:"p99"`
}

// BacklogHealth summarizes backlog item states.
type BacklogHealth struct {
	Total int `json:"total"`
	Ready int `json:"ready"`
	Stale int `json:"stale"`
}

// HealthIndicator is a traffic-light assessment of a metric dimension.
type HealthIndicator struct {
	Dimension       string  `json:"dimension"`
	Status          string  `json:"status"` // green, yellow, red
	Value           float64 `json:"value"`
	ThresholdGreen  float64 `json:"threshold_green"`
	ThresholdYellow float64 `json:"threshold_yellow"`
	Trend           string  `json:"trend"` // improving, stable, declining
}

// VelocityPoint represents velocity for a single sprint.
type VelocityPoint struct {
	Sprint   string  `json:"sprint"`
	Velocity float64 `json:"velocity"`
}

// BottleneckResult identifies a pipeline stage and its wait time.
type BottleneckResult struct {
	Stage       string  `json:"stage"`
	AvgWaitDays float64 `json:"avg_wait_days"`
}
