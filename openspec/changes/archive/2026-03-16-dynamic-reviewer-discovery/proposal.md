## Why

The `/review-council` command hardcodes all five reviewer agent
names (`reviewer-adversary`, `reviewer-architect`, `reviewer-guard`,
`reviewer-testing`, `reviewer-sre`) in its prose instructions. If any
agent file is missing, the Task tool invocation fails at runtime with
no graceful fallback.

This is a real problem because:

1. **`reviewer-testing` is not embedded** in the `unbound` binary.
   Projects scaffolded with `unbound init` receive only four of the
   five reviewers but no orchestration command. If a project adds
   the `review-council.md` command, it will try to invoke an agent
   that does not exist.
2. **The `review-council.md` command itself is not embedded** -- it
   is local to the meta-repo and the Gaze prototype. Any project
   that copies or adapts it inherits the hardcoded five-agent
   assumption.
3. **Future hero repos** may have different subsets of reviewers
   installed depending on which heroes are deployed.

The org constitution's Composability First principle (II) explicitly
states: "Heroes SHOULD auto-detect the presence of other heroes and
activate enhanced functionality when peers are available, without
requiring manual configuration." The review council violates this
by assuming all five agents are always present.

## What Changes

Modify `review-council.md` to add a **discovery step** before
delegating to reviewer agents. The command will:

1. Read the `.opencode/agents/` directory to find files matching
   `reviewer-*.md`
2. Build the invocation list from discovered agents only
3. Skip absent reviewers with an informational note
4. Base the council verdict on discovered reviewers only

This is a single-file change to `.opencode/command/review-council.md`.

## Capabilities

### New Capabilities
- `dynamic-reviewer-discovery`: The review council discovers
  available reviewer agents at runtime by scanning the
  `.opencode/agents/` directory, rather than assuming a fixed set.

### Modified Capabilities
- `review-council verdict`: Verdict policy changes from "all five
  reviewers must APPROVE" to "all discovered reviewers must APPROVE"
  with an informational note listing absent reviewers.
- `review-council final report`: Both Code Review and Spec Review
  mode reports include a discovery summary showing which reviewers
  were invoked and which known reviewers were absent.

### Removed Capabilities
- None.

## Impact

| Item | Impact |
|------|--------|
| `.opencode/command/review-council.md` | Modified: discovery step added, delegation made dynamic, verdict policy updated, final report includes discovery summary |
| Reviewer agent files | No changes |
| `opencode.json` | No changes |
| Scaffold system (`internal/scaffold/`) | No changes |
| `AGENTS.md` | No changes needed |
| Spec 005 (The Divisor) | No changes (still draft; this aligns with its intended direction for pluggable reviewers) |

## Constitution Alignment

Assessed against the Unbound Force org constitution.

### I. Autonomous Collaboration

**Assessment**: PASS

The review council continues to communicate through structured
verdicts (APPROVE / REQUEST CHANGES). Discovery reads agent files
from the filesystem -- no runtime coupling or synchronous
interaction between agents is introduced. Each reviewer remains
independently invocable.

### II. Composability First

**Assessment**: PASS

This change directly advances Composability First. The constitution
states heroes "SHOULD auto-detect the presence of other heroes and
activate enhanced functionality when peers are available." Dynamic
discovery replaces a hardcoded assumption of five agents with
runtime detection, allowing the council to function with any
subset of reviewers. Projects with only four embedded reviewers
will work without modification.

### III. Observable Quality

**Assessment**: PASS

The final report now includes a discovery summary listing which
reviewers were invoked and which were absent. This makes the
review's completeness observable and auditable. The council
verdict format (APPROVE / REQUEST CHANGES) is unchanged.

### IV. Testability

**Assessment**: PASS

The change is testable by:
1. Running `/review-council` with all five agents present
   (should behave identically to current behavior).
2. Temporarily removing an agent file and running again
   (should invoke remaining agents, report the absence).
3. No external services or shared mutable state required.
