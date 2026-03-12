# Oracle Brief Template

Use this template to package a rich, self-contained context dossier for the target model. Fill sections with facts first, then labeled assumptions where facts are missing. Prefer over-inclusion of relevant context to an overly minimal brief, and assume the target model may only see this one document.

## 1. Mission

- Primary objective:
- Why this matters now:
- Desired outcome:
- What the user plans to do with the answer:

## 2. Answer Profile

- Mode: breadth / depth / breadth plus depth
- Deliverable type:
- Audience:
- Preferred structure:
- Tone or style constraints:
- Whether the Oracle should teach, decide, diagnose, synthesize, or propose:

## 3. Scope

- In scope:
- Out of scope:
- Required focus areas:
- Nice-to-have areas:
- Adjacent areas that may still help if relevant:

## 4. Context and Evidence

- Known facts:
- Relevant artifacts:
- Inline excerpts or quoted material the Oracle must see directly:
- Environment or domain background:
- Timeline or sequence of events:
- Prior work already done:
- Failed attempts or dead ends:
- Competing hypotheses or interpretations:
- Supporting context that might be useful even if not central:

## 5. Constraints and Preferences

- Time or budget constraints:
- Tooling or platform constraints:
- Citation or source expectations:
- Risk tolerance:
- Stakeholders or downstream consumers:
- What to avoid:

## 6. Open Questions and Assumptions

- Remaining unknowns:
- Ambiguities that should be preserved rather than forced closed:
- Reasonable assumptions if the Oracle must proceed:
- Areas where confidence should be qualified:

## 7. Instructions to the Oracle

- What to optimize for:
- How to handle uncertainty:
- Whether to compare alternatives, rank hypotheses, or recommend one path:
- Required sections or output schema:
- Verification or evaluation standard:
- Whether to include appendices, source maps, or further-reading branches:

## 8. Optional Context Appendix

Use this section for additional material that may be helpful even if it is not core to the main task.

- Raw notes:
- Extra examples:
- Historical context:
- Stakeholder concerns:
- Related but nonessential details:

## 9. Final Oracle Prompt

```text
You are the Oracle model. Your task is to [objective].

I am giving you an intentionally rich context package in a single self-contained document. Some details may be central and some may be peripheral; use your judgment, but do not ignore potentially relevant context. Do not assume you can open files, follow local file paths, or retrieve missing attachments beyond the text included here.

Mission:
[insert mission]

Answer profile:
[insert answer profile]

Scope:
[insert scope]

Context and evidence:
[insert structured context]

Constraints and preferences:
[insert constraints]

Open questions and assumptions:
[insert open questions]

Your instructions:
- Optimize for [correctness / completeness / practicality / etc.]
- Treat the following as confirmed facts: [facts]
- Treat the following as hypotheses or unresolved areas: [unknowns]
- Stay within scope: [scope]
- Avoid: [avoidances]
- Produce the final deliverable directly rather than only offering suggestions, unless I explicitly ask for brainstorming or critique.
- Produce the answer in this format: [format]
- If uncertainty remains, explicitly state assumptions, confidence, and plausible alternatives.
- If some context appears secondary, retain awareness of it in case it changes the answer.

Success criteria:
[what a useful answer must accomplish]
```

Adapt the final prompt freely. The template is a scaffold, not a rigid format.
