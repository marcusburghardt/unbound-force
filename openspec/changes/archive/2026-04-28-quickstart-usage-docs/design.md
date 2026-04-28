## Context

The project has two documentation extremes: README.md
(75 lines, high-level) and AGENTS.md (800+ lines,
exhaustive). New users fall into the gap between them.
The proposal defines three deliverables: QUICKSTART.md,
USAGE.md, and a README.md update, plus two upstream
issues for agent mode fixes.

## Goals / Non-Goals

### Goals
- 5-minute path from zero to running `/review-council`
- Platform parity: macOS and Fedora treated equally
- Clear orientation for OpenCode modes and agents
- Self-maintaining install instructions (no version rot)
- Serve both maintainer and contributor personas

### Non-Goals
- Replacing AGENTS.md (it remains the deep reference)
- Documenting every slash command (top 15 only)
- Covering Windows (not supported by uf)
- Writing website documentation (separate repo)
- Fixing the agent mode issue in Gaze/Replicator
  (tracked via filed issues, not code in this repo)

## Decisions

### D1: Two files instead of one

**Decision**: Separate QUICKSTART.md (install + first use)
from USAGE.md (workflows + agents).

**Rationale**: Different reading moments. QUICKSTART is
read once during onboarding. USAGE is referenced
repeatedly during daily work. Combining them into one
file would make the quickstart feel too long and the
usage guide hard to find later.

**Composability alignment**: Each file is independently
useful. A maintainer who only wants install instructions
reads QUICKSTART. A contributor who already installed
reads USAGE.

### D2: Recommend Homebrew on Fedora, document dnf as fallback

**Decision**: The primary Fedora path recommends
installing Linuxbrew. A secondary "dnf minimal" path
documents the RPM install + curl for OpenCode.

**Rationale**: `uf setup` relies on Homebrew for 12 of
14 tool installations. Without Homebrew, 5 tools are
skipped entirely (Gaze, Dewey, Replicator, gh, Ollama).
The Homebrew path gives a consistent experience across
platforms. The dnf path is honest about its limitations
but still functional for the minimal `/review-council`
flow.

**Trade-off**: Some Fedora users prefer native package
managers. But the current tooling reality makes Homebrew
the pragmatic choice for the full experience.

### D3: Dynamic RPM version via GitHub API

**Decision**: Use `curl` + GitHub API to resolve the
latest RPM URL instead of hardcoding a version.

**Command**:
```bash
sudo dnf install -y "$(
  curl -fsSL \
    https://api.github.com/repos/unbound-force/unbound-force/releases/latest |
  grep -o 'https://[^"]*linux_amd64\.rpm'
)"
```

**Rationale**: Hardcoded versions go stale. The GitHub
API `/releases/latest` endpoint always resolves to the
current release. The `grep` pattern extracts the RPM
download URL from the JSON response without requiring
`jq` as a dependency.

**Trade-off**: Requires network access to GitHub API at
install time (always true for RPM install anyway).
Rate-limited to 60 requests/hour for unauthenticated
API calls (sufficient for installation).

### D4: /review-council as the "hello world" moment

**Decision**: The quickstart ends with running
`/review-council`, not `/opsx-propose` or `/unleash`.

**Rationale**: `/review-council` requires the fewest
prerequisites (just `uf` + `opencode` + `git`), needs
no external tools (Gaze, Dewey, Replicator all degrade
gracefully), produces immediate visible value (5 AI
reviewers analyzing real code), and works on any existing
codebase without creating spec artifacts first.

### D5: Task-oriented USAGE.md structure

**Decision**: Structure USAGE.md around "I want to..."
recipes rather than tool descriptions.

**Rationale**: Users come with intent ("review my code",
"propose a change"), not tool names. Task-oriented
structure matches their mental model. The agent/mode
section provides orientation, then workflows show how
to accomplish goals.

### D6: README.md pointer instead of duplication

**Decision**: Replace README install section with a
pointer to QUICKSTART.md rather than maintaining both.

**Rationale**: Duplicated install instructions drift.
The README's job is project orientation (what is UF,
constitution, repo contents). Installation details
belong in the quickstart.

### D7: File upstream issues for agent mode fix

**Decision**: File issues in `unbound-force/gaze` and
`unbound-force/replicator` rather than patching agent
files in this repo.

**Rationale**: The agent files are created by `gaze init`
and `replicator init` respectively. Fixing them in this
repo would be overwritten on next init. The fix belongs
at the source. USAGE.md documents the expected behavior
regardless of whether the fix has landed.

## Risks / Trade-offs

### R1: GitHub API rate limiting on RPM install

The dynamic RPM version resolution uses the unauthenticated
GitHub API (60 requests/hour). This is sufficient for
installation but could fail in CI environments that make many
API calls. Mitigation: document that users can replace the
`curl` command with a direct URL if they hit rate limits.

### R2: Linuxbrew recommendation may deter Fedora purists

Some Fedora users strongly prefer native package management.
The documentation is honest about the trade-off: Homebrew
gives the full experience, dnf gives a functional minimum.
Users can adopt Homebrew later without reinstalling `uf`.

### R3: Documentation may drift from tooling

QUICKSTART.md references specific tool names, command syntax,
and the `uf setup` step count (14). If these change,
documentation needs updating. Mitigation: keep references
high-level (e.g., "installs recommended tools" rather than
listing all 14 steps). The AGENTS.md Documentation Validation
Gate already requires documentation impact assessment for
each task.

### R4: Agent mode fix depends on upstream repos

The Gaze and Replicator issues may not be resolved
immediately. USAGE.md must be correct regardless of
whether the fix has landed. The modes/agents section
describes the intended behavior (Build and Plan as primary
modes, everything else as subagents) and does not depend
on the upstream fix being complete.
