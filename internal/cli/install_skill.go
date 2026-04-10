package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
	"github.com/hiromaily/docs-ssot/internal/skillinstall"
)

var installSkillCmd = &cobra.Command{
	Use:   "install-skill",
	Short: "Install the docs-ssot skill for AI coding agents",
	Long: `install-skill writes a SKILL.md file for each specified AI tool so that
the agent knows how to migrate existing documentation to the docs-ssot SSOT
structure, build output files, and validate the result.

Supported tools: claude, cursor, copilot, codex
Default: all four tools

Examples:
  # Install for all tools
  docs-ssot install-skill

  # Install for Claude Code only
  docs-ssot install-skill --tool claude

  # Install for multiple specific tools
  docs-ssot install-skill --tool claude,cursor`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cmd.SilenceUsage = true

		toolStr, _ := cmd.Flags().GetString("tool")

		var tools []agentscan.Tool
		if toolStr == "" {
			tools = agentscan.AllTools()
		} else {
			var err error
			tools, err = ParseToolList(toolStr)
			if err != nil {
				return err
			}
		}

		return skillinstall.Install(tools, os.Stdin, os.Stdout)
	},
}

func init() {
	installSkillCmd.Flags().String("tool", "", "target tool(s): claude, cursor, copilot, codex (comma-separated; default: all)")
}
