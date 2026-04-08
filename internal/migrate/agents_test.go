package migrate_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
	"github.com/hiromaily/docs-ssot/internal/migrate"
)

func TestRunAgents_DryRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupAgentTestFiles(t, dir)

	var buf bytes.Buffer
	cfg := migrate.AgentConfig{
		RootDir:     dir,
		SourceTool:  "auto",
		TargetTools: agentscan.AllTools(),
		OutputDir:   filepath.Join(dir, "template/sections"),
		TemplateDir: filepath.Join(dir, "template/pages"),
		DryRun:      true,
		ConfigFile:  filepath.Join(dir, "docsgen.yaml"),
	}

	if err := migrate.RunAgents(&buf, cfg); err != nil {
		t.Fatalf("RunAgents() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Detected source tool: claude") {
		t.Errorf("expected source tool detection, got:\n%s", output)
	}
	if !strings.Contains(output, "Would create") {
		t.Errorf("expected 'Would create' in dry-run output, got:\n%s", output)
	}

	// Verify no files were written.
	sectionsDir := filepath.Join(dir, "template/sections/ai")
	if _, err := os.Stat(sectionsDir); !os.IsNotExist(err) {
		t.Error("expected no files written in dry-run mode")
	}
}

func TestRunAgents_WritesFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupAgentTestFiles(t, dir)

	var buf bytes.Buffer
	cfg := migrate.AgentConfig{
		RootDir:     dir,
		SourceTool:  "claude",
		TargetTools: []agentscan.Tool{agentscan.ToolClaude, agentscan.ToolCursor},
		OutputDir:   filepath.Join(dir, "template/sections"),
		TemplateDir: filepath.Join(dir, "template/pages"),
		DryRun:      false,
		ConfigFile:  filepath.Join(dir, "docsgen.yaml"),
	}

	if err := migrate.RunAgents(&buf, cfg); err != nil {
		t.Fatalf("RunAgents() error: %v", err)
	}

	// Verify section files were created.
	archSection := filepath.Join(dir, "template/sections/ai/rules/architecture.md")
	if _, err := os.Stat(archSection); err != nil {
		t.Errorf("expected section file %s to exist", archSection)
	}

	// Section content should start with H2 (shifted from H1).
	data, err := os.ReadFile(archSection)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "## Architecture") {
		t.Errorf("section should start with H2, got: %q", string(data)[:30])
	}

	// Verify Claude template was created.
	claudeTpl := filepath.Join(dir, "template/pages/ai-agents/claude/rules/architecture.tpl.md")
	if _, err := os.Stat(claudeTpl); err != nil {
		t.Errorf("expected Claude template %s to exist", claudeTpl)
	}

	// Verify Cursor template was created with .mdc extension.
	cursorTpl := filepath.Join(dir, "template/pages/ai-agents/cursor/rules/architecture.tpl.mdc")
	if _, err := os.Stat(cursorTpl); err != nil {
		t.Errorf("expected Cursor template %s to exist", cursorTpl)
	}

	// Cursor template should contain MDC frontmatter.
	cursorData, err := os.ReadFile(cursorTpl)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cursorData), "alwaysApply: true") {
		t.Errorf("Cursor template should contain alwaysApply, got:\n%s", string(cursorData))
	}

	// Verify docsgen.yaml was created.
	if _, err := os.Stat(filepath.Join(dir, "docsgen.yaml")); err != nil {
		t.Errorf("expected docsgen.yaml to exist")
	}
}

func TestRunAgents_SkillPreservesFields(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupAgentTestFiles(t, dir)

	var buf bytes.Buffer
	cfg := migrate.AgentConfig{
		RootDir:     dir,
		SourceTool:  "claude",
		TargetTools: []agentscan.Tool{agentscan.ToolClaude},
		OutputDir:   filepath.Join(dir, "template/sections"),
		TemplateDir: filepath.Join(dir, "template/pages"),
		DryRun:      false,
		ConfigFile:  filepath.Join(dir, "docsgen.yaml"),
	}

	if err := migrate.RunAgents(&buf, cfg); err != nil {
		t.Fatalf("RunAgents() error: %v", err)
	}

	// Verify skill template preserves Claude-specific fields.
	skillTpl := filepath.Join(dir, "template/pages/ai-agents/claude/skills/deploy/SKILL.tpl.md")
	data, err := os.ReadFile(skillTpl)
	if err != nil {
		t.Fatalf("expected skill template to exist: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "name: deploy") {
		t.Errorf("expected name field, got:\n%s", content)
	}
	if !strings.Contains(content, "description: Deploy to production") {
		t.Errorf("expected description field, got:\n%s", content)
	}
}

func TestRunAgents_NoToolsDetected(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	var buf bytes.Buffer
	cfg := migrate.AgentConfig{
		RootDir:     dir,
		SourceTool:  "auto",
		TargetTools: agentscan.AllTools(),
		OutputDir:   filepath.Join(dir, "template/sections"),
		TemplateDir: filepath.Join(dir, "template/pages"),
		DryRun:      true,
		ConfigFile:  filepath.Join(dir, "docsgen.yaml"),
	}

	err := migrate.RunAgents(&buf, cfg)
	if err == nil {
		t.Fatal("expected error for empty repo")
	}
	if !strings.Contains(err.Error(), "no AI tool configurations detected") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunAgents_SpecificSourceTool(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupAgentTestFiles(t, dir)

	var buf bytes.Buffer
	cfg := migrate.AgentConfig{
		RootDir:     dir,
		SourceTool:  "claude",
		TargetTools: []agentscan.Tool{agentscan.ToolClaude},
		OutputDir:   filepath.Join(dir, "template/sections"),
		TemplateDir: filepath.Join(dir, "template/pages"),
		DryRun:      true,
		ConfigFile:  filepath.Join(dir, "docsgen.yaml"),
	}

	if err := migrate.RunAgents(&buf, cfg); err != nil {
		t.Fatalf("RunAgents() error: %v", err)
	}

	if !strings.Contains(buf.String(), "Detected source tool: claude") {
		t.Errorf("expected explicit source tool, got:\n%s", buf.String())
	}
}

func TestRunAgents_CodexCombinedAGENTS(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupAgentTestFiles(t, dir)

	var buf bytes.Buffer
	cfg := migrate.AgentConfig{
		RootDir:     dir,
		SourceTool:  "claude",
		TargetTools: []agentscan.Tool{agentscan.ToolCodex},
		OutputDir:   filepath.Join(dir, "template/sections"),
		TemplateDir: filepath.Join(dir, "template/pages"),
		DryRun:      false,
		ConfigFile:  filepath.Join(dir, "docsgen.yaml"),
	}

	if err := migrate.RunAgents(&buf, cfg); err != nil {
		t.Fatalf("RunAgents() error: %v", err)
	}

	// Codex rules should be combined into a single AGENTS.tpl.md.
	agentsTpl := filepath.Join(dir, "template/pages/ai-agents/codex/AGENTS.tpl.md")
	data, err := os.ReadFile(agentsTpl)
	if err != nil {
		t.Fatalf("expected AGENTS.tpl.md to exist: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "# Agent Instructions") {
		t.Errorf("expected title in AGENTS template, got:\n%s", content)
	}
	if !strings.Contains(content, "@include:") {
		t.Errorf("expected @include directives, got:\n%s", content)
	}
}

// helpers

func setupAgentTestFiles(t *testing.T, dir string) {
	t.Helper()

	// Claude rules.
	rulesDir := filepath.Join(dir, ".claude", "rules")
	if err := os.MkdirAll(rulesDir, 0o750); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(rulesDir, "architecture.md"), "# Architecture\n\nLayered architecture pattern.\n")
	writeTestFile(t, filepath.Join(rulesDir, "testing.md"), "# Testing\n\nRun go test ./...\n")

	// Claude skill with frontmatter.
	skillDir := filepath.Join(dir, ".claude", "skills", "deploy")
	if err := os.MkdirAll(skillDir, 0o750); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(skillDir, "SKILL.md"), "---\nname: deploy\ndescription: Deploy to production\n---\n\n# Deploy\n\n")
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
