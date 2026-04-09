// Package skillinstall installs docs-ssot SKILL.md files for AI coding agents.
package skillinstall

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
)

// skillPath returns the install path for the SKILL.md file for the given tool.
func skillPath(tool agentscan.Tool) string {
	switch tool {
	case agentscan.ToolClaude:
		return filepath.Join(".claude", "skills", "docs-ssot", "SKILL.md")
	case agentscan.ToolCursor:
		return filepath.Join(".cursor", "skills", "docs-ssot", "SKILL.md")
	case agentscan.ToolCopilot:
		return filepath.Join(".github", "skills", "docs-ssot", "SKILL.md")
	case agentscan.ToolCodex:
		return filepath.Join(".agents", "skills", "docs-ssot", "SKILL.md")
	default:
		return ""
	}
}

// skillFrontmatter returns the YAML frontmatter block for the given tool.
func skillFrontmatter(tool agentscan.Tool) string {
	switch tool {
	case agentscan.ToolClaude:
		return `---
name: docs-ssot
description: Set up docs-ssot SSOT documentation structure — migrate existing docs, build, and validate
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash(docs-ssot *)
  - Bash(make docs*)
  - Bash(git diff *)
---`
	case agentscan.ToolCursor, agentscan.ToolCopilot, agentscan.ToolCodex:
		return `---
name: docs-ssot
description: Set up docs-ssot SSOT documentation structure — migrate existing docs, build, and validate
---`
	default:
		return ""
	}
}

// skillBody returns the shared skill body (same for all tools).
const skillBody = `
# docs-ssot: Documentation SSOT Setup

This skill migrates existing Markdown documentation into a modular Single Source of Truth (SSOT) structure managed by ` + "`docs-ssot`" + `, then builds and validates the output.

## When to use this skill

- You have existing documentation (README.md, CLAUDE.md, AGENTS.md, etc.) to manage as SSOT
- You want to set up ` + "`docs-ssot`" + ` in a new or existing repository
- You want to regenerate documentation after editing source templates

---

## Workflow

### Step 1 — Check prerequisites

Verify ` + "`docs-ssot`" + ` is installed:

` + "```sh" + `
docs-ssot version
` + "```" + `

If not installed:

` + "```sh" + `
# Homebrew (macOS/Linux)
brew tap hiromaily/tap && brew install docs-ssot

# or Go install
go install github.com/hiromaily/docs-ssot/cmd/docs-ssot@latest
` + "```" + `

---

### Step 2 — Identify existing documentation files

List Markdown files in the repository root:

` + "```sh" + `
ls *.md
` + "```" + `

Common candidates: ` + "`README.md`" + `, ` + "`CLAUDE.md`" + `, ` + "`AGENTS.md`" + `, ` + "`CONTRIBUTING.md`" + `

---

### Step 3 — Preview the migration plan

Run a dry-run to see what sections will be created without writing any files:

` + "```sh" + `
docs-ssot migrate --dry-run README.md CLAUDE.md
` + "```" + `

Review the output to understand:
- How many sections will be created and their categories
- Which sections are detected as duplicates (shared across files)
- The proposed template structure under ` + "`template/sections/`" + ` and ` + "`template/pages/`" + `

---

### Step 4 — Run the migration

Migrate the identified files:

` + "```sh" + `
docs-ssot migrate README.md CLAUDE.md AGENTS.md
` + "```" + `

This will:
1. Split each file by H2 headings into section files under ` + "`template/sections/<category>/`" + `
2. Create template files under ` + "`template/pages/`" + ` with ` + "`@include`" + ` directives
3. Create or update ` + "`docsgen.yaml`" + ` with build targets
4. Verify round-trip: build and compare output against originals

---

### Step 5 — Review the generated structure

` + "```sh" + `
docs-ssot index
` + "```" + `

Inspect:
- ` + "`template/sections/`" + ` — modular section files (**edit these, not the generated outputs**)
- ` + "`template/pages/*.tpl.md`" + ` — template files defining document structure
- ` + "`docsgen.yaml`" + ` — build targets mapping templates to output files

---

### Step 6 — Build documentation

Regenerate all output files from the templates:

` + "```sh" + `
docs-ssot build
# or, if a Makefile target exists:
make docs
` + "```" + `

---

### Step 7 — Validate

Check that all include directives resolve correctly:

` + "```sh" + `
docs-ssot validate
` + "```" + `

---

### Step 8 — Check for SSOT violations

Scan for near-duplicate sections that should be merged into a single source:

` + "```sh" + `
docs-ssot check
` + "```" + `

If duplicates are found, consolidate the content into one section file under ` + "`template/sections/`" + ` and update the templates to reference it with ` + "`@include`" + `.

---

## Key Commands Reference

| Command | Purpose |
|---------|---------|
| ` + "`docs-ssot migrate <files>`" + ` | Decompose existing docs into SSOT section structure |
| ` + "`docs-ssot migrate --dry-run <files>`" + ` | Preview migration without writing files |
| ` + "`docs-ssot build`" + ` | Generate all output files from templates |
| ` + "`docs-ssot validate`" + ` | Check all include directives resolve |
| ` + "`docs-ssot check`" + ` | Detect near-duplicate sections (SSOT violations) |
| ` + "`docs-ssot index`" + ` | Show include relationships and orphan detection |
| ` + "`docs-ssot include <template>`" + ` | Expand and print a template to stdout (debugging) |

---

## Editing documentation after migration

After migration, the ongoing workflow is:

1. Edit source files in ` + "`template/sections/`" + `
2. Run ` + "`docs-ssot build`" + ` to regenerate outputs
3. Verify with ` + "`git diff README.md`" + `

**Never edit** ` + "`README.md`" + `, ` + "`CLAUDE.md`" + `, or ` + "`AGENTS.md`" + ` directly — they are overwritten on every build.

---

## docsgen.yaml structure

` + "```yaml" + `
targets:
  - input: template/pages/README.tpl.md
    output: README.md
  - input: template/pages/CLAUDE.tpl.md
    output: CLAUDE.md
` + "```" + `

Add or remove targets to control which output files are generated.

---

## Include directive syntax

Templates use include directives to compose sections:

` + "```markdown" + `
<!-- @include: sections/project/overview.md -->
<!-- @include: sections/development/ -->
<!-- @include: sections/**/*.md level=+1 -->
` + "```" + `

The optional ` + "`level=+N`" + ` parameter adjusts heading depth of the included content.
`

// skillContent returns the complete SKILL.md content for the given tool.
func skillContent(tool agentscan.Tool) string {
	return skillFrontmatter(tool) + "\n" + skillBody
}

// Install installs the docs-ssot skill for each of the specified tools.
// It writes to the current working directory.
// If a skill file already exists, it prompts the user via in/out before overwriting.
func Install(tools []agentscan.Tool, in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)

	for _, tool := range tools {
		path := skillPath(tool)
		if path == "" {
			continue
		}

		content := skillContent(tool)

		if _, err := os.Stat(path); err == nil {
			// File exists — prompt user.
			if _, err := fmt.Fprintf(out, "%s already exists. Overwrite? [y/N]: ", path); err != nil {
				return fmt.Errorf("write prompt: %w", err)
			}
			line, err := reader.ReadString('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("read prompt: %w", err)
			}
			answer := strings.TrimSpace(strings.ToLower(line))
			if answer != "y" && answer != "yes" {
				if _, err := fmt.Fprintf(out, "Skipped %s\n", path); err != nil {
					return fmt.Errorf("write skip message: %w", err)
				}
				continue
			}
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			return fmt.Errorf("create directory for %s: %w", path, err)
		}

		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}

		if _, err := fmt.Fprintf(out, "Installed %s\n", path); err != nil {
			return fmt.Errorf("write install message: %w", err)
		}
	}

	return nil
}
