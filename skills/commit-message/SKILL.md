---
name: commit-message
description: Generate conventional commit messages. Use when asked to create a commit message, write a commit, or summarize changes for version control. Follows conventional commits format with type(scope)!: summary structure.
---

# Commit Message Generator

Generate commit messages following conventional commits format.

## Format

```text
type(scope)!: short summary

Body paragraphs using full sentences. Use imperative, present tense. Use GitHub-flavored markdown. No trailing period or whitespace. Add line breaks only between paragraphs.

BREAKING CHANGE: footer describing breaking change (if applicable)
Closes #123
```

## Process

1. **View the diff** - Run `git diff --staged` or `git diff` to see changes
2. **Analyze changes** - Understand what changed and why (skip extremely large file diffs)
3. **Determine type** - See references/commit-types.md for valid types
4. **Identify scope** - Component, module, or area affected (optional)
5. **Write summary** - Concise imperative description (<50 chars ideal)
6. **Write body** - Explain what and why, not how

## Critical: Base Message on Final Diff Only

**Always ground the commit message on the actual diff, not conversation history.**

During a session, there may be multiple iterations: initial implementation, refactors, bug fixes, style changes. These intermediate steps are transient and irrelevant to the commit message.

The commit describes the **final state** being committed, not the journey to get there. If code was added then refactored before committing, only describe the final refactored version. The git history doesn't know about your conversation—it only sees the diff.

Before writing the message, always run `git diff` or `git diff --staged` to see exactly what will be committed.

## Rules

- **Describe what and why, not how** - Implementation details are in the diff
- **Use imperative mood** - "Add feature" not "Added feature"
- **No bullet points** - Use prose paragraphs in the body
- **Don't assume** - Ask for clarification or explore the codebase if needed
- **Generate only** - Do not commit without user approval
- **Ignore conversation history** - Only the final diff matters

## Examples

### Simple feature
```text
feat(auth): add OAuth2 login support

Users can now authenticate using Google and GitHub OAuth providers. This removes the friction of password-based registration for new users.
```

### Bug fix with issue reference
```text
fix(api): handle null response from external service

The payment gateway occasionally returns null instead of an error object. This caused unhandled exceptions in production.

Closes #456
```

### Breaking change
```text
refactor(config)!: change config file format from YAML to TOML

TOML provides better type safety and clearer syntax for nested configuration. Existing YAML configs must be migrated using the provided migration script.

BREAKING CHANGE: Configuration files must be converted from YAML to TOML format. Run `migrate-config` to convert automatically.
```

### Documentation
```text
docs(readme): add installation instructions for Windows

The README now includes PowerShell commands for Windows users and notes about WSL compatibility.
```
