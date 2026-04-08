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

With --agents, it scans for AI tool configuration files (rules, skills, commands)
and generates multi-tool templates from a single tool's configuration.`,
	Args: func(cmd *cobra.Command, args []string) error {
		agentsMode, _ := cmd.Flags().GetBool("agents")
		if !agentsMode && len(args) == 0 {
			return errors.New("requires at least 1 arg(s) when not using --agents")
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
		agentsMode, _ := flags.GetBool("agents")
		sourceTool, _ := flags.GetString("source-tool")
		targetToolsStr, _ := flags.GetString("target-tools")
		convertCommands, _ := flags.GetBool("convert-commands")
		inferGlobs, _ := flags.GetBool("infer-globs")

		if sectionLevel < 1 || sectionLevel > 6 {
			return fmt.Errorf("--section-level must be between 1 and 6, got %d", sectionLevel)
		}

		// Run agent migration if --agents is set.
		if agentsMode {
			targetTools, err := parseTargetTools(targetToolsStr)
			if err != nil {
				return err
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
	migrateCmd.Flags().String("output-dir", "template/sections", "where to write section files")
	migrateCmd.Flags().String("template-dir", "template/pages", "where to write template files")
	migrateCmd.Flags().Int("section-level", 2, "heading level used as section boundary (1–6)")
	migrateCmd.Flags().Float64("threshold", 0.82, "similarity threshold for duplicate detection (0.0–1.0)")
	migrateCmd.Flags().Bool("dry-run", false, "print the migration plan without writing files")
	migrateCmd.Flags().Bool("agents", false, "enable agent-aware migration mode")
	migrateCmd.Flags().String("source-tool", "auto", "source tool: auto, claude, cursor, copilot")
	migrateCmd.Flags().String("target-tools", "all", "target tools: all or comma-separated (claude,cursor,copilot,codex)")
	migrateCmd.Flags().Bool("convert-commands", false, "convert legacy commands to skills during migration")
	migrateCmd.Flags().Bool("infer-globs", false, "infer path-gated globs from rule slug names")
}

func parseTargetTools(s string) ([]agentscan.Tool, error) {
	if s == "all" || s == "" {
		return agentscan.AllTools(), nil
	}

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
