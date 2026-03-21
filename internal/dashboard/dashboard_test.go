package dashboard

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/unbound-force/unbound-force/internal/metrics"
)

func TestRenderBarChart_4Sprints(t *testing.T) {
	data := []BarChartPoint{
		{Label: "Sprint 1", Value: 12},
		{Label: "Sprint 2", Value: 18},
		{Label: "Sprint 3", Value: 15},
		{Label: "Sprint 4", Value: 20},
	}

	var buf bytes.Buffer
	if err := RenderBarChart("Velocity", data, &buf); err != nil {
		t.Fatalf("RenderBarChart: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "█") {
		t.Error("output should contain bar character █")
	}

	for _, d := range data {
		if !strings.Contains(out, d.Label) {
			t.Errorf("output should contain label %q", d.Label)
		}
	}

	// Values are printed via "%.0f" format
	for _, v := range []string{"12", "18", "15", "20"} {
		if !strings.Contains(out, v) {
			t.Errorf("output should contain value %q", v)
		}
	}
}

func TestRenderBarChart_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderBarChart("Empty", nil, &buf); err != nil {
		t.Fatalf("RenderBarChart: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty data, got %q", buf.String())
	}
}

func TestRenderSparkline_30Days(t *testing.T) {
	data := make([]float64, 30)
	for i := range data {
		data[i] = float64(i) + 1 // 1..30
	}

	var buf bytes.Buffer
	if err := RenderSparkline("Daily Velocity", data, &buf); err != nil {
		t.Fatalf("RenderSparkline: %v", err)
	}

	out := buf.String()

	// Verify at least some sparkline characters appear
	hasSparkChar := false
	for _, ch := range "▁▂▃▄▅▆▇█" {
		if strings.ContainsRune(out, ch) {
			hasSparkChar = true
			break
		}
	}
	if !hasSparkChar {
		t.Error("output should contain sparkline characters (▁▂▃▄▅▆▇█)")
	}

	if !strings.Contains(out, "Min:") {
		t.Error("output should contain Min:")
	}
	if !strings.Contains(out, "Avg:") {
		t.Error("output should contain Avg:")
	}
	if !strings.Contains(out, "Max:") {
		t.Error("output should contain Max:")
	}
}

func TestRenderSparkline_SingleValue(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderSparkline("Single", []float64{42.0}, &buf); err != nil {
		t.Fatalf("RenderSparkline: %v", err)
	}
	// Should not panic and should produce some output
	if buf.Len() == 0 {
		t.Error("expected output for single data point")
	}
}

func TestRenderHealthIndicators_GreenYellowRed(t *testing.T) {
	indicators := []metrics.HealthIndicator{
		{Dimension: "velocity", Status: "green", Value: 18.0, Trend: "improving"},
		{Dimension: "quality", Status: "yellow", Value: 0.25, Trend: "stable"},
		{Dimension: "review", Status: "red", Value: 4.5, Trend: "declining"},
	}

	var buf bytes.Buffer
	if err := RenderHealthIndicators("Health", indicators, &buf); err != nil {
		t.Fatalf("RenderHealthIndicators: %v", err)
	}

	out := buf.String()

	for _, ind := range indicators {
		if !strings.Contains(out, ind.Dimension) {
			t.Errorf("output should contain dimension %q", ind.Dimension)
		}
	}

	// Trend arrows: improving=↑, stable=→, declining=↓
	if !strings.Contains(out, "↑") {
		t.Error("output should contain improving arrow ↑")
	}
	if !strings.Contains(out, "→") {
		t.Error("output should contain stable arrow →")
	}
	if !strings.Contains(out, "↓") {
		t.Error("output should contain declining arrow ↓")
	}
}

func TestRenderHTML_ProducesValidFile(t *testing.T) {
	snap := metrics.MetricsSnapshot{
		Timestamp:        time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
		Velocity:         17.5,
		CycleTime:        metrics.CycleTimeStats{Avg: 24.0, Median: 18.0, P90: 48.0, P99: 72.0},
		LeadTime:         36.0,
		DefectRate:       0.1,
		ReviewIterations: 2.0,
		CIPassRate:       92.5,
		BacklogHealth:    metrics.BacklogHealth{Total: 30, Ready: 20, Stale: 2},
		FlowEfficiency:   65.0,
		SourcesCollected: []string{"github"},
	}

	indicators := []metrics.HealthIndicator{
		{Dimension: "velocity", Status: "green", Value: 17.5, Trend: "improving"},
		{Dimension: "quality", Status: "yellow", Value: 0.1, Trend: "stable"},
	}

	outPath := filepath.Join(t.TempDir(), "dashboard.html")
	if err := RenderHTML(snap, indicators, outPath); err != nil {
		t.Fatalf("RenderHTML: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
	if info.Size() < 100 {
		t.Errorf("file size = %d, want > 100 bytes", info.Size())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "Mx F") {
		t.Error("HTML should contain 'Mx F'")
	}
	if !strings.Contains(content, "17.5") {
		t.Error("HTML should contain velocity value 17.5")
	}
}
