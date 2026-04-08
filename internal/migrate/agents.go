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
	"github.com/hiromaily/docs-ssot/internal/config"
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
	// ConvertCommands converts legacy commands to skills during migration.
	ConvertCommands bool
	// InferGlobs attempts to infer path-gated rules from slug names.
	InferGlobs bool
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
	planSections(sections, cfg)

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
	verifyAgentRoundTrip(w, cfg, sections, templates)

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
		body := frontmatter.ShiftH1ToH2(strings.TrimSpace(parsed.Body))

		sections = append(sections, agentSection{
			Source:  f,
			Content: body,
			Fields:  parsed.Fields,
		})
	}
	return sections, nil
}

func planSections(sections []agentSection, cfg AgentConfig) {
	for i := range sections {
		s := &sections[i]
		effectiveType := s.Source.Type
		// Convert commands to skills if requested.
		if cfg.ConvertCommands && effectiveType == agentscan.FileTypeCommand {
			effectiveType = agentscan.FileTypeSkill
		}
		typeDir := string(effectiveType) + "s" // "rules", "skills", "commands", "subagents"
		s.SectionPath = filepath.Join(cfg.OutputDir, "ai", typeDir, s.Source.Slug+".md")
	}
}

func planTemplates(sections []agentSection, cfg AgentConfig) []templateFile {
	var templates []templateFile

	// Track codex rule includes for combined AGENTS.md.
	var codexRuleIncludes []string

	for _, s := range sections {
		effectiveType := s.Source.Type
		if cfg.ConvertCommands && effectiveType == agentscan.FileTypeCommand {
			effectiveType = agentscan.FileTypeSkill
		}

		for _, tool := range cfg.TargetTools {
			// Original commands (not converted) are Claude-only.
			if effectiveType == agentscan.FileTypeCommand && tool != agentscan.ToolClaude {
				continue
			}

			// Codex rules go into a combined AGENTS.md — skip individual templates.
			if tool == agentscan.ToolCodex && effectiveType == agentscan.FileTypeRule {
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
	if len(codexRuleIncludes) > 0 && slices.Contains(cfg.TargetTools, agentscan.ToolCodex) {
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
	effectiveType := s.Source.Type
	if cfg.ConvertCommands && effectiveType == agentscan.FileTypeCommand {
		effectiveType = agentscan.FileTypeSkill
	}

	// Use effective type for path calculation.
	effectiveSource := s.Source
	effectiveSource.Type = effectiveType

	tplDir := templateDirForTool(cfg.TemplateDir, tool)
	tplPath := templatePathForFile(tplDir, tool, effectiveSource)
	relPath := computeRelPath(filepath.Dir(tplPath), s.SectionPath)

	var content string
	switch effectiveType {
	case agentscan.FileTypeRule:
		var opts frontmatter.RuleTemplateOpts
		if cfg.InferGlobs {
			if globs, ok := agentscan.InferGlobs(s.Source.Slug); ok {
				opts.Globs = globs
			}
		}
		content = frontmatter.GenerateRuleTemplate(tool, s.Source.Slug, relPath, opts)
	case agentscan.FileTypeSkill:
		content = frontmatter.GenerateSkillTemplate(tool, s.Source.Slug, relPath, s.Fields)
	case agentscan.FileTypeCommand:
		content = frontmatter.GenerateRuleTemplate(tool, s.Source.Slug, relPath)
	case agentscan.FileTypeSubagent:
		content = frontmatter.GenerateSubagentTemplate(tool, s.Source.Slug, relPath, s.Fields)
	}

	return templateFile{
		Tool:       tool,
		Type:       effectiveType,
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
	case agentscan.FileTypeSubagent:
		return filepath.Join(tplDir, "agents", f.Slug+".tpl.md")
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
	// Use forward slashes for cross-platform include paths in Markdown.
	return filepath.ToSlash(rel)
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
	newTargets := make([]config.Target, 0, len(templates))
	for _, t := range templates {
		newTargets = append(newTargets, config.Target{
			Input:  filepath.ToSlash(t.OutputPath),
			Output: filepath.ToSlash(resolveOutputPath(t)),
		})
	}

	if len(newTargets) == 0 {
		return nil
	}

	// Load existing config or create a new one.
	existing, loadErr := config.Load(cfg.ConfigFile)
	if loadErr != nil {
		if !os.IsNotExist(loadErr) {
			return fmt.Errorf("load config: %w", loadErr)
		}
		existing = &config.Config{}
	}

	// Deduplicate: skip targets that already exist in the config.
	existingSet := make(map[string]bool, len(existing.Targets))
	for _, t := range existing.Targets {
		existingSet[t.Input] = true
	}
	for _, t := range newTargets {
		if !existingSet[t.Input] {
			existing.Targets = append(existing.Targets, t)
		}
	}

	if err := config.Save(cfg.ConfigFile, existing); err != nil {
		return fmt.Errorf("write %s: %w", cfg.ConfigFile, err)
	}

	_, _ = fmt.Fprintf(w, "Updated %s (%d new targets)\n", cfg.ConfigFile, len(newTargets))
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
		case agentscan.FileTypeSubagent:
			return filepath.Join(".claude", "agents", t.Slug+".md")
		}
	case agentscan.ToolCursor:
		switch t.Type {
		case agentscan.FileTypeRule:
			return filepath.Join(".cursor", "rules", t.Slug+".mdc")
		case agentscan.FileTypeSkill:
			return filepath.Join(".cursor", "skills", t.Slug, "SKILL.md")
		case agentscan.FileTypeCommand:
			return t.OutputPath // commands not supported for Cursor
		case agentscan.FileTypeSubagent:
			return filepath.Join(".cursor", "agents", t.Slug+".md")
		}
	case agentscan.ToolCopilot:
		switch t.Type {
		case agentscan.FileTypeRule:
			return filepath.Join(".github", "instructions", t.Slug+".instructions.md")
		case agentscan.FileTypeSkill:
			return filepath.Join(".github", "skills", t.Slug, "SKILL.md")
		case agentscan.FileTypeCommand:
			return t.OutputPath // commands not supported for Copilot
		case agentscan.FileTypeSubagent:
			return filepath.Join(".github", "agents", t.Slug+".agent.md")
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
		case agentscan.FileTypeSubagent:
			return filepath.Join(".codex", "agents", t.Slug+".md")
		}
	}
	return t.OutputPath
}

// sourceOriginal pairs a source path with its original content and the generated output path
// for deterministic round-trip verification.
type sourceOriginal struct {
	sourcePath string
	outputPath string
	data       []byte
}

func verifyAgentRoundTrip(w io.Writer, cfg AgentConfig, sections []agentSection, templates []templateFile) {
	_, _ = fmt.Fprintln(w, "Verifying round-trip...")

	// Build a map from slug → first resolved output path so we can compare generated files.
	slugToOutput := make(map[string]string, len(templates))
	for _, t := range templates {
		if _, exists := slugToOutput[t.Slug]; !exists {
			slugToOutput[t.Slug] = resolveOutputPath(t)
		}
	}

	// Read original source files before building so we can compare after build.
	// Use a slice (not a map) for deterministic iteration order.
	originals := make([]sourceOriginal, 0, len(sections))
	for _, s := range sections {
		outputPath, ok := slugToOutput[s.Source.Slug]
		if !ok {
			// No standalone output for this section (e.g., Codex rules are embedded in AGENTS.md).
			continue
		}
		srcPath := filepath.Join(cfg.RootDir, s.Source.Path)
		data, err := os.ReadFile(srcPath)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Round-trip verification: SKIP (cannot read %s: %v)\n", s.Source.Path, err)
			return
		}
		originals = append(originals, sourceOriginal{
			sourcePath: s.Source.Path,
			outputPath: outputPath,
			data:       data,
		})
	}

	if err := generator.Build(cfg.ConfigFile); err != nil {
		_, _ = fmt.Fprintf(w, "Round-trip verification: SKIP (build error: %v)\n", err)
		return
	}

	// Compare generated output files against originals.
	allMatch := true
	for _, orig := range originals {
		generated, err := os.ReadFile(orig.outputPath)
		if err != nil {
			_, _ = fmt.Fprintf(w, "  WARNING: cannot read generated %s: %v\n", orig.outputPath, err)
			allMatch = false
			continue
		}

		// Strip frontmatter from both sides for a fair comparison: the original has tool
		// frontmatter, the generated output may have tool-specific frontmatter added by the template.
		originalBody := frontmatter.StripContent(string(orig.data))
		generatedBody := frontmatter.StripContent(string(generated))
		if originalBody != generatedBody {
			_, _ = fmt.Fprintf(w, "  WARNING: %s differs after round-trip\n", orig.sourcePath)
			allMatch = false
		}
	}

	if allMatch {
		_, _ = fmt.Fprintln(w, "Round-trip verification: OK")
	} else {
		_, _ = fmt.Fprintln(w, "Round-trip verification: WARN (some files differ)")
	}
}

func formatTools(tools []agentscan.Tool) string {
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = string(t)
	}
	return strings.Join(names, ", ")
}
