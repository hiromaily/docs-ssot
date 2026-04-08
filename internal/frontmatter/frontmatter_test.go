package frontmatter_test

import (
	"strings"
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

func TestParse_CRLF(t *testing.T) {
	t.Parallel()

	content := "---\r\nname: foo\r\ndescription: bar\r\n---\r\n\r\n# Title\r\n"
	p := frontmatter.Parse(content)

	if p.Fields["name"] != "foo" {
		t.Errorf("expected name='foo', got %q", p.Fields["name"])
	}
	if p.Fields["description"] != "bar" {
		t.Errorf("expected description='bar', got %q", p.Fields["description"])
	}
	if !strings.Contains(p.Body, "# Title") {
		t.Errorf("expected body to contain '# Title', got %q", p.Body)
	}
}

func TestParse_MultiValueField(t *testing.T) {
	t.Parallel()

	content := "---\nallowed-tools:\n  - Read\n  - Edit\n---\n\nBody\n"
	p := frontmatter.Parse(content)

	if _, ok := p.Fields["allowed-tools"]; !ok {
		t.Errorf("expected 'allowed-tools' key in Fields, got %v", p.Fields)
	}
	// With yaml.Unmarshal, multi-line list values are serialised as compact YAML.
	val := p.Fields["allowed-tools"]
	if !strings.Contains(val, "Read") || !strings.Contains(val, "Edit") {
		t.Errorf("expected multi-line value to contain Read and Edit, got %q", val)
	}
}

func TestParse_BoolAndIntFields(t *testing.T) {
	t.Parallel()

	content := "---\nalwaysApply: true\neffort: 3\n---\n\nBody\n"
	p := frontmatter.Parse(content)

	if p.Fields["alwaysApply"] != "true" {
		t.Errorf("expected alwaysApply='true', got %q", p.Fields["alwaysApply"])
	}
	if p.Fields["effort"] != "3" {
		t.Errorf("expected effort='3', got %q", p.Fields["effort"])
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

func TestShiftH1ToH2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple_h1", input: "# Foo", want: "## Foo"},
		{name: "h2_unchanged", input: "## Bar", want: "## Bar"},
		{name: "code_fence_preserved", input: "```\n# not heading\n```", want: "```\n# not heading\n```"},
		{name: "no_headings", input: "just text", want: "just text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := frontmatter.ShiftH1ToH2(tt.input)
			if got != tt.want {
				t.Errorf("ShiftH1ToH2(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
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
	if !strings.Contains(got, "description: Architecture") {
		t.Errorf("expected description in Cursor template, got:\n%s", got)
	}
	if !strings.Contains(got, "alwaysApply: true") {
		t.Errorf("expected alwaysApply in Cursor template, got:\n%s", got)
	}
	if !strings.Contains(got, "@include:") {
		t.Errorf("expected @include in template, got:\n%s", got)
	}
}

func TestGenerateRuleTemplate_Copilot(t *testing.T) {
	t.Parallel()

	got := frontmatter.GenerateRuleTemplate(agentscan.ToolCopilot, "testing", "../../sections/ai/rules/testing.md")

	if !strings.Contains(got, "applyTo: \"**/*\"") {
		t.Errorf("expected applyTo in Copilot template, got:\n%s", got)
	}
}

func TestGenerateRuleTemplate_WithGlobs(t *testing.T) {
	t.Parallel()

	opts := frontmatter.RuleTemplateOpts{Globs: "**/*.go"}

	t.Run("cursor_uses_globs", func(t *testing.T) {
		t.Parallel()
		got := frontmatter.GenerateRuleTemplate(agentscan.ToolCursor, "go", "../../sections/ai/rules/go.md", opts)
		if !strings.Contains(got, "globs: **/*.go") {
			t.Errorf("expected globs in Cursor template, got:\n%s", got)
		}
		if strings.Contains(got, "alwaysApply") {
			t.Errorf("expected no alwaysApply when globs set, got:\n%s", got)
		}
	})

	t.Run("copilot_uses_globs_in_applyTo", func(t *testing.T) {
		t.Parallel()
		got := frontmatter.GenerateRuleTemplate(agentscan.ToolCopilot, "go", "../../sections/ai/rules/go.md", opts)
		if !strings.Contains(got, "applyTo: \"**/*.go\"") {
			t.Errorf("expected globs in applyTo, got:\n%s", got)
		}
	})

	t.Run("claude_ignores_globs", func(t *testing.T) {
		t.Parallel()
		got := frontmatter.GenerateRuleTemplate(agentscan.ToolClaude, "go", "../../sections/ai/rules/go.md", opts)
		if strings.Contains(got, "globs") || strings.Contains(got, "applyTo") {
			t.Errorf("Claude should not have globs/applyTo, got:\n%s", got)
		}
	})
}

func TestGenerateSubagentTemplate(t *testing.T) {
	t.Parallel()

	fields := map[string]string{
		"name":            "critic",
		"description":     "Adversarial critic",
		"disallowedTools": "",
	}

	t.Run("claude_preserves_fields", func(t *testing.T) {
		t.Parallel()
		got := frontmatter.GenerateSubagentTemplate(agentscan.ToolClaude, "critic", "../sections/ai/subagents/critic.md", fields)
		if !strings.Contains(got, "name: critic") {
			t.Errorf("expected name, got:\n%s", got)
		}
		if !strings.Contains(got, "description: Adversarial critic") {
			t.Errorf("expected description, got:\n%s", got)
		}
		if !strings.Contains(got, "disallowedTools:") {
			t.Errorf("expected disallowedTools preserved for Claude, got:\n%s", got)
		}
		if !strings.Contains(got, "@include:") {
			t.Errorf("expected @include, got:\n%s", got)
		}
	})

	t.Run("copilot_only_name_description", func(t *testing.T) {
		t.Parallel()
		got := frontmatter.GenerateSubagentTemplate(agentscan.ToolCopilot, "critic", "../sections/ai/subagents/critic.md", fields)
		if strings.Contains(got, "disallowedTools") {
			t.Errorf("Copilot should not have disallowedTools, got:\n%s", got)
		}
		if !strings.Contains(got, "name: critic") {
			t.Errorf("expected name, got:\n%s", got)
		}
	})
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
		if !strings.Contains(got, "model: opus") {
			t.Errorf("expected model field preserved for Claude, got:\n%s", got)
		}
	})

	t.Run("cursor_only_name_description", func(t *testing.T) {
		t.Parallel()
		got := frontmatter.GenerateSkillTemplate(agentscan.ToolCursor, "deploy", "../sections/ai/skills/deploy.md", fields)
		if strings.Contains(got, "model:") {
			t.Errorf("Cursor template should not include model field, got:\n%s", got)
		}
		if !strings.Contains(got, "name: deploy") {
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

	if !strings.Contains(got, "# Agent Instructions") {
		t.Errorf("expected title in AGENTS template, got:\n%s", got)
	}
	if !strings.Contains(got, "architecture.md") {
		t.Errorf("expected architecture include, got:\n%s", got)
	}
	if !strings.Contains(got, "testing.md") {
		t.Errorf("expected testing include, got:\n%s", got)
	}
}
