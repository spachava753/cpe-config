---
name: oracle
description: Turn a vague, underspecified, or high-stakes user request into a complete Oracle handoff prompt by expanding it with abundant context, clarified scope, evidence, constraints, preferences, and desired output. Use when preparing a prompt for a slow or expensive model whose answer quality improves when it receives a rich, over-complete brief up front.
---

# Oracle

Use this skill when the user wants help preparing a prompt or brief for a highly capable model that should receive as much relevant context as possible in one shot.

## Goal

Interrogate the user to remove ambiguity, surface hidden constraints, uncover adjacent context, and build an Oracle-ready handoff that the user can paste or save directly. The default deliverable is not advice about how to ask the Oracle; it is the finished Oracle prompt or single-file dossier itself.

## Core Principle

Treat context as abundant, not scarce. A strong Oracle handoff should usually contain more context than seems strictly necessary, because the target model is good at separating signal from noise but cannot recover context it was never given.

## Workflow

### 1. Frame the Request

- Infer the task type: debugging, research, design, planning, writing, decision support, or another high-context task.
- Determine the needed answer profile: breadth, depth, or breadth plus depth.
- Restate the current understanding in a short structured form and call out assumptions early.
- Ask what the user ultimately wants to do with the answer so the Oracle can optimize for the right outcome.

### 2. Expand the Context Window

- Ask questions in short batches, but keep gathering context until the handoff feels richly specified.
- Prefer concrete artifacts over summaries whenever possible: code, logs, stack traces, URLs, datasets, notes, prior prompts, screenshots, examples, transcripts, documents, and decision history.
- Pull in surrounding context, not just the immediate question: background, prior attempts, competing hypotheses, stakeholder concerns, relevant timelines, and adjacent constraints.
- If the user is unsure, offer options, defaults, and labeled assumptions, then keep building the dossier.
- When deciding whether to include something, lean toward inclusion unless it is clearly irrelevant or distracting.

Use `references/question-bank.md` for reusable question patterns.

### 3. Resolve the Core Inputs

Before finalizing the prompt, gather or explicitly label assumptions for these fields:

- Objective: what the Oracle should accomplish.
- Scope: what is in bounds, out of bounds, and where to focus.
- Deliverable: the desired output format, structure, and level.
- Audience: who the answer is for and what background they have.
- Constraints: time, risk tolerance, tooling, budget, compliance, citations, or platform limits.
- Evidence: facts, artifacts, prior work, failed attempts, observed behavior, and competing interpretations.
- Success criteria: what would make the answer genuinely useful.
- Background context: why this request exists, what happened before now, and what downstream decisions depend on the answer.

### 4. Tailor the Interrogation to the Task

For debugging and technical diagnosis, gather generously:
- exact symptoms, expected vs actual behavior, reproduction steps, environment, versions, recent changes, logs, stack traces, suspected causes, failed fixes, architecture notes, deployment context, and operational impact.

For research, gather generously:
- topic boundaries, time horizon, geography, audience level, source preferences, controversy handling, desired synthesis style, prior knowledge, and whether the Oracle should optimize for breadth, depth, or both.

For planning, design, or decision support, gather generously:
- goals, constraints, tradeoffs, available resources, timeline, decision criteria, stakeholders, alternatives already considered, political or organizational realities, and acceptable fallback options.

For writing or explanation tasks, gather generously:
- audience, tone, length, examples to emulate or avoid, must-include points, forbidden content, source material, background assumptions, and what action the writing should cause.

### 5. Build the Oracle Dossier

Once enough context has been assembled, package it using `references/oracle-brief-template.md`.

The dossier should clearly separate:
- confirmed facts
- user preferences and constraints
- supporting context that may or may not be central but could still help
- open questions
- assumptions made to unblock progress
- the final task instructions for the Oracle

Assume the target Oracle may not be able to open local files, follow file paths, or retrieve missing attachments. Inline the necessary context directly into the handoff instead of telling the Oracle where to find it.

### 6. Produce the Handoff

Default to delivering a finished Oracle handoff, not just suggestions about how the user could write one.

When useful, return the result in this order:

1. A short list of any remaining context gaps or optional additions.
2. A single self-contained Oracle dossier or prompt that can be pasted directly into the target model or saved as one file.
3. An optional shorter variant only if the user explicitly wants a compact version.

If the user asks for speed, provide a best-effort prompt sooner, but still preserve as much context as is already available.
If the task is high stakes, prefer another clarification round and a more expansive dossier before finalizing.

## Prompt-Writing Rules

- Tell the Oracle exactly what to do, not just the topic.
- Separate facts from hypotheses and guesses.
- State what to optimize for: correctness, practicality, completeness, originality, caution, or speed.
- Specify output structure, formatting, and any required sections.
- Tell the Oracle how to handle uncertainty: note assumptions, confidence levels, alternatives, or missing evidence.
- Instruct the Oracle to produce the final deliverable directly unless the user explicitly wants brainstorming or critique.
- Assume the Oracle only sees the current prompt text. Inline key context, excerpts, and evidence instead of referring to files the Oracle may not be able to read.
- Include scope boundaries, exclusions, and decision criteria.
- Include adjacent context, prior attempts, and potentially relevant background instead of stripping the brief down to a minimal core.
- Preserve ambiguity that matters; do not flatten meaningful uncertainty into false precision.

## Stop Condition

Stop interrogating the user when:
- the objective is crisp
- the answer shape is defined
- major constraints are known
- core evidence is present or explicitly unavailable
- the handoff already contains abundant context and the next questions would be repetitive rather than additive

When in doubt, prefer one more context-rich question and one more useful artifact over an over-trimmed brief.
