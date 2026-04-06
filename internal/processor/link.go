package processor

import (
	"bufio"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hiromaily/docs-ssot/internal/mdutil"
)

// linkPattern matches Markdown inline links and images: [text](url) and ![alt](url).
// Group 1 captures the label (e.g. "[text]" or "![alt]"), group 2 the full destination
// (e.g. `./foo.md` or `./foo.md "My Title"`).
// Known limitations:
//   - The URL pattern [^)]+ does not handle URLs containing unescaped parentheses
//     (e.g. Wikipedia URLs like `.../Markdown_(programming_language)`). Such links are left unchanged.
//   - Link titles that themselves contain a closing parenthesis (e.g. `[label](url "title (info)")`)
//     will be incorrectly truncated at the first `)`. Such links are left unchanged or partially rewritten.
//   - Reference-style links (`[text][ref]`) are not matched and are left unchanged.
//
// A full CommonMark-compliant parser would be required to handle all edge cases; the added complexity
// is not justified for this tool's primary use cases.
var linkPattern = regexp.MustCompile(`(!?\[[^\]]*\])\(([^)]+)\)`)

// absoluteURLPrefixes lists the prefixes that identify non-relative URLs.
var absoluteURLPrefixes = []string{"http://", "https://", "//", "/", "#", "mailto:", "ftp:", "data:", "tel:"}

// LinkTransformer rewrites relative Markdown link and image URLs so they are
// correct relative to OutputDir rather than SourceDir.
// Links inside fenced code blocks are not rewritten.
type LinkTransformer struct {
	SourceDir string
	OutputDir string
}

// Transform implements Transformer.
func (l LinkTransformer) Transform(content string) (string, error) {
	if l.SourceDir == l.OutputDir {
		return content, nil
	}
	var sb strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(nil, 1024*1024)
	fenceType := ""
	for scanner.Scan() {
		line := scanner.Text()
		prevFenceType := fenceType
		fenceType = mdutil.NextFenceType(line, fenceType)
		if prevFenceType == "" && fenceType == "" && strings.ContainsAny(line, "[!") {
			line = rewriteLinksInLine(line, l.SourceDir, l.OutputDir)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String(), scanner.Err()
}

// rewriteLinksInLine rewrites relative Markdown link and image URLs in line so they are
// correct relative to outputDir rather than sourceDir.
// sourceDir and outputDir must be pre-computed absolute paths, and must differ.
// Lines inside code fences must be excluded by the caller.
func rewriteLinksInLine(line, sourceDir, outputDir string) string {
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
