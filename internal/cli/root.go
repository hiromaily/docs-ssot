// Package cli wires together all Cobra subcommands for docs-ssot.
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	appVersion string
)

var rootCmd = &cobra.Command{
	Use:   "docs-ssot",
	Short: "Documentation SSOT generator",
	Long:  "docs-ssot generates documentation files from modular Markdown sources using a template-based composition system.",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "docsgen.yaml", "path to configuration file")
	rootCmd.AddCommand(buildCmd, checkCmd, includeCmd, indexCmd, installSkillCmd, migrateCmd, validateCmd, versionCmd)
}

// Execute initialises the CLI with the given build version and runs the root command.
func Execute(version string) {
	appVersion = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
