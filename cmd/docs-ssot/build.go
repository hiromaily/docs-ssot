package main

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/docs-ssot/internal/generator"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate documentation from templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		return generator.Build(configFile)
	},
}
