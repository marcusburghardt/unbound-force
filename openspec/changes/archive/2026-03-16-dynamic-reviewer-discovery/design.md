## Context

The `/review-council` command (`.opencode/command/review-council.md`)
orchestrates five reviewer agents via prose-directed Task tool
delegation. The five agent names are hardcoded in the command text.
The `unbound` scaffold system embeds only four of the five reviewers
(`reviewer-testing` is excluded), and the `review-council.md` command
itself is not embedded. This means any project that adopts the
command will encounter runtime failures if agents are missing.

The org constitution's Composability First principle states heroes
"SHOULD auto-detect the presence of other heroes and activate
enhanced functionality when peers are available." The current
hardcoded list violates this guidance.

## Goals / Non-Goals

### Goals
- The review council dynamically discovers available reviewer agents
  at runtime by reading the filesystem
- The council works with any subset of reviewer agents (1 to N)
- Absent reviewers are reported as informational, not blocking
- The change is contained to a single file (`review-council.md`)
- The approach uses existing OpenCode capabilities (Read tool on
  directories) with no new tooling dependencies

### Non-Goals
- Prerequisite/capability checking per reviewer (e.g., verifying
  `gaze` binary for `reviewer-testing`). This is a future
  enhancement.
- Agent registry or discovery API in OpenCode runtime. This would
  require upstream OpenCode changes.
- Modifying the scaffold system to embed `review-council.md` or
  `reviewer-testing.md`.
- Implementing the full Divisor framework (Spec 005). This change
  is a tactical improvement to the existing prototype.
- Config-driven agent toggling via `permission.task`. The discovery
  approach is automatic and requires no manual configuration.

## Decisions

### D1: File-system discovery via Read tool

**Decision**: Use the Read tool on `.opencode/agents/` to list
directory contents and filter for `reviewer-*.md` files. Extract
agent names by stripping the `.md` extension.

**Rationale**: This requires zero new tooling. The Read tool can
read directories (returns entries one per line with trailing `/`
for subdirectories). The pattern `reviewer-*.md` is a simple string
match implementable in the command prose. This mirrors the existing
`update-agent-context.sh` pattern that uses `[[ -f "$FILE" ]]` for
dynamic agent detection.

**Alternatives considered**:
- Custom discovery tool (`.opencode/tools/discover-agents.ts`):
  More deterministic but adds Node.js/Bun runtime dependency and
  a new file to maintain. Rejected for over-engineering a single
  use case.
- `permission.task` config: Static, requires manual configuration.
  Rejected because it violates the auto-detect guidance in
  Composability First.

### D2: Informational-only reporting for absent reviewers

**Decision**: Absent reviewers are noted in the final report but
do not block the council verdict.

**Rationale**: The council's purpose is to provide the best
available review. A three-reviewer APPROVE is more useful than
failing because a fifth reviewer is not installed. The user sees
which perspectives were missing and can make an informed decision.

**Trade-off**: A project might merge code without a testing review
because `reviewer-testing` was not installed. This is acceptable
because the absence is reported, and the user retains agency.

### D3: Known-roles reference table, not invocation list

**Decision**: Keep the current role descriptions (adversary,
architect, guard, testing, sre) as a reference table in the
command. Use them for context when delegating to discovered
agents, but derive the actual invocation list solely from
discovery.

**Rationale**: The reference table provides rich context for
delegation prompts (each reviewer's focus area). Without it, the
command would delegate with generic prompts, losing the benefit
of targeted instructions. But the table is documentation, not
the source of truth for what gets invoked.

### D4: Support for unknown reviewer agents

**Decision**: If discovery finds a `reviewer-*.md` file not in
the known-roles table (e.g., `reviewer-performance.md`), the
command still invokes it with a generic review delegation prompt.

**Rationale**: This supports extensibility. Teams can add custom
reviewers without modifying the command. Aligns with Composability
First's extension point guidance.

### D5: Single-file change scope

**Decision**: All modifications are contained to
`.opencode/command/review-council.md`. No changes to reviewer
agent files, `opencode.json`, the scaffold system, or `AGENTS.md`.

**Rationale**: Minimizes risk and review surface. The command is
the orchestration point; the agents are independent. This follows
the existing architecture where the command directs delegation
and agents are self-contained.

## Risks / Trade-offs

### R1: LLM interpretation variability

The discovery logic is expressed in prose, not code. Different
LLM runs might interpret the filtering or matching slightly
differently.

**Mitigation**: The instructions are specific and concrete:
"Read `.opencode/agents/` directory, filter for files matching
`reviewer-*.md`, extract names by stripping `.md`." The pattern
is simple enough that interpretation variance is unlikely.

### R2: False positives from non-reviewer files

A file like `reviewer-notes.md` in `.opencode/agents/` would
match the `reviewer-*.md` pattern and be treated as a reviewer
agent.

**Mitigation**: The convention is that agent files in
`.opencode/agents/` ARE agents. A notes file should not be
placed there. The risk is low and the consequence (an extra
Task tool call that fails) is recoverable.

### R3: No prerequisite checking

A discovered agent might exist as a file but be unable to
function (e.g., `reviewer-testing` exists but `gaze` is not
installed).

**Mitigation**: Individual reviewer agents are responsible for
their own error handling. The `gaze-reporter` agent already has
a 3-tier binary resolution strategy. This is explicitly a
non-goal for this change.
