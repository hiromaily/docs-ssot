package frontmatter_test

import (
	"testing"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
	"github.com/hiromaily/docs-ssot/internal/frontmatter"
)

func TestParse_WithFrontmatter(t *testing.T) {
	t.Parallel()

	content := "---\nname: deploy\ndescription: Deploy to production\n---\n\n# Deploy\n\nSteps...\n"
	p := frontmatter.Parse(content)

	if p.Raw != "name: deploy\ndescription: Deploy to production" {
		t.Errorf("unexpected Raw: %q", p.Raw)
	}
	if p.Fields["name"] != "deploy" {
		t.Errorf("expected name='deploy', got %q", p.Fields["name"])
	}
	if p.Fields["description"] != "Deploy to production" {
		t.Errorf("expected description='Deploy to production', got %q", p.Fields["description"])
	}
	if p.Body != "\n# Deploy\n\nSteps...\n" {
		t.Errorf("unexpected Body: %q", p.Body)
	}
}

func TestParse_WithoutFrontmatter(t *testing.T) {
	t.Parallel()

	content := "# Just Markdown\n\nNo frontmatter here.\n"
	p := frontmatter.Parse(content)

	if p.Raw != "" {
		t.Errorf("expected empty Raw, got %q", p.Raw)
	}
	if len(p.Fields) != 0 {
		t.Errorf("expected empty Fields, got %v", p.Fields)
	}
	if p.Body != content {
		t.Errorf("expected Body to equal input, got %q", p.Body)
	}
}

func TestParse_FrontmatterOnly(t *testing.T) {
	t.Parallel()

	content := "---\nname: test\n---"
	p := frontmatter.Parse(content)

	if p.Fields["name"] != "test" {
		t.Errorf("expected name='test', got %q", p.Fields["name"])
	}
	if p.Body != "" {
		t.Errorf("expected empty Body, got %q", p.Body)
	}
}

func TestStripContent_ShiftsH1ToH2(t *testing.T) {
	t.Parallel()

	content := "---\nname: foo\n---\n\n# Title\n\n## Subsection\n\nContent.\n"
	got := frontmatter.StripContent(content)

	expected := "## Title\n\n## Subsection\n\nContent."
	if got != expected {
		t.Errorf("StripContent() =\n%q\nwant:\n%q", got, expected)
	}
}

func TestStripContent_PreservesCodeFenceHeadings(t *testing.T) {
	t.Parallel()

	content := "# Title\n\n```\n# comment in code\n```\n"
	got := frontmatter.StripContent(content)

	if got != "## Title\n\n```\n# comment in code\n```" {
		t.Errorf("StripContent() should preserve code fence headings, got:\n%q", got)
	}
}

func TestGenerateRuleTemplate_Claude(t *testing.T) {
	t.Parallel()

	got := frontmatter.GenerateRuleTemplate(agentscan.ToolClaude, "architecture", "../../sections/ai/rules/architecture.md")
	want := "<!-- @include: ../../sections/ai/rules/architecture.md level=-1 -->\n"
	if got != want {
		t.Errorf("GenerateRuleTemplate(claude) =\n%q\nwant:\n%q", got, want)
	}
}

func TestGenerateRuleTemplate_Cursor(t *testing.T) {
	t.Parallel()

	got := frontmatter.GenerateRuleTemplate(agentscan.ToolCursor, "architecture", "../../sections/ai/rules/architecture.md")

	if got == "" {
		t.Fatal("empty result")
	}
	// Should contain MDC frontmatter.
	if !contains(got, "description: Architecture") {
		t.Errorf("expected description in Cursor template, got:\n%s", got)
	}
	if !contains(got, "alwaysApply: true") {
		t.Errorf("expected alwaysApply in Cursor template, got:\n%s", got)
	}
	if !contains(got, "@include:") {
		t.Errorf("expected @include in template, got:\n%s", got)
	}
}

func TestGenerateRuleTemplate_Copilot(t *testing.T) {
	t.Parallel()

	got := frontmatter.GenerateRuleTemplate(agentscan.ToolCopilot, "testing", "../../sections/ai/rules/testing.md")

	if !contains(got, "applyTo: \"**/*\"") {
		t.Errorf("expected applyTo in Copilot template, got:\n%s", got)
	}
}

func TestGenerateSkillTemplate(t *testing.T) {
	t.Parallel()

	fields := map[string]string{
		"name":        "deploy",
		"description": "Deploy to production",
		"model":       "opus",
	}

	t.Run("claude_preserves_extra_fields", func(t *testing.T) {
		t.Parallel()
		got := frontmatter.GenerateSkillTemplate(agentscan.ToolClaude, "deploy", "../sections/ai/skills/deploy.md", fields)
		if !contains(got, "model: opus") {
			t.Errorf("expected model field preserved for Claude, got:\n%s", got)
		}
	})

	t.Run("cursor_only_name_description", func(t *testing.T) {
		t.Parallel()
		got := frontmatter.GenerateSkillTemplate(agentscan.ToolCursor, "deploy", "../sections/ai/skills/deploy.md", fields)
		if contains(got, "model:") {
			t.Errorf("Cursor template should not include model field, got:\n%s", got)
		}
		if !contains(got, "name: deploy") {
			t.Errorf("expected name in Cursor template, got:\n%s", got)
		}
	})
}

func TestGenerateCodexAGENTSTemplate(t *testing.T) {
	t.Parallel()

	includes := []string{
		"../../sections/ai/rules/architecture.md",
		"../../sections/ai/rules/testing.md",
	}

	got := frontmatter.GenerateCodexAGENTSTemplate(includes)

	if !contains(got, "# Agent Instructions") {
		t.Errorf("expected title in AGENTS template, got:\n%s", got)
	}
	if !contains(got, "architecture.md") {
		t.Errorf("expected architecture include, got:\n%s", got)
	}
	if !contains(got, "testing.md") {
		t.Errorf("expected testing include, got:\n%s", got)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := range len(s) - len(substr) + 1 {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
