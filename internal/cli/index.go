package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hiromaily/docs-ssot/internal/config"
	"github.com/hiromaily/docs-ssot/internal/index"
)

var indexOutput string

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generate template INDEX.md with include relationships and orphan detection",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		cfg, err := config.Load(configFile)
		if err != nil {
			return err
		}

		templateDir := index.DetectTemplateDir(cfg)
		data, err := index.Generate(templateDir, cfg)
		if err != nil {
			return err
		}

		content := index.Render(data)

		output := indexOutput
		if output == "" && cfg.Index.Output != "" {
			output = cfg.Index.Output
		}

		if output == "" {
			_, _ = fmt.Fprint(os.Stdout, content)
			return nil
		}

		if dir := filepath.Dir(output); dir != "." {
			//nolint:gosec // generated documentation directory
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create output directory %s: %w", dir, err)
			}
		}

		_, _ = fmt.Fprintln(os.Stdout, "Generating:", output)
		//nolint:gosec // generated documentation index is intended to be world-readable
		return os.WriteFile(output, []byte(content), 0o644)
	},
}

func init() {
	indexCmd.Flags().StringVar(&indexOutput, "output", "", "write index to file instead of stdout")
}
