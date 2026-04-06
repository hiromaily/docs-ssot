// Package mdutil provides shared Markdown processing utilities used by
// both the include processor and the index generator.
package mdutil

import (
	"path/filepath"
	"strconv"
	"strings"
)

// NextFenceType returns the updated fence type after processing line.
// fenceType is "" when outside a fence, or the opening fence string (e.g. "```", "~~~~") when inside.
// Per CommonMark:
//   - A backtick fence is only closed by backticks; a tilde fence only by tildes.
//   - A closing fence must have at least as many fence characters as the opening fence.
//   - A closing fence must contain only fence characters and optional trailing spaces.
//   - Up to 3 spaces of indentation are allowed on fence markers.
func NextFenceType(line, fenceType string) string {
	trimmed := strings.TrimLeft(line, " ")
	if len(line)-len(trimmed) > 3 {
		return fenceType // more than 3 spaces of indentation: not a fence marker
	}
	if fenceType == "" {
		// Opening fence: starts with ``` or ~~~
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			char := trimmed[0]
			i := 0
			for i < len(trimmed) && trimmed[i] == char {
				i++
			}
			return trimmed[:i] // store the exact opening fence (length matters for closing)
		}
	} else {
		// Closing fence: same character type, at least as many chars, trailing spaces only
		char := fenceType[0]
		i := 0
		for i < len(trimmed) && trimmed[i] == char {
			i++
		}
		if i >= len(fenceType) && strings.TrimRight(trimmed[i:], " ") == "" {
			return ""
		}
	}
	return fenceType
}

// ResolveIncludePath returns the absolute path for includePath relative to the containing file.
func ResolveIncludePath(absContainingFile, includePath string) string {
	if filepath.IsAbs(includePath) {
		return includePath
	}
	return filepath.Join(filepath.Dir(absContainingFile), includePath)
}

// FindGlobRoot returns the deepest ancestor directory of pattern that contains no glob metacharacters.
func FindGlobRoot(pattern string) string {
	dir := filepath.Dir(pattern)
	for strings.ContainsAny(dir, "*?[") {
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
	return dir
}

// MatchGlobParts recursively matches pattern segments against path segments.
// A "**" segment matches zero or more consecutive path segments.
func MatchGlobParts(pat, path []string) (bool, error) {
	for len(pat) > 0 {
		if pat[0] == "**" {
			if len(pat) == 1 {
				return true, nil // ** at end matches any remaining path
			}
			// Try consuming zero path segments (skip **)
			if ok, err := MatchGlobParts(pat[1:], path); err != nil || ok {
				return ok, err
			}
			// Try consuming one path segment and retrying with the same **
			if len(path) == 0 {
				return false, nil
			}
			return MatchGlobParts(pat, path[1:])
		}
		if len(path) == 0 {
			return false, nil
		}
		matched, err := filepath.Match(pat[0], path[0])
		if err != nil || !matched {
			return false, err
		}
		pat = pat[1:]
		path = path[1:]
	}
	return len(path) == 0, nil
}

// ParseIncludeArgs parses the argument string captured from an include directive.
// The expected form is: <path> [level=<delta>]
// Returns the file path and optional level delta (0 if absent or unparseable).
func ParseIncludeArgs(args string) (string, int) {
	args = strings.TrimSpace(args)
	if args == "" {
		return "", 0
	}
	var level int
	if idx := strings.LastIndex(args, " level="); idx != -1 {
		n, err := strconv.Atoi(args[idx+len(" level="):])
		if err == nil {
			level = n
			args = strings.TrimSpace(args[:idx])
		}
	}
	return args, level
}
