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

//go:embed skills/claude/SKILL.md
var claudeSkill []byte

//go:embed skills/cursor/SKILL.md
var cursorSkill []byte

//go:embed skills/copilot/SKILL.md
var copilotSkill []byte

//go:embed skills/codex/SKILL.md
var codexSkill []byte

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

// skillContent returns the embedded SKILL.md content for the given tool.
func skillContent(tool agentscan.Tool) []byte {
	switch tool {
	case agentscan.ToolClaude:
		return claudeSkill
	case agentscan.ToolCursor:
		return cursorSkill
	case agentscan.ToolCopilot:
		return copilotSkill
	case agentscan.ToolCodex:
		return codexSkill
	default:
		return nil
	}
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

		if err := os.WriteFile(path, content, 0o600); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}

		if _, err := fmt.Fprintf(out, "Installed %s\n", path); err != nil {
			return fmt.Errorf("write install message: %w", err)
		}
	}

	return nil
}
