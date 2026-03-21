package metrics

import (
	"fmt"
	"math"
	"sort"
)

// ComputeVelocity calculates velocity per sprint from snapshot history.
func ComputeVelocity(snapshots []MetricsSnapshot) []VelocityPoint {
	var points []VelocityPoint
	for i, s := range snapshots {
		points = append(points, VelocityPoint{
			Sprint:   fmt.Sprintf("Sprint %d", i+1),
			Velocity: s.Velocity,
		})
	}
	return points
}

// ComputeCycleTimeFromValues calculates cycle time statistics from a slice of hours.
func ComputeCycleTimeFromValues(hours []float64) CycleTimeStats {
	if len(hours) == 0 {
		return CycleTimeStats{}
	}

	sorted := make([]float64, len(hours))
	copy(sorted, hours)
	sort.Float64s(sorted)

	sum := 0.0
	for _, h := range sorted {
		sum += h
	}

	return CycleTimeStats{
		Avg:    sum / float64(len(sorted)),
		Median: percentile(sorted, 50),
		P90:    percentile(sorted, 90),
		P99:    percentile(sorted, 99),
	}
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}
	idx := (p / 100) * float64(len(sorted)-1)
	lower := int(math.Floor(idx))
	upper := int(math.Ceil(idx))
	if lower == upper || upper >= len(sorted) {
		return sorted[lower]
	}
	frac := idx - float64(lower)
	return sorted[lower]*(1-frac) + sorted[upper]*frac
}

// ComputeFlowEfficiency calculates flow efficiency as value-add time / total time.
func ComputeFlowEfficiency(valueAddHours, totalHours float64) float64 {
	if totalHours == 0 {
		return 0
	}
	return (valueAddHours / totalHours) * 100
}

// DetermineTrend computes trend direction from recent values.
func DetermineTrend(values []float64) string {
	if len(values) < 2 {
		return "stable"
	}

	last := values[len(values)-1]
	prev := values[len(values)-2]
	threshold := 0.1 // 10% change threshold

	if prev == 0 {
		if last > 0 {
			return "improving"
		}
		return "stable"
	}

	change := (last - prev) / math.Abs(prev)
	if change > threshold {
		return "improving"
	}
	if change < -threshold {
		return "declining"
	}
	return "stable"
}
