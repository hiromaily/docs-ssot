package cli

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/docs-ssot/internal/generator"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate documentation from templates",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return generator.Build(configFile)
	},
}
