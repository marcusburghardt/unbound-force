# Tasks: Mx F Architecture (Manager)

**Input**: Design documents from `/specs/007-mx-f-architecture/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/mxf-cli.md

**Tests**: Tests are included — the plan specifies 80-90% coverage targets and the constitution's Testability principle mandates coverage strategy.

**Organization**: Tasks are grouped by user story. A foundational phase creates shared infrastructure (CLI skeleton, artifact generalization, storage layer) before user story work begins.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup

**Purpose**: Create the CLI entry point, internal package directories, and project structure.

- [x] T001 Create `cmd/mxf/main.go` with Cobra root command, `MxFParams` struct (per contracts/mxf-cli.md Go API section), `newRootCmdWithParams()` testable constructor, `version`/`commit`/`date` ldflags, persistent `--format` flag, and placeholder subcommand registrations for: `collect`, `metrics`, `impediment`, `dashboard`, `sprint`, `standup`, `retro`.
- [x] T002 [P] Create `internal/metrics/` package directory with placeholder `doc.go`.
- [x] T003 [P] Create `internal/impediment/` package directory with placeholder `doc.go`.
- [x] T004 [P] Create `internal/coaching/` package directory with placeholder `doc.go`.
- [x] T005 [P] Create `internal/dashboard/` package directory with placeholder `doc.go`.
- [x] T006 [P] Create `internal/sprint/` package directory with placeholder `doc.go`.
- [x] T007 Create `cmd/mxf/main_test.go` with `TestMxFParams_Defaults`, `TestNewRootCmd_HasSubcommands` (verify 7 subcommands registered), and `TestRunCollect_NoData` (verify graceful error with no collected data).

---

## Phase 2: Foundational — Artifact Generalization + Storage Layer

**Purpose**: Generalize `internal/artifacts` for multi-hero use (Mx F is both producer and consumer). Create the metrics storage layer used by all user stories.

**CRITICAL**: All user story work depends on the storage layer and artifact reader/writer.

### Artifact Package Generalization

- [x] T008 Refactor `WriteArtifact()` in `internal/artifacts/artifacts.go`: change signature from `writeArtifact(dir, artifactType, id string, payload interface{})` to exported `WriteArtifact(dir, hero, artifactType, id string, payload interface{}) error`. Accept `hero` parameter instead of hardcoding `"muti-mind"`. Update Muti-Mind callers in `internal/artifacts/` to pass `"muti-mind"` explicitly.
- [x] T009 Add `ReadEnvelope(path string) (*Envelope, error)` to `internal/artifacts/artifacts.go`: reads a JSON file, unmarshals into Envelope struct, returns error if file doesn't exist or is malformed.
- [x] T010 Add `FindArtifacts(dir, artifactType string) ([]string, error)` to `internal/artifacts/artifacts.go`: walks the directory tree, returns paths of files whose parsed `artifact_type` field matches the given type. Sorts by timestamp descending (newest first).
- [x] T011 Add tests in `internal/artifacts/artifacts_test.go`: `TestWriteArtifact_CustomHero` (verify hero field is "mx-f"), `TestReadEnvelope_ValidFile`, `TestReadEnvelope_MalformedFile`, `TestFindArtifacts_MultipleTypes`, `TestFindArtifacts_EmptyDir`.

### Metrics Storage Layer

- [x] T012 Create `internal/metrics/models.go` with Go structs per data-model.md: `SourceCollection`, `MetricsSnapshot`, `CycleTimeStats`, `BacklogHealth`, `HealthIndicator`. Include JSON struct tags matching the artifact envelope payload schema from contracts/mxf-cli.md.
- [x] T013 Create `internal/metrics/store.go` with `Store` struct (wraps data dir path), `WriteCollection(source string, data SourceCollection) error` (writes to `.mx-f/data/{source}/{timestamp}.json`), `ReadCollections(source string, since time.Time) ([]SourceCollection, error)` (reads and filters by time), `WriteSnapshot(snapshot MetricsSnapshot) error`, `ReadSnapshots(since time.Time) ([]MetricsSnapshot, error)`.
- [x] T014 Create `internal/metrics/store_test.go` with `TestStore_WriteReadCollection_RoundTrip`, `TestStore_ReadCollections_FilterByTime`, `TestStore_WriteReadSnapshot_RoundTrip`, `TestStore_EmptyDir`. All tests use `t.TempDir()`.

**Checkpoint**: Artifact package supports multi-hero read/write. Metrics storage layer handles JSON file I/O. All foundational tests pass.

---

## Phase 3: User Story 1 — Metrics Collection Platform (Priority: P1) MVP

**Goal**: Implement `mxf collect` with source-specific collectors for GitHub, Gaze, Divisor, and Muti-Mind. Each collector reads from its data source and writes to `.mx-f/data/{source}/`.

**Independent Test**: Configure Mx F for a GitHub repo with 10+ PRs, run collection, verify stored data matches API output.

### GitHub Collector

- [x] T015 [US1] Create `internal/metrics/collect_github.go` with `CollectGitHub(runner sync.GHRunner, repo string, period time.Duration) (*SourceCollection, error)`. Use `gh api` commands via GHRunner to collect: PR count + merge times (`gh api repos/{owner}/{repo}/pulls --json`), review turnaround times, CI pass rate (`gh run list`), commit frequency, issue open/close rates, contributor activity. Parse JSON responses into `SourceCollection.RawData`.
- [x] T016 [US1] Create `internal/metrics/collect_github_test.go` with dispatch-capable `StubGHRunner` that returns different responses for different `gh` commands. Test: `TestCollectGitHub_PRMetrics`, `TestCollectGitHub_CIPassRate`, `TestCollectGitHub_NoAuth` (verify graceful error when `gh` fails), `TestCollectGitHub_RateLimit` (verify detection of rate limit response).

### Hero Artifact Collectors

- [x] T017 [P] [US1] Create `internal/metrics/collect_gaze.go` with `CollectGaze(artifactDir string, since time.Time) (*SourceCollection, error)`. Use `artifacts.FindArtifacts(dir, "quality-report")` to discover Gaze reports. Parse CRAP scores, CRAPload counts, contract coverage, over-specification counts from payloads.
- [x] T018 [P] [US1] Create `internal/metrics/collect_divisor.go` with `CollectDivisor(artifactDir string, since time.Time) (*SourceCollection, error)`. Use `artifacts.FindArtifacts(dir, "review-verdict")` to discover Divisor reports. Parse review iteration counts, finding categories, approval rates, time-to-approval.
- [x] T019 [P] [US1] Create `internal/metrics/collect_mutimind.go` with `CollectMutiMind(artifactDir string, since time.Time) (*SourceCollection, error)`. Use `artifacts.FindArtifacts(dir, "backlog-item")` and `FindArtifacts(dir, "acceptance-decision")`. Parse backlog size, velocity, lead time, acceptance rates.
- [x] T020 [US1] Create `internal/metrics/collect_test.go` with tests for Gaze, Divisor, and Muti-Mind collectors. Use `t.TempDir()` with sample artifact JSON files. Test: `TestCollectGaze_ValidArtifacts`, `TestCollectGaze_NoArtifacts`, `TestCollectDivisor_ValidArtifacts`, `TestCollectMutiMind_ValidArtifacts`.

### Collector Orchestration + CLI Integration

- [x] T021 [US1] Create `internal/metrics/collector.go` with `Collector` struct (holds GHRunner, artifact dirs, store, output writer), `Collect(sources []string, repo string, period time.Duration) error` (orchestrates source-specific collectors, writes results to store, prints summary). Implement graceful degradation (FR-020): skip unavailable sources, report what's missing, exit 0 on partial success.
- [x] T022 [US1] Implement `collect` subcommand in `cmd/mxf/main.go`: register `--source`, `--repo`, `--period` flags per contracts/mxf-cli.md. Delegate to `runCollect(params MxFParams, source, repo, period string) error` which creates a `metrics.Collector` and calls `Collect()`.
- [x] T023 [US1] Create `internal/metrics/collector_test.go` with `TestCollector_AllSources`, `TestCollector_GracefulDegradation` (one source fails, others succeed), `TestCollector_NoSources` (all fail → exit 1).

**Checkpoint**: `mxf collect --source github` retrieves real GitHub data via `gh` CLI. `mxf collect --source all` collects from all available sources with graceful degradation. Data stored in `.mx-f/data/`.

---

## Phase 4: User Story 2 — CLI Metrics Querying (Priority: P1)

**Goal**: Implement `mxf metrics` with subcommands: `summary`, `velocity`, `cycle-time`, `bottlenecks`, `health`. Both `--format text` and `--format json` output.

**Independent Test**: Collect simulated metrics, run queries, verify output matches expected calculations.

### Computation Engine

- [x] T024 [US2] Create `internal/metrics/compute.go` with pure functions: `ComputeVelocity(snapshots []MetricsSnapshot) []VelocityPoint`, `ComputeCycleTime(collections []SourceCollection, period time.Duration) CycleTimeStats`, `ComputeLeadTime(...)`, `ComputeDefectRate(...)`, `ComputeReviewIterations(...)`, `ComputeCIPassRate(...)`, `ComputeBacklogHealth(...)`, `ComputeFlowEfficiency(...)`. All pure — input data in, computed results out.
- [x] T025 [US2] Create `internal/metrics/compute_test.go` with table-driven tests for each computation function using deterministic input data and known expected outputs. Test edge cases: empty input, single data point, zero values.

### Query Engine + Health Indicators

- [x] T026 [US2] Create `internal/metrics/query.go` with `Query` struct (wraps Store), `Summary(period time.Duration) (*MetricsSnapshot, error)` (produces consolidated snapshot), `Velocity(sprints int) ([]VelocityPoint, error)`, `CycleTime(period time.Duration) (*CycleTimeStats, error)`, `Bottlenecks() ([]BottleneckResult, error)` (identifies slowest pipeline stage), `Health() ([]HealthIndicator, error)`.
- [x] T027 [US2] Create `internal/metrics/health.go` with `ComputeHealth(snapshot MetricsSnapshot, previous []MetricsSnapshot) []HealthIndicator`. Compute traffic-light status per dimension (velocity, quality, review, backlog, flow) using configurable thresholds. Compute trend (improving/stable/declining) from previous 3 snapshots.
- [x] T028 [US2] Create `internal/metrics/query_test.go` and `internal/metrics/health_test.go` with tests: `TestQuery_Summary`, `TestQuery_Velocity_MultiSprint`, `TestQuery_Bottlenecks_KnownPlacement` (SC-004), `TestComputeHealth_GreenYellowRed`, `TestComputeHealth_Trend`, `TestQuery_InsufficientData`.

### CLI Integration + JSON Output

- [x] T029 [US2] Implement `metrics` subcommand group in `cmd/mxf/main.go` with `summary`, `velocity`, `cycle-time`, `bottlenecks`, `health` subcommands. Register `--format`, `--sprints`, `--period` flags per contracts/mxf-cli.md. Delegate each to `runMetrics(params, sub, format, sprints, period)`.
- [x] T030 [US2] Implement `--format json` output in `runMetrics`: wrap query results in `artifacts.Envelope` with `hero: "mx-f"`, `artifact_type: "metrics-snapshot"`, produce JSON per contracts/mxf-cli.md JSON example. Implement `--format text` with formatted tabular output per contracts/mxf-cli.md text examples.
- [x] T031 [US2] Create `cmd/mxf/metrics_test.go` with `TestRunMetrics_Summary_Text`, `TestRunMetrics_Summary_JSON` (verify artifact envelope structure), `TestRunMetrics_NoData` (verify error message), `TestRunMetrics_Bottlenecks`.

**Checkpoint**: `mxf metrics summary` produces health snapshot. `mxf metrics bottlenecks` identifies slowest stage. `--format json` produces valid artifact envelopes. SC-003 (5-second query) verified.

---

## Phase 5: User Story 3 — AI Coaching and Retrospective Facilitation (Priority: P2)

**Goal**: Create the `mx-f-coach.md` OpenCode agent and the retrospective record storage engine.

**Independent Test**: Present coaching agent with team scenario, verify reflective questions (not solutions). Run `mxf retro start`, verify structured 5-section output.

### Coaching Agent

- [x] T032 [US3] Create `.opencode/agents/mx-f-coach.md` with: YAML frontmatter (description: "Flow Facilitator and Continuous Improvement Coach", mode: agent, model, temperature: 0.3), H1 Role (coaching philosophy — facilitate, don't prescribe), H2 Source Documents (AGENTS.md, constitution, `.mx-f/data/` metrics, `.mx-f/impediments/`, `.mx-f/retros/`, graphthulhu MCP conditional), H2 Coaching Framework (5 Whys technique, reflective questioning, mirroring, probing — with examples), H2 Retrospective Facilitation Protocol (5-phase format: data presentation → pattern identification → root cause analysis → improvement proposals → action items), H2 Boundary Rules (FR-021: no prescriptions, redirect technical questions to Cobalt-Crush/Gaze/Divisor).

### Retrospective Engine

- [x] T033 [US3] Create `internal/coaching/models.go` with Go structs per data-model.md: `RetroRecord`, `ActionItem`, `CoachingInteraction`. Include YAML struct tags for frontmatter and JSON tags for serialization.
- [x] T034 [US3] Create `internal/coaching/retro.go` with `RetroStore` struct (wraps retro dir path), `StartRetro(date string, metricsData map[string]interface{}) (*RetroRecord, error)` (creates new retro with data presentation phase pre-populated from metrics), `SaveRetro(record *RetroRecord) error` (writes to `.mx-f/retros/YYYY-MM-DD.md`), `LoadRetro(date string) (*RetroRecord, error)`, `ListRetros() ([]RetroRecord, error)`.
- [x] T035 [US3] Create `internal/coaching/actions.go` with `NextActionID(retros []RetroRecord) string` (scans all retros for highest AI-NNN, returns next), `ReviewPreviousActions(retros []RetroRecord) []ActionItem` (finds all non-completed action items from previous retros and marks stale if past deadline).
- [x] T036 [US3] Create `internal/coaching/retro_test.go` with `TestRetroStore_SaveLoad_RoundTrip`, `TestRetroStore_ListRetros`, `TestNextActionID`, `TestReviewPreviousActions_StaleDetection`.

### CLI Integration

- [x] T037 [US3] Implement `retro` subcommand group in `cmd/mxf/main.go` with `start` and `actions` subcommands. `retro start`: loads latest metrics, reviews previous action items, creates new retro record, saves it. `retro actions --status`: lists action items filtered by status. Delegate to `runRetro(params, sub, status)`.

**Checkpoint**: Coaching agent deployed. `mxf retro start` produces structured 5-section record. Previous action items reviewed and stale items flagged.

---

## Phase 6: User Story 4 — Impediment Tracking (Priority: P2)

**Goal**: Implement `mxf impediment` with add/list/resolve/detect subcommands. Proactive detection from metrics trends.

**Independent Test**: Add impediments, list them, resolve one, run detect on simulated metrics with CI spike. Verify round-trip and detection.

### Impediment CRUD

- [x] T038 [US4] Create `internal/impediment/models.go` with Go structs per data-model.md: `Impediment` (YAML frontmatter fields + computed `AgeDays`), `ImpedimentList` sort helpers.
- [x] T039 [US4] Create `internal/impediment/impediment.go` with `Repository` struct (wraps impediment dir path), `Add(title, severity, owner, description string) (*Impediment, error)` (auto-assigns IMP-NNN ID, writes `.mx-f/impediments/IMP-NNN.md`), `List(statusFilter string) ([]Impediment, error)` (reads all, filters by status, sorts by severity), `Resolve(id, resolution string) error` (updates status + resolution + timestamp), `Get(id string) (*Impediment, error)`.
- [x] T040 [US4] Create `internal/impediment/impediment_test.go` with `TestRepository_Add_AutoID`, `TestRepository_List_SortBySeverity`, `TestRepository_Resolve_UpdatesFile`, `TestRepository_List_StaleDetection` (age > 14 days flagged), `TestRepository_Add_NoOwner` (unassigned).

### Proactive Detection

- [x] T041 [US4] Create `internal/impediment/detect.go` with `Detect(metricsStore *metrics.Store, repo *Repository) ([]Impediment, error)`. Analyze metrics trends for: CI failure rate spike (>15% increase over 7 days), review turnaround increase (>50% above rolling avg), velocity drop (>25% below rolling avg). Create draft impediments (source: "detected") for each anomaly.
- [x] T042 [US4] Create `internal/impediment/detect_test.go` with `TestDetect_CISpike` (SC-008), `TestDetect_ReviewTurnaroundIncrease`, `TestDetect_VelocityDrop`, `TestDetect_NoAnomalies`.

### CLI Integration

- [x] T043 [US4] Implement `impediment` subcommand group in `cmd/mxf/main.go` with `add`, `list`, `resolve`, `detect` subcommands. Register flags per contracts/mxf-cli.md. Delegate to `runImpediment(params, sub, flags)`.
- [x] T044 [US4] Create `cmd/mxf/impediment_test.go` with `TestRunImpediment_AddListResolve_RoundTrip` (SC-007), `TestRunImpediment_Detect_CISpike` (SC-008).

**Checkpoint**: Impediment CRUD works. Proactive detection identifies CI spikes from metrics data. Stale impediments flagged.

---

## Phase 7: User Story 5 — Dashboard and Visualization (Priority: P3)

**Goal**: Implement `mxf dashboard` with ASCII charts (sparklines, bar charts) and optional HTML output.

**Independent Test**: Collect simulated metrics, render dashboard, verify charts represent data correctly.

### Text Dashboard

- [x] T045 [P] [US5] Create `internal/dashboard/text.go` with: `RenderBarChart(data []BarChartPoint, w io.Writer) error` (ASCII bar chart with lipgloss styling), `RenderSparkline(data []float64, w io.Writer) error` (Unicode sparkline characters ▁▂▃▄▅▆▇█), `RenderHealthIndicators(indicators []metrics.HealthIndicator, w io.Writer) error` (traffic-light ● with color via lipgloss).
- [x] T046 [P] [US5] Create `internal/dashboard/text_test.go` with `TestRenderBarChart_4Sprints`, `TestRenderSparkline_30Days`, `TestRenderHealthIndicators_GreenYellowRed`. Verify output contains expected characters and values.

### HTML Dashboard

- [x] T047 [US5] Create `internal/dashboard/html.go` with `RenderHTML(snapshot metrics.MetricsSnapshot, indicators []metrics.HealthIndicator, outPath string) error`. Use `html/template` to generate standalone HTML file with embedded Chart.js CDN link for interactive charts.
- [x] T048 [US5] Create `internal/dashboard/html_test.go` with `TestRenderHTML_ProducesValidFile` (verify file exists, contains Chart.js reference, contains data values).

### CLI Integration

- [x] T049 [US5] Implement `dashboard` subcommand group in `cmd/mxf/main.go` with `velocity`, `cycle-time`, `health` subcommands (plus bare `dashboard` for full view). Register `--html`, `--output` flags per contracts/mxf-cli.md. Delegate to `runDashboard(params, sub, html, output)`.

**Checkpoint**: `mxf dashboard` renders ASCII charts in terminal. `mxf dashboard --html` generates standalone HTML file.

---

## Phase 8: User Story 6 — Swarm Coordination (Priority: P3)

**Goal**: Implement `mxf sprint` (plan/review), `mxf standup`, and cross-hero pattern identification.

**Independent Test**: Simulate sprint lifecycle — plan, standup, review — verify outputs at each stage.

### Sprint Lifecycle

- [x] T050 [US6] Create `internal/sprint/models.go` with Go structs per data-model.md: `SprintState`.
- [x] T051 [US6] Create `internal/sprint/sprint.go` with `SprintStore` struct (wraps sprint dir path), `Plan(goal string, velocity float64, backlogItems []string) (*SprintState, error)` (creates sprint, calculates capacity from historical velocity, suggests scope), `Review(sprintName string, metricsStore *metrics.Store) (*SprintState, error)` (aggregates completed items, velocity, quality metrics), `Save/Load` for JSON persistence at `.mx-f/sprints/{name}.json`.
- [x] T052 [US6] Create `internal/sprint/sprint_test.go` with `TestSprintStore_PlanAndReview`, `TestSprintPlan_CapacityCalculation`.

### Standup + Pattern Identification

- [x] T053 [US6] Add `Standup(sprintStore *SprintStore, impedimentRepo *impediment.Repository, metricsStore *metrics.Store, w io.Writer) error` function to `internal/sprint/sprint.go`. Aggregates: items in progress, blocked items (from impediments), CI status, review status, at-risk items. Output per contracts/mxf-cli.md standup format.
- [x] T054 [US6] Add `IdentifyPatterns(snapshots []metrics.MetricsSnapshot, divisorFindings []artifacts.Envelope) ([]string, error)` to `internal/sprint/sprint.go`. Analyzes cross-hero data for process patterns (FR-016): frequent Divisor findings → convention pack recommendation, velocity trends → capacity adjustments.
- [x] T055 [US6] Create `internal/sprint/standup_test.go` with `TestStandup_FullReport`, `TestStandup_NoActiveSprint`, `TestIdentifyPatterns_FrequentDocFindings`.

### CLI Integration

- [x] T056 [US6] Implement `sprint` subcommand group in `cmd/mxf/main.go` with `plan` and `review` subcommands. Register `--goal` flag for plan. Implement `standup` subcommand (no flags). Delegate to `runSprint(params, sub, goal)` and `runStandup(params)`.

**Checkpoint**: `mxf sprint plan` calculates capacity and suggests scope. `mxf standup` produces daily report. `mxf sprint review` summarizes completed sprint. Pattern identification produces cross-hero recommendations.

---

## Phase 9: Scaffold Integration + GoReleaser

**Purpose**: Embed `mx-f-coach.md` in the scaffold engine, add `cmd/mxf/` to GoReleaser, run all tests.

- [x] T057 Copy `.opencode/agents/mx-f-coach.md` to `internal/scaffold/assets/opencode/agents/mx-f-coach.md`.
- [x] T058 Add `"opencode/agents/mx-f-coach.md"` to `expectedAssetPaths` in `internal/scaffold/scaffold_test.go`. Update agent count comment.
- [x] T059 Update `cmd/unbound/main_test.go`: change "46 files processed" to "47 files processed" in `TestRunInit_FreshDir`.
- [x] T060 Add `cmd/mxf/` build entry to `.goreleaser.yaml` per research.md R10: second `builds` entry with same `CGO_ENABLED=0`, `goos`, `goarch`, and ldflags pattern as the `unbound` build.
- [x] T061 Run `go test -race -count=1 ./...` and verify all tests pass including drift detection for the new agent file. Fix any failures.
- [x] T062 Run `go build ./...` and verify both `cmd/unbound` and `cmd/mxf` build successfully.

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, status updates, and final validation.

- [x] T063 [P] Update `AGENTS.md`: change Mx F status in Heroes table from "Spec only (007)" to "Implemented (Spec 007)". Add `mx-f-coach.md` to project structure tree. Add `cmd/mxf/` and `internal/metrics/`, `internal/impediment/`, `internal/coaching/`, `internal/dashboard/`, `internal/sprint/` to project structure. Add Spec 007 to Recent Changes.
- [x] T064 [P] Update `README.md`: change Mx F status from "Spec only" to "Implemented". Update file count from "46 files" to "47 files". Add `mxf` CLI to the available tools section.
- [x] T065 [P] Update `specs/007-mx-f-architecture/spec.md`: change `status: draft` to `status: complete` in frontmatter and body.
- [x] T066 [P] Update `unbound-force.md`: verify Mx F hero description reflects implemented status.
- [x] T067 Run `make check` (or `go build ./... && go test -race -count=1 ./... && go vet ./...`) and verify all checks pass.
- [x] T068 Verify SC-001 through SC-010 success criteria from spec.md are met. Document verification results.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Phase 1 — artifact generalization and storage layer needed by all user stories
- **US1 (Phase 3)**: Depends on Phase 2 — collectors use artifact reader and metrics store
- **US2 (Phase 4)**: Depends on Phase 3 — queries operate on collected data
- **US3 (Phase 5)**: Depends on Phase 2 — retro engine uses metrics store. Can run in parallel with US1/US2.
- **US4 (Phase 6)**: Depends on Phase 2 — impediment repo uses storage pattern. Detect needs metrics store. Can run in parallel with US1/US2.
- **US5 (Phase 7)**: Depends on Phase 4 — dashboard renders query results
- **US6 (Phase 8)**: Depends on Phases 3-6 — sprint uses metrics, impediments, and artifacts
- **Scaffold (Phase 9)**: Depends on Phase 5 — coaching agent must exist before embedding
- **Polish (Phase 10)**: Depends on all phases

### Parallel Opportunities

- **Phase 1**: T002-T006 (directory creation) can run in parallel
- **Phase 2**: T008-T011 (artifacts) and T012-T014 (storage) are sequential within each group but the groups can partially overlap
- **Phase 3**: T017-T019 (hero collectors) can run in parallel
- **Phase 5**: Can run in parallel with Phase 3/4 (different packages)
- **Phase 6**: Can run in parallel with Phase 3/4 (different packages)
- **Phase 7**: T045-T046 (text) and T047-T048 (HTML) can run in parallel
- **Phase 10**: T063-T066 (doc updates) can run in parallel

### User Story Parallel Paths

After Phase 2 completes:
- **Path A**: US1 (collect) → US2 (query) → US5 (dashboard) → US6 (sprint)
- **Path B**: US3 (coaching/retro) — independent
- **Path C**: US4 (impediments) — independent, feeds into US6

---

## Implementation Strategy

### MVP First (Phases 1-4)

1. Phase 1: CLI skeleton + package directories
2. Phase 2: Artifact generalization + storage layer
3. Phase 3: US1 — Metrics collection (GitHub + hero artifacts)
4. Phase 4: US2 — Metrics querying + health indicators
5. **STOP and VALIDATE**: `mxf collect && mxf metrics summary` produces valid output

### Incremental Delivery

1. Phases 1-4 → `mxf collect` + `mxf metrics` (MVP)
2. Phase 5 → `mxf retro` + coaching agent (Mx F's differentiator)
3. Phase 6 → `mxf impediment` (obstacle tracking + detection)
4. Phase 7 → `mxf dashboard` (visualization)
5. Phase 8 → `mxf sprint` + `mxf standup` (swarm coordination)
6. Phase 9 → Scaffold + GoReleaser integration
7. Phase 10 → Documentation + validation

### Estimated Time

| Phase | Tasks | Est. Time |
|-------|-------|-----------|
| Phase 1: Setup | 7 | 30 min |
| Phase 2: Foundational | 6 | 1.5 hrs |
| Phase 3: US1 Collection | 9 | 2.5 hrs |
| Phase 4: US2 Querying | 8 | 2 hrs |
| Phase 5: US3 Coaching | 6 | 1.5 hrs |
| Phase 6: US4 Impediments | 7 | 1.5 hrs |
| Phase 7: US5 Dashboard | 5 | 1.5 hrs |
| Phase 8: US6 Sprint | 7 | 2 hrs |
| Phase 9: Scaffold | 6 | 30 min |
| Phase 10: Polish | 6 | 30 min |
| **Total** | **67** | **~14 hours** |
<!-- scaffolded by unbound vdev -->
