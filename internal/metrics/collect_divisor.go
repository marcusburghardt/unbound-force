package metrics

import (
	"encoding/json"
	"time"

	"github.com/unbound-force/unbound-force/internal/artifacts"
)

// CollectDivisor collects metrics from Divisor review verdict artifacts.
func CollectDivisor(artifactDir string, since time.Time) (*SourceCollection, error) {
	paths, err := artifacts.FindArtifacts(artifactDir, "review-verdict")
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, nil
	}

	now := time.Now().UTC()
	raw := make(map[string]interface{})
	var verdicts []map[string]interface{}

	for _, p := range paths {
		env, err := artifacts.ReadEnvelope(p)
		if err != nil {
			continue
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(env.Payload, &payload); err != nil {
			continue
		}
		verdicts = append(verdicts, payload)
	}

	raw["verdicts"] = verdicts
	raw["verdict_count"] = len(verdicts)

	return &SourceCollection{
		Source:      "divisor",
		CollectedAt: now,
		DataPoints:  len(verdicts),
		RawData:     raw,
	}, nil
}
