---
name: compact
description: Compact or summarize the current conversation into a structured summary that enables a fresh agent to resume work with minimal context loss. Use ONLY when the user explicitly asks to compact, summarize the session, or create a handoff summary. Do NOT read this skill proactively based on conversation length.
---

# Conversation Compaction

Generate a structured compaction summary that preserves the essential state of the current conversation so a fresh agent can resume the task effectively.

## Core Principles

### Preserve what matters, discard what doesn't

A compaction summary is lossy. Accept that. The goal is not a transcript — it is a **launch pad** for a fresh agent to pick up where you left off with minimal ramp-up work and minimal context pollution.

### Don't write what can be discovered

A fresh CPE agent can read files, run `git diff`, check `git log`, browse directories, and read `AGENTS.md`. Never waste compaction space on information the agent can obtain in seconds. Instead, **point** to where information lives.

**Bad:** "Changed files: `auth.go`, `handler.go`, `middleware.go`, `auth_test.go` — added JWT validation, updated handler signatures, added middleware chain..."

**Good:** "Run `git diff main` to see all changes. Key subsystems modified: `internal/auth/` (JWT validation), `internal/api/` (handler signatures)."

### Structure for repeated compaction

Compaction can happen many times in a long session. The summary format must be **stable across rounds** — each compaction replaces the previous one rather than wrapping it in another layer. Use clearly delimited sections so the next compaction can identify what to keep, update, or drop.

## Process

1. **Identify the task type** — Determine which compaction profile fits (see below)
2. **Draft the summary** — Follow the template structure
3. **Trim ruthlessly** — Remove anything a fresh agent can discover or doesn't need
4. **Verify completeness** — Ensure the three mandatory sections are present and sufficient

## Summary Template

Every compaction summary MUST contain these three sections. Additional task-specific sections follow.

```markdown
# Compaction Summary (Round N)

## Objective
[One or two sentences: what is the end goal of this entire session. This rarely changes across compaction rounds.]

## Progress
[What has been accomplished so far. Be specific about outcomes, not process. Reference discoverable artifacts rather than describing them inline.]

## Remaining Work
[What still needs to happen to reach the objective. Ordered by priority or logical sequence. Include decision points if any are pending.]
```

### Optional sections (add when relevant):

```markdown
## Key Decisions & Constraints
[Important design decisions, user preferences, or constraints that a fresh agent cannot infer from code/files alone. Include the *why* behind non-obvious choices.]

## Obstacles & Workarounds
[Problems encountered and how they were resolved — or patterns to avoid. State the obstacle, the resolution, and optionally a warning. Do NOT include the debugging journey.]

## Context Rebuild Instructions
[Explicit steps for the fresh agent to rebuild working context quickly. Point to files, diffs, commands, or specific sections to read. Keep it actionable.]
```

## Task-Specific Profiles

### Codebase / Engineering Tasks

Focus on: subsystems touched, architectural decisions, test status, build state.

- Reference `git diff` and `git log` instead of listing changes
- Name the subsystems/directories modified, not individual files
- Record any non-obvious technical decisions and their rationale
- Note test results: what passes, what's broken, what's not yet written
- Mention any environment setup or tooling quirks

**Example snippet:**
```markdown
## Progress
- JWT authentication flow is fully implemented and tested in `internal/auth/`
- API handlers updated to use new auth middleware — run `git diff main` for details
- All existing tests pass (`go test ./...`); new tests added for token refresh edge cases

## Obstacles & Workarounds
- The `jwt-go` library panics on expired tokens with certain clock skew — switched to `golang-jwt/jwt/v5` which handles this gracefully. Do not use `dgrijalva/jwt-go`.
```

### Research / Writing Tasks

Focus on: thesis/argument structure, sources found, sections drafted, open questions.

- List key sources with URLs or citation keys (these cannot be rediscovered easily)
- Summarize the current outline or argument structure
- Note which sections are drafted vs. outlined vs. not started
- Record any specific findings, data points, or quotes that are hard to re-find

**Example snippet:**
```markdown
## Progress
- Literature review complete. 12 key papers identified — see `references.md` for full list
- Sections 1-3 drafted in `paper.md`; Section 4 outlined only
- Key finding: Chen et al. (2024) contradicts the standard model on protein folding rates (Table 3, p.12) — this is central to our argument

## Remaining Work
- Draft Section 4 (Discussion): synthesize findings from Chen and Park studies
- Write abstract (after all sections complete)
- User wants to limit paper to 4000 words — currently at 2800
```

### Data / Analysis Tasks

Focus on: data sources, transformations applied, results so far, methodology decisions.

- Name input/output file paths
- Summarize methodology choices and their justification
- Record intermediate results that are expensive to recompute
- Note any data quality issues discovered

### General / Mixed Tasks

Use the base template. Add optional sections as needed. When in doubt about what to include, ask: "Would a fresh agent waste more than 30 seconds figuring this out?" If yes, include it.

## What to Exclude

- **The journey** — Don't describe how you arrived at a solution. Just state the solution.
- **Failed attempts** — Unless the failure mode is a trap the next agent could fall into, omit it.
- **Verbatim code or file contents** — Point to the file instead.
- **Conversation dynamics** — "The user asked me to..." is irrelevant. State the requirement directly.
- **Obvious context** — Anything in `AGENTS.md`, README, or standard project structure.
- **Intermediate reasoning** — The next agent doesn't need your thought process, just your conclusions.

## What to Always Include

- **The objective** — Always. Even if it hasn't changed.
- **Non-obvious decisions** — Choices a fresh agent would make differently without guidance.
- **Blockers or pending questions** — Anything requiring user input that hasn't been resolved.
- **Hard-won knowledge** — Things you learned that aren't documented anywhere in the project.
- **Source URLs and citations** — For research tasks, these are expensive to re-find.
- **Current state markers** — Which tests pass, what's deployed, what branch you're on.

## Multi-Round Compaction

When compacting a conversation that already contains a prior compaction summary:

1. **Start from the previous summary**, not from scratch
2. **Merge Progress sections** — Move completed "Remaining Work" items into "Progress"
3. **Update Remaining Work** — Remove completed items, add newly discovered items
4. **Prune Obstacles** — Keep only those still relevant or dangerous; drop resolved ones that are no longer traps
5. **Preserve the Objective** — Change only if the user has explicitly changed the goal
6. **Increment the round number** — `Round N` → `Round N+1`

The summary should never grow unboundedly. If it exceeds ~800 words, aggressively compress the Progress section (summarize groups of completed work) and prune stale obstacles.

## Example: Full Compaction Summary

```markdown
# Compaction Summary (Round 2)

## Objective
Build a CLI tool that syncs local Markdown notes to a remote Notion workspace, supporting incremental updates and conflict detection.

## Progress
- CLI scaffolding complete using `cobra` in `cmd/` — supports `sync`, `status`, and `init` subcommands
- Notion API client implemented in `internal/notion/` — handles page creation, update, and content block conversion
- Markdown parser (`internal/parser/`) converts MD to Notion block format; supports headings, lists, code blocks, and images
- Incremental sync works via content hashing stored in `.notesync/state.json`
- Run `git diff main` for full changeset; primary work in `internal/` and `cmd/`

## Remaining Work
1. Implement conflict detection — compare remote `last_edited_time` with local sync timestamp
2. Add `--dry-run` flag to `sync` command
3. Write integration tests for sync workflow (unit tests all pass currently)
4. User wants a `--force` flag to overwrite remote without conflict check

## Key Decisions & Constraints
- Using content hashing (not timestamps) for change detection because filesystem timestamps are unreliable across OS
- Notion API rate limit is 3 req/sec — client has built-in retry with exponential backoff
- User preference: no external dependencies beyond cobra and the Notion SDK

## Obstacles & Workarounds
- Notion API does not support atomic batch updates — sync writes pages sequentially. If interrupted, `.notesync/state.json` may be partially updated. Added a recovery check in `sync` startup that detects and repairs partial state.

## Context Rebuild Instructions
- Read `AGENTS.md` for build/test commands
- Run `git log --oneline -10` for recent commit history
- Check `internal/notion/client.go` for API client structure
- Run `go test ./...` to verify current test state
```
