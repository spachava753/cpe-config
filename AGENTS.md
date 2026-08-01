# AGENTS.md

This repository is a [CPE](https://github.com/spachava753/cpe) agent configuration + skills workspace. It stores local model/tool config, system prompts, and reusable agent skills.

## Component Index

- `cpe.yaml` - Main CPE config (complete per-model runtime profiles, pricing metadata, provider settings).
- `*_instructions.md` - Golang template instructions for CPE and subagents formatted as markdown.
- `scripts/sprite-dev` - Creates or reconciles Sprites with SSH and CPE, cloning this repository for config.
- `skills/` - Skill library; each skill lives in `skills/<name>/SKILL.md`.
- `auth.json` - Local auth/config state (treat as sensitive).

## Known Config Nuances

- Keep `systemPromptPath` values relative to `cpe.yaml`; CPE resolves them from the config file's directory, which keeps this repository portable between macOS and Linux/Sprites.
- `scripts/sprite-dev` clones `cpe-config` directly into `~/.config/cpe` on each Sprite; only local secrets and OAuth state are uploaded. It installs the latest CPE GitHub release artifact by default. `--cpe-version` or `CPE_INSTALL_VERSION` opts into a source build for a branch, tag, or commit.
- For Cloudflare AI Gateway custom providers, configure the provider's API root, not the final SDK request URL. Include a provider-specific compatibility prefix if that is part of the upstream root, but do not include version or operation suffixes that the CPE client/SDK appends. For `type: anthropic`, CPE appends `/v1/messages`; for Moonshot this means the custom provider base URL is `https://api.moonshot.ai/anthropic`, not `https://api.moonshot.ai/anthropic/v1`.
- Moonshot/Kimi is not a first-class Cloudflare AI Gateway provider in the provider-native list; use the custom provider slug `moonshot-anthropic` for Anthropic-compatible Kimi routes.
- `kimi-k2.7-code` has thinking always on and should not expose selectable `thinkingValues`, but the Anthropic-compatible endpoint rejects omitted thinking in live tests. Keep the CPE request patch that adds `thinking: {type: enabled}` unless the profile is switched to OpenAI chat completions.

## Upstream CPE Reference

For product behavior and CLI usage, refer to the upstream README:
- Raw README: `https://raw.githubusercontent.com/spachava753/cpe/main/README.md`
- Config schema for `cpe.yaml` and `subagent.yaml`: `https://raw.githubusercontent.com/spachava753/cpe/refs/heads/main/schema/cpe-config-schema.json`
