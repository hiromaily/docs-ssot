package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

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
directives that reproduce the original document structure.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		flags := cmd.Flags()
		outputDir, _ := flags.GetString("output-dir")
		templateDir, _ := flags.GetString("template-dir")
		sectionLevel, _ := flags.GetInt("section-level")
		threshold, _ := flags.GetFloat64("threshold")
		dryRun, _ := flags.GetBool("dry-run")

		if sectionLevel < 1 || sectionLevel > 6 {
			return fmt.Errorf("--section-level must be between 1 and 6, got %d", sectionLevel)
		}

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
	},
}

func init() {
	migrateCmd.Flags().String("output-dir", "template/sections", "where to write section files")
	migrateCmd.Flags().String("template-dir", "template/pages", "where to write template files")
	migrateCmd.Flags().Int("section-level", 2, "heading level used as section boundary (1–6)")
	migrateCmd.Flags().Float64("threshold", 0.82, "similarity threshold for duplicate detection (0.0–1.0)")
	migrateCmd.Flags().Bool("dry-run", false, "print the migration plan without writing files")
}
