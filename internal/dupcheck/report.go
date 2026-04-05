package dupcheck

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type jsonOutput struct {
	Root         string       `json:"root"`
	SectionLevel int          `json:"section_level"`
	Threshold    float64      `json:"threshold"`
	MinChars     int          `json:"min_chars"`
	Excludes     []string     `json:"excludes,omitempty"`
	ResultCount  int          `json:"result_count"`
	Results      []jsonResult `json:"results"`
}

type jsonResult struct {
	Score float64  `json:"score"`
	A     Chunk    `json:"a"`
	B     Chunk    `json:"b"`
	Meta  jsonMeta `json:"meta"`
}

type jsonMeta struct {
	SameLastHeading bool `json:"same_last_heading"`
}

func toJSONOutput(cfg Config, results []Result) jsonOutput {
	jResults := make([]jsonResult, 0, len(results))
	for _, r := range results {
		jResults = append(jResults, jsonResult{
			Score: r.Score,
			A:     r.A,
			B:     r.B,
			Meta:  jsonMeta{SameLastHeading: sameLastHeading(r.A, r.B)},
		})
	}
	return jsonOutput{
		Root:         cfg.Root,
		SectionLevel: cfg.SectionLevel,
		Threshold:    cfg.Threshold,
		MinChars:     cfg.MinChars,
		Excludes:     cfg.Excludes,
		ResultCount:  len(results),
		Results:      jResults,
	}
}

func writeJSON(w io.Writer, out jsonOutput) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

var resultSeparator = strings.Repeat("-", 100)

func writeText(w io.Writer, results []Result) {
	for _, r := range results {
		_, _ = fmt.Fprintf(w, "score=%.3f\n", r.Score)
		_, _ = fmt.Fprintf(w, "A: %s [%s]\n", r.A.File, strings.Join(r.A.Headings, " > "))
		_, _ = fmt.Fprintf(w, "B: %s [%s]\n", r.B.File, strings.Join(r.B.Headings, " > "))
		_, _ = fmt.Fprintf(w, "A title: %s\n", r.A.Title)
		_, _ = fmt.Fprintf(w, "B title: %s\n", r.B.Title)
		_, _ = fmt.Fprintf(w, "A snippet: %s\n", truncate(r.A.Text, 160))
		_, _ = fmt.Fprintf(w, "B snippet: %s\n", truncate(r.B.Text, 160))
		_, _ = fmt.Fprintln(w, resultSeparator)
	}
}

func truncate(s string, limit int) string {
	r := []rune(s)
	if len(r) <= limit {
		return s
	}
	return string(r[:limit]) + "..."
}
