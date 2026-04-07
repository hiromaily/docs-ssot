// Package categoriser assigns Markdown sections to directory categories
// based on heading keyword heuristics.
package categoriser

import (
	"strings"
)

// rule maps keyword patterns to a category directory name.
type rule struct {
	keywords []string
	category string
}

// rules are evaluated in order; the first match wins.
var rules = []rule{
	{keywords: []string{"architecture", "design", "system", "pipeline"}, category: "architecture"},
	{keywords: []string{"overview", "about", "introduction", "background"}, category: "project"},
	{keywords: []string{"install", "setup", "getting started", "prerequisites"}, category: "development"},
	{keywords: []string{"test", "testing", "ci"}, category: "development"},
	{keywords: []string{"lint", "format", "code quality"}, category: "development"},
	{keywords: []string{"api", "commands", "cli", "reference"}, category: "reference"},
	{keywords: []string{"contributing", "contribute"}, category: "development"},
	{keywords: []string{"license", "changelog", "roadmap"}, category: "project"},
	{keywords: []string{"faq", "troubleshooting"}, category: "product"},
}

// Categorise returns the category directory name for a given section heading title.
// If no heuristic matches, it returns "misc".
func Categorise(title string) string {
	lower := strings.ToLower(title)

	for _, r := range rules {
		for _, kw := range r.keywords {
			if strings.Contains(lower, kw) {
				return r.category
			}
		}
	}

	return "misc"
}
