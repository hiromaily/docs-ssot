package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintln(os.Stdout, "docs-ssot", appVersion)
	},
}
