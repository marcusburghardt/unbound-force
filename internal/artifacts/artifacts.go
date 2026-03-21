package artifacts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/unbound-force/unbound-force/internal/backlog"
)

// Envelope represents the Spec 002 Artifact Envelope
type Envelope struct {
	Hero          string          `json:"hero"`
	Version       string          `json:"version"`
	Timestamp     string          `json:"timestamp"`
	ArtifactType  string          `json:"artifact_type"`
	SchemaVersion string          `json:"schema_version"`
	Context       json.RawMessage `json:"context,omitempty"`
	Payload       json.RawMessage `json:"payload"`
}

// AcceptanceDecision represents the payload for an acceptance-decision artifact
type AcceptanceDecision struct {
	ItemID         string   `json:"item_id"`
	Decision       string   `json:"decision"` // accept, reject, conditional
	Rationale      string   `json:"rationale"`
	CriteriaMet    []string `json:"criteria_met"`
	CriteriaFailed []string `json:"criteria_failed"`
	GazeReportRef  string   `json:"gaze_report_ref"`
	DecidedAt      string   `json:"decided_at"`
}

const Version = "1.0.0"

// WriteArtifact writes a JSON artifact envelope to the given directory.
// The hero parameter identifies the producing hero (e.g., "mx-f", "muti-mind").
func WriteArtifact(dir, hero, artifactType, id string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	envelope := Envelope{
		Hero:          hero,
		Version:       Version,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		ArtifactType:  artifactType,
		SchemaVersion: "1.0.0",
		Payload:       payloadBytes,
	}

	envelopeBytes, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal envelope: %w", err)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create artifacts directory %q: %w", dir, err)
	}

	filename := fmt.Sprintf("%s-%s.json", id, artifactType)
	targetPath := filepath.Join(dir, filename)
	if err := os.WriteFile(targetPath, envelopeBytes, 0644); err != nil {
		return fmt.Errorf("write artifact %q: %w", targetPath, err)
	}
	return nil
}

// ReadEnvelope reads a JSON artifact file and returns the parsed Envelope.
func ReadEnvelope(path string) (*Envelope, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read artifact %q: %w", path, err)
	}
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse artifact %q: %w", path, err)
	}
	return &env, nil
}

// FindArtifacts discovers artifact files of the given type in a directory tree.
// Returns file paths sorted by filename descending (newest timestamp first).
func FindArtifacts(dir, artifactType string) ([]string, error) {
	var matches []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		env, readErr := ReadEnvelope(path)
		if readErr != nil {
			return nil // skip non-artifact JSON files
		}
		if env.ArtifactType == artifactType {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %q: %w", dir, err)
	}
	// Sort descending by filename (timestamps sort naturally)
	sort.Sort(sort.Reverse(sort.StringSlice(matches)))
	return matches, nil
}

// GenerateBacklogItemArtifact generates a backlog-item JSON artifact
func GenerateBacklogItemArtifact(dir string, item *backlog.Item) error {
	return WriteArtifact(dir, "muti-mind", "backlog-item", item.ID, item)
}

// GenerateAcceptanceDecision generates an acceptance-decision JSON artifact
func GenerateAcceptanceDecision(dir string, decision *AcceptanceDecision) error {
	return WriteArtifact(dir, "muti-mind", "acceptance-decision", decision.ItemID, decision)
}
