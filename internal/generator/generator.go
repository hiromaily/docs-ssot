// Package generator builds documentation output files from template files.
package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hiromaily/docs-ssot/internal/config"
	"github.com/hiromaily/docs-ssot/internal/index"
	"github.com/hiromaily/docs-ssot/internal/processor"
)

// ErrValidationFailed is returned by Validate when one or more templates contain unresolvable includes.
var ErrValidationFailed = errors.New("validation failed")

// Build generates all output files defined in the given config file.
func Build(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	for _, t := range cfg.Targets {
		_, _ = fmt.Fprintln(os.Stdout, "Generating:", t.Output)

		content, err := processor.ProcessFile(t.Input, t.Output)
		if err != nil {
			return err
		}

		if dir := filepath.Dir(t.Output); dir != "." {
			//nolint:gosec // generated documentation directory
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create output directory %s: %w", dir, err)
			}
		}

		//nolint:gosec // generated documentation files are intended to be world-readable
		if err := os.WriteFile(t.Output, []byte(content), 0o644); err != nil {
			return err
		}
	}

	// Generate INDEX.md if configured
	if cfg.Index.Output != "" {
		if err := generateIndex(cfg); err != nil {
			return err
		}
	}

	return nil
}

// generateIndex generates the INDEX.md file from the template directory.
func generateIndex(cfg *config.Config) error {
	templateDir := index.DetectTemplateDir(cfg)
	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		return fmt.Errorf("index generation: %w", err)
	}

	content := index.Render(data)

	_, _ = fmt.Fprintln(os.Stdout, "Generating:", cfg.Index.Output)

	if dir := filepath.Dir(cfg.Index.Output); dir != "." {
		//nolint:gosec // generated documentation directory
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create index output directory %s: %w", dir, err)
		}
	}

	//nolint:gosec // generated index file is intended to be world-readable
	return os.WriteFile(cfg.Index.Output, []byte(content), 0o644)
}

// Validate performs a dry run over all templates in the given config file, checking that all
// include directives can be resolved without errors. No output files are written.
// It prints "OK" on success or one "ERROR: ..." line per failing template, then returns
// ErrValidationFailed so the caller can exit with a non-zero status code.
func Validate(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	ok := true
	for _, t := range cfg.Targets {
		if _, err := processor.ProcessFile(t.Input, t.Output); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "ERROR:", err)
			ok = false
		}
	}

	if ok {
		_, _ = fmt.Fprintln(os.Stdout, "OK")
		return nil
	}
	return ErrValidationFailed
}
