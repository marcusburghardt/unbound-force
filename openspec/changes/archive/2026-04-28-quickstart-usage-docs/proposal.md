## Why

The project lacks short, actionable onboarding documentation.
The README.md is high-level (75 lines) and points to AGENTS.md
(800+ lines, aimed at AI agents). There is no QUICKSTART.md,
no CONTRIBUTING.md, and no guide for using OpenCode after
`uf init`.

A new user -- whether a project maintainer adopting Unbound
Force or a contributor joining a UF-enabled project -- must
reverse-engineer the onboarding flow (`brew install` then
`uf setup` then `uf init` then `uf doctor`) from Go source
code. The "hello world" moment (`/review-council`) is buried.

Additionally, `uf init` scaffolds 18 agents and 44 slash
commands into a project. After initialization, users see
unexpected extra modes in OpenCode (agents from Gaze and
Replicator that lack `mode: subagent` and appear as
Tab-cyclable primary modes). There is no documentation
explaining the agent landscape, when to use Plan vs Build,
or how slash commands relate to agents.

This change creates two new files and updates one existing
file, plus files two upstream issues.

## What Changes

1. **QUICKSTART.md** (new, ~100 lines): Installation and
   first-use guide covering macOS (Homebrew) and
   Fedora/RHEL (Homebrew recommended, dnf minimal path).
   Covers both maintainer (`uf init`) and contributor
   (`uf setup`) journeys. Ends with running
   `/review-council` as the first tangible value moment.

2. **USAGE.md** (new, ~100 lines): Post-initialization
   guide explaining OpenCode modes and agents, common
   workflows (review, propose, feature, unleash, quality),
   a decision table for workflow selection, customization
   via convention packs, and a quick-reference command
   table.

3. **README.md** (update): Replace the install section
   (lines 30-41) with a pointer to QUICKSTART.md to
   avoid duplication.

4. **Upstream issues**: File issues in `unbound-force/gaze`
   and `unbound-force/replicator` to add `mode: subagent`
   to agents that incorrectly appear as primary modes.

## Capabilities

### New Capabilities
- `quickstart-guide`: QUICKSTART.md with platform-specific
  install instructions (macOS Homebrew, Fedora Homebrew,
  Fedora dnf), maintainer and contributor paths, and
  first-review walkthrough.
- `usage-guide`: USAGE.md with OpenCode modes/agents
  orientation, 5 workflow recipes, decision table,
  customization section, and command quick reference.
- `auto-latest-rpm`: RPM install command uses `curl` +
  GitHub API to resolve latest version dynamically, so
  the documentation never goes stale.

### Modified Capabilities
- `readme-install`: README.md install section replaced
  with pointer to QUICKSTART.md to eliminate duplication.

### Removed Capabilities
- None.

## Impact

- `QUICKSTART.md`: New file in repo root.
- `USAGE.md`: New file in repo root.
- `README.md`: Lines 30-41 replaced with QUICKSTART.md
  pointer. Rest of file unchanged.
- `unbound-force/gaze` (external): Issue filed to add
  `mode: subagent` to `gaze-reporter.md` and
  `gaze-test-generator.md`.
- `unbound-force/replicator` (external): Issue filed to
  add `mode: subagent` to `coordinator.md`, `worker.md`,
  and `background-worker.md`.
- No Go code changes. No scaffold asset changes.
  Documentation only.

## Constitution Alignment

Assessed against the Unbound Force org constitution.

### I. Autonomous Collaboration

**Assessment**: N/A

This change creates documentation files only. It does not
affect inter-hero artifact formats, communication protocols,
or metadata. The upstream agent mode fixes improve agent
discoverability but do not change artifact exchange.

### II. Composability First

**Assessment**: PASS

The documentation explicitly supports standalone adoption.
QUICKSTART.md documents `uf init --divisor` (deploy only
review agents) as a valid subset. The usage guide explains
that most tools are optional with graceful degradation.
This reinforces composability by making the standalone
value of each hero visible to new users.

### III. Observable Quality

**Assessment**: N/A

This change does not produce machine-parseable output or
modify artifact schemas. The RPM install command uses the
GitHub API for version resolution, which is observable and
reproducible.

### IV. Testability

**Assessment**: N/A

Documentation-only change. No components to test in
isolation. The auto-latest RPM command is a shell one-liner
that can be validated manually.
