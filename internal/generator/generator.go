package generator

import (
	"fmt"
	"os"

	"github.com/hiromaily/docs-ssot/internal/config"
	"github.com/hiromaily/docs-ssot/internal/include"
)

func Build(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	for _, t := range cfg.Targets {
		fmt.Println("Generating:", t.Output)

		content, err := include.ProcessFile(t.Input, t.Output)
		if err != nil {
			return err
		}

		if err := os.WriteFile(t.Output, []byte(content), 0o644); err != nil {
			return err
		}
	}

	return nil
}
