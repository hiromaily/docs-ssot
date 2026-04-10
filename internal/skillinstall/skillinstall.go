// Package skillinstall installs docs-ssot SKILL.md files for AI coding agents.
package skillinstall

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
)

//go:embed skills/SKILL.md
var skillBody []byte

// claudeFrontmatter restricts allowed-tools so the skill cannot escape its
// documented workflow. Other tools do not yet support tool restrictions.
const claudeFrontmatter = `---
name: docs-ssot
description: Set up docs-ssot SSOT documentation structure — migrate existing docs, build, and validate
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash(docs-ssot *)
  - Bash(make docs*)
  - Bash(git diff *)
---
`

const minimalFrontmatter = `---
name: docs-ssot
description: Set up docs-ssot SSOT documentation structure — migrate existing docs, build, and validate
---
`

// Pre-computed per-tool content: frontmatter + shared body, assembled once at init.
var (
	claudeContent  []byte
	minimalContent []byte
)

func init() {
	claudeContent = append([]byte(claudeFrontmatter), skillBody...)
	minimalContent = append([]byte(minimalFrontmatter), skillBody...)
}

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

func skillContent(tool agentscan.Tool) []byte {
	if tool == agentscan.ToolClaude {
		return claudeContent
	}
	return minimalContent
}

// Install writes a SKILL.md file for each tool into the appropriate skills directory.
// If a file already exists, the user is prompted before overwriting.
func Install(tools []agentscan.Tool, in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)

	for _, tool := range tools {
		path := skillPath(tool)
		if path == "" {
			continue
		}

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

		if err := os.WriteFile(path, skillContent(tool), 0o600); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}

		if _, err := fmt.Fprintf(out, "Installed %s\n", path); err != nil {
			return fmt.Errorf("write install message: %w", err)
		}
	}

	return nil
}
