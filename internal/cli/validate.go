package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/hiromaily/docs-ssot/internal/generator"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check that all include directives can be resolved",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := generator.Validate(configFile)
		if errors.Is(err, generator.ErrValidationFailed) {
			// Per-template errors are already printed by Validate itself.
			// Silence Cobra's "Error: validation failed" to avoid double-reporting.
			cmd.Root().SilenceErrors = true
		}

		return err
	},
}
