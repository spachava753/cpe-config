---
name: subagent
description: Instructions for invoking CPE as a sub-agent via CLI for parallel and decomposed task execution. Use when tasks can be parallelized (research from multiple sources, processing multiple files, fan-out operations), decomposed into independent subtasks, when offloading work to separate agent processes would improve efficiency, when you want to preserve context window space, or when you only need final results without intermediate steps cluttering the conversation.
---

# Subagent Invocation

Invoke CPE as a sub-agent using:

```bash
cpe --skip-stdin -m model-ref-id -n -G -i resource "Prompt as single argument"
```

## Flags

| Flag | Purpose |
|------|---------|
| `--skip-stdin` | Skip reading from stdin |
| `-m` | Model: `opus` (smart/slow/expensive), `sonnet` (balanced), `haiku` (fast/cheap/follows instructions well) |
| `-n` | Start new conversation |
| `-G` | Don't save conversation |
| `-i` | Input resource (file path or URL). Supports text, images, PDFs. Prefer raw content URLs over HTML pages |

**Stderr handling:** Use `cmd.Output()` to capture stdout. On failure, extract stderr from `*exec.ExitError` for debugging.

## When to Use Subagents

**Parallel execution:** Query multiple sources or process multiple files simultaneously.

**Context window preservation:** Offload work to subagents to avoid filling up your context with intermediate steps, tool calls, and verbose outputs.

**Results-only tasks:** When you only need the final output, not the reasoning or intermediate steps that led to it.

**Decomposed workflows:** Break complex tasks into independent subtasks (analyze N items → synthesize results).

**Batch processing:** Process multiple files/URLs with the same operation.

**Model selection:**
- `haiku`: Repetitive tasks, simple transformations, following explicit instructions
- `sonnet`: Analysis, summarization, moderate complexity
- `opus`: Complex reasoning, nuanced decisions, high-stakes outputs

## Code Mode Examples

### Parallel File Processing

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
)

func Run(ctx context.Context) error {
	files := []string{"auth.go", "api.go", "db.go"}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[string]string)
	
	for _, f := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := exec.CommandContext(ctx, "cpe",
				"--skip-stdin", "-m", "haiku", "-n", "-G",
				"-i", f,
				"Add comprehensive error handling. Output only the modified code.")
			out, err := cmd.Output()
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					results[f] = fmt.Sprintf("error: %s", string(exitErr.Stderr))
				} else {
					results[f] = fmt.Sprintf("error: %v", err)
				}
				return
			}
			results[f] = string(out)
		}()
	}
	wg.Wait()
	
	for f, r := range results {
		fmt.Printf("=== %s ===\n%s\n", f, r)
	}
	return nil
}
```

### Parallel Research from URLs

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context) error {
	sources := []struct {
		name string
		url  string
	}{
		{"Go", "https://raw.githubusercontent.com/golang/go/master/README.md"},
		{"Rust", "https://raw.githubusercontent.com/rust-lang/rust/master/README.md"},
		{"Zig", "https://raw.githubusercontent.com/ziglang/zig/master/README.md"},
	}

	g, ctx := errgroup.WithContext(ctx)
	results := make([]string, len(sources))

	for i, s := range sources {
		g.Go(func() error {
			cmd := exec.CommandContext(ctx, "cpe",
				"--skip-stdin", "-m", "sonnet", "-n", "-G",
				"-i", s.url,
				fmt.Sprintf("Summarize key features of %s from this README in 3 bullets", s.name))
			out, err := cmd.Output()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					return fmt.Errorf("%s: %s", s.name, string(exitErr.Stderr))
				}
				return fmt.Errorf("%s: %w", s.name, err)
			}
			results[i] = fmt.Sprintf("## %s\n%s", s.name, out)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	for _, r := range results {
		fmt.Println(r)
	}
	return nil
}
```

### Analyze Then Synthesize (Two-Phase)

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context) error {
	// Phase 1: Parallel analysis
	modules := []string{"src/auth", "src/api", "src/storage"}
	g, ctx := errgroup.WithContext(ctx)
	analyses := make([]string, len(modules))

	for i, mod := range modules {
		g.Go(func() error {
			cmd := exec.CommandContext(ctx, "cpe",
				"--skip-stdin", "-m", "haiku", "-n", "-G",
				"-i", mod,
				"Analyze architecture and list potential issues. Output as markdown.")
			out, err := cmd.Output()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					return fmt.Errorf("%s: %s", mod, string(exitErr.Stderr))
				}
				return fmt.Errorf("%s: %w", mod, err)
			}
			analyses[i] = fmt.Sprintf("## %s\n%s", filepath.Base(mod), out)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	// Phase 2: Synthesize
	tmp, _ := os.CreateTemp("", "analysis-*.md")
	tmp.WriteString(strings.Join(analyses, "\n\n"))
	tmp.Close()
	defer os.Remove(tmp.Name())

	cmd := exec.CommandContext(ctx, "cpe",
		"--skip-stdin", "-m", "sonnet", "-n", "-G",
		"-i", tmp.Name(),
		"Synthesize these analyses into an architecture review with prioritized recommendations.")
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("synthesis: %s", string(exitErr.Stderr))
		}
		return err
	}
	fmt.Println(string(out))
	return nil
}
```

### Batch Test Generation

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context) error {
	entries, _ := filepath.Glob("pkg/*/*.go")
	var targets []string
	for _, e := range entries {
		if !strings.HasSuffix(e, "_test.go") {
			testFile := strings.TrimSuffix(e, ".go") + "_test.go"
			if _, err := os.Stat(testFile); os.IsNotExist(err) {
				targets = append(targets, e)
			}
		}
	}

	g, _ := errgroup.WithContext(ctx)
	for _, t := range targets {
		g.Go(func() error {
			testPath := strings.TrimSuffix(t, ".go") + "_test.go"
			cmd := exec.CommandContext(ctx, "cpe",
				"--skip-stdin", "-m", "haiku", "-n", "-G",
				"-i", t,
				"Generate unit tests for this file. Output only valid Go test code.")
			out, err := cmd.Output()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					return fmt.Errorf("%s: %s", t, string(exitErr.Stderr))
				}
				return fmt.Errorf("%s: %w", t, err)
			}
			return os.WriteFile(testPath, out, 0644)
		})
	}
	return g.Wait()
}
```

### Image Analysis Fan-Out

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context) error {
	images, _ := filepath.Glob("screenshots/*.png")
	g, ctx := errgroup.WithContext(ctx)
	results := make([]string, len(images))

	for i, img := range images {
		g.Go(func() error {
			cmd := exec.CommandContext(ctx, "cpe",
				"--skip-stdin", "-m", "haiku", "-n", "-G",
				"-i", img,
				"Describe this UI screenshot: layout, key elements, any UX issues.")
			out, err := cmd.Output()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					return fmt.Errorf("%s: %s", img, string(exitErr.Stderr))
				}
				return fmt.Errorf("%s: %w", img, err)
			}
			results[i] = fmt.Sprintf("### %s\n%s", filepath.Base(img), out)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	for _, r := range results {
		fmt.Println(r)
	}
	return nil
}
```
