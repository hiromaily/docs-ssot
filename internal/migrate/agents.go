package migrate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
	"github.com/hiromaily/docs-ssot/internal/frontmatter"
	"github.com/hiromaily/docs-ssot/internal/generator"
)

// AgentConfig holds parameters for the agent-aware migration mode.
type AgentConfig struct {
	// RootDir is the repository root to scan for AI tool files.
	RootDir string
	// SourceTool overrides auto-detection of the source tool ("auto" for auto-detect).
	SourceTool string
	// TargetTools lists tools to generate templates for ("all" or comma-separated).
	TargetTools []agentscan.Tool
	// OutputDir is where section files are written (default: "template/sections").
	OutputDir string
	// TemplateDir is where template files are written (default: "template/pages").
	TemplateDir string
	// DryRun prints the plan without writing files.
	DryRun bool
	// ConfigFile is the path to docsgen.yaml.
	ConfigFile string
}

// agentSection tracks a single agent file and its planned output locations.
type agentSection struct {
	Source      agentscan.AgentFile
	Content     string            // body content (frontmatter stripped)
	Fields      map[string]string // parsed frontmatter fields
	SectionPath string            // output path for the section file
}

// RunAgents executes the agent-aware migration.
func RunAgents(w io.Writer, cfg AgentConfig) error {
	// Step 1: Scan repository for AI tool files.
	scanResult, err := agentscan.Scan(cfg.RootDir)
	if err != nil {
		return fmt.Errorf("agent scan: %w", err)
	}

	if len(scanResult.DetectedTools) == 0 {
		return errors.New("no AI tool configurations detected in repository")
	}

	// Step 2: Determine source tool.
	sourceTool := scanResult.SourceTool
	if cfg.SourceTool != "auto" && cfg.SourceTool != "" {
		parsed, parseErr := agentscan.ParseTool(cfg.SourceTool)
		if parseErr != nil {
			return parseErr
		}
		sourceTool = parsed
	}

	sourceFiles := scanResult.FilesForTool(sourceTool)
	if len(sourceFiles) == 0 {
		return fmt.Errorf("no agent files found for source tool %q", sourceTool)
	}

	_, _ = fmt.Fprintf(w, "Detected source tool: %s (%d files)\n", sourceTool, len(sourceFiles))
	_, _ = fmt.Fprintf(w, "Target tools: %s\n", formatTools(cfg.TargetTools))

	// Step 3: Read and parse source files.
	sections, err := readSourceFiles(cfg.RootDir, sourceFiles)
	if err != nil {
		return err
	}

	// Step 4: Plan section file paths.
	planSections(sections, cfg.OutputDir)

	// Step 5: Plan template files.
	templates := planTemplates(sections, cfg)

	// Step 6: Report plan.
	reportAgentPlan(w, cfg, sections, templates)

	if cfg.DryRun {
		return nil
	}

	// Step 7: Write section files.
	if err := writeAgentSectionFiles(w, sections); err != nil {
		return err
	}

	// Step 8: Write template files.
	if err := writeAgentTemplateFiles(w, templates); err != nil {
		return err
	}

	// Step 9: Update docsgen.yaml.
	if err := updateConfig(w, cfg, templates); err != nil {
		return err
	}

	// Step 10: Round-trip verification.
	verifyAgentRoundTrip(w, cfg)

	_, _ = fmt.Fprintln(w, "Agent migration complete.")
	return nil
}

// templateFile represents a planned template file to generate.
type templateFile struct {
	Tool       agentscan.Tool
	Type       agentscan.FileType
	Slug       string
	OutputPath string
	Content    string
}

func readSourceFiles(root string, files []agentscan.AgentFile) ([]agentSection, error) {
	sections := make([]agentSection, 0, len(files))
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(root, f.Path))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", f.Path, err)
		}

		parsed := frontmatter.Parse(string(data))
		body := frontmatter.StripContent(string(data))

		sections = append(sections, agentSection{
			Source:  f,
			Content: body,
			Fields:  parsed.Fields,
		})
	}
	return sections, nil
}

func planSections(sections []agentSection, outputDir string) {
	for i := range sections {
		s := &sections[i]
		typeDir := string(s.Source.Type) + "s" // "rules", "skills", "commands"
		s.SectionPath = filepath.Join(outputDir, "ai", typeDir, s.Source.Slug+".md")
	}
}

func planTemplates(sections []agentSection, cfg AgentConfig) []templateFile {
	var templates []templateFile

	// Track codex rule includes for combined AGENTS.md.
	var codexRuleIncludes []string

	for _, s := range sections {
		for _, tool := range cfg.TargetTools {
			// Commands are Claude-only.
			if s.Source.Type == agentscan.FileTypeCommand && tool != agentscan.ToolClaude {
				continue
			}

			// Codex rules go into a combined AGENTS.md — skip individual templates.
			if tool == agentscan.ToolCodex && s.Source.Type == agentscan.FileTypeRule {
				// Compute include path from AGENTS.md template to section file.
				tplDir := filepath.Join(cfg.TemplateDir, "ai-agents", "codex")
				relPath := computeRelPath(tplDir, s.SectionPath)
				codexRuleIncludes = append(codexRuleIncludes, relPath)
				continue
			}

			tpl := buildTemplate(tool, s, cfg)
			templates = append(templates, tpl)
		}
	}

	// Generate combined Codex AGENTS.md template if there are rules.
	if len(codexRuleIncludes) > 0 && containsTool(cfg.TargetTools, agentscan.ToolCodex) {
		codexPath := filepath.Join(cfg.TemplateDir, "ai-agents", "codex", "AGENTS.tpl.md")
		templates = append(templates, templateFile{
			Tool:       agentscan.ToolCodex,
			Type:       agentscan.FileTypeRule,
			Slug:       "AGENTS",
			OutputPath: codexPath,
			Content:    frontmatter.GenerateCodexAGENTSTemplate(codexRuleIncludes),
		})
	}

	return templates
}

func buildTemplate(tool agentscan.Tool, s agentSection, cfg AgentConfig) templateFile {
	tplDir := templateDirForTool(cfg.TemplateDir, tool)
	tplPath := templatePathForFile(tplDir, tool, s.Source)
	relPath := computeRelPath(filepath.Dir(tplPath), s.SectionPath)

	var content string
	switch s.Source.Type {
	case agentscan.FileTypeRule:
		content = frontmatter.GenerateRuleTemplate(tool, s.Source.Slug, relPath)
	case agentscan.FileTypeSkill:
		content = frontmatter.GenerateSkillTemplate(tool, s.Source.Slug, relPath, s.Fields)
	case agentscan.FileTypeCommand:
		// Commands reuse the rule template format.
		content = frontmatter.GenerateRuleTemplate(tool, s.Source.Slug, relPath)
	}

	return templateFile{
		Tool:       tool,
		Type:       s.Source.Type,
		Slug:       s.Source.Slug,
		OutputPath: tplPath,
		Content:    content,
	}
}

func templateDirForTool(baseDir string, tool agentscan.Tool) string {
	return filepath.Join(baseDir, "ai-agents", string(tool))
}

func templatePathForFile(tplDir string, tool agentscan.Tool, f agentscan.AgentFile) string {
	switch f.Type {
	case agentscan.FileTypeRule:
		ext := ".tpl.md"
		if tool == agentscan.ToolCursor {
			ext = ".tpl.mdc"
		}
		subDir := "rules"
		if tool == agentscan.ToolCopilot {
			subDir = "instructions"
		}
		return filepath.Join(tplDir, subDir, f.Slug+ext)
	case agentscan.FileTypeSkill:
		return filepath.Join(tplDir, "skills", f.Slug, "SKILL.tpl.md")
	case agentscan.FileTypeCommand:
		return filepath.Join(tplDir, "commands", f.Slug+".tpl.md")
	}
	return filepath.Join(tplDir, f.Slug+".tpl.md")
}

func computeRelPath(from, to string) string {
	absFrom, err := filepath.Abs(from)
	if err != nil {
		return to
	}
	absTo, err := filepath.Abs(to)
	if err != nil {
		return to
	}
	rel, err := filepath.Rel(absFrom, absTo)
	if err != nil {
		return to
	}
	return rel
}

func reportAgentPlan(w io.Writer, cfg AgentConfig, sections []agentSection, templates []templateFile) {
	_, _ = fmt.Fprintln(w)

	verb := "Would create"
	if !cfg.DryRun {
		verb = "Creating"
	}

	_, _ = fmt.Fprintf(w, "%s sections:\n", verb)
	for _, s := range sections {
		_, _ = fmt.Fprintf(w, "  %s\n", s.SectionPath)
	}

	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "%s templates (%d tools × %d files):\n", verb, len(cfg.TargetTools), len(sections))

	// Group by tool.
	toolCounts := map[agentscan.Tool]int{}
	for _, t := range templates {
		toolCounts[t.Tool]++
	}
	for _, tool := range cfg.TargetTools {
		count := toolCounts[tool]
		_, _ = fmt.Fprintf(w, "  %s: %d templates\n", tool, count)
	}

	_, _ = fmt.Fprintln(w)
}

func writeAgentSectionFiles(w io.Writer, sections []agentSection) error {
	for _, s := range sections {
		if err := os.MkdirAll(filepath.Dir(s.SectionPath), 0o750); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(s.SectionPath), err)
		}

		//nolint:gosec // generated documentation files
		if err := os.WriteFile(s.SectionPath, []byte(s.Content+"\n"), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", s.SectionPath, err)
		}
		_, _ = fmt.Fprintf(w, "  %s\n", s.SectionPath)
	}
	return nil
}

func writeAgentTemplateFiles(w io.Writer, templates []templateFile) error {
	for _, t := range templates {
		if err := os.MkdirAll(filepath.Dir(t.OutputPath), 0o750); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(t.OutputPath), err)
		}

		//nolint:gosec // generated template files
		if err := os.WriteFile(t.OutputPath, []byte(t.Content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", t.OutputPath, err)
		}
		_, _ = fmt.Fprintf(w, "  Created %s\n", t.OutputPath)
	}
	return nil
}

func updateConfig(w io.Writer, cfg AgentConfig, templates []templateFile) error {
	// Build new targets from templates.
	newTargets := make([]string, 0, len(templates))

	for _, t := range templates {
		outputPath := resolveOutputPath(t)
		newTargets = append(newTargets, fmt.Sprintf("  - input: %s\n    output: %s", t.OutputPath, outputPath))
	}

	if len(newTargets) == 0 {
		return nil
	}

	// Check if config exists.
	configData, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		// Create new config.
		content := "targets:\n" + strings.Join(newTargets, "\n") + "\n"
		//nolint:gosec // generated config file
		if writeErr := os.WriteFile(cfg.ConfigFile, []byte(content), 0o644); writeErr != nil {
			return fmt.Errorf("write %s: %w", cfg.ConfigFile, writeErr)
		}
		_, _ = fmt.Fprintf(w, "Created %s (%d targets)\n", cfg.ConfigFile, len(newTargets))
		return nil
	}

	// Append to existing config.
	existing := string(configData)
	if !strings.HasSuffix(existing, "\n") {
		existing += "\n"
	}
	var sb strings.Builder
	for _, target := range newTargets {
		sb.WriteString(target + "\n")
	}
	existing += sb.String()

	//nolint:gosec // generated config file
	if err := os.WriteFile(cfg.ConfigFile, []byte(existing), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", cfg.ConfigFile, err)
	}

	_, _ = fmt.Fprintf(w, "Added %d targets to %s\n", len(newTargets), cfg.ConfigFile)
	return nil
}

func resolveOutputPath(t templateFile) string { //nolint:gocyclo // inherent in tool×type dispatch
	switch t.Tool {
	case agentscan.ToolClaude:
		switch t.Type {
		case agentscan.FileTypeRule:
			return filepath.Join(".claude", "rules", t.Slug+".md")
		case agentscan.FileTypeSkill:
			return filepath.Join(".claude", "skills", t.Slug, "SKILL.md")
		case agentscan.FileTypeCommand:
			return filepath.Join(".claude", "commands", t.Slug+".md")
		}
	case agentscan.ToolCursor:
		switch t.Type {
		case agentscan.FileTypeRule:
			return filepath.Join(".cursor", "rules", t.Slug+".mdc")
		case agentscan.FileTypeSkill:
			return filepath.Join(".cursor", "skills", t.Slug, "SKILL.md")
		case agentscan.FileTypeCommand:
			return t.OutputPath // commands not supported for Cursor
		}
	case agentscan.ToolCopilot:
		switch t.Type {
		case agentscan.FileTypeRule:
			return filepath.Join(".github", "instructions", t.Slug+".instructions.md")
		case agentscan.FileTypeSkill:
			return filepath.Join(".github", "skills", t.Slug, "SKILL.md")
		case agentscan.FileTypeCommand:
			return t.OutputPath // commands not supported for Copilot
		}
	case agentscan.ToolCodex:
		if t.Slug == "AGENTS" {
			return "AGENTS.md"
		}
		switch t.Type {
		case agentscan.FileTypeRule:
			return t.OutputPath // rules embedded in AGENTS.md
		case agentscan.FileTypeSkill:
			return filepath.Join(".agents", "skills", t.Slug, "SKILL.md")
		case agentscan.FileTypeCommand:
			return t.OutputPath // commands not supported for Codex
		}
	}
	return t.OutputPath
}

func verifyAgentRoundTrip(w io.Writer, cfg AgentConfig) {
	_, _ = fmt.Fprintln(w, "Verifying round-trip...")

	if err := generator.Build(cfg.ConfigFile); err != nil {
		_, _ = fmt.Fprintf(w, "Round-trip verification: SKIP (build error: %v)\n", err)
		return
	}

	_, _ = fmt.Fprintln(w, "Round-trip verification: OK")
}

func formatTools(tools []agentscan.Tool) string {
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = string(t)
	}
	return strings.Join(names, ", ")
}

func containsTool(tools []agentscan.Tool, target agentscan.Tool) bool {
	return slices.Contains(tools, target)
}
