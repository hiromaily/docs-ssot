package agentscan_test

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
)

func TestScan_DetectsClaude(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupClaudeFiles(t, dir)

	result, err := agentscan.Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if !containsTool(result.DetectedTools, agentscan.ToolClaude) {
		t.Errorf("expected Claude in detected tools, got %v", result.DetectedTools)
	}

	rules := result.FilesByType(agentscan.ToolClaude, agentscan.FileTypeRule)
	if len(rules) != 2 {
		t.Errorf("expected 2 Claude rules, got %d", len(rules))
	}

	skills := result.FilesByType(agentscan.ToolClaude, agentscan.FileTypeSkill)
	if len(skills) != 1 {
		t.Errorf("expected 1 Claude skill, got %d", len(skills))
	}
}

func TestScan_DetectsSubagents(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupClaudeFiles(t, dir)

	// Add subagents.
	agentsDir := filepath.Join(dir, ".claude", "agents")
	mkdirAll(t, agentsDir)
	writeFile(t, filepath.Join(agentsDir, "critic.md"), "---\nname: critic\ndescription: Adversarial critic\n---\n\n# Critic\n")
	writeFile(t, filepath.Join(agentsDir, "debugger.md"), "---\nname: debugger\ndescription: Debugger\n---\n\n# Debugger\n")

	result, err := agentscan.Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	subagents := result.FilesByType(agentscan.ToolClaude, agentscan.FileTypeSubagent)
	if len(subagents) != 2 {
		t.Errorf("expected 2 Claude subagents, got %d", len(subagents))
	}

	// Total files should include subagents.
	claudeFiles := result.FilesForTool(agentscan.ToolClaude)
	if len(claudeFiles) != 5 { // 2 rules + 1 skill + 2 subagents
		t.Errorf("expected 5 Claude files, got %d", len(claudeFiles))
	}
}

func TestScan_DetectsCursor(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Create Cursor rules.
	rulesDir := filepath.Join(dir, ".cursor", "rules")
	mkdirAll(t, rulesDir)
	writeFile(t, filepath.Join(rulesDir, "architecture.mdc"), "---\ndescription: Arch\nalwaysApply: true\n---\n\n# Architecture\n")

	result, err := agentscan.Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if !containsTool(result.DetectedTools, agentscan.ToolCursor) {
		t.Errorf("expected Cursor in detected tools, got %v", result.DetectedTools)
	}

	rules := result.FilesByType(agentscan.ToolCursor, agentscan.FileTypeRule)
	if len(rules) != 1 {
		t.Errorf("expected 1 Cursor rule, got %d", len(rules))
	}
	if rules[0].Slug != "architecture" {
		t.Errorf("expected slug 'architecture', got %q", rules[0].Slug)
	}
}

func TestScan_DetectsCopilot(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Create Copilot instructions.
	instrDir := filepath.Join(dir, ".github", "instructions")
	mkdirAll(t, instrDir)
	writeFile(t, filepath.Join(instrDir, "go.instructions.md"), "---\napplyTo: \"**/*.go\"\n---\n\n# Go Rules\n")
	writeFile(t, filepath.Join(instrDir, "testing.instructions.md"), "---\napplyTo: \"**/*_test.go\"\n---\n\n# Testing\n")

	result, err := agentscan.Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if !containsTool(result.DetectedTools, agentscan.ToolCopilot) {
		t.Errorf("expected Copilot in detected tools, got %v", result.DetectedTools)
	}

	rules := result.FilesByType(agentscan.ToolCopilot, agentscan.FileTypeRule)
	if len(rules) != 2 {
		t.Errorf("expected 2 Copilot rules, got %d", len(rules))
	}
	if rules[0].Slug != "go" {
		t.Errorf("expected slug 'go', got %q", rules[0].Slug)
	}
}

func TestScan_DetectsCopilotViaInstructionsFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Create only copilot-instructions.md (no instructions/ dir).
	githubDir := filepath.Join(dir, ".github")
	mkdirAll(t, githubDir)
	writeFile(t, filepath.Join(githubDir, "copilot-instructions.md"), "# Instructions\n")

	result, err := agentscan.Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if !containsTool(result.DetectedTools, agentscan.ToolCopilot) {
		t.Errorf("expected Copilot detected via copilot-instructions.md, got %v", result.DetectedTools)
	}
}

func TestScan_SourceToolAutoDetect(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupClaudeFiles(t, dir) // 2 rules + 1 skill = 3 files

	// Create 1 Cursor rule.
	rulesDir := filepath.Join(dir, ".cursor", "rules")
	mkdirAll(t, rulesDir)
	writeFile(t, filepath.Join(rulesDir, "core.mdc"), "---\n---\n\n# Core\n")

	result, err := agentscan.Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if result.SourceTool != agentscan.ToolClaude {
		t.Errorf("expected source tool 'claude', got %q", result.SourceTool)
	}
}

func TestScan_EmptyRepo(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	result, err := agentscan.Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if len(result.DetectedTools) != 0 {
		t.Errorf("expected no detected tools, got %v", result.DetectedTools)
	}
}

func TestParseTool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    agentscan.Tool
		wantErr bool
	}{
		{name: "claude", input: "claude", want: agentscan.ToolClaude},
		{name: "cursor", input: "Cursor", want: agentscan.ToolCursor},
		{name: "copilot", input: "COPILOT", want: agentscan.ToolCopilot},
		{name: "codex", input: "codex", want: agentscan.ToolCodex},
		{name: "unknown", input: "vim", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := agentscan.ParseTool(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTool(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseTool(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllTools(t *testing.T) {
	t.Parallel()

	tools := agentscan.AllTools()
	want := []agentscan.Tool{agentscan.ToolClaude, agentscan.ToolCursor, agentscan.ToolCopilot, agentscan.ToolCodex}
	if !reflect.DeepEqual(tools, want) {
		t.Errorf("AllTools() = %v, want %v", tools, want)
	}
}

func TestFilesForTool(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	setupClaudeFiles(t, dir)

	// Also create Cursor files.
	rulesDir := filepath.Join(dir, ".cursor", "rules")
	mkdirAll(t, rulesDir)
	writeFile(t, filepath.Join(rulesDir, "core.mdc"), "# Core\n")

	result, err := agentscan.Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	claudeFiles := result.FilesForTool(agentscan.ToolClaude)
	cursorFiles := result.FilesForTool(agentscan.ToolCursor)

	if len(claudeFiles) != 3 {
		t.Errorf("expected 3 Claude files, got %d", len(claudeFiles))
	}
	if len(cursorFiles) != 1 {
		t.Errorf("expected 1 Cursor file, got %d", len(cursorFiles))
	}
}

// helpers

func setupClaudeFiles(t *testing.T, dir string) {
	t.Helper()

	rulesDir := filepath.Join(dir, ".claude", "rules")
	mkdirAll(t, rulesDir)
	writeFile(t, filepath.Join(rulesDir, "architecture.md"), "# Architecture\n\nLayered architecture.\n")
	writeFile(t, filepath.Join(rulesDir, "testing.md"), "# Testing\n\nRun go test.\n")

	skillDir := filepath.Join(dir, ".claude", "skills", "deploy")
	mkdirAll(t, skillDir)
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), "---\nname: deploy\ndescription: Deploy to production\n---\n\n# Deploy\n\nSteps...\n")
}

func containsTool(tools []agentscan.Tool, target agentscan.Tool) bool {
	return slices.Contains(tools, target)
}

func mkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o750); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
