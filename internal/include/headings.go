package include

import (
	"bufio"
	"regexp"
	"strings"
)

// headingPattern matches CommonMark ATX headings: up to 3 optional leading spaces,
// then 1–6 '#' characters, followed by a space/tab or end of line.
// Group 1 captures leading indent; group 2 captures the '#' run.
var headingPattern = regexp.MustCompile(`^( {0,3})(#{1,6})(?:[ \t]|$)`)

// adjustHeadingLevels shifts every ATX heading in content by delta levels.
// Headings inside fenced code blocks are not adjusted.
// Positive delta deepens headings (# → ## for delta=1).
// Negative delta shallows headings (## → # for delta=-1).
// Levels are clamped to [1, 6].
func adjustHeadingLevels(content string, delta int) string {
	if delta == 0 {
		return content
	}
	var sb strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	fenceType := ""
	for scanner.Scan() {
		line := scanner.Text()
		prevFenceType := fenceType
		fenceType = nextFenceType(line, fenceType)
		if prevFenceType == "" && fenceType == "" {
			line = shiftHeading(line, delta)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

// shiftHeading adjusts the ATX heading level of line by delta.
// Returns line unchanged if it is not a heading.
func shiftHeading(line string, delta int) string {
	m := headingPattern.FindStringSubmatchIndex(line)
	if m == nil {
		return line
	}
	indent := line[m[2]:m[3]]
	hashes := line[m[4]:m[5]]
	rest := line[m[5]:]

	newLevel := len(hashes) + delta
	if newLevel < 1 {
		newLevel = 1
	} else if newLevel > 6 {
		newLevel = 6
	}
	return indent + strings.Repeat("#", newLevel) + rest
}
