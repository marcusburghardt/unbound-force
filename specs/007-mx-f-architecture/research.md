# Research: Mx F Architecture

**Spec**: 007-mx-f-architecture
**Date**: 2026-03-20

## R1: CLI Architecture — AppParams Pattern

**Decision**: Follow the Muti-Mind `AppParams +
newRootCmdWithParams()` pattern. `MxFParams` struct carries
`io.Writer`, data dir paths, output format, and `GHRunner`
interface.

**Rationale**: Proven pattern, enables test injection.
Muti-Mind has 12 subcommands using this exact pattern.
All CLI commands delegate to `runXxx(params)` functions,
making unit testing straightforward without subprocess
execution or `os.Stdout` mocking.

**Alternatives considered**:
- Flag-based (less testable): Parsing flags directly in
  command handlers couples CLI parsing to business logic.
- Global state (violates conventions): AGENTS.md explicitly
  requires "No global state: Prefer functional style and
  dependency injection."

## R2: GitHub API via GHRunner

**Decision**: Import `internal/sync.GHRunner` directly.
Use `gh api --json` commands for data retrieval.

**Rationale**: Already proven in Muti-Mind's GitHub sync
implementation. Handles auth (via `gh auth login`),
pagination via `--paginate`. Zero new auth code required.

**Alternatives considered**:
- `go-github` library: Adds a large external dependency
  and duplicates auth handling already solved by `gh`.
- Direct HTTP: Reinvents authentication, token management,
  and pagination logic.

**Note**: Mx F needs a dispatch-capable test stub that
returns different responses for different `gh` commands.
Muti-Mind's `GHRunner` mock returns a single response;
Mx F's collectors call multiple distinct `gh api` endpoints
in a single `collect` run (PRs, issues, CI runs, etc.).

## R3: Metrics Storage — JSON Files

**Decision**: JSON files at
`.mx-f/data/{source}/{timestamp}.json`. One file per
collection run.

**Rationale**: `CGO_ENABLED=0` constraint rules out SQLite.
JSON is simple, inspectable, consistent with project
patterns (Muti-Mind uses JSON for artifact output). Files
are individually readable by `graphthulhu` for knowledge
graph integration.

**Alternatives considered**:
- SQLite: Requires CGO, violates the zero-CGO build
  constraint.
- BoltDB: Pure Go but requires custom query logic for
  time-range and aggregation queries.
- Markdown+YAML: Awkward for numeric data and time-series
  aggregation. Better suited for narrative content.

**Directory structure**:
```
.mx-f/data/
├── github/       # PR, issue, CI, commit metrics
├── gaze/         # Quality report metrics
├── divisor/      # Review verdict metrics
├── muti-mind/    # Backlog and acceptance metrics
└── snapshots/    # Consolidated metrics snapshots
```

## R4: Impediment Storage — Markdown+YAML

**Decision**: Impediments stored as
`.mx-f/impediments/IMP-NNN.md` with YAML frontmatter.
Same pattern as Muti-Mind's backlog items.

**Rationale**: Human-readable, git-trackable, consistent
with the existing Muti-Mind backlog pattern
(`.muti-mind/backlog/*.md`). Auto-incrementing `IMP-NNN`
IDs derived by scanning existing files in the directory.

**Alternatives considered**:
- JSON files: Less human-readable for narrative
  descriptions (impediments include context, resolution
  notes, and escalation history that benefit from
  Markdown formatting).

## R5: Retrospective Records

**Decision**: Retrospective records stored as
`.mx-f/retros/YYYY-MM-DD.md` with YAML frontmatter for
structured data and Markdown body for notes.

**Rationale**: Date-keyed for natural chronological
ordering. Human-readable and git-trackable. Action items
embedded in frontmatter as an array, enabling programmatic
parsing by the CLI while keeping the full record readable
in any Markdown viewer.

**Example frontmatter**:
```yaml
---
date: 2026-03-20
participants:
  - "@dev1"
  - "@dev2"
action_items:
  - id: AI-001
    description: "Update CI config to pin base image"
    owner: "@dev1"
    deadline: 2026-03-27
    status: pending
---
```

## R6: Artifact Envelope Generalization

**Decision**: Refactor `internal/artifacts.writeArtifact()`
to accept a `hero` parameter instead of hardcoding
`"muti-mind"`. Mx F passes `"mx-f"`. Add
`ReadEnvelope(path)` and `FindArtifacts(dir, artifactType)`
functions for artifact consumption.

**Rationale**: Mx F is both a producer (`metrics-snapshot`,
`coaching-record`) and consumer (`quality-report`,
`review-verdict`, `backlog-item`) of artifacts. The current
`internal/artifacts` package is write-only and hardcoded
to Muti-Mind. Generalizing it benefits all current and
future heroes.

**Alternatives considered**:
- Duplicate `writeArtifact` in Mx F's package: Less DRY,
  creates two implementations of the same envelope logic
  that could drift.

**New functions**:
- `WriteArtifact(hero, artifactType, payload, outDir)` —
  generalized version of current function.
- `ReadEnvelope(path) (*Envelope, error)` — parse a JSON
  artifact file into the standard envelope struct.
- `FindArtifacts(dir, artifactType) ([]string, error)` —
  discover artifact files by type in a directory tree.

## R7: ASCII Chart Rendering

**Decision**: Custom ASCII rendering using `lipgloss`
(already in dependency tree) for styling + a minimal
sparkline/bar chart module in `internal/dashboard/text.go`.
Consider `guptarohit/asciigraph` for line charts if custom
rendering is too complex.

**Rationale**: `lipgloss` provides colors and styling for
terminal output. Bar charts and sparklines are simple
enough to implement inline (~100-150 lines each). Avoids
new dependencies for US5 P3 work. If line chart rendering
proves complex during implementation,
`guptarohit/asciigraph` is a lightweight, well-maintained
option.

**Alternatives considered**:
- `termui`: Full TUI framework, far heavier than needed
  for static chart output.
- `asciigraph` from the start: Adds a dependency for
  something that may be simple enough to build in-house.
- No charts (tables only): Loses the visual trend
  comprehension that makes dashboards useful.

## R8: HTML Dashboard

**Decision**: Go's `html/template` stdlib package generates
standalone HTML files. Embed a lightweight JS charting
library (Chart.js via CDN link) for interactive charts.
No server required.

**Rationale**: Static HTML generation is a Go strength.
Chart.js CDN link keeps the generated file small (~1KB
of template code). No Node.js runtime needed. The
generated HTML is a single self-contained file that opens
in any browser.

**Alternatives considered**:
- Embedded Chart.js (no CDN): Increases binary size by
  ~200KB for a P3 feature. CDN is acceptable since HTML
  dashboards are for stakeholder sharing (implies network).
- SVG generation: More complex, less interactive, no
  tooltips or hover effects.
- Web server (`net/http`): Overkill for a report file.
  Violates the "no runtime dependency" design goal.

## R9: Coaching Agent Design

**Decision**: Single Markdown file (`mx-f-coach.md`)
deployed via `unbound init`. Structure:
1. YAML frontmatter (description, mode, model, temperature,
   tools)
2. H1: Role — coaching philosophy, flow facilitation
3. H2: Source Documents — AGENTS.md, constitution, specs,
   metrics data at `.mx-f/data/`, impediments at
   `.mx-f/impediments/`, retro records at `.mx-f/retros/`
4. H2: Coaching Framework — 5 Whys, reflective questioning,
   mirroring, probing techniques
5. H2: Retrospective Facilitation Protocol — structured
   5-section format (data, patterns, root causes,
   proposals, action items)
6. H2: Boundary Rules — no prescriptions, redirect
   technical questions to appropriate heroes

**Rationale**: Same pattern as `cobalt-crush-dev.md` and
`divisor-*.md`. The coaching persona is instruction-based —
the CLI provides data, the agent provides facilitation.
Separation of concerns: CLI owns computation, agent owns
conversation.

**Alternatives considered**:
- Multiple agent files (coach, retro-facilitator,
  impediment-analyst): Over-engineering for what is
  fundamentally one persona with different modes.
- Agent-only (no CLI): Cannot compute metrics, render
  charts, or manage structured data reliably.

## R10: GoReleaser Multi-Binary

**Decision**: Add a second `builds` entry for `cmd/mxf/`
in `.goreleaser.yaml`. Both binaries ship in the same
archive. Homebrew cask updated to include both.

**Rationale**: Follows the `unbound` pattern. Both binaries
share the same Go module and dependency tree. Single
release pipeline, single version number.

**Note**: `cmd/mutimind` is also not yet in GoReleaser —
both `mutimind` and `mxf` could be added together as a
single GoReleaser configuration update.

**Alternatives considered**:
- Separate releases per binary: Fragments the release
  pipeline and version tracking.
- Single binary with subcommands: Would merge Mx F into
  the `unbound` CLI, conflating the scaffold tool with
  hero functionality.

## R11: Embedded Asset Count

**Decision**: Add 1 new embedded file (`mx-f-coach.md`),
bringing total from 46 to 47.

**Impact on tests**:
- `expectedAssetPaths` in `internal/scaffold/scaffold_test.go`
  gains 1 entry for `opencode/agents/mx-f-coach.md`.
- `TestRunInit_FreshDir` file count assertion changes
  46 → 47.
- Drift detection test must include the new file to ensure
  the embedded asset matches its canonical source.
<!-- scaffolded by unbound vdev -->
