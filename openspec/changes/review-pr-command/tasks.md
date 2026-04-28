## 1. Create Standardized Command File

- [ ] 1.1 Create `.opencode/command/review-pr.md` as the
  canonical source (live copy). Adapt from org-infra's
  `review_pr.md` with all standardization deltas:
  - Kebab-case filename (`review-pr.md`)
  - Optional PR number with auto-detection preamble
  - Generic constitution reference (no hardcoded
    principle names)
  - Convention pack awareness with graceful degradation
  - Severity pack reference with inline fallback
  - Proper YAML frontmatter with description
  - All original operational steps preserved (CI
    causality, local tools, scoped diff, spec awareness,
    fix-branch, in-line comments)

## 2. Scaffold Asset Integration

- [ ] 2.1 Copy `.opencode/command/review-pr.md` to
  `internal/scaffold/assets/opencode/command/review-pr.md`
  (embedded scaffold asset copy).
- [ ] 2.2 Update `expectedAssetPaths` in
  `internal/scaffold/scaffold_test.go` — add
  `"opencode/command/review-pr.md"` to the slice
  (alphabetically after `review-council.md`).

## 3. Documentation

- [ ] 3.1 Update AGENTS.md — add `/review-pr` to the
  Project Structure tree under `.opencode/command/`.
- [ ] 3.2 Update AGENTS.md — add a "PR Review" section or
  table documenting when to use `/review-pr` vs
  `/review-council`, including the command's capabilities
  and `gh` CLI prerequisite.

## 4. Verification

- [ ] 4.1 Run `make check` — verify build, test, vet, and
  lint all pass with the new asset file.
- [ ] 4.2 Verify `TestAssetPaths_MatchExpected` passes
  (confirms `expectedAssetPaths` matches actual embedded
  files).
- [ ] 4.3 Verify `TestEmbeddedAssets_MatchSource` passes
  (confirms live copy and scaffold asset are byte-identical
  after version marker insertion).
- [ ] 4.4 Verify `TestRun_CreatesFiles` passes (confirms
  file count includes the new asset).
- [ ] 4.5 Verify `TestIsToolOwned` recognizes the new
  command as tool-owned (no code change needed — all
  `opencode/command/` files are tool-owned by prefix
  match).
- [ ] 4.6 Verify constitution alignment: Composability
  (command works standalone without other UF tools),
  Observable Quality (structured output with severity
  levels), Testability (scaffold test suite covers
  deployment).
