package migrate_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hiromaily/docs-ssot/internal/migrate"
)

func TestRun_DryRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	readme := filepath.Join(dir, "README.md")
	err := os.WriteFile(readme, []byte(`# My Project

Welcome to the project.

## Setup

Install with go install.

## Testing

Run go test ./...

## Architecture

Layered architecture.
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cfg := migrate.Config{
		InputFiles:   []string{readme},
		OutputDir:    filepath.Join(dir, "template/sections"),
		TemplateDir:  filepath.Join(dir, "template/pages"),
		SectionLevel: 2,
		Threshold:    0.82,
		DryRun:       true,
		ConfigFile:   filepath.Join(dir, "docsgen.yaml"),
	}

	if err := migrate.Run(&buf, cfg); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Parsed") {
		t.Errorf("expected 'Parsed' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "4 sections") {
		t.Errorf("expected '4 sections' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "Would create") {
		t.Errorf("expected 'Would create' in dry-run output, got:\n%s", output)
	}

	// Verify no files were written.
	if _, err := os.Stat(filepath.Join(dir, "template/sections")); !os.IsNotExist(err) {
		t.Error("expected no files written in dry-run mode")
	}
}

func TestRun_WritesFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	readme := filepath.Join(dir, "README.md")
	err := os.WriteFile(readme, []byte(`# My Project

Welcome.

## Setup

Install deps.

## Testing

Run tests.
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cfg := migrate.Config{
		InputFiles:   []string{readme},
		OutputDir:    filepath.Join(dir, "template/sections"),
		TemplateDir:  filepath.Join(dir, "template/pages"),
		SectionLevel: 2,
		Threshold:    0.82,
		DryRun:       false,
		ConfigFile:   filepath.Join(dir, "docsgen.yaml"),
	}

	if err := migrate.Run(&buf, cfg); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	// Verify section files were created.
	setupPath := filepath.Join(dir, "template/sections/development/setup.md")
	if _, err := os.Stat(setupPath); err != nil {
		t.Errorf("expected section file %s to exist", setupPath)
	}

	// Verify section content preserves heading.
	setupContent, err := os.ReadFile(setupPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(setupContent), "## Setup") {
		t.Errorf("section file should start with H2 heading, got: %q", string(setupContent)[:30])
	}

	testingPath := filepath.Join(dir, "template/sections/development/testing.md")
	if _, err := os.Stat(testingPath); err != nil {
		t.Errorf("expected section file %s to exist", testingPath)
	}

	// Verify template file was created with @include directives.
	tplPath := filepath.Join(dir, "template/pages/README.tpl.md")
	tplContent, err := os.ReadFile(tplPath)
	if err != nil {
		t.Fatalf("expected template file: %v", err)
	}
	if !strings.Contains(string(tplContent), "<!-- @include:") {
		t.Error("expected @include directives in template")
	}
	// Template should NOT contain --- separators.
	if strings.Contains(string(tplContent), "\n---\n") {
		t.Error("template should not contain --- separators")
	}

	// Verify docsgen.yaml was created.
	configPath := filepath.Join(dir, "docsgen.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("expected docsgen.yaml to exist")
	}

	// Verify output mentions "Creating" (not "Would create").
	output := buf.String()
	if !strings.Contains(output, "Creating") {
		t.Errorf("expected 'Creating' in non-dry-run output, got:\n%s", output)
	}
}

func TestRun_DuplicateDetection(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	readme := filepath.Join(dir, "README.md")
	err := os.WriteFile(readme, []byte(`## Setup

To set up this project, you need to install Go 1.26 or later.
Then run make install to install all dependencies.
After that, run make build to compile the project.
Finally, verify the installation by running make test.

## Architecture

The system uses a layered architecture pattern.
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	claude := filepath.Join(dir, "CLAUDE.md")
	err = os.WriteFile(claude, []byte(`## Setup

To set up this project, you need to install Go 1.26 or later.
Then run make install to install all dependencies.
After that, run make build to compile the project.
Finally, verify the installation by running make test.

## Contributing

Please submit a PR.
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cfg := migrate.Config{
		InputFiles:   []string{readme, claude},
		OutputDir:    filepath.Join(dir, "template/sections"),
		TemplateDir:  filepath.Join(dir, "template/pages"),
		SectionLevel: 2,
		Threshold:    0.82,
		DryRun:       true,
		ConfigFile:   filepath.Join(dir, "docsgen.yaml"),
	}

	if err := migrate.Run(&buf, cfg); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "duplicate") {
		t.Errorf("expected duplicate detection in output, got:\n%s", output)
	}
	if !strings.Contains(output, "merged into") {
		t.Errorf("expected 'merged into' in output, got:\n%s", output)
	}
}

func TestRun_ConfigAlreadyExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	readme := filepath.Join(dir, "README.md")
	err := os.WriteFile(readme, []byte("## Hello\n\nWorld.\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Pre-create docsgen.yaml.
	configPath := filepath.Join(dir, "docsgen.yaml")
	err = os.WriteFile(configPath, []byte("targets: []\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cfg := migrate.Config{
		InputFiles:   []string{readme},
		OutputDir:    filepath.Join(dir, "template/sections"),
		TemplateDir:  filepath.Join(dir, "template/pages"),
		SectionLevel: 2,
		Threshold:    0.82,
		DryRun:       false,
		ConfigFile:   configPath,
	}

	if err := migrate.Run(&buf, cfg); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if !strings.Contains(buf.String(), "already exists") {
		t.Errorf("expected 'already exists' in output, got:\n%s", buf.String())
	}
}

func TestToSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple", input: "Setup", want: "setup"},
		{name: "multi_word", input: "Getting Started", want: "getting-started"},
		{name: "special_chars", input: "API Reference (v2)", want: "api-reference-v2"},
		{name: "trailing_special", input: "FAQ!", want: "faq"},
		{name: "numbers", input: "Step 1 Guide", want: "step-1-guide"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := migrate.ToSlug(tt.input)
			if got != tt.want {
				t.Errorf("ToSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
