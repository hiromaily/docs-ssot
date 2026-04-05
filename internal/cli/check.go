package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hiromaily/docs-ssot/internal/dupcheck"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check docs for SSOT violations by detecting near-duplicate sections",
	Long: `check scans Markdown files under the docs directory for sections with similar content.

It uses TF-IDF cosine similarity to find near-duplicate sections across different files,
which indicates potential SSOT (Single Source of Truth) violations where the same information
exists in multiple places.

A similarity score of 1.0 means identical content; the default threshold of 0.82 catches
near-duplicates while filtering out loosely related content.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cmd.SilenceUsage = true

		flags := cmd.Flags()
		root, _ := flags.GetString("root")
		threshold, _ := flags.GetFloat64("threshold")
		minChars, _ := flags.GetInt("min-chars")
		sectionLevel, _ := flags.GetInt("section-level")
		format, _ := flags.GetString("format")
		excludes, _ := flags.GetStringArray("exclude")

		if sectionLevel < 1 || sectionLevel > 6 {
			return fmt.Errorf("--section-level must be between 1 and 6, got %d", sectionLevel)
		}
		if format != "text" && format != "json" {
			return fmt.Errorf("--format must be text or json, got %q", format)
		}

		cfg := dupcheck.Config{
			Root:         root,
			Threshold:    threshold,
			MinChars:     minChars,
			SectionLevel: sectionLevel,
			Format:       format,
			Excludes:     excludes,
		}
		return dupcheck.Run(os.Stdout, cfg)
	},
}

func init() {
	checkCmd.Flags().String("root", "docs", "root directory to scan for Markdown files")
	checkCmd.Flags().Float64("threshold", 0.82, "similarity threshold (0.0–1.0); sections scoring above this are reported")
	checkCmd.Flags().Int("min-chars", 120, "minimum character count for a section to be included in comparison")
	checkCmd.Flags().Int("section-level", 2, "heading level used as section boundary (1–6)")
	checkCmd.Flags().String("format", "text", "output format: text or json")
	checkCmd.Flags().StringArray("exclude", nil, "exclude path pattern (repeatable, e.g. --exclude docs/changelog/**)")
}
