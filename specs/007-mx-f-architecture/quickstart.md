# Quickstart: Mx F Architecture (Manager)

**Spec**: 007-mx-f-architecture
**Date**: 2026-03-20

## Prerequisites

- `unbound` CLI installed:
  ```bash
  brew install unbound-force/tap/unbound
  ```
  This installs `unbound`, `mxf`, and `graphthulhu`
  automatically.

- GitHub CLI authenticated (required for GitHub metrics):
  ```bash
  gh auth login
  ```

- OpenCode running in the target project (required for
  the coaching agent).

## Deploy Mx F

### Option 1: Full scaffold (recommended for new projects)

```bash
unbound init
```

This deploys everything: Speckit templates, OpenSpec
schema, Divisor agents, convention packs, the Mx F
coaching agent, and all supporting files.

### Option 2: Coaching agent only

The `mx-f-coach.md` agent is deployed by `unbound init`
alongside all other agents. The `mxf` binary is installed
separately via Homebrew as part of the `unbound` package.

## Quick Start Workflow

The standard Mx F workflow follows four steps:
**collect → metrics → dashboard → retro**.

### 1. Collect metrics

```bash
mxf collect
```

Collects from all available sources (GitHub, Gaze,
Divisor, Muti-Mind). Missing sources are skipped
gracefully.

To collect from a specific source:

```bash
mxf collect --source github --period 30d
```

To collect from a different repository:

```bash
mxf collect --source github --repo unbound-force/gaze
```

### 2. Query metrics

View a consolidated summary:

```bash
mxf metrics summary
```

Check velocity over the last 4 sprints:

```bash
mxf metrics velocity --sprints 4
```

Analyze cycle time for the last 30 days:

```bash
mxf metrics cycle-time --period 30d
```

Find the bottleneck in your pipeline:

```bash
mxf metrics bottlenecks
```

View health indicators:

```bash
mxf metrics health
```

Get JSON output for any query (conforms to the
`metrics-snapshot` artifact envelope):

```bash
mxf metrics summary --format json
```

### 3. Visualize trends

Full dashboard in the terminal:

```bash
mxf dashboard
```

Specific chart:

```bash
mxf dashboard velocity
mxf dashboard cycle-time
mxf dashboard health
```

Generate an HTML report for stakeholders:

```bash
mxf dashboard --html --output sprint-15-report.html
```

### 4. Run a retrospective

Start a structured retrospective:

```bash
mxf retro start
```

This walks through five phases: data presentation,
pattern identification, root cause analysis, improvement
proposals, and action items. Previous action items are
reviewed automatically at the start.

Review action items between retrospectives:

```bash
mxf retro actions
mxf retro actions --status stale
```

## Additional Commands

### Track impediments

```bash
# Log a new impediment
mxf impediment add \
  --title "CI pipeline flaky on Linux" \
  --severity high \
  --owner "@dev"

# List active impediments
mxf impediment list

# Auto-detect impediments from metrics trends
mxf impediment detect

# Resolve an impediment
mxf impediment resolve IMP-007 \
  --resolution "Pinned CI base image to specific version"
```

### Sprint lifecycle

```bash
# Begin sprint planning
mxf sprint plan --goal "Implement user authentication"

# Daily standup report
mxf standup

# Sprint review summary
mxf sprint review
```

## What Gets Deployed

```text
.opencode/
└── agents/
    └── mx-f-coach.md              # Coaching persona (user-owned)

.mx-f/                              # Created at runtime by mxf
├── data/                           # Metrics storage
│   ├── github/
│   │   └── {timestamp}.json
│   ├── gaze/
│   │   └── {timestamp}.json
│   ├── divisor/
│   │   └── {timestamp}.json
│   └── mutimind/
│       └── {timestamp}.json
├── impediments/                    # Impediment records
│   └── IMP-NNN.md
└── retros/                         # Retrospective records
    └── {date}-retro.json
```

The `mx-f-coach.md` agent file is deployed by
`unbound init`. The `.mx-f/` data directory is created
automatically by the `mxf` binary on first use — no
initialization command is required.

## Working Without Other Heroes (Standalone Mode)

Mx F is designed for standalone operation per Org
Constitution Principle II (Composability First). It
works with GitHub as the sole data source when other
heroes are not deployed:

```bash
# Collect only GitHub metrics
mxf collect --source github

# All metrics commands work with GitHub data alone
mxf metrics summary
mxf metrics velocity --sprints 4
mxf dashboard
```

When a source is unavailable, Mx F reports what is
missing and continues with available data:

```text
Collecting metrics (period: 90d)...

  github     42 data points
  gaze       --  no quality-report artifacts found
  divisor    --  no review-verdict artifacts found
  muti-mind  --  no backlog-item artifacts found

Total: 42 data points collected from 1/4 sources.
```

Sprint planning and standup reports adapt to use GitHub
issues instead of Muti-Mind backlog items when
Muti-Mind is not deployed.

## Customizing the Coaching Agent

The coaching agent at `.opencode/agents/mx-f-coach.md`
is user-owned. You can customize it to fit your team's
process and culture.

### Coaching Philosophy

By default, Mx F uses reflective questioning (5 Whys,
mirroring, probing) rather than prescribing solutions.
This is a core design principle (FR-008) — the agent
guides the team to discover root causes themselves.

### Editing the Agent

Open `.opencode/agents/mx-f-coach.md` and modify:

- **Retrospective format**: Adjust the five-phase
  structure if your team prefers a different format
  (e.g., Start/Stop/Continue, 4Ls, sailboat).

- **Coaching style**: Adjust the tone and questioning
  approach. The default is Socratic — you can make it
  more direct for teams that prefer it.

- **Metric thresholds**: Customize the health indicator
  thresholds (green/yellow/red) for your team's context.

- **Process focus**: Emphasize specific areas (e.g.,
  review efficiency, CI reliability) depending on your
  team's current challenges.

### Example: Adding a Custom Focus Area

Add a section to the agent file to highlight a specific
concern:

```markdown
## Current Team Focus

The team is currently focused on reducing review
turnaround time. When facilitating retrospectives or
presenting metrics, give extra attention to review
iteration counts and time-to-approval trends.
```

### Restoring Defaults

If you want to restore the canonical coaching agent,
delete the file and re-run:

```bash
rm .opencode/agents/mx-f-coach.md
unbound init
```

The agent file will be re-created from the embedded
canonical version.
<!-- scaffolded by unbound vdev -->
