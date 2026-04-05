package main

import (
	"os"

	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "docs-ssot",
	Short: "Documentation SSOT generator",
	Long:  "docs-ssot generates documentation files from modular Markdown sources using a template-based composition system.",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "docsgen.yaml", "path to configuration file")
	rootCmd.AddCommand(buildCmd, includeCmd, validateCmd, versionCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
