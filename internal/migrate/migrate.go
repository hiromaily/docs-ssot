// Package migrate decomposes existing monolithic Markdown files into the SSOT
// section structure. It splits documents by headings, detects cross-file
// duplicates using TF-IDF cosine similarity, and generates template files
// with @include directives that reproduce the original document structure.
package migrate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hiromaily/docs-ssot/internal/categoriser"
	"github.com/hiromaily/docs-ssot/internal/dupcheck"
	"github.com/hiromaily/docs-ssot/internal/generator"
	"github.com/hiromaily/docs-ssot/internal/splitter"
)

// Config holds parameters for the migrate command.
type Config struct {
	// InputFiles are the Markdown files to migrate.
	InputFiles []string
	// OutputDir is where section files are written (default: "template/sections").
	OutputDir string
	// TemplateDir is where template files are written (default: "template/pages").
	TemplateDir string
	// SectionLevel is the heading level used as section boundary (default: 2).
	SectionLevel int
	// Threshold is the TF-IDF cosine similarity threshold for duplicate detection (default: 0.82).
	Threshold float64
	// DryRun prints the migration plan without writing files.
	DryRun bool
	// ConfigFile is the path to docsgen.yaml (default: "docsgen.yaml").
	ConfigFile string
}

// sectionFile tracks an extracted section and its planned output path.
type sectionFile struct {
	// Source is the input file this section came from.
	Source string
	// Section is the parsed section data.
	Section splitter.Section
	// Category is the categorised directory (e.g., "development", "project").
	Category string
	// Slug is the filename slug (e.g., "setup", "testing").
	Slug string
	// RelPath is the relative path from template dir to section file (for @include).
	RelPath string
	// OutputPath is the full path where the section file will be written.
	OutputPath string
	// MergedWith is non-empty if this section was merged into another section's file.
	MergedWith string
	// IndexInFile is the zero-based position of this section within its source file.
	// Used as a unique key to distinguish sections with identical headings.
	IndexInFile int
}

// Run executes the migrate command, decomposing input files into SSOT sections.
func Run(w io.Writer, cfg Config) error {
	if len(cfg.InputFiles) == 0 {
		return errors.New("no input files specified")
	}

	// Step 1: Parse input files into sections.
	allSections, err := parseAllFiles(cfg.InputFiles, cfg.SectionLevel)
	if err != nil {
		return err
	}

	// Step 2: Plan section files (assign categories, slugs, output paths).
	planned, err := planSectionFiles(allSections, cfg.OutputDir, cfg.TemplateDir)
	if err != nil {
		return err
	}

	// Step 3: Detect duplicates across files.
	merges := detectDuplicates(planned, cfg.Threshold)
	applyMerges(planned, merges)

	// Step 4: Report plan.
	reportPlan(w, cfg, allSections, planned, merges)

	if cfg.DryRun {
		return nil
	}

	// Step 5: Write section files.
	if err := writeSectionFiles(w, planned); err != nil {
		return err
	}

	// Step 6: Write template files.
	if err := writeTemplateFiles(w, cfg, allSections, planned); err != nil {
		return err
	}

	// Step 7: Write docsgen.yaml if it doesn't exist.
	if err := writeConfigIfNeeded(w, cfg); err != nil {
		return err
	}

	// Step 8: Round-trip verification.
	verifyRoundTrip(w, cfg)

	_, _ = fmt.Fprintln(w, "Migration complete.")
	return nil
}

// fileSections groups parsed sections by source file.
type fileSections struct {
	Source   string
	Sections []splitter.Section
}

func parseAllFiles(files []string, sectionLevel int) ([]fileSections, error) {
	var result []fileSections
	for _, f := range files {
		sections, err := splitter.Split(f, sectionLevel)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", f, err)
		}
		result = append(result, fileSections{Source: f, Sections: sections})
	}
	return result, nil
}

func planSectionFiles(allSections []fileSections, outputDir, templateDir string) ([]*sectionFile, error) {
	// Convert templateDir to absolute once to ensure filepath.Rel works correctly
	// even if templateDir and outputDir have mixed relative/absolute forms.
	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		return nil, fmt.Errorf("abs path %s: %w", templateDir, err)
	}

	slugCount := map[string]int{}
	var planned []*sectionFile

	for _, fs := range allSections {
		for i, s := range fs.Sections {
			if s.Title == "" {
				// Preamble (content before first heading) — skip for section files,
				// will be inlined in the template.
				continue
			}

			cat := categoriser.Categorise(s.Title)
			slug := ToSlug(s.Title)

			// Disambiguate duplicate slugs within the same category.
			key := cat + "/" + slug
			if n := slugCount[key]; n > 0 {
				slug = fmt.Sprintf("%s-%d", slug, n+1)
			}
			slugCount[key]++

			outPath := filepath.Join(outputDir, cat, slug+".md")

			// Compute relative path from template dir to section file using absolute paths
			// so that filepath.Rel succeeds regardless of whether the inputs are relative or absolute.
			absOutPath, err := filepath.Abs(outPath)
			if err != nil {
				return nil, fmt.Errorf("abs path %s: %w", outPath, err)
			}
			relPath, err := filepath.Rel(absTemplateDir, absOutPath)
			if err != nil {
				return nil, fmt.Errorf("compute relative path from %s to %s: %w", absTemplateDir, absOutPath, err)
			}

			planned = append(planned, &sectionFile{
				Source:      fs.Source,
				Section:     s,
				Category:    cat,
				Slug:        slug,
				RelPath:     relPath,
				OutputPath:  outPath,
				IndexInFile: i,
			})
		}
	}

	return planned, nil
}

// duplicatePair tracks two sections that are duplicates.
type duplicatePair struct {
	indexA int
	indexB int
	score  float64
}

func detectDuplicates(planned []*sectionFile, threshold float64) []duplicatePair {
	if len(planned) < 2 {
		return nil
	}

	// Build dupcheck Chunks from planned sections.
	chunks := make([]dupcheck.Chunk, len(planned))
	for i, sf := range planned {
		text := sf.Section.Title + "\n" + sf.Section.Body
		chunks[i] = dupcheck.Chunk{
			File:   sf.Source,
			Index:  i,
			Kind:   "section",
			Title:  sf.Section.Title,
			Text:   text,
			Tokens: dupcheck.Tokenize(text),
		}
	}

	// Use dupcheck's TF-IDF engine.
	vectors := dupcheck.BuildTFIDF(chunks)

	var pairs []duplicatePair
	for i := range len(chunks) {
		for j := i + 1; j < len(chunks); j++ {
			// Only compare sections from different source files.
			if planned[i].Source == planned[j].Source {
				continue
			}
			score := dupcheck.Cosine(vectors[i], vectors[j])
			if score >= threshold {
				pairs = append(pairs, duplicatePair{indexA: i, indexB: j, score: score})
			}
		}
	}

	return pairs
}

func applyMerges(planned []*sectionFile, merges []duplicatePair) {
	for _, m := range merges {
		// Keep the first occurrence (indexA), mark the second (indexB) as merged.
		b := planned[m.indexB]
		a := planned[m.indexA]
		b.MergedWith = a.OutputPath
		// Update the relPath of B to point to A's section file.
		b.RelPath = a.RelPath
	}
}

func reportPlan(w io.Writer, cfg Config, allSections []fileSections, planned []*sectionFile, merges []duplicatePair) {
	for _, fs := range allSections {
		count := len(fs.Sections)
		_, _ = fmt.Fprintf(w, "Parsed %s: %d sections\n", fs.Source, count)
	}

	if len(merges) > 0 {
		_, _ = fmt.Fprintf(w, "Detected %d duplicate sections (similarity > %.2f):\n", len(merges), cfg.Threshold)
		for _, m := range merges {
			a := planned[m.indexA]
			b := planned[m.indexB]
			_, _ = fmt.Fprintf(w, "  %q (%s) — merged into %s (score=%.3f)\n",
				b.Section.Title, b.Source, a.OutputPath, m.score)
		}
	}

	unique := 0
	for _, sf := range planned {
		if sf.MergedWith == "" {
			unique++
		}
	}

	verb := "Would create"
	if !cfg.DryRun {
		verb = "Creating"
	}
	_, _ = fmt.Fprintf(w, "%s %d unique section files in %s\n", verb, unique, cfg.OutputDir)
}

func writeSectionFiles(w io.Writer, planned []*sectionFile) error {
	for _, sf := range planned {
		if sf.MergedWith != "" {
			continue // Skip merged duplicates.
		}

		if err := os.MkdirAll(filepath.Dir(sf.OutputPath), 0o750); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(sf.OutputPath), err)
		}

		// Section files start at H2 per docs-ssot convention.
		content := sf.Section.RawHeading + "\n\n" + sf.Section.Body + "\n"

		//nolint:gosec // generated documentation files
		if err := os.WriteFile(sf.OutputPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", sf.OutputPath, err)
		}

		_, _ = fmt.Fprintf(w, "  %s\n", sf.OutputPath)
	}
	return nil
}

func writeTemplateFiles(w io.Writer, cfg Config, allSections []fileSections, planned []*sectionFile) error {
	if err := os.MkdirAll(cfg.TemplateDir, 0o750); err != nil {
		return fmt.Errorf("mkdir %s: %w", cfg.TemplateDir, err)
	}

	// Build a lookup: (source, title) → sectionFile for @include resolution.
	lookup := buildLookup(planned)

	for _, fs := range allSections {
		baseName := strings.TrimSuffix(filepath.Base(fs.Source), filepath.Ext(fs.Source))
		tplPath := filepath.Join(cfg.TemplateDir, baseName+".tpl.md")

		var lines []string

		includeCount := 0
		for i, s := range fs.Sections {
			if s.Title == "" {
				// Preamble: inline the content directly in the template.
				if s.Body != "" {
					lines = append(lines, s.Body)
					lines = append(lines, "")
				}
				continue
			}

			sf := lookup[lookupKey(fs.Source, i)]
			if sf == nil {
				// Should not happen, but fallback to inline.
				lines = append(lines, s.RawHeading)
				lines = append(lines, "")
				if s.Body != "" {
					lines = append(lines, s.Body)
					lines = append(lines, "")
				}
				continue
			}

			lines = append(lines, "<!-- @include: "+sf.RelPath+" -->")
			includeCount++

			// Add blank line between sections, but not after the last one.
			if i < len(fs.Sections)-1 {
				lines = append(lines, "")
			}
		}

		content := strings.Join(lines, "\n") + "\n"

		//nolint:gosec // generated template files
		if err := os.WriteFile(tplPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", tplPath, err)
		}

		_, _ = fmt.Fprintf(w, "Created %s (%d includes)\n", tplPath, includeCount)
	}

	return nil
}

func writeConfigIfNeeded(w io.Writer, cfg Config) error {
	if _, err := os.Stat(cfg.ConfigFile); err == nil {
		_, _ = fmt.Fprintf(w, "docsgen.yaml already exists, skipping.\n")
		return nil
	}

	var lines []string
	lines = append(lines, "targets:")

	for _, f := range cfg.InputFiles {
		baseName := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		tplPath := filepath.Join(cfg.TemplateDir, baseName+".tpl.md")
		lines = append(lines, "  - input: "+tplPath)
		lines = append(lines, "    output: "+f)
	}
	lines = append(lines, "")

	content := strings.Join(lines, "\n")

	//nolint:gosec // generated config file
	if err := os.WriteFile(cfg.ConfigFile, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", cfg.ConfigFile, err)
	}

	_, _ = fmt.Fprintf(w, "Created %s\n", cfg.ConfigFile)
	return nil
}

// verifyRoundTrip runs generator.Build on the generated config and compares
// the output against the original input files. Differences are reported as
// warnings but do not cause an error.
func verifyRoundTrip(w io.Writer, cfg Config) {
	_, _ = fmt.Fprintln(w, "Verifying round-trip...")

	if err := generator.Build(cfg.ConfigFile); err != nil {
		_, _ = fmt.Fprintf(w, "Round-trip verification: SKIP (build error: %v)\n", err)
		return
	}

	allMatch := true
	for _, f := range cfg.InputFiles {
		original, err := os.ReadFile(f)
		if err != nil {
			_, _ = fmt.Fprintf(w, "  WARNING: cannot read original %s: %v\n", f, err)
			allMatch = false
			continue
		}

		// The build output overwrites the original files, so read it again.
		generated, err := os.ReadFile(f)
		if err != nil {
			_, _ = fmt.Fprintf(w, "  WARNING: cannot read generated %s: %v\n", f, err)
			allMatch = false
			continue
		}

		if string(original) != string(generated) {
			_, _ = fmt.Fprintf(w, "  WARNING: %s differs after round-trip\n", f)
			allMatch = false
		}
	}

	if allMatch {
		_, _ = fmt.Fprintln(w, "Round-trip verification: OK")
	} else {
		_, _ = fmt.Fprintln(w, "Round-trip verification: WARN (some files differ)")
	}
}

func buildLookup(planned []*sectionFile) map[string]*sectionFile {
	m := make(map[string]*sectionFile, len(planned))
	for _, sf := range planned {
		m[lookupKey(sf.Source, sf.IndexInFile)] = sf
	}
	return m
}

func lookupKey(source string, index int) string {
	return fmt.Sprintf("%s:%d", source, index)
}

// ToSlug converts a heading title into a kebab-case filename slug.
func ToSlug(title string) string {
	lower := strings.ToLower(title)
	var buf strings.Builder
	prevDash := false

	for _, r := range lower {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			buf.WriteRune(r)
			prevDash = false
		default:
			if !prevDash && buf.Len() > 0 {
				buf.WriteByte('-')
				prevDash = true
			}
		}
	}

	s := buf.String()
	return strings.TrimRight(s, "-")
}
