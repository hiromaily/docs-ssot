// Package index generates a template INDEX.md by scanning the template directory,
// parsing include directives, building a reverse include map, and detecting orphans.
package index

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/hiromaily/docs-ssot/internal/config"
	"github.com/hiromaily/docs-ssot/internal/mdutil"
)

// includePattern matches <!-- @include: <path> [level=±N] -->.
var includePattern = regexp.MustCompile(`^\s*<!--\s*@include:\s*(.*?)\s*-->\s*$`)

// TemplateInfo holds a template file and the docs files it includes.
type TemplateInfo struct {
	Path     string   // relative to template dir (e.g. "pages/README.tpl.md")
	Output   string   // output file (e.g. "README.md")
	Includes []string // resolved docs files it references (relative to template dir)
}

// IndexData contains all data needed to render an INDEX.md.
type IndexData struct {
	Pages    []TemplateInfo      // page templates (pages/*.tpl.md)
	Sections map[string][]string // sections/** file -> list of referencing template labels
	Rules    map[string][]string // sections/ai/rules/** file -> list of referencing template labels
	Commands map[string][]string // sections/ai/commands/** file -> list of referencing template labels
	Orphans  []string            // sections/** files not referenced by any template
}

// Generate scans the template directory, parses include directives, builds
// the reverse include map, and returns the structured index data.
func Generate(templateDir string, cfg *config.Config) (*IndexData, error) {
	// 1. Collect all template files (.tpl.md, .tpl.mdc) and their include relationships
	templates, err := scanTemplates(templateDir, cfg)
	if err != nil {
		return nil, fmt.Errorf("scanning templates: %w", err)
	}

	// 2. Collect all section files
	sectionsDir := filepath.Join(templateDir, "sections")
	allDocs, err := scanDocsFiles(sectionsDir, templateDir)
	if err != nil {
		return nil, fmt.Errorf("scanning docs: %w", err)
	}

	// 3. Build reverse map: docs file -> list of referencing template labels
	reverseMap := buildReverseMap(templates)

	// 4. Classify docs files into sections/rules/commands and detect orphans
	data := classify(allDocs, reverseMap, templates)

	return data, nil
}

// Render produces the INDEX.md content from the given data.
func Render(data *IndexData) string {
	var sb strings.Builder

	sb.WriteString("<!-- AUTO-GENERATED FILE — DO NOT EDIT -->\n")
	sb.WriteString("<!-- Regenerate with: docs-ssot index -->\n")
	sb.WriteString("# Template Index\n\n")

	// Pages
	sb.WriteString("## Pages\n\n")
	sb.WriteString("| Template | Output | Sections included |\n")
	sb.WriteString("| --- | --- | --- |\n")
	for _, p := range data.Pages {
		fmt.Fprintf(&sb, "| %s | %s | %d |\n", p.Path, p.Output, len(p.Includes))
	}
	sb.WriteByte('\n')

	// Sections
	renderRefTable(&sb, "Sections", data.Sections)

	// Rules
	renderRefTable(&sb, "Rules", data.Rules)

	// Commands
	renderRefTable(&sb, "Commands", data.Commands)

	// Orphans
	sb.WriteString("## Orphans\n\n")
	if len(data.Orphans) == 0 {
		sb.WriteString("| File | Note |\n")
		sb.WriteString("| --- | --- |\n")
		sb.WriteString("| (none) | All files are referenced |\n")
	} else {
		sb.WriteString("| File | Note |\n")
		sb.WriteString("| --- | --- |\n")
		for _, o := range data.Orphans {
			fmt.Fprintf(&sb, "| %s | Not referenced by any template |\n", o)
		}
	}

	return sb.String()
}

// renderRefTable writes a "Referenced by" table section.
func renderRefTable(sb *strings.Builder, title string, refs map[string][]string) {
	fmt.Fprintf(sb, "## %s\n\n", title)
	sb.WriteString("| File | Referenced by |\n")
	sb.WriteString("| --- | --- |\n")

	keys := sortedKeys(refs)
	for _, file := range keys {
		labels := refs[file]
		fmt.Fprintf(sb, "| %s | %s |\n", file, strings.Join(labels, ", "))
	}
	sb.WriteByte('\n')
}

// scanTemplates discovers all template files and resolves their include directives
// to build a list of TemplateInfo.
func scanTemplates(templateDir string, cfg *config.Config) ([]TemplateInfo, error) {
	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		return nil, err
	}

	var templates []TemplateInfo

	for _, t := range cfg.Targets {
		absInput, err := filepath.Abs(t.Input)
		if err != nil {
			return nil, err
		}

		includes, err := extractIncludes(absInput, templateDir)
		if err != nil {
			return nil, fmt.Errorf("parsing includes in %s: %w", t.Input, err)
		}

		relInput, err := filepath.Rel(absTemplateDir, absInput)
		if err != nil {
			relInput = t.Input
		}

		templates = append(templates, TemplateInfo{
			Path:     relInput,
			Output:   t.Output,
			Includes: includes,
		})
	}

	return templates, nil
}

// extractIncludes parses a template file and returns all resolved docs file paths
// (relative to templateDir) that it references via include directives.
// It does NOT recurse into included files — it only looks at the template itself.
func extractIncludes(absPath, templateDir string) ([]string, error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		return nil, err
	}

	var includes []string
	fenceType := ""

	for line := range strings.SplitSeq(string(data), "\n") {
		prevFence := fenceType
		fenceType = mdutil.NextFenceType(line, fenceType)
		if prevFence != "" || fenceType != "" {
			continue
		}

		matches := includePattern.FindStringSubmatch(line)
		if len(matches) <= 1 {
			continue
		}

		includePath, _ := mdutil.ParseIncludeArgs(matches[1])
		absInclude := mdutil.ResolveIncludePath(absPath, includePath)

		// Resolve the include to actual files
		resolved, err := resolveToFiles(absInclude, includePath, absTemplateDir)
		if err != nil {
			// Skip unresolvable includes (they may fail at build time, not index time)
			continue
		}
		includes = append(includes, resolved...)
	}

	return dedupSorted(includes), nil
}

// resolveToFiles resolves an include path to a list of actual file paths (relative to templateDir).
func resolveToFiles(absInclude, includePath, absTemplateDir string) ([]string, error) {
	switch {
	case strings.HasSuffix(includePath, "/"):
		return resolveDirInclude(absInclude, absTemplateDir)
	case strings.Contains(includePath, "**"):
		return resolveRecursiveGlob(absInclude, absTemplateDir)
	case strings.ContainsAny(includePath, "*?["):
		return resolveGlobInclude(absInclude, absTemplateDir)
	default:
		return resolveSingleFile(absInclude, absTemplateDir)
	}
}

func resolveDirInclude(absDir, absTemplateDir string) ([]string, error) {
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		absFile := filepath.Join(absDir, entry.Name())
		rel, relErr := filepath.Rel(absTemplateDir, absFile)
		if relErr != nil {
			continue
		}
		files = append(files, rel)
	}
	return files, nil
}

func resolveRecursiveGlob(pattern, absTemplateDir string) ([]string, error) {
	root := mdutil.FindGlobRoot(pattern)
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}
	patParts := strings.Split(filepath.ToSlash(pattern), "/")
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		pathParts := strings.Split(filepath.ToSlash(path), "/")
		matched, matchErr := mdutil.MatchGlobParts(patParts, pathParts)
		if matchErr != nil || !matched {
			return matchErr
		}
		rel, relErr := filepath.Rel(absTemplateDir, path)
		if relErr != nil {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	return files, err
}

func resolveGlobInclude(pattern, absTemplateDir string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, match := range matches {
		info, statErr := os.Stat(match)
		if statErr != nil || info.IsDir() {
			continue
		}
		rel, relErr := filepath.Rel(absTemplateDir, match)
		if relErr != nil {
			continue
		}
		files = append(files, rel)
	}
	return files, nil
}

func resolveSingleFile(absInclude, absTemplateDir string) ([]string, error) {
	if _, err := os.Stat(absInclude); err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(absTemplateDir, absInclude)
	if err != nil {
		return nil, err
	}
	return []string{rel}, nil
}

// scanDocsFiles returns all .md files under docsDir, relative to templateDir.
func scanDocsFiles(docsDir, templateDir string) ([]string, error) {
	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		return nil, err
	}
	absDocsDir, err := filepath.Abs(docsDir)
	if err != nil {
		return nil, err
	}

	if _, statErr := os.Stat(absDocsDir); statErr != nil {
		return nil, nil // no docs dir → no files
	}

	var files []string
	walkErr := filepath.WalkDir(absDocsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		rel, relErr := filepath.Rel(absTemplateDir, path)
		if relErr != nil {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	slices.Sort(files)
	return files, nil
}

// buildReverseMap builds a map from docs file path to list of template labels that reference it.
func buildReverseMap(templates []TemplateInfo) map[string][]string {
	rm := make(map[string][]string)

	for _, t := range templates {
		label := templateLabel(t)
		for _, inc := range t.Includes {
			rm[inc] = append(rm[inc], label)
		}
	}

	// Deduplicate labels
	for k, v := range rm {
		rm[k] = dedupSorted(v)
	}

	return rm
}

// templateLabel returns a short human-readable label for a template.
func templateLabel(t TemplateInfo) string {
	path := t.Path

	// AI agent templates: use tool name (e.g. "claude", "cursor")
	if after, ok := strings.CutPrefix(path, "pages/ai-agents/"); ok {
		parts := strings.SplitN(after, "/", 2)
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// Page templates: use output name without extension
	if strings.HasPrefix(path, "pages/") {
		name := strings.TrimSuffix(filepath.Base(t.Output), filepath.Ext(t.Output))
		return name
	}

	return filepath.Base(path)
}

// classify categorizes docs files into sections/rules/commands and identifies orphans.
func classify(allDocs []string, reverseMap map[string][]string, templates []TemplateInfo) *IndexData {
	data := &IndexData{
		Sections: make(map[string][]string),
		Rules:    make(map[string][]string),
		Commands: make(map[string][]string),
	}

	// Extract page templates
	for _, t := range templates {
		if strings.HasPrefix(t.Path, "pages/") {
			data.Pages = append(data.Pages, t)
		}
	}

	// Classify docs files
	for _, f := range allDocs {
		refs, referenced := reverseMap[f]
		if !referenced {
			data.Orphans = append(data.Orphans, f)
			continue
		}

		switch {
		case strings.HasPrefix(f, "sections/ai/rules/"):
			data.Rules[f] = refs
		case strings.HasPrefix(f, "sections/ai/commands/"):
			data.Commands[f] = refs
		default:
			data.Sections[f] = refs
		}
	}

	return data
}

// dedupSorted returns a sorted, deduplicated copy of the string slice.
func dedupSorted(s []string) []string {
	if len(s) == 0 {
		return s
	}
	sorted := make([]string, len(s))
	copy(sorted, s)
	sort.Strings(sorted)
	result := sorted[:1]
	for _, v := range sorted[1:] {
		if v != result[len(result)-1] {
			result = append(result, v)
		}
	}
	return result
}

// DetectTemplateDir infers the template root directory from the config targets.
// It looks for the first target whose input path contains "template/" and returns
// the path up to and including that segment.
func DetectTemplateDir(cfg *config.Config) string {
	for _, t := range cfg.Targets {
		if idx := strings.Index(t.Input, "template/"); idx >= 0 {
			return t.Input[:idx+len("template")]
		}
	}
	return "template"
}

// sortedKeys returns the keys of a map in sorted order.
func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
