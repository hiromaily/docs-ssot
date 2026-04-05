package include

import (
	"path/filepath"
	"regexp"
	"strings"
)

// linkPattern matches Markdown inline links and images: [text](url) and ![alt](url).
// Group 1 captures the label (e.g. "[text]" or "![alt]"), group 2 the full destination
// (e.g. `./foo.md` or `./foo.md "My Title"`).
// Limitation: the URL pattern [^)]+ does not handle URLs that contain unescaped parentheses
// (e.g. Wikipedia URLs like `.../Markdown_(programming_language)`). Such links are left unchanged.
var linkPattern = regexp.MustCompile(`(!?\[[^\]]*\])\(([^)]+)\)`)

// absoluteURLPrefixes lists the prefixes that identify non-relative URLs.
var absoluteURLPrefixes = []string{"http://", "https://", "//", "/", "#", "mailto:", "ftp:", "data:", "tel:"}

// rewriteLinksInDirs rewrites relative Markdown link and image URLs in line so they are
// correct relative to outputDir rather than sourceDir.
// sourceDir and outputDir must be pre-computed absolute paths by the caller.
// Lines inside code fences must be excluded by the caller.
func rewriteLinksInDirs(line, sourceDir, outputDir string) string {
	if sourceDir == outputDir {
		return line
	}
	return linkPattern.ReplaceAllStringFunc(line, func(match string) string {
		// Re-extract submatches from the already-matched substring.
		// FindStringSubmatch on a string that already matched the pattern always succeeds.
		parts := linkPattern.FindStringSubmatch(match)
		label := parts[1] // e.g. "[text]" or "![alt]"
		dest := parts[2]  // e.g. `./foo.md` or `./foo.md "My Title"`

		// Split optional link title: `./foo.md "My Title"` → url=`./foo.md`, title=` "My Title"`
		url, title := dest, ""
		if idx := strings.IndexByte(dest, ' '); idx != -1 {
			url, title = dest[:idx], dest[idx:]
		}

		if !isRelativeURL(url) {
			return match
		}

		// Split off fragment: "file.md#section" → urlPath="file.md", fragment="#section"
		urlPath, fragment := url, ""
		if idx := strings.Index(url, "#"); idx != -1 {
			urlPath, fragment = url[:idx], url[idx:]
		}
		if urlPath == "" {
			return match // pure anchor (#section) — no path to rewrite
		}

		absTarget := filepath.Join(sourceDir, urlPath)
		newRel, err := filepath.Rel(outputDir, absTarget)
		if err != nil {
			return match // fallback: leave unchanged
		}
		newRel = filepath.ToSlash(newRel)
		if !strings.HasPrefix(newRel, ".") {
			newRel = "./" + newRel
		}
		return label + "(" + newRel + fragment + title + ")"
	})
}

// isRelativeURL reports whether url is a relative path that should be rewritten.
// Absolute URLs, absolute paths, and pure anchors are left unchanged.
func isRelativeURL(url string) bool {
	for _, prefix := range absoluteURLPrefixes {
		if strings.HasPrefix(url, prefix) {
			return false
		}
	}
	return true
}
