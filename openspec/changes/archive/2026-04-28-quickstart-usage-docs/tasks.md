## 1. QUICKSTART.md

- [x] 1.1 Create `QUICKSTART.md` in repo root with
  document header and "What is Unbound Force?" section
  (3 bullet points: review council, spec workflow,
  quality analysis; one-line OpenCode + portability note)
- [x] 1.2 Add "Prerequisites" section: git, LLM API key
  (link to OpenCode provider docs), Go 1.24+ (only for
  Go projects)
- [x] 1.3 Add "Install" section with macOS subsection:
  `brew install unbound-force/tap/unbound-force` one-liner
- [x] 1.4 Add "Install" Fedora subsection (recommended):
  Homebrew install one-liner then `brew install` (same as
  macOS)
- [x] 1.5 Add "Install" Fedora subsection (dnf minimal):
  dynamic RPM install via `curl` + GitHub API (design D3),
  plus OpenCode install via `curl -fsSL
  https://opencode.ai/install | bash`. Include ARM64 note.
- [x] 1.6 Add "For Project Maintainers" section:
  `uf init`, brief summary of what it creates (~5 lines),
  git add/commit/push commands, mention `uf init --divisor`
  for review-only subset
- [x] 1.7 Add "For Contributors" section: `uf setup`,
  `uf doctor`, note about `--dry-run`, note that most
  tools are optional
- [x] 1.8 Add "Your First Review" section: `opencode`
  then `/review-council`, 2-3 lines on what to expect
  (5 AI reviewers, APPROVE/REQUEST CHANGES verdict)
- [x] 1.9 Add "Next Steps" section: pointer to USAGE.md,
  tiered progression (reviews → specs → unleash → full
  setup), pointer to AGENTS.md for deep reference

## 2. USAGE.md

- [x] 2.1 Create `USAGE.md` in repo root with document
  header and "Start OpenCode" section (`opencode` command)
- [x] 2.2 Add "Modes and Agents" section: primary modes
  table (Build, Plan with descriptions), subagents table
  (agent name → invoking command mapping), brief
  explanation of Tab switching vs @mention vs slash
  commands
- [x] 2.3 Add "Common Workflows" section with "Review
  Code" recipe: `/review-council`, what happens, what to
  expect
- [x] 2.4 Add "Common Workflows" recipe for "Propose a
  Change (Small)": `/opsx-propose <desc>` →
  `/opsx-apply` → `/finale`
- [x] 2.5 Add "Common Workflows" recipe for "Build a
  Feature": Speckit pipeline with ASCII diagram
  (`/speckit.specify` → `/speckit.plan` →
  `/speckit.tasks` → `/speckit.implement` → `/finale`)
- [x] 2.6 Add "Common Workflows" recipe for "Go Fully
  Autonomous": `/unleash`, when to use, what it handles
- [x] 2.7 Add "Common Workflows" recipe for "Check Code
  Quality": `/gaze` → `/gaze-fix` (Go projects only)
- [x] 2.8 Add "When to Use What" decision table:
  situation (bug fix, new feature, autonomous, review,
  quality) → workflow (OpenSpec, Speckit, standalone)
  → starting command
- [x] 2.9 Add "Customization" section: convention packs
  at `.opencode/uf/packs/*-custom.md`, brief explanation
  of tool-owned vs user-owned, one-line portability note
- [x] 2.10 Add "Quick Reference" table: 10-15 most-used
  commands with one-line descriptions
- [x] 2.11 Add "See Also" section: Muti-Mind backlog,
  Workflow orchestration, Forge parallel execution, link
  to AGENTS.md

## 3. README.md Update

- [x] 3.1 Replace README.md install section (lines 30-41,
  "Specification Framework" opening with brew/dnf commands)
  with a "Getting Started" pointer to QUICKSTART.md that
  retains the `brew install` one-liner for visibility
- [x] 3.2 Verify remaining README content is unchanged
  (constitution, repo contents, knowledge layer, license
  sections)

## 4. Upstream Issues

- [x] 4.1 File issue in `unbound-force/gaze`: add
  `mode: subagent` to `gaze-reporter.md` and
  `gaze-test-generator.md` agent frontmatter so they
  do not appear as Tab-cyclable primary modes in OpenCode
  (unbound-force/gaze#91)
- [x] 4.2 File issue in `unbound-force/replicator`: add
  `mode: subagent` (or `hidden: true`) to
  `coordinator.md`, `worker.md`, and
  `background-worker.md` agent frontmatter
  (unbound-force/replicator#12)

## 5. Verification

- [x] 5.1 Verify QUICKSTART.md follows the project's
  content convention pack (`.opencode/uf/packs/content.md`)
  for voice and formatting
- [x] 5.2 Verify USAGE.md command references are accurate
  against current `.opencode/command/` directory contents
- [x] 5.3 Verify all internal links resolve (QUICKSTART →
  USAGE, QUICKSTART → AGENTS.md, README → QUICKSTART)
- [x] 5.4 Verify constitution alignment: documentation
  supports Composability First by documenting standalone
  `--divisor` subset and optional tool degradation
- [x] 5.5 Update AGENTS.md "Recent Changes" section with
  a summary of this change
