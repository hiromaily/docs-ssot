package dupcheck

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strings"
)

// Config holds parameters for the duplicate section checker.
type Config struct {
	Root         string
	Threshold    float64
	MinChars     int
	SectionLevel int
	Format       string
	Excludes     []string
}

// Result holds a pair of similar sections and their similarity score.
type Result struct {
	A     Chunk
	B     Chunk
	Score float64
}

// Run scans the docs directory for near-duplicate sections and writes results to w.
// It returns an error if file scanning or JSON encoding fails.
func Run(w io.Writer, cfg Config) error {
	files, err := markdownFiles(cfg.Root, cfg.Excludes)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	var chunks []Chunk
	for _, path := range files {
		extracted, err := extractSectionChunks(path, cfg.MinChars, cfg.SectionLevel)
		if err != nil {
			_, _ = fmt.Fprintf(w, "warning: parse failed: %s: %v\n", path, err)
			continue
		}
		chunks = append(chunks, extracted...)
	}

	if len(chunks) == 0 {
		if cfg.Format == "json" {
			return writeJSON(w, toJSONOutput(cfg, nil))
		}
		_, _ = fmt.Fprintln(w, "no section chunks found")
		return nil
	}

	vectors := buildTFIDF(chunks)

	var results []Result
	for i := range len(chunks) {
		for j := i + 1; j < len(chunks); j++ {
			if chunks[i].File == chunks[j].File {
				continue
			}
			score := cosine(vectors[i], vectors[j])
			if sameLastHeading(chunks[i], chunks[j]) {
				score = min(score+0.03, 1.0)
			}
			if score >= cfg.Threshold {
				results = append(results, Result{A: chunks[i], B: chunks[j], Score: score})
			}
		}
	}

	slices.SortFunc(results, func(a, b Result) int {
		if s := cmp.Compare(b.Score, a.Score); s != 0 { // descending by score
			return s
		}
		if s := strings.Compare(a.A.File, b.A.File); s != 0 {
			return s
		}
		return cmp.Compare(a.A.Index, b.A.Index)
	})

	switch cfg.Format {
	case "json":
		return writeJSON(w, toJSONOutput(cfg, results))
	default:
		writeText(w, results)
		return nil
	}
}

func sameLastHeading(a, b Chunk) bool {
	if len(a.Headings) == 0 || len(b.Headings) == 0 {
		return false
	}
	return strings.EqualFold(a.Headings[len(a.Headings)-1], b.Headings[len(b.Headings)-1])
}
