package metrics

// DefaultThresholds defines traffic-light thresholds for health indicators.
var DefaultThresholds = map[string][2]float64{
	"velocity": {8.0, 5.0},   // green >= 8, yellow >= 5, red < 5
	"quality":  {0.10, 0.20}, // green <= 0.10, yellow <= 0.20, red > 0.20 (inverted — lower is better)
	"review":   {2.0, 3.0},   // green <= 2, yellow <= 3, red > 3 (inverted)
	"backlog":  {80.0, 60.0}, // green >= 80% ready, yellow >= 60%
	"flow":     {70.0, 50.0}, // green >= 70% efficiency, yellow >= 50%
}

// ComputeHealth computes traffic-light health indicators from a snapshot.
func ComputeHealth(current MetricsSnapshot, history []MetricsSnapshot) []HealthIndicator {
	var indicators []HealthIndicator

	// Velocity
	indicators = append(indicators, computeIndicator(
		"velocity", current.Velocity,
		DefaultThresholds["velocity"][0], DefaultThresholds["velocity"][1],
		false, extractValues(history, func(s MetricsSnapshot) float64 { return s.Velocity }),
	))

	// Quality (defect rate — inverted: lower is better)
	indicators = append(indicators, computeIndicator(
		"quality", current.DefectRate,
		DefaultThresholds["quality"][0], DefaultThresholds["quality"][1],
		true, extractValues(history, func(s MetricsSnapshot) float64 { return s.DefectRate }),
	))

	// Review (iterations — inverted: lower is better)
	indicators = append(indicators, computeIndicator(
		"review", current.ReviewIterations,
		DefaultThresholds["review"][0], DefaultThresholds["review"][1],
		true, extractValues(history, func(s MetricsSnapshot) float64 { return s.ReviewIterations }),
	))

	// Backlog health (% ready)
	readyPct := 0.0
	if current.BacklogHealth.Total > 0 {
		readyPct = float64(current.BacklogHealth.Ready) / float64(current.BacklogHealth.Total) * 100
	}
	indicators = append(indicators, computeIndicator(
		"backlog", readyPct,
		DefaultThresholds["backlog"][0], DefaultThresholds["backlog"][1],
		false, nil,
	))

	// Flow efficiency
	indicators = append(indicators, computeIndicator(
		"flow", current.FlowEfficiency,
		DefaultThresholds["flow"][0], DefaultThresholds["flow"][1],
		false, extractValues(history, func(s MetricsSnapshot) float64 { return s.FlowEfficiency }),
	))

	return indicators
}

func computeIndicator(dimension string, value, greenThreshold, yellowThreshold float64, inverted bool, history []float64) HealthIndicator {
	status := computeStatus(value, greenThreshold, yellowThreshold, inverted)
	trend := "stable"
	if len(history) >= 2 {
		if inverted {
			// For inverted metrics, negate values for trend computation
			negated := make([]float64, len(history))
			for i, v := range history {
				negated[i] = -v
			}
			trend = DetermineTrend(negated)
		} else {
			trend = DetermineTrend(history)
		}
	}

	return HealthIndicator{
		Dimension:       dimension,
		Status:          status,
		Value:           value,
		ThresholdGreen:  greenThreshold,
		ThresholdYellow: yellowThreshold,
		Trend:           trend,
	}
}

func computeStatus(value, green, yellow float64, inverted bool) string {
	if inverted {
		// Lower is better (e.g., defect rate, review iterations)
		if value <= green {
			return "green"
		}
		if value <= yellow {
			return "yellow"
		}
		return "red"
	}
	// Higher is better (e.g., velocity, flow efficiency)
	if value >= green {
		return "green"
	}
	if value >= yellow {
		return "yellow"
	}
	return "red"
}

func extractValues(snapshots []MetricsSnapshot, fn func(MetricsSnapshot) float64) []float64 {
	var values []float64
	for _, s := range snapshots {
		values = append(values, fn(s))
	}
	return values
}
