package agentscan

import "strings"

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
var knownPaths = map[string]string{
	"frontend":  "frontend/**",
	"backend":   "backend/**",
	"app-web":   "frontend/app-web/**",
	"api":       "backend/servers/**",
	"migration": "backend/**/migrations/**",
	"ci":        ".github/**",
}

// InferGlobs attempts to infer a file glob pattern from a rule slug.
// Returns the pattern and true if inference succeeded, or ("", false) if unknown.
func InferGlobs(slug string) (string, bool) {
	lower := strings.ToLower(slug)

	// Check exact match first.
	if pattern, ok := knownExtensions[lower]; ok {
		return pattern, true
	}
	if pattern, ok := knownPaths[lower]; ok {
		return pattern, true
	}

	// Check if slug contains a known keyword.
	for keyword, pattern := range knownPaths {
		if strings.Contains(lower, keyword) {
			return pattern, true
		}
	}
	for keyword, pattern := range knownExtensions {
		if strings.Contains(lower, keyword) {
			return pattern, true
		}
	}

	return "", false
}
