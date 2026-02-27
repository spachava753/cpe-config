# AGENTS.md

This repository is a [CPE](https://github.com/spachava753/cpe) agent configuration + skills workspace. It stores local model/tool config, system prompts, and reusable agent skills.

## Component Index

- `cpe.yaml` - Main CPE config (models, defaults, pricing metadata, provider settings).
- `subagent.yaml` - Subagent runtime config
- `*_instructions.md` - Golang template instructions for CPE and subagents formatted as markdown.
- `skills/` - Skill library; each skill lives in `skills/<name>/SKILL.md`.
- `auth.json` - Local auth/config state (treat as sensitive).

## Upstream CPE Reference

For product behavior and CLI usage, refer to the upstream README:
- Raw README: `https://raw.githubusercontent.com/spachava753/cpe/main/README.md`
- Config schema for `cpe.yaml` and `subagent.yaml`: `https://raw.githubusercontent.com/spachava753/cpe/refs/heads/main/schema/cpe-config-schema.json`
