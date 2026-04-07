package splitter_test

import (
	"testing"

	"github.com/hiromaily/docs-ssot/internal/splitter"
)

func TestSplitString_BasicSections(t *testing.T) {
	t.Parallel()

	content := `# Title

Intro text.

## Setup

Install instructions here.

## Testing

Test instructions.

### Unit Tests

Details about unit tests.

## Architecture

Arch overview.
`

	sections, err := splitter.SplitString(content, 2)
	if err != nil {
		t.Fatalf("SplitString() error: %v", err)
	}

	// Expect: "Title" (H1, treated as section boundary since level<=2),
	// "Setup", "Testing" (includes H3 subsection), "Architecture".
	if len(sections) != 4 {
		t.Fatalf("got %d sections, want 4", len(sections))
	}

	if sections[0].Title != "Title" {
		t.Errorf("section 0 title = %q, want %q", sections[0].Title, "Title")
	}
	if sections[0].Level != 1 {
		t.Errorf("section 0 level = %d, want 1", sections[0].Level)
	}
	if sections[0].Body != "Intro text." {
		t.Errorf("section 0 body = %q, want %q", sections[0].Body, "Intro text.")
	}

	if sections[1].Title != "Setup" {
		t.Errorf("section 1 title = %q, want %q", sections[1].Title, "Setup")
	}

	if sections[2].Title != "Testing" {
		t.Errorf("section 2 title = %q, want %q", sections[2].Title, "Testing")
	}

	if sections[3].Title != "Architecture" {
		t.Errorf("section 3 title = %q, want %q", sections[3].Title, "Architecture")
	}
}

func TestSplitString_ContentBeforeFirstHeading(t *testing.T) {
	t.Parallel()

	content := `Some preamble text.

## First Section

Content.
`

	sections, err := splitter.SplitString(content, 2)
	if err != nil {
		t.Fatalf("SplitString() error: %v", err)
	}

	if len(sections) != 2 {
		t.Fatalf("got %d sections, want 2", len(sections))
	}

	if sections[0].Title != "" {
		t.Errorf("preamble section title = %q, want empty", sections[0].Title)
	}
	if sections[0].Body != "Some preamble text." {
		t.Errorf("preamble body = %q, want %q", sections[0].Body, "Some preamble text.")
	}

	if sections[1].Title != "First Section" {
		t.Errorf("section 1 title = %q, want %q", sections[1].Title, "First Section")
	}
}

func TestSplitString_CodeFencesPreserved(t *testing.T) {
	t.Parallel()

	content := "## Example\n\n```markdown\n## Not a heading\n```\n\n## Next\n\nContent.\n"

	sections, err := splitter.SplitString(content, 2)
	if err != nil {
		t.Fatalf("SplitString() error: %v", err)
	}

	if len(sections) != 2 {
		t.Fatalf("got %d sections, want 2", len(sections))
	}

	if sections[0].Title != "Example" {
		t.Errorf("section 0 title = %q, want %q", sections[0].Title, "Example")
	}
	if sections[1].Title != "Next" {
		t.Errorf("section 1 title = %q, want %q", sections[1].Title, "Next")
	}
}

func TestSplitString_NestedCodeFences(t *testing.T) {
	t.Parallel()

	// 4-backtick fence should not be closed by 3-backtick line.
	content := "## Example\n\n````markdown\n```\n## Not a heading\n```\n````\n\n## Next\n\nContent.\n"

	sections, err := splitter.SplitString(content, 2)
	if err != nil {
		t.Fatalf("SplitString() error: %v", err)
	}

	if len(sections) != 2 {
		t.Fatalf("got %d sections, want 2", len(sections))
	}

	if sections[0].Title != "Example" {
		t.Errorf("section 0 title = %q, want %q", sections[0].Title, "Example")
	}
	if sections[1].Title != "Next" {
		t.Errorf("section 1 title = %q, want %q", sections[1].Title, "Next")
	}
}

func TestSplitString_Empty(t *testing.T) {
	t.Parallel()

	sections, err := splitter.SplitString("", 2)
	if err != nil {
		t.Fatalf("SplitString() error: %v", err)
	}

	if len(sections) != 0 {
		t.Fatalf("got %d sections, want 0", len(sections))
	}
}
