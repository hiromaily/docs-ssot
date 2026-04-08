// Package frontmatter handles parsing, stripping, and generating YAML frontmatter
// for AI tool configuration files across different tools (Claude, Cursor, Copilot, Codex).
package frontmatter

import (
	"strings"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
)

// Parsed holds the parsed frontmatter and body of a Markdown file.
type Parsed struct {
	// Raw is the raw YAML frontmatter string (without delimiters).
	Raw string
	// Body is the Markdown content after the frontmatter.
	Body string
	// Fields holds parsed key-value pairs from the frontmatter.
	Fields map[string]string
}

// Parse splits a Markdown file into frontmatter and body.
// If no frontmatter is present, Raw and Fields are empty and Body contains the full content.
func Parse(content string) Parsed {
	if !strings.HasPrefix(content, "---\n") {
		return Parsed{Body: content}
	}

	// Find closing delimiter.
	rest := content[4:] // skip "---\n"
	raw, body, found := strings.Cut(rest, "\n---\n")
	if !found {
		// Check for "---\n" at the very end.
		if strings.HasSuffix(rest, "\n---") {
			raw = rest[:len(rest)-4]
			return Parsed{
				Raw:    raw,
				Body:   "",
				Fields: parseFields(raw),
			}
		}
		return Parsed{Body: content}
	}

	return Parsed{
		Raw:    raw,
		Body:   body,
		Fields: parseFields(raw),
	}
}

// StripContent returns the body content with frontmatter removed and
// H1 headings shifted to H2 (section file convention).
func StripContent(content string) string {
	p := Parse(content)
	return shiftH1ToH2(strings.TrimSpace(p.Body))
}

// GenerateRuleTemplate generates a tool-specific rule template with appropriate
// frontmatter and an @include directive for the given section path.
func GenerateRuleTemplate(tool agentscan.Tool, slug, includePath string) string {
	var b strings.Builder

	switch tool {
	case agentscan.ToolClaude:
		// Claude rules have no frontmatter.
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")

	case agentscan.ToolCursor:
		b.WriteString("---\n")
		b.WriteString("description: " + slugToDescription(slug) + "\n")
		b.WriteString("alwaysApply: true\n")
		b.WriteString("---\n")
		b.WriteString("\n")
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")

	case agentscan.ToolCopilot:
		b.WriteString("---\n")
		b.WriteString("applyTo: \"**/*\"\n")
		b.WriteString("---\n")
		b.WriteString("\n")
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")

	case agentscan.ToolCodex:
		// Codex rules are embedded in AGENTS.md via @include.
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")
	}

	return b.String()
}

// GenerateSkillTemplate generates a tool-specific skill template.
func GenerateSkillTemplate(tool agentscan.Tool, slug, includePath string, sourceFields map[string]string) string {
	var b strings.Builder

	name := slug
	description := slugToDescription(slug)
	if v, ok := sourceFields["name"]; ok {
		name = v
	}
	if v, ok := sourceFields["description"]; ok {
		description = v
	}

	switch tool {
	case agentscan.ToolClaude:
		// Preserve all original frontmatter fields for Claude.
		b.WriteString("---\n")
		b.WriteString("name: " + name + "\n")
		b.WriteString("description: " + description + "\n")
		// Preserve Claude-specific fields.
		for _, key := range []string{"argument-hint", "disable-model-invocation", "user-invocable", "allowed-tools", "model", "effort", "context"} {
			if v, ok := sourceFields[key]; ok {
				b.WriteString(key + ": " + v + "\n")
			}
		}
		b.WriteString("---\n")
		b.WriteString("\n")
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")

	default:
		// Other tools: name + description only.
		b.WriteString("---\n")
		b.WriteString("name: " + name + "\n")
		b.WriteString("description: " + description + "\n")
		b.WriteString("---\n")
		b.WriteString("\n")
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")
	}

	return b.String()
}

// GenerateCodexAGENTSTemplate generates a combined AGENTS.md template
// that includes all rules via @include directives.
func GenerateCodexAGENTSTemplate(includes []string) string {
	var b strings.Builder
	b.WriteString("# Agent Instructions\n\n")
	for _, inc := range includes {
		b.WriteString("<!-- @include: " + inc + " -->\n\n")
	}
	return b.String()
}

// parseFields parses simple key: value pairs from YAML frontmatter.
// This is intentionally simple — it handles single-line values only.
func parseFields(raw string) map[string]string {
	fields := make(map[string]string)
	for line := range strings.SplitSeq(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		fields[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return fields
}

// shiftH1ToH2 shifts all H1 headings (# Foo) to H2 (## Foo) for section file convention.
// Headings inside code fences are left unchanged.
func shiftH1ToH2(content string) string {
	lines := strings.Split(content, "\n")
	inFence := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "## ") {
			lines[i] = "#" + line
		}
	}
	return strings.Join(lines, "\n")
}

// slugToDescription converts a slug like "architecture" to "Architecture" for use as a description.
func slugToDescription(slug string) string {
	if slug == "" {
		return ""
	}
	parts := strings.Split(slug, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
