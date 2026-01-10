//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const skillTemplate = `---
name: %s
description: [TODO: Complete explanation of what the skill does and when to use it. Include WHEN to use - specific scenarios, file types, or tasks that trigger it.]
---

# %s

## Overview

[TODO: 1-2 sentences explaining what this skill enables]

## [TODO: Main Section]

[TODO: Add content. See examples in existing skills:
- Code samples for technical skills
- Decision trees for complex workflows
- References to scripts/templates/references as needed]

## Resources

This skill includes example resource directories:

### scripts/
Executable Go code that performs specific operations.

### references/
Documentation loaded into context as needed.

### assets/
Files used in output (templates, images, fonts).

Any unneeded directories can be deleted.
`

const exampleScript = `//go:build ignore

package main

import "fmt"

// Example helper script for %s
// Replace with actual implementation or delete if not needed.

func main() {
	fmt.Println("This is an example script for %s")
	// TODO: Add actual script logic here
}
`

const exampleReference = `# Reference Documentation for %s

This is a placeholder for detailed reference documentation.
Replace with actual content or delete if not needed.

## When Reference Docs Are Useful

- Comprehensive API documentation
- Detailed workflow guides
- Complex multi-step processes
- Information too lengthy for main SKILL.md
`

const exampleAsset = `# Example Asset File

This placeholder represents where asset files would be stored.
Replace with actual files (templates, images, fonts) or delete if not needed.

Asset files are used in output, not loaded into context.
`

func titleCase(s string) string {
	words := strings.Split(s, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func main() {
	if len(os.Args) < 4 || os.Args[2] != "--path" {
		fmt.Println("Usage: go run init_skill.go <skill-name> --path <path>")
		fmt.Println("\nSkill name requirements:")
		fmt.Println("  - Hyphen-case identifier (e.g., 'data-analyzer')")
		fmt.Println("  - Lowercase letters, digits, and hyphens only")
		fmt.Println("  - Max 64 characters")
		fmt.Println("\nExamples:")
		fmt.Println("  go run init_skill.go my-new-skill --path ./skills")
		os.Exit(1)
	}

	skillName := os.Args[1]
	basePath := os.Args[3]
	skillDir := filepath.Join(basePath, skillName)
	skillTitle := titleCase(skillName)

	fmt.Printf("🚀 Initializing skill: %s\n", skillName)
	fmt.Printf("   Location: %s\n\n", basePath)

	// Check if exists
	if _, err := os.Stat(skillDir); err == nil {
		fmt.Printf("❌ Error: Skill directory already exists: %s\n", skillDir)
		os.Exit(1)
	}

	// Create directories
	dirs := []string{
		skillDir,
		filepath.Join(skillDir, "scripts"),
		filepath.Join(skillDir, "references"),
		filepath.Join(skillDir, "assets"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			fmt.Printf("❌ Error creating directory: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("✅ Created skill directory: %s\n", skillDir)

	// Write SKILL.md
	content := fmt.Sprintf(skillTemplate, skillName, skillTitle)
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
		fmt.Printf("❌ Error creating SKILL.md: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Created SKILL.md")

	// Write example script
	scriptContent := fmt.Sprintf(exampleScript, skillName, skillName)
	if err := os.WriteFile(filepath.Join(skillDir, "scripts", "example.go"), []byte(scriptContent), 0644); err != nil {
		fmt.Printf("❌ Error creating example script: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Created scripts/example.go")

	// Write example reference
	refContent := fmt.Sprintf(exampleReference, skillTitle)
	if err := os.WriteFile(filepath.Join(skillDir, "references", "api_reference.md"), []byte(refContent), 0644); err != nil {
		fmt.Printf("❌ Error creating example reference: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Created references/api_reference.md")

	// Write example asset
	if err := os.WriteFile(filepath.Join(skillDir, "assets", "example_asset.txt"), []byte(exampleAsset), 0644); err != nil {
		fmt.Printf("❌ Error creating example asset: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Created assets/example_asset.txt")

	fmt.Printf("\n✅ Skill '%s' initialized successfully at %s\n", skillName, skillDir)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Edit SKILL.md to complete the TODO items")
	fmt.Println("2. Customize or delete example files in scripts/, references/, assets/")
	fmt.Println("3. Run the validator when ready")
}
