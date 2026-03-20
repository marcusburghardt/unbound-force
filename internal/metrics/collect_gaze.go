package metrics

import (
	"encoding/json"
	"time"

	"github.com/unbound-force/unbound-force/internal/artifacts"
)

// CollectGaze collects metrics from Gaze quality report artifacts.
func CollectGaze(artifactDir string, since time.Time) (*SourceCollection, error) {
	paths, err := artifacts.FindArtifacts(artifactDir, "quality-report")
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, nil // No artifacts found — not an error
	}

	now := time.Now().UTC()
	raw := make(map[string]interface{})
	var reports []map[string]interface{}

	for _, p := range paths {
		env, err := artifacts.ReadEnvelope(p)
		if err != nil {
			continue
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(env.Payload, &payload); err != nil {
			continue
		}
		reports = append(reports, payload)
	}

	raw["reports"] = reports
	raw["report_count"] = len(reports)

	return &SourceCollection{
		Source:      "gaze",
		CollectedAt: now,
		DataPoints:  len(reports),
		RawData:     raw,
	}, nil
}
