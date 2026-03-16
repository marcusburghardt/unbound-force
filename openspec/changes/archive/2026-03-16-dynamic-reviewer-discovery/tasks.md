## 1. Add Discovery Step

- [x] 1.1 Insert a new "Discover Available Reviewers" section
  between the "Determine Review Mode" section and the mode-specific
  instructions in `review-council.md`. This section instructs the
  agent to read the `.opencode/agents/` directory using the Read
  tool and filter entries for files matching `reviewer-*.md`.
- [x] 1.2 Add instructions to extract agent names by stripping the
  `.md` extension from discovered filenames (e.g.,
  `reviewer-adversary.md` -> `reviewer-adversary`).
- [x] 1.3 Add a guard clause: if zero reviewer agents are
  discovered, report to the user that no reviewer agents were
  found in `.opencode/agents/` and stop without delegating.

## 2. Convert Hardcoded Lists to Dynamic Delegation

- [x] 2.1 In Code Review Mode (step 2), replace the hardcoded
  five-agent invocation list with an instruction to delegate to
  all *discovered* reviewer agents in parallel.
- [x] 2.2 In Spec Review Mode (step 1), apply the same replacement
  -- delegate to all *discovered* reviewer agents in parallel.
- [x] 2.3 Add a known-roles reference table that maps agent names
  to their focus areas. Instruct the agent to use role descriptions
  from this table when delegating to known agents, and to use a
  generic review prompt for agents not in the table.

## 3. Update Verdict Policy

- [x] 3.1 Update the "Verdict" section to change "all five
  reviewers" to "all discovered reviewers." Absent reviewers
  MUST NOT affect the verdict.
- [x] 3.2 Update the "APPROVE WITH ADVISORIES" text in Spec Review
  Mode to reference discovered reviewers instead of five.

## 4. Add Discovery Summary to Final Reports

- [x] 4.1 In Code Review Mode's final report (step 6), add a
  discovery summary listing: (a) how many reviewers were
  discovered, (b) which reviewers were invoked, (c) which known
  reviewer roles were absent (informational).
- [x] 4.2 In Spec Review Mode's final report (step 6), add the
  same discovery summary.

## 5. Update Iteration Loop References

- [x] 5.1 In Code Review Mode (step 4), change "re-run all five
  reviewers" to "re-run all discovered reviewers."
- [x] 5.2 In Spec Review Mode (step 4), apply the same change.

## 6. Verification

- [x] 6.1 Verify constitution alignment: confirm the updated
  command maintains Autonomous Collaboration (artifact-based
  communication only), Composability First (auto-detects available
  reviewers), Observable Quality (discovery summary in reports),
  and Testability (testable with agent file presence/absence).
- [x] 6.2 Read the final `review-council.md` end-to-end to confirm
  no hardcoded agent count ("five") remains in the command text
  and all references use "discovered" language.
