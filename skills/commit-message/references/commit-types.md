# Commit Types

| Type | Description |
|------|-------------|
| feat | New feature for the user |
| fix | Bug fix for the user |
| docs | Documentation only changes |
| style | Formatting, missing semicolons, etc (no code change) |
| refactor | Code change that neither fixes a bug nor adds a feature |
| perf | Code change that improves performance |
| test | Adding or correcting tests |
| build | Changes to build system or external dependencies |
| ci | Changes to CI configuration files and scripts |
| chore | Other changes that don't modify src or test files |
| revert | Reverts a previous commit |

## Scope

The scope is optional and indicates the section of the codebase affected:
- Component name (e.g., `auth`, `api`, `ui`)
- Module name (e.g., `parser`, `database`)
- File or directory (e.g., `readme`, `config`)

## Breaking Changes

Add `!` after the scope for breaking changes:
```
feat(api)!: remove deprecated endpoints
```

Always include a `BREAKING CHANGE:` footer explaining the impact.
