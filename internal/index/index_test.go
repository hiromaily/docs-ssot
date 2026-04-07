package index_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hiromaily/docs-ssot/internal/config"
	"github.com/hiromaily/docs-ssot/internal/index"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// setupTemplateTree creates a minimal template directory structure for testing.
func setupTemplateTree(t *testing.T) (templateDir string, cfg *config.Config) {
	t.Helper()
	dir := t.TempDir()
	templateDir = filepath.Join(dir, "template")

	// sections/
	writeFile(t, filepath.Join(templateDir, "sections", "project", "overview.md"), "# Overview\n")
	writeFile(t, filepath.Join(templateDir, "sections", "project", "vision.md"), "# Vision\n")
	writeFile(t, filepath.Join(templateDir, "sections", "architecture", "system.md"), "# System\n")

	// sections/ai/rules/
	writeFile(t, filepath.Join(templateDir, "sections", "ai", "rules", "general.md"), "# General Rules\n")
	writeFile(t, filepath.Join(templateDir, "sections", "ai", "rules", "git.md"), "# Git Rules\n")

	// sections/ai/commands/
	writeFile(t, filepath.Join(templateDir, "sections", "ai", "commands", "fix-pr-reviews.md"), "# Fix PR Reviews\n")

	// pages/
	writeFile(t, filepath.Join(templateDir, "pages", "README.tpl.md"),
		"# README\n<!-- @include: ../sections/project/overview.md -->\n<!-- @include: ../sections/architecture/system.md -->\n")

	writeFile(t, filepath.Join(templateDir, "pages", "CLAUDE.tpl.md"),
		"# CLAUDE\n<!-- @include: ../sections/project/overview.md -->\n<!-- @include: ../sections/project/vision.md -->\n<!-- @include: ../sections/architecture/system.md -->\n")

	// pages/ai-agents/claude/rules/
	writeFile(t, filepath.Join(templateDir, "pages", "ai-agents", "claude", "rules", "general.tpl.md"),
		"<!-- @include: ../../../../sections/ai/rules/general.md -->\n")

	writeFile(t, filepath.Join(templateDir, "pages", "ai-agents", "claude", "rules", "git.tpl.md"),
		"<!-- @include: ../../../../sections/ai/rules/git.md -->\n")

	// pages/ai-agents/claude/commands/
	writeFile(t, filepath.Join(templateDir, "pages", "ai-agents", "claude", "commands", "fix-pr-reviews.tpl.md"),
		"<!-- @include: ../../../../sections/ai/commands/fix-pr-reviews.md -->\n")

	// pages/ai-agents/cursor/rules/
	writeFile(t, filepath.Join(templateDir, "pages", "ai-agents", "cursor", "rules", "general.tpl.mdc"),
		"---\ndescription: general\nalwaysApply: true\n---\n<!-- @include: ../../../../sections/ai/rules/general.md -->\n")

	cfg = &config.Config{
		Targets: []config.Target{
			{Input: filepath.Join(templateDir, "pages", "README.tpl.md"), Output: "README.md"},
			{Input: filepath.Join(templateDir, "pages", "CLAUDE.tpl.md"), Output: "CLAUDE.md"},
			{Input: filepath.Join(templateDir, "pages", "ai-agents", "claude", "rules", "general.tpl.md"), Output: ".claude/rules/general.md"},
			{Input: filepath.Join(templateDir, "pages", "ai-agents", "claude", "rules", "git.tpl.md"), Output: ".claude/rules/git.md"},
			{Input: filepath.Join(templateDir, "pages", "ai-agents", "claude", "commands", "fix-pr-reviews.tpl.md"), Output: ".claude/commands/fix-pr-reviews.md"},
			{Input: filepath.Join(templateDir, "pages", "ai-agents", "cursor", "rules", "general.tpl.mdc"), Output: ".cursor/rules/general.mdc"},
		},
	}

	return templateDir, cfg
}

func TestGenerate_BasicStructure(t *testing.T) {
	t.Parallel()
	templateDir, cfg := setupTemplateTree(t)

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Should have 6 page templates (2 regular + 4 ai-agents)
	if len(data.Pages) != 6 {
		t.Errorf("pages count = %d, want 6", len(data.Pages))
	}

	// Sections should include project and architecture files
	if len(data.Sections) == 0 {
		t.Error("expected sections to be non-empty")
	}

	// Rules should include general.md and git.md
	if len(data.Rules) != 2 {
		t.Errorf("rules count = %d, want 2", len(data.Rules))
	}

	// Commands should include fix-pr-reviews.md
	if len(data.Commands) != 1 {
		t.Errorf("commands count = %d, want 1", len(data.Commands))
	}

	// No orphans
	if len(data.Orphans) != 0 {
		t.Errorf("orphans = %v, want none", data.Orphans)
	}
}

func TestGenerate_OrphanDetection(t *testing.T) {
	t.Parallel()
	templateDir, cfg := setupTemplateTree(t)

	// Add an orphan file
	writeFile(t, filepath.Join(templateDir, "sections", "project", "orphan.md"), "# Orphan\n")

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(data.Orphans) != 1 {
		t.Fatalf("orphans count = %d, want 1; got %v", len(data.Orphans), data.Orphans)
	}

	if !strings.HasSuffix(data.Orphans[0], "orphan.md") {
		t.Errorf("orphan = %q, want to end with orphan.md", data.Orphans[0])
	}
}

func TestGenerate_MultipleReferences(t *testing.T) {
	t.Parallel()
	templateDir, cfg := setupTemplateTree(t)

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// docs/sections/project/overview.md is referenced by both README and CLAUDE pages
	overviewKey := "sections/project/overview.md"
	refs, ok := data.Sections[overviewKey]
	if !ok {
		t.Fatalf("expected %q in sections", overviewKey)
	}

	if len(refs) < 2 {
		t.Errorf("refs for %q = %v, want at least 2 (README and CLAUDE)", overviewKey, refs)
	}
}

func TestGenerate_RulesReferencedByMultipleTools(t *testing.T) {
	t.Parallel()
	templateDir, cfg := setupTemplateTree(t)

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// docs/rules/general.md should be referenced by both claude and cursor
	generalKey := "sections/ai/rules/general.md"
	refs, ok := data.Rules[generalKey]
	if !ok {
		t.Fatalf("expected %q in rules", generalKey)
	}

	if len(refs) < 2 {
		t.Errorf("refs for %q = %v, want at least 2 (claude and cursor)", generalKey, refs)
	}
}

func TestGenerate_GlobInclude(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	templateDir := filepath.Join(dir, "template")

	writeFile(t, filepath.Join(templateDir, "sections", "ai", "claude.md"), "# Claude\n")
	writeFile(t, filepath.Join(templateDir, "sections", "ai", "cursor.md"), "# Cursor\n")

	writeFile(t, filepath.Join(templateDir, "pages", "README.tpl.md"),
		"<!-- @include: ../sections/ai/*.md -->\n")

	cfg := &config.Config{
		Targets: []config.Target{
			{Input: filepath.Join(templateDir, "pages", "README.tpl.md"), Output: "README.md"},
		},
	}

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Both ai files should be referenced
	if len(data.Sections) != 2 {
		t.Errorf("sections count = %d, want 2; got keys: %v", len(data.Sections), keysOf(data.Sections))
	}

	if len(data.Orphans) != 0 {
		t.Errorf("orphans = %v, want none", data.Orphans)
	}
}

func TestGenerate_DirectoryInclude(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	templateDir := filepath.Join(dir, "template")

	writeFile(t, filepath.Join(templateDir, "sections", "product", "concept.md"), "# Concept\n")
	writeFile(t, filepath.Join(templateDir, "sections", "product", "features.md"), "# Features\n")

	writeFile(t, filepath.Join(templateDir, "pages", "README.tpl.md"),
		"<!-- @include: ../sections/product/ -->\n")

	cfg := &config.Config{
		Targets: []config.Target{
			{Input: filepath.Join(templateDir, "pages", "README.tpl.md"), Output: "README.md"},
		},
	}

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(data.Sections) != 2 {
		t.Errorf("sections count = %d, want 2; got keys: %v", len(data.Sections), keysOf(data.Sections))
	}

	if len(data.Orphans) != 0 {
		t.Errorf("orphans = %v, want none", data.Orphans)
	}
}

func TestGenerate_RecursiveGlobInclude(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	templateDir := filepath.Join(dir, "template")

	writeFile(t, filepath.Join(templateDir, "sections", "a.md"), "a\n")
	writeFile(t, filepath.Join(templateDir, "sections", "sub", "b.md"), "b\n")

	writeFile(t, filepath.Join(templateDir, "pages", "README.tpl.md"),
		"<!-- @include: ../sections/**/*.md -->\n")

	cfg := &config.Config{
		Targets: []config.Target{
			{Input: filepath.Join(templateDir, "pages", "README.tpl.md"), Output: "README.md"},
		},
	}

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(data.Sections) != 2 {
		t.Errorf("sections count = %d, want 2; got keys: %v", len(data.Sections), keysOf(data.Sections))
	}

	if len(data.Orphans) != 0 {
		t.Errorf("orphans = %v, want none", data.Orphans)
	}
}

func TestRender_DeterministicOutput(t *testing.T) {
	t.Parallel()
	templateDir, cfg := setupTemplateTree(t)

	data1, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	output1 := index.Render(data1)
	output2 := index.Render(data2)

	if output1 != output2 {
		t.Error("two consecutive renders produced different output")
	}
}

func TestRender_ContainsExpectedSections(t *testing.T) {
	t.Parallel()
	templateDir, cfg := setupTemplateTree(t)

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	output := index.Render(data)

	expected := []string{
		"<!-- AUTO-GENERATED FILE — DO NOT EDIT -->",
		"# Template Index",
		"## Pages",
		"## Sections",
		"## Rules",
		"## Commands",
		"## Orphans",
		"| (none) | All files are referenced |",
	}

	for _, s := range expected {
		if !strings.Contains(output, s) {
			t.Errorf("output missing expected string: %q", s)
		}
	}
}

func TestRender_OrphansListed(t *testing.T) {
	t.Parallel()
	templateDir, cfg := setupTemplateTree(t)

	writeFile(t, filepath.Join(templateDir, "sections", "orphan.md"), "# Orphan\n")

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	output := index.Render(data)

	if !strings.Contains(output, "orphan.md") {
		t.Error("orphan file not listed in rendered output")
	}

	if strings.Contains(output, "All files are referenced") {
		t.Error("should not show 'All files are referenced' when orphans exist")
	}
}

func TestRender_OutputFlag(t *testing.T) {
	t.Parallel()
	templateDir, cfg := setupTemplateTree(t)

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	content := index.Render(data)

	outputPath := filepath.Join(t.TempDir(), "INDEX.md")
	if err := os.WriteFile(outputPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	written, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(written) != content {
		t.Error("written file content doesn't match rendered content")
	}
}

func TestDetectTemplateDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		inputs []string
		want   string
	}{
		{
			name:   "standard",
			inputs: []string{"template/pages/README.tpl.md"},
			want:   "template",
		},
		{
			name:   "nested",
			inputs: []string{"poc/docs-ssot/template/pages/README.tpl.md"},
			want:   "poc/docs-ssot/template",
		},
		{
			name:   "no_template",
			inputs: []string{"src/main.go"},
			want:   "template",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var targets []config.Target
			for _, input := range tc.inputs {
				targets = append(targets, config.Target{Input: input})
			}
			cfg := &config.Config{Targets: targets}
			got := index.DetectTemplateDir(cfg)
			if got != tc.want {
				t.Errorf("DetectTemplateDir() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGenerate_IncludeInsideCodeFenceIgnored(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	templateDir := filepath.Join(dir, "template")

	writeFile(t, filepath.Join(templateDir, "sections", "guide.md"), "# Guide\n")

	// Template has an include inside a code fence — should NOT count as a reference
	writeFile(t, filepath.Join(templateDir, "pages", "README.tpl.md"),
		"```md\n<!-- @include: ../sections/guide.md -->\n```\n")

	cfg := &config.Config{
		Targets: []config.Target{
			{Input: filepath.Join(templateDir, "pages", "README.tpl.md"), Output: "README.md"},
		},
	}

	data, err := index.Generate(templateDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// guide.md is not actually included (it's in a code fence), so it should be an orphan
	if len(data.Orphans) != 1 {
		t.Errorf("orphans count = %d, want 1; orphans = %v", len(data.Orphans), data.Orphans)
	}
}

func keysOf(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
