## ADDED Requirements

### Requirement: QUICKSTART.md existence

The repository MUST contain a `QUICKSTART.md` file in the
repo root providing installation and first-use instructions.

#### Scenario: New user finds quickstart

- **GIVEN** a user visits the repository for the first time
- **WHEN** they look for onboarding documentation
- **THEN** `QUICKSTART.md` exists in the repo root and is
  linked from `README.md`

### Requirement: Platform coverage

QUICKSTART.md MUST document installation for macOS
(Homebrew) and Fedora/RHEL (Homebrew recommended, dnf
fallback). Both paths MUST lead to a functional
`uf` + `opencode` installation.

#### Scenario: macOS user installs via Homebrew

- **GIVEN** a macOS user with Homebrew installed
- **WHEN** they follow the QUICKSTART.md macOS section
- **THEN** `uf` and `opencode` are available in PATH

#### Scenario: Fedora user installs via Homebrew

- **GIVEN** a Fedora user without Homebrew
- **WHEN** they follow the QUICKSTART.md Fedora (recommended)
  section
- **THEN** Homebrew is installed and `uf` is available in PATH

#### Scenario: Fedora user installs via dnf

- **GIVEN** a Fedora user who prefers native packages
- **WHEN** they follow the QUICKSTART.md Fedora (dnf) section
- **THEN** `uf` is installed via RPM and `opencode` is
  installed via curl, both available in PATH

### Requirement: Self-maintaining RPM version

The dnf install command MUST resolve the latest RPM version
dynamically via the GitHub API. The command MUST NOT contain
a hardcoded version string.

#### Scenario: New release published

- **GIVEN** a new unbound-force release is published
- **WHEN** a user runs the dnf install command from the docs
- **THEN** the latest RPM is installed without documentation
  changes

### Requirement: Dual persona coverage

QUICKSTART.md MUST document both the maintainer journey
(`uf init` to add UF to a project) and the contributor
journey (`uf setup` + `uf doctor` to set up a development
environment).

#### Scenario: Maintainer adopts UF

- **GIVEN** a project maintainer with `uf` installed
- **WHEN** they follow the maintainer section
- **THEN** `uf init` scaffolds files and the section
  explains what to commit

#### Scenario: Contributor sets up environment

- **GIVEN** a contributor cloning a UF-enabled project
- **WHEN** they follow the contributor section
- **THEN** `uf setup` installs recommended tools and
  `uf doctor` verifies the environment

### Requirement: First-value endpoint

QUICKSTART.md MUST end with a "Your First Review" section
that walks the user through starting OpenCode and running
`/review-council`.

#### Scenario: User runs first review

- **GIVEN** a user who completed the install and init steps
- **WHEN** they follow the "Your First Review" section
- **THEN** they start `opencode` and run `/review-council`
  successfully

### Requirement: USAGE.md existence

The repository MUST contain a `USAGE.md` file in the repo
root explaining how to use OpenCode with the scaffolded
agents and commands.

#### Scenario: User looks for workflow guidance

- **GIVEN** a user who completed the quickstart
- **WHEN** they want to know how to use UF for daily work
- **THEN** `USAGE.md` exists and is linked from
  `QUICKSTART.md`

### Requirement: Agent and mode orientation

USAGE.md MUST contain a section explaining OpenCode primary
modes (Build, Plan) and subagents (Divisor, Cobalt-Crush,
Muti-Mind, etc.), including which slash commands invoke which
agents.

#### Scenario: User sees unexpected modes

- **GIVEN** a user pressing Tab in OpenCode after `uf init`
- **WHEN** they see modes beyond Build and Plan
- **THEN** USAGE.md explains what each mode/agent is and
  when to use it

### Requirement: Workflow recipes

USAGE.md MUST contain at least 5 task-oriented workflow
recipes: code review, small change, new feature, autonomous
pipeline, and quality check.

#### Scenario: User wants to review code

- **GIVEN** a user reading the "Review Code" recipe
- **WHEN** they follow the steps
- **THEN** they successfully run `/review-council`

#### Scenario: User wants to propose a small change

- **GIVEN** a user reading the "Propose a Change" recipe
- **WHEN** they follow the steps
- **THEN** they run `/opsx-propose`, `/opsx-apply`, and
  `/finale` in sequence

### Requirement: Decision table

USAGE.md MUST contain a decision table mapping situations
(bug fix, new feature, autonomous) to workflows (OpenSpec,
Speckit) and starting commands.

#### Scenario: User unsure which workflow to use

- **GIVEN** a user with a task to complete
- **WHEN** they consult the decision table
- **THEN** they can identify the correct workflow and
  starting command

### Requirement: Command quick reference

USAGE.md MUST contain a quick-reference table of the 10-15
most commonly used slash commands with one-line descriptions.

#### Scenario: User looking up a command

- **GIVEN** a user who knows they need a command but not
  the exact name
- **WHEN** they consult the quick reference table
- **THEN** they find the command and its description

## MODIFIED Requirements

### Requirement: README.md install section

README.md install section (currently lines 30-41) MUST be
replaced with a pointer to QUICKSTART.md. The pointer MUST
include the `brew install` one-liner for immediate
visibility but defer detailed platform instructions to the
quickstart.

Previously: README contained full install instructions
inline with a version placeholder for RPM.

#### Scenario: User reads README

- **GIVEN** a user reading README.md
- **WHEN** they reach the install section
- **THEN** they see a brief install command and a link to
  QUICKSTART.md for full instructions

## REMOVED Requirements

None.
