package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// stubGHRunner satisfies sync.GHRunner for CLI-level tests.
type stubGHRunner struct {
	out []byte
	err error
}

func (s *stubGHRunner) Run(args ...string) ([]byte, error) {
	return s.out, s.err
}

func TestMxFParams_Defaults(t *testing.T) {
	p := &MxFParams{}
	p.defaults()

	if p.DataDir != ".mx-f" {
		t.Errorf("DataDir = %q, want .mx-f", p.DataDir)
	}
	if p.Stdout == nil {
		t.Error("Stdout should not be nil after defaults()")
	}
	if p.Stderr == nil {
		t.Error("Stderr should not be nil after defaults()")
	}
	if p.Now == nil {
		t.Error("Now should not be nil after defaults()")
	}
	if p.GHRunner == nil {
		t.Error("GHRunner should not be nil after defaults()")
	}
}

func TestNewRootCmd_HasSubcommands(t *testing.T) {
	cmd := newRootCmd()

	expected := []string{
		"collect", "metrics", "impediment",
		"dashboard", "sprint", "standup", "retro",
	}
	found := make(map[string]bool)
	for _, c := range cmd.Commands() {
		found[c.Name()] = true
	}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("missing subcommand %q", name)
		}
	}
}

func TestNewRootCmdWithParams_InjectsStubs(t *testing.T) {
	stub := &stubGHRunner{out: []byte("ok"), err: nil}
	var buf bytes.Buffer
	p := &MxFParams{
		Stdout:   &buf,
		Stderr:   &buf,
		GHRunner: stub,
		DataDir:  t.TempDir(),
		Now:      func() time.Time { return time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC) },
	}
	cmd := newRootCmdWithParams(p)
	if cmd == nil {
		t.Fatal("newRootCmdWithParams returned nil")
	}
}

func TestRunCollect_NoData(t *testing.T) {
	var buf bytes.Buffer
	p := MxFParams{
		DataDir: t.TempDir(),
		Stdout:  &buf,
		Stderr:  &buf,
		Now:     func() time.Time { return time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC) },
	}
	p.defaults()
	// Collecting with no repo and no artifacts should succeed
	// (graceful degradation — 0/4 sources is not an error)
	err := runCollect(p, "all", "", "30d")
	if err != nil {
		t.Errorf("expected graceful degradation, got error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "0/4 sources") {
		t.Errorf("expected 0/4 sources in output, got:\n%s", output)
	}
}
