// Package processor resolves include directives and transforms Markdown content.
// The public API is ProcessFile; extensible content transformations are defined
// via the Transformer interface and applied through Apply.
package processor

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/hiromaily/docs-ssot/internal/mdutil"
)

// includePattern matches an include directive that occupies its own line (with optional surrounding whitespace).
var includePattern = regexp.MustCompile(`^\s*<!--\s*@include:\s*(.*?)\s*-->\s*$`)

// ProcessFile processes a template file, recursively resolving all include directives.
// Relative include paths are resolved relative to the directory of the file containing the directive.
// Relative Markdown links and image URLs are rewritten to be correct relative to outputPath.
// Directives inside fenced code blocks are treated as literal text and not expanded.
func ProcessFile(path, outputPath string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path (%s): %w", path, err)
	}
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve output path (%s): %w", outputPath, err)
	}
	return processFile(absPath, []string{}, absOutputPath)
}

// processFile recursively resolves include directives.
// ancestors holds the absolute paths of files in the current include chain for circular detection.
// absOutputPath is the absolute path of the final output file, used for link rewriting.
func processFile(absPath string, ancestors []string, absOutputPath string) (string, error) {
	if slices.Contains(ancestors, absPath) {
		return "", fmt.Errorf("circular include detected: %s -> %s", strings.Join(ancestors, " -> "), absPath)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("include error (%s): %w", absPath, err)
	}
	defer func() { _ = file.Close() }()

	// Build a new chain slice with its own backing array so recursive calls cannot
	// accidentally modify each other's ancestor lists via shared capacity.
	chain := slices.Concat(ancestors, []string{absPath})

	// Precompute directories once for link rewriting on every line of this file.
	sourceDir := filepath.Dir(absPath)
	outputDir := filepath.Dir(absOutputPath)

	var sb strings.Builder

	// localBuf accumulates raw lines from the current file between include directives.
	// It is flushed (with link rewriting applied) whenever an include directive is encountered
	// and at the end of the file, ensuring link rewriting is scoped to each file's own lines.
	var localBuf strings.Builder

	flushLocalBuf := func() error {
		if localBuf.Len() == 0 {
			return nil
		}
		rewritten, applyErr := Apply(localBuf.String(), LinkTransformer{SourceDir: sourceDir, OutputDir: outputDir})
		if applyErr != nil {
			return applyErr
		}
		sb.WriteString(rewritten)
		localBuf.Reset()
		return nil
	}

	scanner := bufio.NewScanner(file)
	// Allow lines up to 1 MB; start with a nil buffer so the scanner allocates only as needed.
	scanner.Buffer(nil, 1024*1024)

	// fenceType is "" when outside a code fence, or "```"/"~~~" when inside one.
	// Per CommonMark: backtick fences are closed only by backticks, tilde fences only by tildes.
	fenceType := ""

	for scanner.Scan() {
		line := scanner.Text()

		// Detect fence open/close. CommonMark allows up to 3 spaces of indentation.
		prevFenceType := fenceType
		fenceType = mdutil.NextFenceType(line, fenceType)

		// Only process include directives when the line is normal text
		// (not a fence marker and not inside a fence block).
		if prevFenceType == "" && fenceType == "" {
			matches := includePattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				// Flush local lines (with link rewriting) before inserting included content.
				if err := flushLocalBuf(); err != nil {
					return "", err
				}

				includePath, levelDelta := mdutil.ParseIncludeArgs(matches[1])
				content, includeErr := resolveInclude(absPath, includePath, levelDelta, chain, absOutputPath)
				if includeErr != nil {
					return "", includeErr
				}

				sb.WriteString(content)
				if content != "" && !strings.HasSuffix(content, "\n") {
					sb.WriteByte('\n')
				}
				continue
			}
		}
		localBuf.WriteString(line + "\n")
	}

	if err := flushLocalBuf(); err != nil {
		return "", err
	}

	return sb.String(), scanner.Err()
}

// resolveInclude dispatches an include directive to the appropriate handler based on the path,
// applies optional heading-level adjustment, and returns the assembled content.
func resolveInclude(absContainingFile, includePath string, levelDelta int, chain []string, absOutputPath string) (string, error) {
	absInclude := mdutil.ResolveIncludePath(absContainingFile, includePath)

	var (
		content string
		err     error
	)
	switch {
	case strings.HasSuffix(includePath, "/"):
		content, err = processDirectory(absInclude, chain, absOutputPath)
	case strings.Contains(includePath, "**"):
		content, err = processRecursiveGlob(absInclude, chain, absOutputPath)
	case strings.ContainsAny(includePath, "*?["):
		content, err = processGlob(absInclude, chain, absOutputPath)
	default:
		content, err = processFile(absInclude, chain, absOutputPath)
	}
	if err != nil {
		return "", err
	}

	if levelDelta != 0 {
		content, err = HeadingTransformer{Delta: levelDelta}.Transform(content)
		if err != nil {
			return "", err
		}
	}

	return content, nil
}

// processGlob includes all files matched by the glob pattern (sorted lexically) by processing each in order.
// Only regular files (not directories) are included. No error is returned when the pattern matches nothing.
func processGlob(pattern string, ancestors []string, absOutputPath string) (string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("include error (glob %s): %w", pattern, err)
	}
	var sb strings.Builder
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return "", fmt.Errorf("include error (glob match %s): %w", match, err)
		}
		if info.IsDir() {
			continue
		}
		content, err := processFile(match, ancestors, absOutputPath)
		if err != nil {
			return "", err
		}
		sb.WriteString(content)
		if content != "" && !strings.HasSuffix(content, "\n") {
			sb.WriteByte('\n')
		}
	}
	return sb.String(), nil
}

// processRecursiveGlob handles patterns containing "**" by walking the directory tree recursively.
// Files are included in lexical path order. No error is returned when the root directory does not exist
// or no files match the pattern.
func processRecursiveGlob(pattern string, ancestors []string, absOutputPath string) (string, error) {
	root := mdutil.FindGlobRoot(pattern)

	if _, statErr := os.Stat(root); statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("include error (recursive glob %s): %w", pattern, statErr)
	}

	patParts := strings.Split(filepath.ToSlash(pattern), "/")

	var matchedPaths []string
	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		pathParts := strings.Split(filepath.ToSlash(path), "/")
		matched, matchErr := mdutil.MatchGlobParts(patParts, pathParts)
		if matchErr != nil {
			return fmt.Errorf("include error (recursive glob %s): %w", pattern, matchErr)
		}
		if matched {
			matchedPaths = append(matchedPaths, path)
		}
		return nil
	})
	if walkErr != nil {
		return "", fmt.Errorf("include error (recursive glob %s): %w", pattern, walkErr)
	}

	slices.Sort(matchedPaths)

	var sb strings.Builder
	for _, match := range matchedPaths {
		content, err := processFile(match, ancestors, absOutputPath)
		if err != nil {
			return "", err
		}
		sb.WriteString(content)
		if content != "" && !strings.HasSuffix(content, "\n") {
			sb.WriteByte('\n')
		}
	}
	return sb.String(), nil
}

// processDirectory includes all .md files in absDir (sorted by filename) by processing each in order.
func processDirectory(absDir string, ancestors []string, absOutputPath string) (string, error) {
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return "", fmt.Errorf("include error (directory %s): %w", absDir, err)
	}

	var sb strings.Builder
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		absFile := filepath.Join(absDir, entry.Name())
		content, err := processFile(absFile, ancestors, absOutputPath)
		if err != nil {
			return "", err
		}
		sb.WriteString(content)
		if content != "" && !strings.HasSuffix(content, "\n") {
			sb.WriteByte('\n')
		}
	}
	return sb.String(), nil
}
