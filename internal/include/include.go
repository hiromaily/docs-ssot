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

var includePattern = regexp.MustCompile(`<!--\s*@include:\s*(.*?)\s*-->`)

// ProcessFile processes a template file, recursively resolving all include directives.
// Include paths are resolved relative to the current working directory.
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
	defer file.Close()

	chain := append(ancestors, absPath)

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	inCodeFence := false

	for scanner.Scan() {
		line := scanner.Text()

		// Track fenced code blocks so include directives inside them are treated as literal text.
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			inCodeFence = !inCodeFence
		}

		if !inCodeFence {
			matches := includePattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				includePath := matches[1]
				absInclude, err := filepath.Abs(includePath)
				if err != nil {
					return "", fmt.Errorf("failed to resolve include path (%s): %w", includePath, err)
				}

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
