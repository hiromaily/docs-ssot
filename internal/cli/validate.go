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
		if err != nil && !errors.Is(err, generator.ErrValidationFailed) {
			// ErrValidationFailed is already reported line-by-line by Validate itself.
			return err
		}

		return err
	},
}
