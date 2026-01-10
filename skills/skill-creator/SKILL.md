---
name: skill-creator
description: Guide for creating effective skills using Go scripts. Use when users want to create a new skill (or update an existing skill) that extends AI agent capabilities with specialized knowledge, workflows, or tool integrations. This variant uses Go for all scripts.
---

# Go Skill Creator

This skill provides guidance for creating effective skills, with Go-based tooling.

## About Skills

Skills are modular, self-contained packages that extend an AI agent's capabilities by providing specialized knowledge, workflows, and tools. They transform a general-purpose agent into a specialized agent equipped with procedural knowledge.

### What Skills Provide

1. **Specialized workflows** - Multi-step procedures for specific domains
2. **Tool integrations** - Instructions for working with specific file formats or APIs
3. **Domain expertise** - Company-specific knowledge, schemas, business logic
4. **Bundled resources** - Scripts, references, and assets for complex tasks

## Core Principles

### Concise is Key

The context window is a shared resource. Only add context the agent doesn't already have. Challenge each piece of information: "Is this really needed?"

Prefer concise examples over verbose explanations.

### Set Appropriate Degrees of Freedom

- **High freedom (text instructions)**: Multiple valid approaches, context-dependent decisions
- **Medium freedom (pseudocode/parameterized scripts)**: Preferred pattern exists with acceptable variation
- **Low freedom (specific scripts)**: Fragile operations, consistency critical

### Anatomy of a Skill

```
skill-name/
├── SKILL.md (required)
│   ├── YAML frontmatter (name, description - required)
│   └── Markdown instructions
└── Bundled Resources (optional)
    ├── scripts/     - Executable code (Go/Bash/etc.)
    ├── references/  - Documentation loaded into context as needed
    └── assets/      - Files used in output (templates, icons, fonts)
```

#### SKILL.md Frontmatter

- **name**: Skill identifier (hyphen-case, max 64 chars)
- **description**: What it does AND when to use it (max 1024 chars). This is the primary triggering mechanism.

#### Bundled Resources

- **scripts/**: Deterministic, reusable code. May be executed without loading into context.
- **references/**: Documentation read while working. Keep SKILL.md lean, move details here.
- **assets/**: Templates, images, fonts used in output—not loaded into context.

### Progressive Disclosure

1. **Metadata** - Always in context (~100 words)
2. **SKILL.md body** - When skill triggers (<5k words)
3. **Bundled resources** - As needed (unlimited)

Keep SKILL.md under 500 lines. Split into references/ when approaching this limit.

## Skill Creation Process

1. Understand the skill with concrete examples
2. Plan reusable contents (scripts, references, assets)
3. Initialize the skill (run init_skill.go)
4. Edit the skill (implement resources, write SKILL.md)
5. Package the skill (run package_skill.go)
6. Iterate based on real usage

### Step 1: Understanding with Concrete Examples

Ask users for concrete usage examples:
- "What functionality should this skill support?"
- "Can you give examples of how it would be used?"
- "What should trigger this skill?"

### Step 2: Planning Contents

For each example, identify what scripts, references, and assets would help when executing repeatedly.

### Step 3: Initializing the Skill

Run the init script:

```bash
go run scripts/init_skill.go <skill-name> --path <output-directory>
```

This creates a template skill directory with SKILL.md and example resource directories.

### Step 4: Edit the Skill

1. Start with reusable resources (scripts/, references/, assets/)
2. Test any added scripts by running them
3. Delete unneeded example files
4. Update SKILL.md with clear instructions

For design patterns, see:
- **references/workflows.md** - Sequential workflows and conditional logic
- **references/output-patterns.md** - Template and example patterns

#### SKILL.md Writing Guidelines

- Use imperative/infinitive form
- Include WHEN to use in description (not body)
- Keep body focused on HOW to execute

### Step 5: Packaging

```bash
go run scripts/package_skill.go <path/to/skill-folder> [output-directory]
```

This validates the skill and creates a distributable .skill file.

### Step 6: Iterate

Test the skill on real tasks, identify improvements, update, and test again.
