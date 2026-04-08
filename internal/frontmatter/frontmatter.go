// Package frontmatter handles parsing, stripping, and generating YAML frontmatter
// for AI tool configuration files across different tools (Claude, Cursor, Copilot, Codex).
package frontmatter

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

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
// Handles both LF and CRLF line endings.
func Parse(content string) Parsed {
	// Normalize CRLF to LF for consistent parsing.
	content = normalizeCRLF(content)

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

// normalizeCRLF replaces all \r\n with \n.
func normalizeCRLF(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// StripContent returns the body content with frontmatter removed and
// H1 headings shifted to H2 (section file convention).
func StripContent(content string) string {
	return ShiftH1ToH2(strings.TrimSpace(Parse(content).Body))
}

// ShiftH1ToH2 shifts all H1 headings (# Foo) to H2 (## Foo) for section file convention.
// Headings inside code fences are left unchanged.
func ShiftH1ToH2(content string) string {
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

// RuleTemplateOpts configures rule template generation.
type RuleTemplateOpts struct {
	// Globs is an optional inferred glob pattern for path-gated rules.
	// When set, Cursor uses globs instead of alwaysApply, and Copilot uses it for applyTo.
	Globs string
}

// GenerateRuleTemplate generates a tool-specific rule template with appropriate
// frontmatter and an @include directive for the given section path.
func GenerateRuleTemplate(tool agentscan.Tool, slug, includePath string, opts ...RuleTemplateOpts) string {
	var o RuleTemplateOpts
	if len(opts) > 0 {
		o = opts[0]
	}

	var b strings.Builder

	switch tool {
	case agentscan.ToolClaude:
		// Claude rules have no frontmatter.
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")

	case agentscan.ToolCursor:
		b.WriteString("---\n")
		b.WriteString("description: " + slugToDescription(slug) + "\n")
		if o.Globs != "" {
			b.WriteString("globs: " + o.Globs + "\n")
		} else {
			b.WriteString("alwaysApply: true\n")
		}
		b.WriteString("---\n")
		b.WriteString("\n")
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")

	case agentscan.ToolCopilot:
		b.WriteString("---\n")
		if o.Globs != "" {
			b.WriteString("applyTo: \"" + o.Globs + "\"\n")
		} else {
			b.WriteString("applyTo: \"**/*\"\n")
		}
		b.WriteString("---\n")
		b.WriteString("\n")
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")

	case agentscan.ToolCodex:
		// Codex rules are embedded in AGENTS.md via @include.
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")

	default:
		// Fallback: plain include for unknown tools.
		b.WriteString("<!-- @include: " + includePath + " level=-1 -->\n")
	}

	return b.String()
}

// GenerateSubagentTemplate generates a tool-specific subagent template.
// Note: multi-line YAML fields (e.g., disallowedTools with list values) are preserved
// as single-line key: value pairs. The simple parseFields parser only captures the
// first line of list values, so the value may be empty for multi-line fields.
func GenerateSubagentTemplate(tool agentscan.Tool, slug, includePath string, sourceFields map[string]string) string {
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
		// Preserve all original frontmatter for Claude agents.
		b.WriteString("---\n")
		b.WriteString("name: " + name + "\n")
		b.WriteString("description: " + description + "\n")
		for _, key := range []string{"disallowedTools", "allowedTools", "model", "effort"} {
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

// parseFields parses YAML frontmatter into string key-value pairs.
// Uses yaml.Unmarshal to handle multi-line values (lists, nested maps).
// Complex values are serialised back to inline YAML strings.
func parseFields(raw string) map[string]string {
	// First try proper YAML parsing.
	var parsed map[string]any
	if err := yaml.Unmarshal([]byte(raw), &parsed); err == nil && len(parsed) > 0 {
		return flattenYAML(parsed)
	}

	// Fallback to simple line-by-line parsing for non-standard YAML.
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

// flattenYAML converts a parsed YAML map to string key-value pairs.
// Scalar values are converted directly. Lists and maps are serialised
// to compact YAML strings for preservation in generated templates.
func flattenYAML(m map[string]any) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case string:
			result[k] = val
		case bool:
			if val {
				result[k] = "true"
			} else {
				result[k] = "false"
			}
		case int:
			result[k] = strconv.Itoa(val)
		case float64:
			result[k] = fmt.Sprintf("%g", val)
		default:
			// Lists, maps — serialise to compact YAML.
			data, err := yaml.Marshal(v)
			if err != nil {
				result[k] = fmt.Sprintf("%v", v)
				continue
			}
			result[k] = strings.TrimSpace(string(data))
		}
	}
	return result
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
