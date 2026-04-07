// Package splitter parses Markdown files and splits them into sections by heading level.
// Unlike dupcheck (which uses goldmark AST for text extraction), splitter preserves
// the raw Markdown content of each section for faithful round-trip reproduction.
package splitter

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hiromaily/docs-ssot/internal/mdutil"
)

// headingRe matches ATX-style headings: up to 3 spaces of leading indentation
// (per CommonMark spec), one or more # followed by a space.
var headingRe = regexp.MustCompile(`^\s{0,3}(#{1,6})\s+(.+)$`)

// Section represents a single section extracted from a Markdown file.
type Section struct {
	// Title is the heading text (without the # prefix).
	Title string
	// Level is the heading level (1–6).
	Level int
	// Body is the raw Markdown content below the heading (excluding the heading line itself).
	Body string
	// RawHeading is the original heading line including the # prefix.
	RawHeading string
}

// Split reads a Markdown file and splits it into sections at the given heading level.
// Content before the first heading at sectionLevel is returned as a section with an empty title.
// Each section includes all content until the next heading at sectionLevel or higher.
func Split(path string, sectionLevel int) ([]Section, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	sections, err := SplitString(string(data), sectionLevel)
	if err != nil {
		return nil, fmt.Errorf("split %s: %w", path, err)
	}

	return sections, nil
}

// SplitString splits Markdown content (as a string) into sections at the given heading level.
func SplitString(content string, sectionLevel int) ([]Section, error) {
	var sections []Section
	var current *Section
	var bodyLines []string
	fenceType := "" // "" = outside fence, otherwise the opening fence marker

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Track code fence state using the shared CommonMark-compliant utility.
		fenceType = mdutil.NextFenceType(line, fenceType)
		inCodeFence := fenceType != ""

		if inCodeFence {
			bodyLines = append(bodyLines, line)
			continue
		}

		m := headingRe.FindStringSubmatch(line)
		if m != nil {
			level := len(m[1])
			title := strings.TrimSpace(m[2])

			if level <= sectionLevel {
				// Flush the current section.
				if current != nil {
					current.Body = joinBody(bodyLines)
					sections = append(sections, *current)
				} else if len(bodyLines) > 0 {
					// Content before the first section heading.
					sections = append(sections, Section{
						Title: "",
						Level: 0,
						Body:  joinBody(bodyLines),
					})
				}

				current = &Section{
					Title:      title,
					Level:      level,
					RawHeading: line,
				}
				bodyLines = nil
				continue
			}
		}

		bodyLines = append(bodyLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	// Flush the last section.
	if current != nil {
		current.Body = joinBody(bodyLines)
		sections = append(sections, *current)
	} else if len(bodyLines) > 0 {
		sections = append(sections, Section{
			Title: "",
			Level: 0,
			Body:  joinBody(bodyLines),
		})
	}

	return sections, nil
}

// joinBody joins lines into a single string, removing leading/trailing blank lines
// while preserving internal blank lines and per-line indentation.
func joinBody(lines []string) string {
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}
