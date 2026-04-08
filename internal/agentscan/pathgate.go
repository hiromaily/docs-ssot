package agentscan

import (
	"slices"
	"strings"
)

// knownExtensions maps slug keywords to file glob patterns.
// Used to infer path-gated rules for Cursor (globs) and Copilot (applyTo).
var knownExtensions = map[string]string{
	"go":         "**/*.go",
	"golang":     "**/*.go",
	"typescript": "**/*.{ts,tsx}",
	"ts":         "**/*.{ts,tsx}",
	"javascript": "**/*.{js,jsx}",
	"js":         "**/*.{js,jsx}",
	"python":     "**/*.py",
	"py":         "**/*.py",
	"rust":       "**/*.rs",
	"java":       "**/*.java",
	"ruby":       "**/*.rb",
	"css":        "**/*.css",
	"html":       "**/*.html",
	"sql":        "**/*.sql",
	"proto":      "**/*.proto",
	"yaml":       "**/*.{yaml,yml}",
	"json":       "**/*.json",
	"markdown":   "**/*.md",
	"md":         "**/*.md",
	"docker":     "**/Dockerfile*",
	"terraform":  "**/*.tf",
	"shell":      "**/*.sh",
	"sh":         "**/*.sh",
	"test":       "**/*_test.*",
	"testing":    "**/*_test.*",
}

// knownPaths maps slug keywords to directory path patterns.
// Path patterns take priority over extension patterns in substring matching.
var knownPaths = map[string]string{
	"frontend":  "frontend/**",
	"backend":   "backend/**",
	"api":       "backend/servers/**",
	"migration": "backend/**/migrations/**",
	"ci":        ".github/**",
}

// InferGlobs attempts to infer a file glob pattern from a rule slug.
// Returns the pattern and true if inference succeeded, or ("", false) if unknown.
//
// Matching priority:
//  1. Exact match in knownExtensions
//  2. Exact match in knownPaths
//  3. Longest substring match in knownPaths (deterministic)
//  4. Longest substring match in knownExtensions (deterministic)
func InferGlobs(slug string) (string, bool) {
	lower := strings.ToLower(slug)

	// Phase 1: exact match.
	if pattern, ok := knownExtensions[lower]; ok {
		return pattern, true
	}
	if pattern, ok := knownPaths[lower]; ok {
		return pattern, true
	}

	// Phase 2: longest substring match (deterministic via sorted keys).
	if pattern, ok := longestSubstringMatch(lower, knownPaths); ok {
		return pattern, true
	}
	if pattern, ok := longestSubstringMatch(lower, knownExtensions); ok {
		return pattern, true
	}

	return "", false
}

// longestSubstringMatch finds the keyword with the longest match in slug.
// When multiple keywords match with the same length, the lexicographically
// first keyword wins (deterministic via sorted iteration).
func longestSubstringMatch(slug string, m map[string]string) (string, bool) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	bestPattern := ""
	bestLen := 0
	for _, keyword := range keys {
		if strings.Contains(slug, keyword) && len(keyword) > bestLen {
			bestLen = len(keyword)
			bestPattern = m[keyword]
		}
	}

	if bestLen > 0 {
		return bestPattern, true
	}
	return "", false
}
