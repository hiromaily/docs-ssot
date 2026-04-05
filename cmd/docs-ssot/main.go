package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hiromaily/docs-ssot/internal/generator"
	"github.com/hiromaily/docs-ssot/internal/include"
)

// appVersion is set at build time via -ldflags "-X main.appVersion=<version>".
var appVersion = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "build":
		if err := generator.Build("docsgen.yaml"); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "include":
		if len(os.Args) < 3 {
			_, _ = fmt.Fprintln(os.Stderr, "Usage: docs-ssot include <file>")
			os.Exit(1)
		}
		content, err := include.ProcessFile(os.Args[2], os.Args[2])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		_, _ = fmt.Fprint(os.Stdout, content)
	case "validate":
		err := generator.Validate("docsgen.yaml")
		if err != nil && !errors.Is(err, generator.ErrValidationFailed) {
			// ErrValidationFailed is already reported line-by-line by Validate itself.
			_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		}
		if err != nil {
			os.Exit(1)
		}
	case "version":
		_, _ = fmt.Fprintln(os.Stdout, "docs-ssot", appVersion)
	case "help", "--help", "-h":
		printUsage()
	default:
		_, _ = fmt.Fprintln(os.Stderr, "Unknown command:", cmd)
		os.Exit(1)
	}
}

func printUsage() {
	_, _ = fmt.Fprintln(os.Stdout, "Usage: docs-ssot <command>")
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintln(os.Stdout, "Commands:")
	_, _ = fmt.Fprintln(os.Stdout, "  build             Generate documentation from templates")
	_, _ = fmt.Fprintln(os.Stdout, "  include <file>    Expand includes in <file> and print to stdout")
	_, _ = fmt.Fprintln(os.Stdout, "  validate          Check that all include directives can be resolved")
	_, _ = fmt.Fprintln(os.Stdout, "  version           Print version")
	_, _ = fmt.Fprintln(os.Stdout, "  help              Show this help message")
}
