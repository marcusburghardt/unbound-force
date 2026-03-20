# Data Model: Mx F Architecture (Manager)

**Spec**: 007-mx-f-architecture
**Date**: 2026-03-20

## Entities

### Source Collection

Raw data from a single collection run. One file per source
per run. The lowest-level data entity — all computed
metrics derive from source collections.

**Storage**: `.mx-f/data/{source}/{timestamp}.json`

| Attribute | Type | Notes |
|-----------|------|-------|
| source | string | `github`, `gaze`, `divisor`, `muti-mind` |
| collected_at | ISO 8601 | When collection ran (UTC) |
| data_points | int | Number of items collected |
| raw_data | object | Source-specific JSON payload |

**Identity**: Unique by `(source, collected_at)`.

**Lifecycle**: Created by `mx-f collect --source {name}`.
Immutable after creation — collections are append-only.
Old collections MAY be pruned by retention policy.

### Metrics Snapshot

Point-in-time collection of all computed metrics.
Produced by aggregating one or more Source Collections.
Conforms to the `metrics-snapshot` artifact type
(Spec 002 artifact envelope).

**Storage**: `.mx-f/data/snapshots/{timestamp}.json`

| Attribute | Type | Notes |
|-----------|------|-------|
| timestamp | ISO 8601 | When snapshot was taken (UTC) |
| velocity | float | Items completed per sprint |
| cycle_time | CycleTimeStats | avg/median/p90/p99 in hours |
| lead_time | float | Hours from backlog entry to done |
| defect_rate | float | Defects per completed item |
| review_iterations | float | Average review iterations per PR |
| ci_pass_rate | float | Percentage (0-100) |
| backlog_health | BacklogHealth | total/ready/stale counts |
| flow_efficiency | float | Percentage (0-100) |
| sources_collected | []string | Which sources contributed data |

**Identity**: Unique by `timestamp`.

**Lifecycle**: Created by `mx-f metrics summary`. Each
snapshot is immutable once written. Snapshots accumulate
over time to enable trend analysis.

#### CycleTimeStats (embedded)

Statistical breakdown of cycle time distribution.

| Attribute | Type |
|-----------|------|
| avg | float |
| median | float |
| p90 | float |
| p99 | float |

All values are in hours.

#### BacklogHealth (embedded)

Summary counts of backlog item states.

| Attribute | Type |
|-----------|------|
| total | int |
| ready | int |
| stale | int |

`stale` = items in `draft` or `ready` status for more
than 14 days without update.

### Health Indicator

Traffic-light assessment of a single metric dimension.
Computed from a Metrics Snapshot. Not persisted
independently — embedded in snapshots and sprint state.

| Attribute | Type | Notes |
|-----------|------|-------|
| dimension | string | `velocity`, `quality`, `review`, `backlog`, `flow` |
| status | string | `green`, `yellow`, `red` |
| value | float | Current metric value |
| threshold_green | float | Value at or above which status is green |
| threshold_yellow | float | Value at or above which status is yellow (below green) |
| trend | string | `improving`, `stable`, `declining` |

**Status computation**:
- `value >= threshold_green` → green
- `threshold_yellow <= value < threshold_green` → yellow
- `value < threshold_yellow` → red

**Trend computation**: Derived from comparison of the
current value against the previous 3 snapshots.

### Impediment

Tracked blocker affecting team flow. Stored as Markdown
with YAML frontmatter (same pattern as Muti-Mind backlog
items).

**Storage**: `.mx-f/impediments/IMP-NNN.md`

| Attribute | Type | Notes |
|-----------|------|-------|
| id | string | `IMP-NNN` auto-incrementing |
| title | string | Short description |
| description | string | Markdown body (not in frontmatter) |
| severity | string | `critical`, `high`, `medium`, `low` |
| owner | string | Assigned person (may be `unassigned`) |
| status | string | `open`, `in-progress`, `resolved`, `escalated` |
| created_at | ISO 8601 | When impediment was logged |
| resolved_at | ISO 8601 (optional) | When impediment was resolved |
| resolution | string (optional) | How it was resolved |
| age_days | int | Computed: days since `created_at` |
| source | string | `manual` or `detected` |

**Identity**: Unique by `id`. IDs are auto-incremented
from the highest existing `IMP-NNN` in the directory.

**State transitions**:
```text
  ┌────────┐
  │  open  │
  └───┬────┘
      │
      ▼
┌─────────────┐
│ in-progress │
└──┬──────┬───┘
   │      │
   ▼      ▼
┌────────┐ ┌───────────┐
│resolved│ │ escalated │
└────────┘ └───────────┘
```

**Stale rule**: An impediment with `status = open` and
`age_days > 14` is flagged as stale. Mx F recommends
escalation for stale impediments.

**Lifecycle**: Created by `mx-f impediment add` (manual)
or `mx-f impediment detect` (detected from metrics).
Resolved by `mx-f impediment resolve IMP-NNN`.

### Retrospective Record

Structured record of a retrospective session. Follows
a five-phase facilitation format: data presentation,
pattern identification, root cause analysis, improvement
proposals, and action items.

**Storage**: `.mx-f/retros/YYYY-MM-DD.md`

| Attribute | Type | Notes |
|-----------|------|-------|
| date | date | Session date (YYYY-MM-DD) |
| participants | []string | Team members present |
| data_presented | map | Metrics shown to the team (key-value pairs) |
| patterns_identified | []string | Observed patterns from the data |
| root_causes | []string | Identified root causes |
| improvement_proposals | []string | Proposed improvements |
| action_items | []ActionItem | Tracked commitments |

**Identity**: Unique by `date`. One retrospective per day.

**Lifecycle**: Created by `mx-f retro`. Previous action
items are reviewed at the start of the next retrospective.

#### Action Item (embedded)

Tracked improvement commitment. Embedded in retrospective
records. Action items carry forward across retrospectives
until completed or marked stale.

| Attribute | Type | Notes |
|-----------|------|-------|
| id | string | `AI-NNN` auto-incrementing |
| description | string | What to do |
| owner | string | Responsible person |
| deadline | date | Due date (YYYY-MM-DD) |
| status | string | `pending`, `in-progress`, `completed`, `stale` |
| retrospective_id | string | `YYYY-MM-DD` of source retrospective |

**Stale rule**: An action item with `status != completed`
and `deadline` in the past is marked `stale`.

### Sprint State

Sprint lifecycle tracking. Captures planned scope,
completed work, and end-of-sprint health assessment.

**Storage**: `.mx-f/sprints/{sprint-name}.json`

| Attribute | Type | Notes |
|-----------|------|-------|
| sprint_name | string | Unique identifier (e.g., `sprint-2026-03-20`) |
| goal | string | Sprint goal statement |
| start_date | date | Sprint start (YYYY-MM-DD) |
| end_date | date | Sprint end (YYYY-MM-DD) |
| planned_items | []string | Muti-Mind backlog item IDs (e.g., `BI-001`) |
| completed_items | []string | Items completed during the sprint |
| velocity | float | Computed at sprint end: `len(completed_items)` |
| health_indicators | []HealthIndicator | End-of-sprint health assessment |

**Identity**: Unique by `sprint_name`.

**Lifecycle**: Created by `mx-f sprint plan`. Updated
during the sprint as items complete. Finalized by
`mx-f sprint review`.

### Coaching Interaction

Record of a coaching session between Mx F and the team.
Captures the reflective questioning process and outcome.

**Storage**: `.mx-f/coaching/{timestamp}.json`

| Attribute | Type | Notes |
|-----------|------|-------|
| topic | string | What was discussed |
| questions_asked | []string | Reflective questions posed by Mx F |
| insights_surfaced | []string | Insights the team discovered |
| outcome | string | `action_item`, `escalation`, `resolved`, `deferred` |
| timestamp | ISO 8601 | When session occurred (UTC) |

**Identity**: Unique by `timestamp`.

**Lifecycle**: Created when a coaching session completes.
Coaching interactions feed into the `coaching-record`
artifact type (Spec 002 artifact envelope).

## Relationships

```text
Source Collection
  (github, gaze, divisor, muti-mind)
        │
        │ computation
        ▼
  Metrics Snapshot
        │
        ├──────────────────────────────────┐
        │                                  │
        │ computes                         │ detects
        ▼                                  ▼
  Health Indicator                    Impediment
        │                               (IMP-NNN)
        │ embedded in
        ▼
  Sprint State ──── planned/completed ───► Muti-Mind
  (sprint-name)     items reference        Backlog Items
                                           (BI-NNN)

  Metrics Snapshot ── reviewed at ──► Sprint State
                      sprint end


  Coaching Interaction
        │
        │ may produce
        ▼
  Action Item
        ▲
        │ contains
        │
  Retrospective Record
  (YYYY-MM-DD)
```

**Cross-hero references**:
- `Sprint State.planned_items` and
  `Sprint State.completed_items` reference Muti-Mind
  `backlog-item` IDs (`BI-NNN`).
- `Source Collection.raw_data` contains data ingested from
  Gaze `quality-report`, Divisor `review-verdict`, and
  Muti-Mind `backlog-item`/`acceptance-decision` artifacts.
- `Coaching Interaction` and `Retrospective Record`
  produce `coaching-record` artifacts for consumption by
  other heroes (per Spec 009).

## Validation Rules

1. Source Collection `source` MUST be one of: `github`,
   `gaze`, `divisor`, `muti-mind`.
2. Source Collection `collected_at` MUST be a valid
   ISO 8601 timestamp in UTC.
3. Metrics Snapshot `ci_pass_rate` MUST be in the range
   [0, 100].
4. Metrics Snapshot `flow_efficiency` MUST be in the
   range [0, 100].
5. Metrics Snapshot `velocity`, `lead_time`,
   `defect_rate`, and `review_iterations` MUST be
   non-negative.
6. Metrics Snapshot `sources_collected` MUST contain at
   least one entry.
7. CycleTimeStats values MUST satisfy:
   `avg >= 0`, `median >= 0`, `p90 >= median`,
   `p99 >= p90`.
8. BacklogHealth `total >= ready` and `total >= stale`.
9. Health Indicator `dimension` MUST be one of:
   `velocity`, `quality`, `review`, `backlog`, `flow`.
10. Health Indicator `status` MUST be one of: `green`,
    `yellow`, `red`.
11. Health Indicator `trend` MUST be one of: `improving`,
    `stable`, `declining`.
12. Health Indicator `threshold_green > threshold_yellow`.
13. Impediment `id` MUST match the pattern `IMP-\d{3,}`.
14. Impediment `severity` MUST be one of: `critical`,
    `high`, `medium`, `low`.
15. Impediment `status` MUST be one of: `open`,
    `in-progress`, `resolved`, `escalated`.
16. Impediment `resolved_at` and `resolution` MUST be
    present when `status = resolved` and MUST NOT be
    present when `status = open`.
17. Impediment `age_days` MUST equal the number of days
    between `created_at` and now (or `resolved_at`).
18. Action Item `id` MUST match the pattern `AI-\d{3,}`.
19. Action Item `status` MUST be one of: `pending`,
    `in-progress`, `completed`, `stale`.
20. Action Item `retrospective_id` MUST match a valid
    retrospective date (`YYYY-MM-DD`).
21. Sprint State `end_date >= start_date`.
22. Sprint State `completed_items` MUST be a subset of
    `planned_items` (no item can be completed that was
    not planned).
23. Sprint State `velocity` MUST equal
    `len(completed_items)` when computed at sprint end.
24. Retrospective Record `date` MUST be unique across all
    records (one retrospective per day).
25. Retrospective Record `action_items` entries MUST each
    have a unique `id` within the record.
26. Coaching Interaction `outcome` MUST be one of:
    `action_item`, `escalation`, `resolved`, `deferred`.
27. All timestamps MUST be normalized to UTC. Data sources
    with clock drift exceeding 5 minutes MUST be reported.
<!-- scaffolded by unbound vdev -->
