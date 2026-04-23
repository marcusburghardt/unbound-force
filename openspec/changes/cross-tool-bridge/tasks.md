## Tasks

### Part 1: Implement ensureAGENTSmdPackSection()

- [x] Add `agentsmdPackMarker` constant (heading `## Convention Packs`)
- [x] Implement `ensureAGENTSmdPackSection()` that appends Convention Packs section to AGENTS.md listing deployed packs -- idempotent via heading detection, uses opts.ReadFile/WriteFile
- [x] Add `collectDeployedPacks()` helper that enumerates which pack files were deployed based on resolved language (reuse shouldDeployPack logic)

### Part 2: Implement ensureCLAUDEmd()

- [x] Add `claudemdMarker` constant (`# Unbound Force — managed by uf init`)
- [x] Implement `ensureCLAUDEmd()` that creates or appends managed block to CLAUDE.md with @imports for AGENTS.md, cobalt-crush agent, deployed convention packs, and on-demand review agent references -- idempotent via marker detection, uses opts.ReadFile/WriteFile

### Part 3: Implement ensureCursorrules()

- [x] Add `cursorrulesMarker` constant (same marker pattern)
- [x] Implement `ensureCursorrules()` that creates or appends managed block to .cursorrules with pack reference instructions and agent file references -- idempotent via marker detection, uses opts.ReadFile/WriteFile

### Part 4: Wire into Run() and printSummary()

- [x] Call ensureAGENTSmdPackSection(), ensureCLAUDEmd(), ensureCursorrules() after ensureGitignore() in Run(), append results to subResults
- [x] Verify printSummary() handles new subToolResult entries (existing format should work)

### Part 5: Agent references in bridge files

- [x] Add cobalt-crush-dev.md @import and Divisor review agent on-demand list to ensureCLAUDEmd() generated content
- [x] Add cobalt-crush-dev.md reference and Divisor review agent list to ensureCursorrules() generated content
- [x] Update proposal.md with agent reference design rationale and command file exclusion

### Part 6: Tests

- [x] Add TestEnsureAGENTSmdPackSection (fresh, existing without section, existing with section, idempotent)
- [x] Add TestEnsureCLAUDEmd (fresh, existing without marker, existing with marker, idempotent)
- [x] Add TestEnsureCursorrules (fresh, existing without marker, existing with marker, idempotent)
- [x] Add TestCollectDeployedPacks for different languages (go, typescript, default)
- [x] Verify agent references in CLAUDE.md tests (cobalt-crush @import, divisor-guard reference)
- [x] Verify agent references in .cursorrules tests (cobalt-crush reference, divisor-guard reference)

### Part 7: Verification

- [x] Run `go test -race -count=1 ./...` to verify all tests pass
- [x] Run `golangci-lint run` to verify no lint issues
