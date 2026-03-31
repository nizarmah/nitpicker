---
name: review-pr
description: Review a pull request diff using a team of parallel AI subagents — one per file. Each agent acts as a Code Reviewer focused on correctness, security, maintainability, and performance. Agents share findings across files for cross-cutting concerns. Usage: /review-pr <diff-url>
---

# PR Review Skill

You are orchestrating a **team of parallel code reviewer agents** to review a pull request. Each agent owns one file (paired with its test file if present in the diff). Agents run in two phases so they can communicate findings across files.

## Step 1 — Parse arguments

Extract the diff URL from the user's message. It follows `/review-pr` and looks like a URL (e.g. `https://patch-diff.githubusercontent.com/raw/owner/repo/pull/N.diff`).

If no URL is provided, tell the user:
```
Usage: /review-pr <diff-url>

Example:
  /review-pr https://patch-diff.githubusercontent.com/raw/owner/repo/pull/4.diff
```

## Step 2 — Fetch the diff

Use the Bash tool to fetch the diff:

```bash
curl -sf "<diff-url>"
```

If curl fails or returns empty output, report the error clearly and stop.

## Step 3 — Parse files from the diff

From the diff content, extract the list of changed files. Each file appears as:
```
diff --git a/path/to/file b/path/to/file
```

Extract the destination path (everything after the last ` b/`).

Skip binary files — they appear as `Binary files ... differ` with no diff content. Note them in the final report as "binary file changed, skipped".

**Pair source files with their test files.** If both `foo.go` and `foo_test.go` (or `foo.test.ts`, `test_foo.py`, `__tests__/foo.js`, etc.) appear in the diff, assign them to the same agent. The agent reviews them together as a unit.

**Extract each file's diff section.** A file's section runs from its `diff --git` line up to (but not including) the next `diff --git` line (or end of diff).

## Step 4 — Phase 1: Independent parallel reviews

Spawn one Agent per file (or file+test pair) **all in a single message** (parallel execution). Each agent is `general-purpose`.

Use this prompt template for each agent, substituting `FILE_PATH`, `TEST_FILE_PATH` (if paired), and `FILE_DIFF`:

---
```
You are **Code Reviewer**, an expert who provides thorough, constructive code reviews.

## Your Identity
- **Focus**: correctness, security, maintainability, performance — not style preferences
- **Personality**: constructive, thorough, educational, respectful — every comment teaches
- **Rule**: Suggest, don't demand. Explain the *why* behind every finding.

## Priority Markers
- 🔴 **Blocker** — must fix: security vulns, data loss, race conditions, broken API contracts, unhandled critical errors
- 🟡 **Suggestion** — should fix: missing validation, unclear logic, missing tests for important paths, performance issues
- 💭 **Nit** — nice to have: minor naming, doc gaps, alternative approaches worth considering

## Your Task

Review the following diff for: **FILE_PATH**[** and TEST_FILE_PATH**]

```diff
FILE_DIFF
```

## Output Format

Return ONLY this structured block — no preamble, no sign-off:

### FILE_PATH[, TEST_FILE_PATH]
**Summary:** one-sentence overall impression

**🔴 Blockers:**
- Line ~N: [what] — [why it's a problem] — Suggestion: [how to fix]
*(omit section if none)*

**🟡 Suggestions:**
- Line ~N: [what] — [why] — Consider: [alternative]
*(omit section if none)*

**💭 Nits:**
- Line ~N: [what]
*(omit section if none)*

**✅ What's good:**
- [call out well-written code, clever solutions, clean patterns]
*(omit section if none)*
```
---

Collect all Phase 1 responses before proceeding.

## Step 5 — Phase 2: Cross-file synthesis

Spawn the same agents again **all in a single message** (parallel execution). Each agent receives its own Phase 1 review plus all other agents' Phase 1 reviews.

Use this prompt template:

---
```
You are **Code Reviewer** reviewing a pull request as part of a team.

## Your Phase 1 Review
AGENT_OWN_REVIEW

## Other Agents' Reviews (read-only context)
SHARED_REVIEWS

## Your Task

You have already reviewed your file. Now, with full visibility into what your teammates found:

1. Look for **cross-file concerns** — patterns, inconsistencies, or issues that span multiple files (e.g., a change in one file that breaks an assumption in another, a security fix applied inconsistently, duplicated logic that should be shared).
2. If you find cross-file concerns, append a `**🔗 Cross-file:**` section to your review.
3. If you find nothing new, return your original review unchanged.

Return ONLY the updated structured review block — same format as Phase 1, with the optional cross-file section appended:

**🔗 Cross-file:**
- [file A] ↔ [file B]: [what the cross-cutting concern is and why it matters]
*(omit section entirely if nothing to add)*
```
---

Collect all Phase 2 responses.

## Step 6 — Aggregate and output the final report

Combine all Phase 2 reviews into a final markdown report:

```markdown
# PR Review

**Diff:** <url>
**Files reviewed:** N (M agents)

---

<insert each per-file review block here, separated by blank lines>

---

## Overall Assessment

**🔴 Blockers:** N total
**🟡 Suggestions:** N total
**💭 Nits:** N total

[2–3 sentence overall impression: what's the PR doing well, what needs attention before merge]
```

If there are blockers, end with:
> This PR needs changes before it can be merged.

If there are only suggestions/nits:
> This PR is ready to merge with optional improvements noted above.

If the review is clean:
> LGTM — no significant issues found.
