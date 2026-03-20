package metrics

import "testing"

func TestComputeStatus_NonInverted_Green(t *testing.T) {
	// Higher is better: value=10 >= green=8 → green
	got := computeStatus(10, 8, 5, false)
	if got != "green" {
		t.Errorf("status = %q, want green", got)
	}
}

func TestComputeStatus_NonInverted_Yellow(t *testing.T) {
	// Higher is better: 5 <= value=6 < green=8 → yellow
	got := computeStatus(6, 8, 5, false)
	if got != "yellow" {
		t.Errorf("status = %q, want yellow", got)
	}
}

func TestComputeStatus_NonInverted_Red(t *testing.T) {
	// Higher is better: value=3 < yellow=5 → red
	got := computeStatus(3, 8, 5, false)
	if got != "red" {
		t.Errorf("status = %q, want red", got)
	}
}

func TestComputeStatus_Inverted_Green(t *testing.T) {
	// Lower is better: value=0.05 <= green=0.10 → green
	got := computeStatus(0.05, 0.10, 0.20, true)
	if got != "green" {
		t.Errorf("status = %q, want green", got)
	}
}

func TestComputeStatus_Inverted_Yellow(t *testing.T) {
	// Lower is better: green=0.10 < value=0.15 <= yellow=0.20 → yellow
	got := computeStatus(0.15, 0.10, 0.20, true)
	if got != "yellow" {
		t.Errorf("status = %q, want yellow", got)
	}
}

func TestComputeStatus_Inverted_Red(t *testing.T) {
	// Lower is better: value=0.30 > yellow=0.20 → red
	got := computeStatus(0.30, 0.10, 0.20, true)
	if got != "red" {
		t.Errorf("status = %q, want red", got)
	}
}

func TestComputeHealth_ProducesAllDimensions(t *testing.T) {
	snap := MetricsSnapshot{
		Velocity:         10.0,
		DefectRate:       0.05,
		ReviewIterations: 1.5,
		BacklogHealth:    BacklogHealth{Total: 100, Ready: 85, Stale: 2},
		FlowEfficiency:   75.0,
	}

	indicators := ComputeHealth(snap, nil)
	if len(indicators) != 5 {
		t.Fatalf("got %d indicators, want 5", len(indicators))
	}

	expectedDimensions := map[string]bool{
		"velocity": false,
		"quality":  false,
		"review":   false,
		"backlog":  false,
		"flow":     false,
	}

	for _, ind := range indicators {
		if _, ok := expectedDimensions[ind.Dimension]; !ok {
			t.Errorf("unexpected dimension %q", ind.Dimension)
		}
		expectedDimensions[ind.Dimension] = true
	}

	for dim, found := range expectedDimensions {
		if !found {
			t.Errorf("missing dimension %q", dim)
		}
	}
}

func TestComputeHealth_TrendComputation(t *testing.T) {
	// Build 3 snapshots with improving velocity: 5 → 7 → 10
	history := []MetricsSnapshot{
		{Velocity: 5.0},
		{Velocity: 7.0},
		{Velocity: 10.0},
	}
	current := history[len(history)-1]

	indicators := ComputeHealth(current, history)

	var velocityIndicator *HealthIndicator
	for i := range indicators {
		if indicators[i].Dimension == "velocity" {
			velocityIndicator = &indicators[i]
			break
		}
	}

	if velocityIndicator == nil {
		t.Fatal("velocity indicator not found")
	}
	if velocityIndicator.Trend != "improving" {
		t.Errorf("velocity trend = %q, want improving", velocityIndicator.Trend)
	}
	if velocityIndicator.Status != "green" {
		t.Errorf("velocity status = %q, want green (value=10 >= threshold=8)", velocityIndicator.Status)
	}
}
