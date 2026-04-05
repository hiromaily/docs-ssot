// Package include resolves include directives in Markdown files.
package include

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// includePattern matches an include directive that occupies its own line (with optional surrounding whitespace).
var includePattern = regexp.MustCompile(`^\s*<!--\s*@include:\s*(.*?)\s*-->\s*$`)

// ProcessFile processes a template file, recursively resolving all include directives.
// Relative include paths are resolved relative to the directory of the file containing the directive.
// Directives inside fenced code blocks are treated as literal text and not expanded.
func ProcessFile(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path (%s): %w", path, err)
	}
	return processFile(absPath, []string{})
}

// processFile recursively resolves include directives.
// ancestors holds the absolute paths of files in the current include chain for circular detection.
func processFile(absPath string, ancestors []string) (string, error) {
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

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	// Use a 1 MB buffer to handle files with long lines (e.g. large embedded diagrams).
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	// fenceType is "" when outside a code fence, or "```"/"~~~" when inside one.
	// Per CommonMark: backtick fences are closed only by backticks, tilde fences only by tildes.
	fenceType := ""

	for scanner.Scan() {
		line := scanner.Text()

		// Detect fence open/close. CommonMark allows up to 3 spaces of indentation.
		prevFenceType := fenceType
		fenceType = nextFenceType(line, fenceType)

		// Only process include directives when the line is normal text (not a fence
		// marker and not inside a fence block).
		if prevFenceType == "" && fenceType == "" {
			matches := includePattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				absInclude := resolveIncludePath(absPath, matches[1])

				content, err := processFile(absInclude, chain)
				if err != nil {
					return "", err
				}

				sb.WriteString(content)
				if !strings.HasSuffix(content, "\n") {
					sb.WriteByte('\n')
				}
				continue
			}
		}

		sb.WriteString(line + "\n")
	}

	return sb.String(), scanner.Err()
}

// nextFenceType returns the updated fence type after processing line.
// Per CommonMark, a backtick fence is only closed by backticks, and a tilde fence only by tildes.
// Up to 3 spaces of indentation are allowed on fence markers.
func nextFenceType(line, fenceType string) string {
	trimmed := strings.TrimLeft(line, " ")
	if len(line)-len(trimmed) > 3 {
		return fenceType // more than 3 spaces of indentation: not a fence marker
	}
	switch fenceType {
	case "":
		if strings.HasPrefix(trimmed, "```") {
			return "```"
		}
		if strings.HasPrefix(trimmed, "~~~") {
			return "~~~"
		}
	case "```":
		if strings.HasPrefix(trimmed, "```") {
			return ""
		}
	case "~~~":
		if strings.HasPrefix(trimmed, "~~~") {
			return ""
		}
	}
	return fenceType
}

// resolveIncludePath returns the absolute path for includePath relative to the containing file.
func resolveIncludePath(absContainingFile, includePath string) string {
	if filepath.IsAbs(includePath) {
		return includePath
	}
	return filepath.Join(filepath.Dir(absContainingFile), includePath)
}
