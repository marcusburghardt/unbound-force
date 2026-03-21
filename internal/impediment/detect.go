package impediment

import (
	"fmt"
	"time"

	"github.com/unbound-force/unbound-force/internal/metrics"
)

// Detect analyzes metrics trends and creates draft impediments for anomalies.
func Detect(store *metrics.Store, repo *Repository, now time.Time) ([]Impediment, error) {
	snapshots, err := store.ReadSnapshots(time.Time{})
	if err != nil {
		return nil, err
	}
	if len(snapshots) < 2 {
		return nil, fmt.Errorf("insufficient data for trend analysis (need at least 2 snapshots)")
	}

	var detected []Impediment
	latest := snapshots[len(snapshots)-1]

	// Check for CI failure rate spike
	if len(snapshots) >= 2 {
		prev := snapshots[len(snapshots)-2]
		ciDrop := prev.CIPassRate - latest.CIPassRate
		if ciDrop > 15 {
			imp, err := repo.Add(
				fmt.Sprintf("CI pass rate dropped %.0f%% (from %.1f%% to %.1f%%)", ciDrop, prev.CIPassRate, latest.CIPassRate),
				"high", "", "Detected automatically from metrics trend analysis.", now,
			)
			if err == nil {
				imp.Source = "detected"
				if saveErr := repo.save(imp); saveErr != nil {
					return detected, fmt.Errorf("update impediment source: %w", saveErr)
				}
				detected = append(detected, *imp)
			}
		}
	}

	// Check for review turnaround increase
	if len(snapshots) >= 3 {
		var avgIter float64
		for _, s := range snapshots[:len(snapshots)-1] {
			avgIter += s.ReviewIterations
		}
		avgIter /= float64(len(snapshots) - 1)
		if latest.ReviewIterations > avgIter*1.5 {
			imp, err := repo.Add(
				fmt.Sprintf("Review iterations increased %.0f%% above average", (latest.ReviewIterations-avgIter)/avgIter*100),
				"medium", "", "Detected automatically from metrics trend analysis.", now,
			)
			if err == nil {
				imp.Source = "detected"
				if saveErr := repo.save(imp); saveErr != nil {
					return detected, fmt.Errorf("update impediment source: %w", saveErr)
				}
				detected = append(detected, *imp)
			}
		}
	}

	// Check for velocity drop
	if len(snapshots) >= 3 {
		var avgVel float64
		for _, s := range snapshots[:len(snapshots)-1] {
			avgVel += s.Velocity
		}
		avgVel /= float64(len(snapshots) - 1)
		if avgVel > 0 && latest.Velocity < avgVel*0.75 {
			imp, err := repo.Add(
				fmt.Sprintf("Velocity dropped %.0f%% below average", (avgVel-latest.Velocity)/avgVel*100),
				"medium", "", "Detected automatically from metrics trend analysis.", now,
			)
			if err == nil {
				imp.Source = "detected"
				if saveErr := repo.save(imp); saveErr != nil {
					return detected, fmt.Errorf("update impediment source: %w", saveErr)
				}
				detected = append(detected, *imp)
			}
		}
	}

	return detected, nil
}
