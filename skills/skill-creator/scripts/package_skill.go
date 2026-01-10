//go:build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io"
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

	nameRe := regexp.MustCompile(`(?m)^name:\s*(.+)$`)
	if m := nameRe.FindStringSubmatch(frontmatter); m != nil {
		name = strings.TrimSpace(m[1])
	}

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

	if name == "" {
		return false, "Missing 'name' in frontmatter"
	}
	if description == "" {
		return false, "Missing 'description' in frontmatter"
	}

	nameRe := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !nameRe.MatchString(name) {
		return false, fmt.Sprintf("Name '%s' should be hyphen-case", name)
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") || strings.Contains(name, "--") {
		return false, fmt.Sprintf("Name '%s' has invalid hyphens", name)
	}
	if len(name) > 64 {
		return false, "Name too long (max 64 chars)"
	}
	if strings.Contains(description, "<") || strings.Contains(description, ">") {
		return false, "Description cannot contain angle brackets"
	}
	if len(description) > 1024 {
		return false, "Description too long (max 1024 chars)"
	}

	return true, "Skill is valid!"
}

func packageSkill(skillPath, outputDir string) error {
	fmt.Println("🔍 Validating skill...")
	valid, msg := validateSkill(skillPath)
	if !valid {
		return fmt.Errorf("validation failed: %s", msg)
	}
	fmt.Printf("✅ %s\n\n", msg)

	skillName := filepath.Base(skillPath)
	if outputDir == "" {
		outputDir = "."
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	outputFile := filepath.Join(outputDir, skillName+".skill")
	zipFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("creating skill file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	baseDir := filepath.Dir(skillPath)
	err = filepath.Walk(skillPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		fmt.Printf("  Added: %s\n", relPath)
		return nil
	})

	if err != nil {
		return fmt.Errorf("adding files to archive: %w", err)
	}

	fmt.Printf("\n✅ Successfully packaged skill to: %s\n", outputFile)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run package_skill.go <path/to/skill-folder> [output-directory]")
		fmt.Println("\nExample:")
		fmt.Println("  go run package_skill.go ./my-skill")
		fmt.Println("  go run package_skill.go ./my-skill ./dist")
		os.Exit(1)
	}

	skillPath := os.Args[1]
	outputDir := ""
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	fmt.Printf("📦 Packaging skill: %s\n", skillPath)
	if outputDir != "" {
		fmt.Printf("   Output directory: %s\n", outputDir)
	}
	fmt.Println()

	absPath, err := filepath.Abs(skillPath)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	if err := packageSkill(absPath, outputDir); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}
