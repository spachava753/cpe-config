//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func parseFrontmatter(content string) (name, description string, err error) {
	if !strings.HasPrefix(content, "---") {
		return "", "", fmt.Errorf("no YAML frontmatter found")
	}

	re := regexp.MustCompile(`(?s)^---\n(.*?)\n---`)
	match := re.FindStringSubmatch(content)
	if match == nil {
		return "", "", fmt.Errorf("invalid frontmatter format")
	}

	frontmatter := match[1]

	// Parse name
	nameRe := regexp.MustCompile(`(?m)^name:\s*(.+)$`)
	if m := nameRe.FindStringSubmatch(frontmatter); m != nil {
		name = strings.TrimSpace(m[1])
	}

	// Parse description (may be multi-line)
	descRe := regexp.MustCompile(`(?m)^description:\s*(.+)$`)
	if m := descRe.FindStringSubmatch(frontmatter); m != nil {
		description = strings.TrimSpace(m[1])
	}

	return name, description, nil
}

func validateSkill(skillPath string) (bool, string) {
	skillMDPath := filepath.Join(skillPath, "SKILL.md")
	content, err := os.ReadFile(skillMDPath)
	if err != nil {
		return false, "SKILL.md not found"
	}

	name, description, err := parseFrontmatter(string(content))
	if err != nil {
		return false, err.Error()
	}

	// Check required fields
	if name == "" {
		return false, "Missing 'name' in frontmatter"
	}
	if description == "" {
		return false, "Missing 'description' in frontmatter"
	}

	// Validate name format
	nameRe := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !nameRe.MatchString(name) {
		return false, fmt.Sprintf("Name '%s' should be hyphen-case (lowercase letters, digits, and hyphens only)", name)
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") || strings.Contains(name, "--") {
		return false, fmt.Sprintf("Name '%s' cannot start/end with hyphen or contain consecutive hyphens", name)
	}
	if len(name) > 64 {
		return false, fmt.Sprintf("Name is too long (%d characters). Maximum is 64 characters.", len(name))
	}

	// Validate description
	if strings.Contains(description, "<") || strings.Contains(description, ">") {
		return false, "Description cannot contain angle brackets (< or >)"
	}
	if len(description) > 1024 {
		return false, fmt.Sprintf("Description is too long (%d characters). Maximum is 1024 characters.", len(description))
	}

	return true, "Skill is valid!"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run quick_validate.go <skill_directory>")
		os.Exit(1)
	}

	valid, message := validateSkill(os.Args[1])
	fmt.Println(message)
	if !valid {
		os.Exit(1)
	}
}
