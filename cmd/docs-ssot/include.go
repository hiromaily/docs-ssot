package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hiromaily/docs-ssot/internal/include"
)

var includeCmd = &cobra.Command{
	Use:   "include <file>",
	Short: "Expand include directives in <file> and print to stdout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := include.ProcessFile(args[0], args[0])
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(os.Stdout, content)

		return err
	},
}
