package dupcheck

import (
	"bytes"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (
	multiSpaceRe = regexp.MustCompile(`\s+`)
	// mdParser is a single shared parser instance; goldmark.Parser is safe to reuse across calls.
	mdParser = goldmark.New().Parser()
)

// extractSectionChunks parses a Markdown file and extracts sections at the
// specified heading level. Each section includes all content until the next
// heading at the same level or higher. Sections shorter than minChars are skipped.
func extractSectionChunks(path string, minChars, sectionLevel int) ([]Chunk, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	doc := mdParser.Parse(text.NewReader(src))

	var chunks []Chunk
	var headingStack []string
	index := 0

	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		h, ok := n.(*ast.Heading)
		if !ok {
			continue
		}

		title := normalizeText(extractInlineText(h, src))
		if title == "" {
			continue
		}

		headingStack = adjustHeadings(headingStack, h.Level, title)

		if h.Level != sectionLevel {
			continue
		}

		sectionHeadings := slices.Clone(headingStack)
		parts := []string{title}

		for cur := n.NextSibling(); cur != nil; cur = cur.NextSibling() {
			if nextH, ok := cur.(*ast.Heading); ok && nextH.Level <= sectionLevel {
				break
			}
			if isSkippableBlock(cur) {
				continue
			}
			txt := normalizeText(extractBlockText(cur, src))
			if txt != "" {
				parts = append(parts, txt)
			}
		}

		fullText := normalizeText(strings.Join(parts, "\n"))
		if len([]rune(fullText)) < minChars {
			continue
		}

		tokens := tokenize(fullText)
		if len(tokens) == 0 {
			continue
		}

		chunks = append(chunks, Chunk{
			File:     path,
			Index:    index,
			Kind:     "section",
			Title:    title,
			Text:     fullText,
			Tokens:   tokens,
			Headings: sectionHeadings,
			Level:    h.Level,
		})
		index++
	}

	return chunks, nil
}

func isSkippableBlock(n ast.Node) bool {
	switch n.(type) {
	case *ast.FencedCodeBlock, *ast.CodeBlock, *ast.HTMLBlock:
		return true
	default:
		return false
	}
}

func extractBlockText(n ast.Node, src []byte) string {
	switch node := n.(type) {
	case *ast.Heading:
		return extractInlineText(node, src)
	case *ast.Paragraph:
		return extractInlineText(node, src)
	case *ast.Blockquote:
		var parts []string
		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			if isSkippableBlock(c) {
				continue
			}
			txt := normalizeText(extractBlockText(c, src))
			if txt != "" {
				parts = append(parts, txt)
			}
		}
		return strings.Join(parts, "\n")
	case *ast.List:
		var parts []string
		for item := node.FirstChild(); item != nil; item = item.NextSibling() {
			txt := normalizeText(extractListItemText(item, src))
			if txt != "" {
				parts = append(parts, txt)
			}
		}
		return strings.Join(parts, "\n")
	case *ast.ListItem:
		return extractListItemText(node, src)
	default:
		var parts []string
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if isSkippableBlock(c) {
				continue
			}
			txt := normalizeText(extractBlockText(c, src))
			if txt != "" {
				parts = append(parts, txt)
			}
		}
		return strings.Join(parts, "\n")
	}
}

func extractInlineText(n ast.Node, src []byte) string {
	var buf bytes.Buffer

	var visit func(ast.Node)
	visit = func(node ast.Node) {
		switch t := node.(type) {
		case *ast.Text:
			buf.Write(t.Segment.Value(src))
			if t.HardLineBreak() || t.SoftLineBreak() {
				buf.WriteByte(' ')
			}
		case *ast.CodeSpan:
			for c := t.FirstChild(); c != nil; c = c.NextSibling() {
				if tx, ok := c.(*ast.Text); ok {
					buf.Write(tx.Segment.Value(src))
				}
			}
			buf.WriteByte(' ')
		case *ast.Link, *ast.Emphasis:
			for c := node.FirstChild(); c != nil; c = c.NextSibling() {
				visit(c)
			}
		default:
			for c := node.FirstChild(); c != nil; c = c.NextSibling() {
				visit(c)
			}
		}
	}

	visit(n)
	return buf.String()
}

func extractListItemText(n ast.Node, src []byte) string {
	var parts []string
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if isSkippableBlock(c) {
			continue
		}
		txt := normalizeText(extractBlockText(c, src))
		if txt != "" {
			parts = append(parts, txt)
		}
	}
	return strings.Join(parts, " ")
}

func adjustHeadings(curr []string, level int, title string) []string {
	if level <= 0 {
		return curr
	}
	keep := min(level-1, len(curr))
	return append(slices.Clone(curr[:keep]), title)
}

func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\u00A0", " ")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = multiSpaceRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func tokenize(s string) []string {
	s = strings.ToLower(s)

	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r):
			b.WriteRune(r)
		case unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Han):
			b.WriteRune(r)
		default:
			b.WriteRune(' ')
		}
	}

	// strings.Fields handles consecutive whitespace, so no regex pass needed.
	fields := strings.Fields(b.String())

	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if utf8.RuneCountInString(f) <= 1 {
			continue
		}
		out = append(out, f)
	}
	return out
}

// Ensure the goldmark parser interface is satisfied at compile time.
var _ parser.Parser = mdParser
