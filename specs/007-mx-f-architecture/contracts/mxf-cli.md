# Contract: Mx F CLI (mxf)

**Spec**: 007-mx-f-architecture
**Date**: 2026-03-20
**Type**: CLI command schema

## Overview

The `mxf` CLI is a Go binary that provides metrics
collection, querying, impediment tracking, dashboard
visualization, sprint lifecycle management, and
retrospective facilitation for the Mx F Manager hero.
It follows the same hybrid pattern as Muti-Mind
(Spec 004): Go CLI backend for data operations, OpenCode
agent (`mx-f-coach.md`) for AI coaching.

## Installation

```bash
# Homebrew (recommended — installs unbound + mxf + graphthulhu)
brew install unbound-force/tap/unbound

# Go install
go install github.com/unbound-force/unbound-force/\
cmd/mxf@latest

# Build from source
git clone https://github.com/unbound-force/unbound-force
cd unbound-force
go build -o mxf ./cmd/mxf
```

## Commands

### `mxf collect`

Collects metrics from one or more data sources and stores
them as JSON files in `.mx-f/data/{source}/{timestamp}.json`.

```bash
mxf collect [flags]
```

**Flags**:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--source` | string | `all` | Data source to collect from: `github`, `gaze`, `divisor`, `muti-mind`, or `all` |
| `--repo` | string | current | GitHub repository (owner/repo format). Defaults to the repository in the current working directory. |
| `--period` | duration | `90d` | Time range to collect. Accepts Go-style durations: `30d`, `7d`, `90d`. |

**Behavior**:

1. Resolve target repository from `--repo` flag or current
   working directory (via `gh repo view`).
2. For each selected source:
   - **github**: Query GitHub API via `gh` CLI (`GHRunner`
     interface). Collect PR count, PR merge time, review
     turnaround, CI pass rate, commit frequency, issue
     open/close rates, contributor activity.
   - **gaze**: Scan for `quality-report` artifact files in
     the project. Collect CRAP scores, CRAPload counts,
     contract coverage, over-specification counts.
   - **divisor**: Scan for `review-verdict` artifact files.
     Collect review iteration counts, finding categories,
     approval rates, time-to-approval.
   - **muti-mind**: Scan for `backlog-item` and
     `acceptance-decision` artifact files. Collect backlog
     size, velocity, lead time, acceptance rates.
3. Write collected data to
   `.mx-f/data/{source}/{timestamp}.json`.
4. Print summary of collected data points.

**Output** (text):

```text
Collecting metrics (period: 90d)...

  github     42 data points  .mx-f/data/github/2026-03-20T14:30:00Z.json
  gaze       18 data points  .mx-f/data/gaze/2026-03-20T14:30:00Z.json
  divisor    23 data points  .mx-f/data/divisor/2026-03-20T14:30:00Z.json
  muti-mind  31 data points  .mx-f/data/mutimind/2026-03-20T14:30:00Z.json

Total: 114 data points collected from 4 sources.
```

**Graceful Degradation** (FR-020):

When a source is unavailable (no artifacts found, no `gh`
auth, etc.), the collector MUST:
1. Report the unavailable source with reason.
2. Continue collecting from remaining sources.
3. Exit 0 (partial collection is not an error).

```text
Collecting metrics (period: 90d)...

  github     42 data points  .mx-f/data/github/2026-03-20T14:30:00Z.json
  gaze       --              no quality-report artifacts found
  divisor    23 data points  .mx-f/data/divisor/2026-03-20T14:30:00Z.json
  muti-mind  --              no backlog-item artifacts found

Total: 65 data points collected from 2/4 sources.
```

**Error Cases**:

| Scenario | Behavior |
|----------|----------|
| `gh` not authenticated | Report "GitHub: gh auth login required". Skip github source. Exit 0 if other sources succeed. |
| GitHub API rate limit hit | Report rate limit, use cached data if available (with staleness warning). |
| Invalid `--source` value | Error: "unknown source '{value}'. Valid: github, gaze, divisor, muti-mind, all". Exit 1. |
| `.mx-f/data/` not writable | Error with `os.MkdirAll` context. Exit 1. |
| No sources available | Report all sources unavailable with reasons. Exit 1. |

---

### `mxf metrics`

Queries collected metrics and produces reports.

```bash
mxf metrics <subcommand> [flags]
```

**Subcommands**:

| Subcommand | Description |
|------------|-------------|
| `summary` | Consolidated health snapshot across all dimensions |
| `velocity` | Velocity (items completed per sprint) with trend |
| `cycle-time` | Cycle time statistics (avg, median, P90, P99) |
| `bottlenecks` | Identifies the pipeline stage with longest wait time |
| `health` | Traffic-light health dashboard (green/yellow/red) |

**Flags** (shared across all subcommands):

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `text` | Output format: `text` or `json` |
| `--sprints` | int | 0 | Number of sprints to include (0 = all available) |
| `--period` | duration | `30d` | Time range to query. Accepts: `7d`, `30d`, `90d`. |

**Behavior per Subcommand**:

**`mxf metrics summary`**:
Produces a consolidated snapshot: velocity, cycle time,
defect rate, test quality trend, review efficiency, and
overall flow health. With `--format json`, output conforms
to the `metrics-snapshot` artifact envelope (Spec 002).

```text
Metrics Summary (last 30d)
─────────────────────────────
Velocity:         8.2 items/sprint (stable)
Cycle Time:       18.4h avg / 14.2h median / 42.1h P90
Lead Time:        3.2d avg
Defect Rate:      0.12 defects/item
Review Iters:     1.8 avg
CI Pass Rate:     94.3%
Backlog Health:   32 total / 12 ready / 3 stale
Flow Efficiency:  68.4%
```

**`mxf metrics velocity --sprints N`**:
Reports velocity for each of the last N sprints with
trend indicator (increasing/stable/decreasing).

**`mxf metrics cycle-time --period Nd`**:
Reports average, median, P90, and P99 cycle times for
items completed within the period.

**`mxf metrics bottlenecks`**:
Identifies the stage with the longest average wait time.
Reports comparative wait times across all stages.

```text
Bottleneck Analysis (last 30d)
──────────────────────────────
  Review       2.3d avg wait  ████████████████  ← bottleneck
  Testing      0.6d avg wait  ████
  Implementation 0.5d avg wait  ███
  Planning     0.2d avg wait  █

Review is 4x slower than the next stage.
```

**`mxf metrics health`**:
Produces traffic-light indicators for each dimension:
velocity, quality, review efficiency, backlog health.

**JSON Output** (`--format json`):

All subcommands produce a `metrics-snapshot` artifact
envelope when `--format json` is specified:

```json
{
  "hero": "mx-f",
  "version": "0.4.0",
  "timestamp": "2026-03-20T14:30:00Z",
  "artifact_type": "metrics-snapshot",
  "schema_version": "1.0.0",
  "context": {
    "repository": "unbound-force/unbound-force",
    "period": "30d"
  },
  "payload": {
    "velocity": { "value": 8.2, "unit": "items/sprint", "trend": "stable" },
    "cycle_time": { "avg_hours": 18.4, "median_hours": 14.2, "p90_hours": 42.1, "p99_hours": 71.3 },
    "lead_time_days": 3.2,
    "defect_rate": 0.12,
    "review_iterations_avg": 1.8,
    "ci_pass_rate": 0.943,
    "backlog_health": { "total": 32, "ready": 12, "stale": 3 },
    "flow_efficiency": 0.684,
    "health_indicators": [
      { "dimension": "velocity", "status": "green", "value": 8.2, "trend": "stable" },
      { "dimension": "quality", "status": "yellow", "value": 0.12, "trend": "declining" }
    ]
  }
}
```

**Error Cases**:

| Scenario | Behavior |
|----------|----------|
| No collected data exists | Error: "no metrics data found. Run `mxf collect` first." Exit 1. |
| Insufficient data for trend analysis | Report "insufficient data" with required data point count. Exit 0. |
| Invalid `--format` value | Error: "unknown format '{value}'. Valid: text, json". Exit 1. |
| Invalid subcommand | Error: "unknown subcommand '{value}'. Valid: summary, velocity, cycle-time, bottlenecks, health". Exit 1. |

---

### `mxf impediment`

Tracks impediments that block the team's flow.

```bash
mxf impediment <subcommand> [flags]
```

**Subcommands**:

| Subcommand | Description |
|------------|-------------|
| `add` | Log a new impediment |
| `list` | List impediments by status |
| `resolve` | Mark an impediment as resolved |
| `detect` | Proactively detect impediments from metrics trends |

**Flags by Subcommand**:

**`mxf impediment add`**:

| Flag | Type | Required | Default | Description |
|------|------|:--------:|---------|-------------|
| `--title` | string | Yes | — | Short description of the impediment |
| `--severity` | string | No | `medium` | Severity: `critical`, `high`, `medium`, `low` |
| `--owner` | string | No | — | Owner responsible for resolution (e.g., `@dev`). Omitting creates an "unassigned" impediment. |
| `--description` | string | No | — | Detailed description of the impediment |

**Behavior**: Creates `.mx-f/impediments/IMP-NNN.md` with
YAML frontmatter containing all attributes. Auto-assigns
the next sequential IMP-NNN ID.

```text
Created impediment IMP-007: "CI pipeline flaky on Linux" (severity: high, owner: @dev)
```

**`mxf impediment list`**:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | `open` | Filter: `open`, `resolved`, `all` |
| `--format` | string | `text` | Output format: `text` or `json` |

**Behavior**: Lists impediments sorted by severity
(critical first). Impediments open > 14 days are flagged
as "stale" with an escalation recommendation.

```text
Active Impediments
──────────────────
ID       Severity  Age   Owner  Title
IMP-007  high      3d    @dev   CI pipeline flaky on Linux
IMP-004  medium    16d   —      Slow test suite (stale — consider escalation)
IMP-009  low       1d    @lead  Missing API docs

3 open impediments (1 stale, 1 unassigned)
```

**`mxf impediment resolve`**:

```bash
mxf impediment resolve <IMP-NNN> --resolution "<text>"
```

| Argument/Flag | Type | Required | Description |
|---------------|------|:--------:|-------------|
| `<IMP-NNN>` | positional | Yes | Impediment ID to resolve |
| `--resolution` | string | Yes | Description of how the impediment was resolved |

**Behavior**: Updates the impediment file with resolution
text and timestamp. Sets status to `resolved`.

```text
Resolved IMP-007: "CI pipeline flaky on Linux"
Resolution: Pinned CI base image to specific version
```

**`mxf impediment detect`**:

```bash
mxf impediment detect
```

No flags. Reads collected metrics data and identifies
anomalies that suggest emerging impediments. Creates draft
impediments (source: `detected`) for review.

**Behavior**: Analyzes trends for:
- CI failure rate spikes (>15% increase over 7 days)
- Review turnaround increases (>50% above rolling average)
- Velocity drops (>25% below rolling average)
- Stale backlog items (>30 days without movement)

```text
Detected 2 potential impediments:

  IMP-010 (draft)  severity: high
    CI failure rate increased from 5% to 25% over 7 days
    Supporting data: .mx-f/data/github/2026-03-20T14:30:00Z.json

  IMP-011 (draft)  severity: medium
    Review turnaround increased 62% above rolling average
    Supporting data: .mx-f/data/divisor/2026-03-20T14:30:00Z.json

Run `mxf impediment list` to review. Edit draft impediments to confirm.
```

**Error Cases**:

| Scenario | Behavior |
|----------|----------|
| `add` without `--title` | Error: "required flag --title not set". Exit 1. |
| `resolve` with unknown ID | Error: "impediment '{ID}' not found". Exit 1. |
| `resolve` without `--resolution` | Error: "required flag --resolution not set". Exit 1. |
| `detect` with no metrics data | Error: "no metrics data found. Run `mxf collect` first." Exit 1. |
| `add` without `--owner` | Accepted. Impediment flagged as "unassigned" in list. Exit 0. |

---

### `mxf dashboard`

Renders trend visualizations for key metrics.

```bash
mxf dashboard [subcommand] [flags]
```

**Subcommands**:

| Subcommand | Description |
|------------|-------------|
| (none) | Full dashboard: velocity + cycle time + health |
| `velocity` | Velocity bar chart with trend |
| `cycle-time` | Cycle time sparkline with annotations |
| `health` | Traffic-light health indicators with sparklines |

**Flags**:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--html` | bool | false | Generate standalone HTML dashboard |
| `--output` | string | `dashboard.html` | Output file path (only with `--html`) |

**Behavior (text mode)**:

Renders ASCII bar charts, sparklines, and traffic-light
indicators directly in the terminal using `lipgloss` for
styling.

```text
Velocity (last 4 sprints)
─────────────────────────
Sprint 12  ████████████████████  10
Sprint 13  ██████████████████    9
Sprint 14  ██████████████████    9
Sprint 15  ████████████████████████  12
Trend: increasing ↑

Cycle Time (last 30d)
─────────────────────
▂▃▅▇▆▄▃▂▃▃▄▅▃▂▂▃▄▃▂▃▃▂▃▄▃▂▂▃▃▂
Min: 4.2h  Avg: 18.4h  Max: 71.3h

Health
──────
● Velocity     12 items/sprint     improving ↑
● Quality      0.12 defects/item   declining ↓
● Reviews      1.8 iterations avg  stable →
● Backlog      32 items (3 stale)  stable →
```

**Behavior (HTML mode)**:

Generates a standalone HTML file with embedded CSS and
JavaScript (no external dependencies, no server required).
Uses a lightweight charting library vendored inline.

```bash
mxf dashboard --html --output report.html
```

```text
Generated HTML dashboard: report.html (42KB)
```

**Error Cases**:

| Scenario | Behavior |
|----------|----------|
| No collected data | Error: "no metrics data found. Run `mxf collect` first." Exit 1. |
| `--output` without `--html` | Warning: "--output ignored without --html". Render text. Exit 0. |
| Output path not writable | Error with filesystem context. Exit 1. |

---

### `mxf sprint`

Manages sprint lifecycle: planning and review.

```bash
mxf sprint <subcommand> [flags]
```

**Subcommands**:

| Subcommand | Description |
|------------|-------------|
| `plan` | Begin sprint planning with capacity calculation |
| `review` | Summarize completed sprint |

**Flags by Subcommand**:

**`mxf sprint plan`**:

| Flag | Type | Required | Default | Description |
|------|------|:--------:|---------|-------------|
| `--goal` | string | No | — | Sprint goal description |

**Behavior**: Pulls the prioritized backlog from Muti-Mind
artifacts, calculates team capacity from historical
velocity, and suggests a sprint scope.

```text
Sprint Planning: Sprint 16
──────────────────────────
Goal: Implement user authentication
Historical velocity: 10.2 items/sprint (avg last 3)
Suggested capacity: 10 items

Prioritized items from backlog:
  1. BI-042  P1  Implement login endpoint
  2. BI-043  P1  Add session management
  3. BI-038  P1  Create auth middleware
  4. BI-051  P2  Password reset flow
  ...

Suggested scope: items 1-10 (matches historical capacity)
```

**`mxf sprint review`**:

No flags. Summarizes the completed sprint.

**Behavior**: Aggregates data from the current sprint
period: items completed vs. planned, velocity, quality
metrics (from Gaze), acceptance decisions (from
Muti-Mind), and review efficiency (from Divisor).

```text
Sprint Review: Sprint 15
────────────────────────
Completed: 12/10 items (120% of plan)
Velocity:  12 items (up from 9)
Quality:   0.08 defects/item (improved)
Reviews:   1.6 iterations avg (improved)
Accepted:  11/12 items accepted by Muti-Mind

Goal: "Implement search functionality" — ACHIEVED
```

**Error Cases**:

| Scenario | Behavior |
|----------|----------|
| No sprint data available | Error: "no sprint data found. Run `mxf collect` and `mxf sprint plan` first." Exit 1. |
| No Muti-Mind data for planning | Warning: "Muti-Mind backlog not available. Using GitHub issues as backlog source." Exit 0. |

---

### `mxf standup`

Produces a daily standup report.

```bash
mxf standup
```

No flags or subcommands.

**Behavior**: Aggregates current state from all available
sources: items in progress, blocked items (from impediment
tracker), CI/test status (from Gaze), review status (from
Divisor), and items at risk of missing the sprint goal.

```text
Daily Standup — 2026-03-20
──────────────────────────

In Progress (4):
  BI-042  Implement login endpoint       @dev1  day 2
  BI-043  Add session management         @dev2  day 1
  BI-038  Create auth middleware         @dev1  day 3
  BI-051  Password reset flow            @dev3  day 1

Blocked (1):
  BI-038  Create auth middleware
    └─ IMP-007: CI pipeline flaky on Linux (high, 3d open)

CI Status:
  Last run: PASS (94.3% pass rate, 7d rolling)

Reviews Pending (2):
  PR #142  auth middleware (waiting 1.2d, reviewer: @lead)
  PR #145  session storage (waiting 0.3d, reviewer: @dev2)

At Risk:
  BI-038  auth middleware — blocked 3d, sprint ends in 4d
```

**Error Cases**:

| Scenario | Behavior |
|----------|----------|
| No data collected | Error: "no metrics data found. Run `mxf collect` first." Exit 1. |
| No active sprint | Warning: "no active sprint. Showing all in-progress items." Exit 0. |

---

### `mxf retro`

Facilitates structured retrospectives and tracks action
items from previous sessions.

```bash
mxf retro <subcommand> [flags]
```

**Subcommands**:

| Subcommand | Description |
|------------|-------------|
| `start` | Begin a structured retrospective session |
| `actions` | List and track retrospective action items |

**Flags by Subcommand**:

**`mxf retro start`**:

No flags. Begins a structured retrospective following
the five-phase format.

**Behavior**: Produces a structured retrospective record:

1. **Data Presentation**: Pulls metrics from the completed
   sprint and presents key trends.
2. **Pattern Identification**: Highlights recurring themes
   across sprints (e.g., "review iterations increasing
   for 3 consecutive sprints").
3. **Root Cause Analysis**: Presents 5 Whys prompts based
   on identified patterns.
4. **Improvement Proposals**: Captures proposed process
   changes.
5. **Action Items**: Records commitments with owners and
   deadlines (auto-assigned AI-NNN IDs).

Reviews previous action items at the start of each
session and reports status (completed/in-progress/stale).

Output is saved to `.mx-f/retros/{date}-retro.md`.

```text
Retrospective — Sprint 15 (2026-03-20)
───────────────────────────────────────

Previous Action Items:
  AI-003  completed   Add PR template checklist (@lead)
  AI-004  stale       Reduce CI suite to <10min (@dev1, due 2026-03-10)

Sprint Data:
  Velocity:  12 items (↑ from 9)
  Quality:   0.08 defects/item (↑ improved)
  Reviews:   1.6 iterations avg (↑ improved)
  CI:        94.3% pass rate (→ stable)

Patterns Identified:
  • Review iterations decreased after adding PR template (AI-003 impact)
  • CI suite still >15min (AI-004 stale — 10d overdue)

Saved to .mx-f/retros/2026-03-20-retro.json
```

**`mxf retro actions`**:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | `all` | Filter: `pending`, `completed`, `stale`, `all` |

**Behavior**: Lists action items from all retrospectives.
Items past their deadline and not completed are marked
"stale".

```text
Retrospective Action Items
──────────────────────────
ID      Status     Owner   Deadline    Description
AI-001  completed  @dev1   2026-02-28  Automate deployment pipeline
AI-002  completed  @lead   2026-03-05  Write contribution guide
AI-003  completed  @lead   2026-03-15  Add PR template checklist
AI-004  stale      @dev1   2026-03-10  Reduce CI suite to <10min
AI-005  pending    @dev2   2026-03-25  Add integration test for auth flow

5 action items (3 completed, 1 stale, 1 pending)
```

**Error Cases**:

| Scenario | Behavior |
|----------|----------|
| `start` with no metrics data | Warning: "no metrics data available. Retrospective will proceed without data presentation phase." Exit 0. |
| `actions` with no retrospectives | "No action items found. Run `mxf retro start` to begin a retrospective." Exit 0. |
| Invalid `--status` value | Error: "unknown status '{value}'. Valid: pending, completed, stale, all". Exit 1. |

---

## Exit Codes

All commands share a consistent exit code convention:

| Code | Meaning |
|------|---------|
| 0 | Success (including partial collection with graceful degradation) |
| 1 | Error (invalid input, missing required data, filesystem failure) |

## Go API: `MxFParams`

All subcommands delegate to testable functions that accept
a params struct, following the testable CLI pattern
established in the `unbound` binary:

```go
// MxFParams provides dependency injection for all mxf
// subcommands. Enables unit testing without subprocess
// execution or os.Stdout mocking.
type MxFParams struct {
    DataDir      string    // Root data dir (default: ".mx-f")
    Stdout       io.Writer // Summary output writer
    Stderr       io.Writer // Error output writer
    GHRunner     GHRunner  // GitHub API interface (injectable)
    Now          func() time.Time // Clock injection for testing
}
```

Each subcommand has a corresponding `runXxx` function:

```go
func runCollect(params MxFParams, source string, repo string, period string) error
func runMetrics(params MxFParams, sub string, format string, sprints int, period string) error
func runImpediment(params MxFParams, sub string, flags ImpedimentFlags) error
func runDashboard(params MxFParams, sub string, html bool, output string) error
func runSprint(params MxFParams, sub string, goal string) error
func runStandup(params MxFParams) error
func runRetro(params MxFParams, sub string, status string) error
```

## File Ownership

| Pattern | Ownership | Auto-update? |
|---------|-----------|:---:|
| `.opencode/agents/mx-f-coach.md` | User-owned | No |
| `.mx-f/data/**/*.json` | Runtime-generated | N/A |
| `.mx-f/impediments/*.md` | Runtime-generated | N/A |
| `.mx-f/retros/*.json` | Runtime-generated | N/A |

## Storage Layout

```text
.mx-f/
├── data/                          # Metrics storage (runtime)
│   ├── github/
│   │   └── {timestamp}.json       # GitHub API metrics
│   ├── gaze/
│   │   └── {timestamp}.json       # Gaze quality report metrics
│   ├── divisor/
│   │   └── {timestamp}.json       # Divisor review metrics
│   └── mutimind/
│       └── {timestamp}.json       # Muti-Mind backlog metrics
├── impediments/                   # Impediment records (runtime)
│   └── IMP-NNN.md                 # YAML frontmatter + description
└── retros/                        # Retrospective records (runtime)
    └── {date}-retro.json          # Structured retro output
```

All directories are created on first use by the relevant
subcommand. No initialization command is required.
<!-- scaffolded by unbound vdev -->
