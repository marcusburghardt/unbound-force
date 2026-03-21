package metrics

import (
	"encoding/json"
	"time"

	"github.com/unbound-force/unbound-force/internal/artifacts"
)

// CollectMutiMind collects metrics from Muti-Mind backlog artifacts.
func CollectMutiMind(artifactDir string, since time.Time) (*SourceCollection, error) {
	blPaths, err := artifacts.FindArtifacts(artifactDir, "backlog-item")
	if err != nil {
		return nil, err
	}
	adPaths, err := artifacts.FindArtifacts(artifactDir, "acceptance-decision")
	if err != nil {
		return nil, err
	}

	if len(blPaths) == 0 && len(adPaths) == 0 {
		return nil, nil
	}

	now := time.Now().UTC()
	raw := make(map[string]interface{})

	var backlogItems []map[string]interface{}
	for _, p := range blPaths {
		env, err := artifacts.ReadEnvelope(p)
		if err != nil {
			continue
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(env.Payload, &payload); err != nil {
			continue
		}
		backlogItems = append(backlogItems, payload)
	}
	raw["backlog_items"] = backlogItems
	raw["backlog_size"] = len(backlogItems)

	var decisions []map[string]interface{}
	for _, p := range adPaths {
		env, err := artifacts.ReadEnvelope(p)
		if err != nil {
			continue
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(env.Payload, &payload); err != nil {
			continue
		}
		decisions = append(decisions, payload)
	}
	raw["acceptance_decisions"] = decisions
	raw["acceptance_count"] = len(decisions)

	accepted := 0
	for _, d := range decisions {
		if dec, ok := d["decision"].(string); ok && dec == "accept" {
			accepted++
		}
	}
	if len(decisions) > 0 {
		raw["acceptance_rate"] = float64(accepted) / float64(len(decisions))
	}

	return &SourceCollection{
		Source:      "muti-mind",
		CollectedAt: now,
		DataPoints:  len(backlogItems) + len(decisions),
		RawData:     raw,
	}, nil
}
