// Package agentscan detects AI tool configuration files in a repository and
// collects agent files (rules, skills, commands) for migration.
package agentscan

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Tool represents a supported AI coding tool.
type Tool string

const (
	ToolClaude  Tool = "claude"
	ToolCursor  Tool = "cursor"
	ToolCopilot Tool = "copilot"
	ToolCodex   Tool = "codex"
)

// AllTools returns all supported tools in canonical order.
func AllTools() []Tool {
	return []Tool{ToolClaude, ToolCursor, ToolCopilot, ToolCodex}
}

// ParseTool parses a tool name string into a Tool value.
func ParseTool(s string) (Tool, error) {
	switch strings.ToLower(s) {
	case "claude":
		return ToolClaude, nil
	case "cursor":
		return ToolCursor, nil
	case "copilot":
		return ToolCopilot, nil
	case "codex":
		return ToolCodex, nil
	default:
		return "", fmt.Errorf("unknown tool: %q (supported: claude, cursor, copilot, codex)", s)
	}
}

// FileType categorises an agent configuration file.
type FileType string

const (
	FileTypeRule     FileType = "rule"
	FileTypeSkill    FileType = "skill"
	FileTypeCommand  FileType = "command"
	FileTypeSubagent FileType = "subagent"
)

// AgentFile represents a single agent configuration file discovered in the repository.
type AgentFile struct {
	// Tool is the AI tool this file belongs to.
	Tool Tool
	// Type categorises the file (rule, skill, command).
	Type FileType
	// Path is the relative path from the repository root.
	Path string
	// Slug is the short identifier derived from the filename or directory name.
	Slug string
}

// ScanResult holds the detected AI tools and their agent files.
type ScanResult struct {
	// DetectedTools lists tools that have configuration directories present.
	DetectedTools []Tool
	// Files contains all discovered agent files, grouped by tool.
	Files []AgentFile
	// SourceTool is the tool selected as the source of truth (most files).
	SourceTool Tool
}

// Scan detects AI tool configuration in the given root directory and collects agent files.
func Scan(root string) (*ScanResult, error) {
	result := &ScanResult{}

	toolFiles := map[Tool]int{}

	// Detect Claude.
	claudeDir := filepath.Join(root, ".claude")
	if isDir(claudeDir) {
		result.DetectedTools = append(result.DetectedTools, ToolClaude)
		files, err := collectClaude(root)
		if err != nil {
			return nil, fmt.Errorf("scan claude: %w", err)
		}
		result.Files = append(result.Files, files...)
		toolFiles[ToolClaude] = len(files)
	}

	// Detect Cursor.
	cursorDir := filepath.Join(root, ".cursor")
	if isDir(cursorDir) {
		result.DetectedTools = append(result.DetectedTools, ToolCursor)
		files, err := collectCursor(root)
		if err != nil {
			return nil, fmt.Errorf("scan cursor: %w", err)
		}
		result.Files = append(result.Files, files...)
		toolFiles[ToolCursor] = len(files)
	}

	// Detect Copilot.
	copilotInstructions := filepath.Join(root, ".github", "copilot-instructions.md")
	copilotInstructionsDir := filepath.Join(root, ".github", "instructions")
	if fileExists(copilotInstructions) || isDir(copilotInstructionsDir) {
		result.DetectedTools = append(result.DetectedTools, ToolCopilot)
		files, err := collectCopilot(root)
		if err != nil {
			return nil, fmt.Errorf("scan copilot: %w", err)
		}
		result.Files = append(result.Files, files...)
		toolFiles[ToolCopilot] = len(files)
	}

	// Detect Codex.
	codexDir := filepath.Join(root, ".codex")
	agentsMD := filepath.Join(root, "AGENTS.md")
	if isDir(codexDir) || fileExists(agentsMD) {
		result.DetectedTools = append(result.DetectedTools, ToolCodex)
		// Codex uses AGENTS.md — not individual rule files to collect.
		toolFiles[ToolCodex] = 0
	}

	// Select source tool: tool with most files.
	result.SourceTool = selectSourceTool(toolFiles)

	return result, nil
}

// FilesForTool returns only files belonging to the given tool.
func (r *ScanResult) FilesForTool(tool Tool) []AgentFile {
	var result []AgentFile
	for _, f := range r.Files {
		if f.Tool == tool {
			result = append(result, f)
		}
	}
	return result
}

// FilesByType returns files of the given tool filtered by type.
func (r *ScanResult) FilesByType(tool Tool, ft FileType) []AgentFile {
	var result []AgentFile
	for _, f := range r.Files {
		if f.Tool == tool && f.Type == ft {
			result = append(result, f)
		}
	}
	return result
}

func collectClaude(root string) ([]AgentFile, error) {
	var files []AgentFile

	rules, err := collectMDFiles(root, filepath.Join(".claude", "rules"), ToolClaude, FileTypeRule, ".md")
	if err != nil {
		return nil, err
	}
	files = append(files, rules...)

	skills, err := collectSkillDirs(root, filepath.Join(".claude", "skills"), ToolClaude)
	if err != nil {
		return nil, err
	}
	files = append(files, skills...)

	agents, err := collectMDFiles(root, filepath.Join(".claude", "agents"), ToolClaude, FileTypeSubagent, ".md")
	if err != nil {
		return nil, err
	}
	files = append(files, agents...)

	cmds, err := collectMDFiles(root, filepath.Join(".claude", "commands"), ToolClaude, FileTypeCommand, ".md")
	if err != nil {
		return nil, err
	}
	files = append(files, cmds...)

	return files, nil
}

func collectCursor(root string) ([]AgentFile, error) {
	var files []AgentFile

	rules, err := collectMDFiles(root, filepath.Join(".cursor", "rules"), ToolCursor, FileTypeRule, ".mdc")
	if err != nil {
		return nil, err
	}
	files = append(files, rules...)

	skills, err := collectSkillDirs(root, filepath.Join(".cursor", "skills"), ToolCursor)
	if err != nil {
		return nil, err
	}
	files = append(files, skills...)

	return files, nil
}

func collectCopilot(root string) ([]AgentFile, error) {
	var files []AgentFile

	// Primary instructions file (detected by Scan but not previously collected).
	if fileExists(filepath.Join(root, ".github", "copilot-instructions.md")) {
		files = append(files, AgentFile{
			Tool: ToolCopilot,
			Type: FileTypeRule,
			Path: filepath.Join(".github", "copilot-instructions.md"),
			Slug: "copilot-instructions",
		})
	}

	rules, err := collectMDFiles(root, filepath.Join(".github", "instructions"), ToolCopilot, FileTypeRule, ".instructions.md")
	if err != nil {
		return nil, err
	}
	files = append(files, rules...)

	skills, err := collectSkillDirs(root, filepath.Join(".github", "skills"), ToolCopilot)
	if err != nil {
		return nil, err
	}
	files = append(files, skills...)

	return files, nil
}

// collectMDFiles collects files with the given extension from a directory.
func collectMDFiles(root, relDir string, tool Tool, ft FileType, ext string) ([]AgentFile, error) {
	absDir := filepath.Join(root, relDir)
	if !isDir(absDir) {
		return nil, nil
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}

	var files []AgentFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ext) {
			continue
		}
		slug := strings.TrimSuffix(e.Name(), ext)
		files = append(files, AgentFile{
			Tool: tool,
			Type: ft,
			Path: filepath.Join(relDir, e.Name()),
			Slug: slug,
		})
	}
	return files, nil
}

// collectSkillDirs collects skills from subdirectories containing SKILL.md.
func collectSkillDirs(root, relDir string, tool Tool) ([]AgentFile, error) {
	absDir := filepath.Join(root, relDir)
	if !isDir(absDir) {
		return nil, nil
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}

	var files []AgentFile
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillFile := filepath.Join(absDir, e.Name(), "SKILL.md")
		if fileExists(skillFile) {
			files = append(files, AgentFile{
				Tool: tool,
				Type: FileTypeSkill,
				Path: filepath.Join(relDir, e.Name(), "SKILL.md"),
				Slug: e.Name(),
			})
		}
	}
	return files, nil
}

func selectSourceTool(toolFiles map[Tool]int) Tool {
	if len(toolFiles) == 0 {
		return "" // no tools detected
	}

	type toolCount struct {
		tool  Tool
		count int
	}

	var counts []toolCount
	for t, c := range toolFiles {
		counts = append(counts, toolCount{tool: t, count: c})
	}

	sort.Slice(counts, func(i, j int) bool {
		if counts[i].count != counts[j].count {
			return counts[i].count > counts[j].count
		}
		return counts[i].tool < counts[j].tool // stable tie-break
	})

	return counts[0].tool
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
