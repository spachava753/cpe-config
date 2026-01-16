---
name: init-agents-md
description: Create an AGENTS.md file that serves as a minimal index/README for coding agents new to the repository. Use when a repository needs an AGENTS.md file to help coding agents quickly understand the project structure, tech stack, and how to work with the codebase.
---

# AGENTS.md Creator

This skill guides the creation of an AGENTS.md file—a minimal index for coding agents entering a repository.

## Purpose

AGENTS.md is a concise navigation guide for coding agents. It should:
- Index where components and systems are located
- Reference documentation files the agent can retrieve independently
- Provide essential technical context (tech stack, build/test/lint commands)
- Remain minimal—defer details to referenced documents

## When to Create AGENTS.md

- Repository lacks an AGENTS.md file
- Existing AGENTS.md is outdated or insufficient
- New agent needs orientation before working on the codebase

## How to Create AGENTS.md

### 1. Gather Context

Analyze the repository structure to understand:
- Programming language(s) and frameworks
- Build system (Makefile, package.json, go.mod, Cargo.toml, etc.)
- Documentation files (README.md, docs/, architecture.md, etc.)
- Directory organization and where key components live

### 2. Document Essential Information

Include only what a coding agent needs immediately:

**Tech Stack**
- Language(s), major frameworks/libraries
- Build tools and package managers

**Essential Commands**
- How to build/compile the project
- How to run tests
- How to lint/format code
- Any other critical development commands

**Documentation References**
- Link to architecture.md (if exists)
- Link to any design docs or technical specifications
- Note where project structure is documented

**Component Index**
- Where core systems are located (e.g., "Authentication: src/auth/")
- Where configuration lives
- Where entry points are (main.go, index.ts, etc.)

### 3. Keep It Minimal

Follow these principles:
- **Index, don't duplicate**: Reference existing docs rather than repeating them
- **What, not how**: State what something is, not how it works (that's in the referenced doc)
- **Agent-centric**: Focus on what helps an agent navigate and work effectively
- **One screen**: Aim for ~50-100 lines max

### 4. Use Your Judgment

The exact content depends on the repository. Apply discretion:
- Small/simple repo: AGENTS.md may be very brief
- Large/complex repo: Include more navigation help
- Well-documented repo: Emphasize links to docs
- Poorly-documented repo: Include more essential context

### 5. Format

Use clear sections with headings. Example structure:

```
# AGENTS.md

## Tech Stack
...

## Essential Commands
...

## Documentation
- [Architecture](docs/architecture.md)
- [API Reference](docs/api.md)

## Component Index
- Auth system: src/auth/
- Database: src/db/
```