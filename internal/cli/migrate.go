package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
	"github.com/hiromaily/docs-ssot/internal/migrate"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [files...]",
	Short: "Decompose existing Markdown files into SSOT section structure",
	Long: `migrate converts existing monolithic Markdown files (e.g., README.md, CLAUDE.md)
into the docs-ssot section structure.

It splits each file by H2 headings into candidate sections, detects duplicates
across files using TF-IDF cosine similarity, creates section files under
template/sections/<category>/, and generates template files with @include
directives that reproduce the original document structure.

With --from, it scans AI tool configuration files (rules, skills, commands)
from the specified tool and generates multi-tool templates for the target tools.

Examples:
  # Migrate Claude configs to all other tools
  docs-ssot migrate --from claude

  # Migrate Claude configs to specific tools
  docs-ssot migrate --from claude --to cursor,codex

  # Migrate with path inference and command conversion
  docs-ssot migrate --from claude --to cursor --infer-globs --convert-commands

  # Combine agent migration with file migration
  docs-ssot migrate --from claude --to cursor README.md CLAUDE.md

  # Migrate existing docs only (no agent migration)
  docs-ssot migrate README.md CLAUDE.md`,
	Args: func(cmd *cobra.Command, args []string) error {
		from, _ := cmd.Flags().GetString("from")
		agentsMode, _ := cmd.Flags().GetBool("agents")
		if from == "" && !agentsMode && len(args) == 0 {
			return errors.New("requires --from <tool> or at least 1 file argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		flags := cmd.Flags()
		outputDir, _ := flags.GetString("output-dir")
		templateDir, _ := flags.GetString("template-dir")
		sectionLevel, _ := flags.GetInt("section-level")
		threshold, _ := flags.GetFloat64("threshold")
		dryRun, _ := flags.GetBool("dry-run")
		convertCommands, _ := flags.GetBool("convert-commands")
		inferGlobs, _ := flags.GetBool("infer-globs")

		// Resolve source tool: --from takes priority over legacy --source-tool.
		from, _ := flags.GetString("from")
		if from == "" {
			from, _ = flags.GetString("source-tool")
		}

		// Resolve target tools: --to takes priority over legacy --target-tools.
		toString, _ := flags.GetString("to")
		if toString == "" {
			toString, _ = flags.GetString("target-tools")
		}

		// Legacy --agents flag: --from implicitly enables agent mode.
		agentsMode, _ := flags.GetBool("agents")
		if from != "" {
			agentsMode = true
		}

		if sectionLevel < 1 || sectionLevel > 6 {
			return fmt.Errorf("--section-level must be between 1 and 6, got %d", sectionLevel)
		}

		// Run agent migration if enabled.
		if agentsMode {
			targetTools, err := resolveTargetTools(toString, from)
			if err != nil {
				return err
			}

			sourceTool := from
			if sourceTool == "" {
				sourceTool = "auto"
			}

			agentCfg := migrate.AgentConfig{
				RootDir:         ".",
				SourceTool:      sourceTool,
				TargetTools:     targetTools,
				OutputDir:       outputDir,
				TemplateDir:     templateDir,
				DryRun:          dryRun,
				ConfigFile:      configFile,
				ConvertCommands: convertCommands,
				InferGlobs:      inferGlobs,
			}

			if err := migrate.RunAgents(os.Stdout, agentCfg); err != nil {
				return err
			}
		}

		// Run regular file migration if files are provided.
		if len(args) > 0 {
			cfg := migrate.Config{
				InputFiles:   args,
				OutputDir:    outputDir,
				TemplateDir:  templateDir,
				SectionLevel: sectionLevel,
				Threshold:    threshold,
				DryRun:       dryRun,
				ConfigFile:   configFile,
			}
			return migrate.Run(os.Stdout, cfg)
		}

		return nil
	},
}

func init() {
	// File migration flags.
	migrateCmd.Flags().String("output-dir", "template/sections", "where to write section files")
	migrateCmd.Flags().String("template-dir", "template/pages", "where to write template files")
	migrateCmd.Flags().Int("section-level", 2, "heading level used as section boundary (1–6)")
	migrateCmd.Flags().Float64("threshold", 0.82, "similarity threshold for duplicate detection (0.0–1.0)")
	migrateCmd.Flags().Bool("dry-run", false, "print the migration plan without writing files")

	// Agent migration flags.
	migrateCmd.Flags().String("from", "", "source AI tool to migrate from (claude, cursor, copilot)")
	migrateCmd.Flags().String("to", "", "target AI tools, comma-separated (default: all except --from)")
	migrateCmd.Flags().Bool("convert-commands", false, "convert legacy commands to skills during migration")
	migrateCmd.Flags().Bool("infer-globs", false, "infer path-gated globs from rule slug names")

	// Legacy flags (backward-compatible, hidden).
	migrateCmd.Flags().Bool("agents", false, "enable agent-aware migration mode (use --from instead)")
	migrateCmd.Flags().String("source-tool", "", "source tool (use --from instead)")
	migrateCmd.Flags().String("target-tools", "", "target tools (use --to instead)")
	_ = migrateCmd.Flags().MarkHidden("agents")
	_ = migrateCmd.Flags().MarkHidden("source-tool")
	_ = migrateCmd.Flags().MarkHidden("target-tools")
}

// resolveTargetTools determines the target tools based on --to and --from flags.
// When --to is empty, all tools except the source tool are used.
func resolveTargetTools(toString, from string) ([]agentscan.Tool, error) {
	if toString != "" && toString != "all" {
		return parseToolList(toString)
	}

	// Default: all tools except the source.
	all := agentscan.AllTools()
	if from == "" || from == "auto" {
		return all, nil
	}

	sourceTool, err := agentscan.ParseTool(from)
	if err != nil {
		return nil, err
	}

	tools := make([]agentscan.Tool, 0, len(all)-1)
	for _, t := range all {
		if t != sourceTool {
			tools = append(tools, t)
		}
	}
	return tools, nil
}

func parseToolList(s string) ([]agentscan.Tool, error) {
	seen := map[agentscan.Tool]bool{}
	var tools []agentscan.Tool
	for name := range strings.SplitSeq(s, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		tool, err := agentscan.ParseTool(name)
		if err != nil {
			return nil, err
		}
		if !seen[tool] {
			seen[tool] = true
			tools = append(tools, tool)
		}
	}
	return tools, nil
}
