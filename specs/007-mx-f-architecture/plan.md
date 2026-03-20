# Implementation Plan: Mx F Architecture (Manager)

**Branch**: `007-mx-f-architecture` | **Date**: 2026-03-20 | **Spec**: [[specs/007-mx-f-architecture/spec]]
**Input**: Feature specification from `/specs/007-mx-f-architecture/spec.md`

## Summary

Mx F is the Manager hero — the Flow Facilitator and
Continuous Improvement Coach. It has two components: a Go
CLI backend (`cmd/mxf/`) for data operations (metrics
collection, querying, impediment tracking, dashboards,
sprint management) and an OpenCode agent
(`mx-f-coach.md`) for AI coaching and retrospective
facilitation. Same hybrid pattern as Muti-Mind (Spec 004).

## Technical Context

**Language/Version**: Go 1.24+ (CLI backend), Markdown (coaching agent)
**Primary Dependencies**: `github.com/spf13/cobra` (CLI), `github.com/charmbracelet/log` (logging), `github.com/charmbracelet/lipgloss` (terminal styling), `embed.FS` (agent embedding)
**Storage**: JSON files in `.mx-f/data/{source}/{timestamp}.json` for metrics, Markdown+YAML frontmatter in `.mx-f/impediments/` for impediments, `.mx-f/retros/` for retrospective records
**Testing**: Standard library `testing`, `-race -count=1`, `t.TempDir()`, `GHRunner` stubs. **Coverage Strategy**: 80% global unit test coverage minimum. 90% unit test coverage for `internal/metrics` (collector, store, compute, query) and `internal/impediment` parsing logic. 80% unit test coverage for GitHub collection via `GHRunner` interface mocking. 100% integration test coverage for collect/query/impediment round-trip happy paths. Functional tests via OpenCode scenarios for coaching agent interactions.
**Target Platform**: macOS, Linux (cross-compiled via GoReleaser)
**Project Type**: CLI tool (Go binary at `cmd/mxf/` + OpenCode agent)
**Performance Goals**: SC-003 specifies 5-second query target for `mx-f metrics summary`. Sub-second local data reads for impediment and sprint queries.
**Constraints**: `CGO_ENABLED=0`, no SQLite, no external services beyond `gh` CLI. All metrics storage is JSON-file-based.
**Scale/Scope**: Project-level metrics and process management, single repository scope

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Autonomous Collaboration**: PASS. Mx F communicates
  through artifacts (`metrics-snapshot`, `coaching-record`).
  Consumes artifacts from Gaze (`quality-report`), Divisor
  (`review-verdict`), and Muti-Mind (`backlog-item`,
  `acceptance-decision`) via file reads. No runtime coupling
  — all hero interaction is through standard artifact
  envelopes on the filesystem.
- **II. Composability First**: PASS. Mx F works standalone
  with just GitHub data (SC-010). Missing heroes = graceful
  degradation (FR-020). The coaching agent works without the
  CLI backend (it reads AGENTS.md and project files
  directly). Each CLI subcommand functions independently.
- **III. Observable Quality**: PASS. All queries support
  `--format json` (FR-005). Produces `metrics-snapshot` and
  `coaching-record` artifacts with provenance metadata
  conforming to the Spec 002 artifact envelope. Traffic-light
  health indicators provide machine-parseable quality signals.
- **IV. Testability**: PASS. `GHRunner` interface for GitHub
  API injection (reuses pattern from `internal/sync`). JSON
  file storage testable with `t.TempDir()`. Domain packages
  (`metrics`, `impediment`, `coaching`, `dashboard`, `sprint`)
  independently testable. Pure computation functions (velocity,
  cycle time, percentiles) are stateless and trivially
  unit-testable. Coverage strategy defined above.

## Project Structure

### Documentation (this feature)

```text
specs/007-mx-f-architecture/
├── spec.md              # Feature specification
├── plan.md              # This file (/speckit.plan command output)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
# CLI entry point
cmd/mxf/
└── main.go                    # Cobra root + 7 subcommands
                               # (collect, metrics, impediment,
                               #  dashboard, sprint, standup, retro)

# Internal packages (domain-oriented)
internal/
├── metrics/                   # US1+US2: Collection + computation + querying
│   ├── collector.go           # Source-specific collectors (github, gaze,
│   │                          #   divisor, mutimind)
│   ├── store.go               # JSON file read/write in .mx-f/data/
│   ├── compute.go             # Metric calculations (velocity, cycle time,
│   │                          #   lead time, defect rate, percentiles)
│   ├── query.go               # Query engine (filtering, aggregation, trends)
│   └── health.go              # Health indicators (traffic-light computation)
├── impediment/                # US4: Obstacle tracking
│   ├── impediment.go          # CRUD operations, IMP-NNN auto-ID assignment
│   ├── detect.go              # Proactive detection from metrics trends
│   └── models.go              # Impediment, Resolution structs
├── coaching/                  # US3: Retrospective facilitation
│   ├── retro.go               # Structured retrospective engine
│   ├── actions.go             # Action item tracking (AI-NNN IDs)
│   └── models.go              # RetroRecord, ActionItem structs
├── dashboard/                 # US5: Visualization
│   ├── text.go                # ASCII sparklines, bar charts (lipgloss)
│   └── html.go                # HTML dashboard generation (standalone file)
└── sprint/                    # US6: Sprint lifecycle
    ├── sprint.go              # Planning, review, standup
    └── models.go              # SprintState struct

# OpenCode agent (coaching persona)
.opencode/agents/
└── mx-f-coach.md              # NEW: coaching persona agent
                               # (canonical source)

# Scaffold embedded copy
internal/scaffold/assets/opencode/agents/
└── mx-f-coach.md              # NEW: embedded copy for unbound init

# Local data directories (created at runtime)
.mx-f/
├── data/                      # Metrics storage
│   ├── github/                # GitHub-sourced metrics
│   │   └── {timestamp}.json
│   ├── gaze/                  # Gaze quality report metrics
│   │   └── {timestamp}.json
│   ├── divisor/               # Divisor review metrics
│   │   └── {timestamp}.json
│   └── mutimind/              # Muti-Mind backlog metrics
│       └── {timestamp}.json
├── impediments/               # Impediment records (MD+YAML frontmatter)
│   └── IMP-NNN.md
└── retros/                    # Retrospective records
    └── {date}-retro.json
```

**Structure Decision**: The project uses the same hybrid
layout established by Muti-Mind (Spec 004). Go application
logic lives under `cmd/mxf/` and `internal/` (domain-oriented
packages mirroring the user story structure). The coaching
persona lives in `.opencode/agents/mx-f-coach.md`. This
separation keeps data operations (collection, computation,
storage, charts) in compiled Go code for performance and
testability, while coaching/facilitation runs as an AI agent
with access to the collected data.

## Complexity Tracking

No constitution violations. The most complex area is the
metrics computation engine (velocity, cycle time,
percentiles, trends, bottleneck detection). This is pure
computation on JSON data — no external dependencies, fully
testable with deterministic inputs and expected outputs.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
<!-- scaffolded by unbound vdev -->
