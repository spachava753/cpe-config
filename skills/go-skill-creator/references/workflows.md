# Workflow Patterns

## Sequential Workflows

For complex tasks, break operations into clear, sequential steps:

```markdown
Filling a PDF form involves these steps:

1. Analyze the form (run analyze_form.go)
2. Create field mapping (edit fields.json)
3. Validate mapping (run validate_fields.go)
4. Fill the form (run fill_form.go)
5. Verify output (run verify_output.go)
```

## Conditional Workflows

For tasks with branching logic, guide through decision points:

```markdown
1. Determine the modification type:
   **Creating new content?** → Follow "Creation workflow" below
   **Editing existing content?** → Follow "Editing workflow" below

2. Creation workflow: [steps]
3. Editing workflow: [steps]
```
