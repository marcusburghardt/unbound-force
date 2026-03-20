package metrics

import (
	"math"
	"testing"
)

func TestComputeVelocity_MultipleSnapshots(t *testing.T) {
	snapshots := []MetricsSnapshot{
		{Velocity: 5.0},
		{Velocity: 8.0},
		{Velocity: 12.0},
	}

	points := ComputeVelocity(snapshots)
	if len(points) != 3 {
		t.Fatalf("got %d points, want 3", len(points))
	}

	want := []struct {
		sprint   string
		velocity float64
	}{
		{"Sprint 1", 5.0},
		{"Sprint 2", 8.0},
		{"Sprint 3", 12.0},
	}

	for i, w := range want {
		if points[i].Sprint != w.sprint {
			t.Errorf("points[%d].Sprint = %q, want %q", i, points[i].Sprint, w.sprint)
		}
		if points[i].Velocity != w.velocity {
			t.Errorf("points[%d].Velocity = %f, want %f", i, points[i].Velocity, w.velocity)
		}
	}
}

func TestComputeVelocity_Empty(t *testing.T) {
	points := ComputeVelocity(nil)
	if len(points) != 0 {
		t.Errorf("got %d points, want 0", len(points))
	}

	points = ComputeVelocity([]MetricsSnapshot{})
	if len(points) != 0 {
		t.Errorf("got %d points for empty slice, want 0", len(points))
	}
}

func TestComputeCycleTimeFromValues_MultipleValues(t *testing.T) {
	hours := []float64{10, 20, 30, 40, 50}
	stats := ComputeCycleTimeFromValues(hours)

	if math.Abs(stats.Avg-30.0) > 0.1 {
		t.Errorf("Avg = %f, want 30.0", stats.Avg)
	}
	if math.Abs(stats.Median-30.0) > 0.1 {
		t.Errorf("Median = %f, want 30.0", stats.Median)
	}
	if math.Abs(stats.P90-46.0) > 0.1 {
		t.Errorf("P90 = %f, want ~46.0", stats.P90)
	}
	if math.Abs(stats.P99-49.6) > 0.1 {
		t.Errorf("P99 = %f, want ~49.6", stats.P99)
	}
}

func TestComputeCycleTimeFromValues_SingleValue(t *testing.T) {
	stats := ComputeCycleTimeFromValues([]float64{42})

	if stats.Avg != 42 {
		t.Errorf("Avg = %f, want 42", stats.Avg)
	}
	if stats.Median != 42 {
		t.Errorf("Median = %f, want 42", stats.Median)
	}
	if stats.P90 != 42 {
		t.Errorf("P90 = %f, want 42", stats.P90)
	}
	if stats.P99 != 42 {
		t.Errorf("P99 = %f, want 42", stats.P99)
	}
}

func TestComputeCycleTimeFromValues_Empty(t *testing.T) {
	stats := ComputeCycleTimeFromValues(nil)

	if stats.Avg != 0 {
		t.Errorf("Avg = %f, want 0", stats.Avg)
	}
	if stats.Median != 0 {
		t.Errorf("Median = %f, want 0", stats.Median)
	}
	if stats.P90 != 0 {
		t.Errorf("P90 = %f, want 0", stats.P90)
	}
	if stats.P99 != 0 {
		t.Errorf("P99 = %f, want 0", stats.P99)
	}
}

func TestPercentile_BoundaryIndices(t *testing.T) {
	// p0 on sorted data returns the first element
	sorted := []float64{1, 2, 3, 4, 5}
	got := percentile(sorted, 0)
	if got != 1.0 {
		t.Errorf("percentile(0) = %f, want 1.0", got)
	}

	// p100 returns the last element
	got = percentile(sorted, 100)
	if got != 5.0 {
		t.Errorf("percentile(100) = %f, want 5.0", got)
	}

	// p50 on even count: [1,2,3,4] → index 1.5 → interpolate 2 and 3 → 2.5
	even := []float64{1, 2, 3, 4}
	got = percentile(even, 50)
	if math.Abs(got-2.5) > 0.01 {
		t.Errorf("percentile(50) on [1,2,3,4] = %f, want 2.5", got)
	}

	// Empty returns 0
	got = percentile(nil, 50)
	if got != 0 {
		t.Errorf("percentile on empty = %f, want 0", got)
	}
}

func TestComputeFlowEfficiency_Normal(t *testing.T) {
	got := ComputeFlowEfficiency(50, 100)
	if math.Abs(got-50.0) > 0.1 {
		t.Errorf("FlowEfficiency = %f, want 50.0", got)
	}
}

func TestComputeFlowEfficiency_ZeroDenominator(t *testing.T) {
	got := ComputeFlowEfficiency(50, 0)
	if got != 0 {
		t.Errorf("FlowEfficiency = %f, want 0 (zero denominator)", got)
	}
}

func TestDetermineTrend_Improving(t *testing.T) {
	// 5 → 8 is +60%, which exceeds the 10% threshold
	trend := DetermineTrend([]float64{5, 8})
	if trend != "improving" {
		t.Errorf("trend = %q, want improving", trend)
	}
}

func TestDetermineTrend_Declining(t *testing.T) {
	// 10 → 5 is -50%, which is below -10% threshold
	trend := DetermineTrend([]float64{10, 5})
	if trend != "declining" {
		t.Errorf("trend = %q, want declining", trend)
	}
}

func TestDetermineTrend_Stable(t *testing.T) {
	// 10 → 10.5 is +5%, within the ±10% threshold
	trend := DetermineTrend([]float64{10, 10.5})
	if trend != "stable" {
		t.Errorf("trend = %q, want stable", trend)
	}
}

func TestDetermineTrend_SingleValue(t *testing.T) {
	trend := DetermineTrend([]float64{10})
	if trend != "stable" {
		t.Errorf("trend = %q, want stable (single value)", trend)
	}
}

func TestDetermineTrend_ZeroPrev(t *testing.T) {
	// 0 → 5: prev is zero, last > 0 → improving
	trend := DetermineTrend([]float64{0, 5})
	if trend != "improving" {
		t.Errorf("trend = %q, want improving (zero prev)", trend)
	}

	// 0 → 0: prev is zero, last is not > 0 → stable
	trend = DetermineTrend([]float64{0, 0})
	if trend != "stable" {
		t.Errorf("trend = %q, want stable (zero→zero)", trend)
	}
}
