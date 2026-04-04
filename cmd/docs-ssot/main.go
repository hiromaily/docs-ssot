package main

import (
	"fmt"
	"os"

	"github.com/hiromaily/docs-ssot/internal/generator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: docs-ssot build")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "build":
		if err := generator.Build("docsgen.yaml"); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", cmd)
		os.Exit(1)
	}
}
