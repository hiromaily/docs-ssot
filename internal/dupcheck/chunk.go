// Package dupcheck detects near-duplicate sections in Markdown documentation,
// helping identify SSOT violations where the same content exists in multiple places.
package dupcheck

// Chunk represents a single extracted section from a Markdown file.
type Chunk struct {
	File     string   `json:"file"`
	Index    int      `json:"index"`
	Kind     string   `json:"kind"`
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	Tokens   []string `json:"-"`
	Headings []string `json:"headings"`
	Level    int      `json:"level"`
}
