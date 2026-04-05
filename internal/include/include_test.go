package include_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hiromaily/docs-ssot/internal/include"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestProcessFile_SingleInclude(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	childPath := filepath.Join(dir, "child.md")
	writeFile(t, childPath, "child content\n")
	writeFile(t, filepath.Join(dir, "root.md"), fmt.Sprintf("before\n<!-- @include: %s -->\nafter\n", childPath))

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "before\nchild content\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_RecursiveInclude(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	cPath := filepath.Join(dir, "c.md")
	bPath := filepath.Join(dir, "b.md")
	writeFile(t, cPath, "level-c\n")
	writeFile(t, bPath, fmt.Sprintf("level-b\n<!-- @include: %s -->\n", cPath))
	writeFile(t, filepath.Join(dir, "a.md"), fmt.Sprintf("level-a\n<!-- @include: %s -->\n", bPath))

	got, err := include.ProcessFile(filepath.Join(dir, "a.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "level-a\nlevel-b\nlevel-c\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_CircularInclude(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	aPath := filepath.Join(dir, "a.md")
	bPath := filepath.Join(dir, "b.md")
	writeFile(t, aPath, fmt.Sprintf("<!-- @include: %s -->\n", bPath))
	writeFile(t, bPath, fmt.Sprintf("<!-- @include: %s -->\n", aPath))

	_, err := include.ProcessFile(aPath, filepath.Join(dir, "output.md"))
	if err == nil {
		t.Fatal("expected circular include error, got nil")
	}
}

func TestProcessFile_MissingFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	missingPath := filepath.Join(dir, "missing.md")
	writeFile(t, filepath.Join(dir, "root.md"), fmt.Sprintf("<!-- @include: %s -->\n", missingPath))

	_, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err == nil {
		t.Fatal("expected error for missing include, got nil")
	}
}

func TestProcessFile_NoIncludes(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "plain.md"), "just text\nno includes\n")

	got, err := include.ProcessFile(filepath.Join(dir, "plain.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "just text\nno includes\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_IncludeInsideCodeFence(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	childPath := filepath.Join(dir, "child.md")
	// child.md does NOT exist — if the include inside the fence were expanded, it would error
	writeFile(t, filepath.Join(dir, "root.md"), fmt.Sprintf("```md\n<!-- @include: %s -->\n```\n", childPath))

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatalf("include inside code fence should not be expanded, got error: %v", err)
	}

	want := fmt.Sprintf("```md\n<!-- @include: %s -->\n```\n", childPath)
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_RelativeIncludePath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "sub", "child.md"), "child content\n")
	// Use a relative path from the root file's directory
	writeFile(t, filepath.Join(dir, "root.md"), "before\n<!-- @include: sub/child.md -->\nafter\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "before\nchild content\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_MixedFenceTypes(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// A tilde fence should NOT be closed by backticks — the include inside remains literal.
	// child.md does NOT exist; if the include were expanded it would error.
	writeFile(t, filepath.Join(dir, "root.md"), "~~~\n<!-- @include: missing.md -->\n~~~\nafter\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatalf("include inside tilde fence should not be expanded, got error: %v", err)
	}

	want := "~~~\n<!-- @include: missing.md -->\n~~~\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_IncludeInlineNotExpanded(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// Directive embedded within other text should NOT be expanded (regex is anchored).
	writeFile(t, filepath.Join(dir, "root.md"), "Note: <!-- @include: missing.md --> end\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatalf("inline directive should not be expanded, got error: %v", err)
	}

	want := "Note: <!-- @include: missing.md --> end\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestProcessFile_LevelAdjustment_RecursiveExpansion verifies that level=+N shifts headings
// in the entire recursively-expanded content of the included file, not just its top level.
func TestProcessFile_LevelAdjustment_RecursiveExpansion(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// c.md is included by b.md (no level param), b.md is included by root with level=+1.
	// All headings in the expanded B+C content should be shifted by +1.
	writeFile(t, filepath.Join(dir, "c.md"), "### Deep\n")
	writeFile(t, filepath.Join(dir, "b.md"), "## Section\n<!-- @include: c.md -->\n")
	writeFile(t, filepath.Join(dir, "root.md"), "<!-- @include: b.md level=+1 -->\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "### Section\n#### Deep\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestProcessFile_LevelAdjustment_WithLinkRewrite verifies that heading level adjustment
// and link rewriting do not interfere with each other.
func TestProcessFile_LevelAdjustment_WithLinkRewrite(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	docsDir := filepath.Join(dir, "docs")
	writeFile(t, filepath.Join(docsDir, "guide.md"), "## Title\n\nSee [ref](./other.md).\n")
	writeFile(t, filepath.Join(dir, "root.md"),
		fmt.Sprintf("<!-- @include: %s level=+1 -->\n", filepath.Join(docsDir, "guide.md")))

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	// heading shifted +1; link rewritten relative to output dir
	want := "### Title\n\nSee [ref](./docs/other.md).\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestProcessFile_LinkRewrite_DirectFile tests that a link inside the template file itself
// is rewritten when the template and output live in different directories.
func TestProcessFile_LinkRewrite_DirectFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// template in dir/template/, output at dir/output.md
	templateDir := filepath.Join(dir, "template")
	writeFile(t, filepath.Join(templateDir, "root.tpl.md"), "[guide](./guide.md)\n")

	got, err := include.ProcessFile(
		filepath.Join(templateDir, "root.tpl.md"),
		filepath.Join(dir, "output.md"),
	)
	if err != nil {
		t.Fatal(err)
	}

	// "./guide.md" relative to template/ → dir/template/guide.md → relative to dir/ → "./template/guide.md"
	want := "[guide](./template/guide.md)\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestProcessFile_LinkRewrite_IncludedFile tests that links inside an included file
// are rewritten to be correct relative to the output file.
func TestProcessFile_LinkRewrite_IncludedFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// included file in dir/docs/sub/, output at dir/output.md
	docsDir := filepath.Join(dir, "docs", "sub")
	includedPath := filepath.Join(docsDir, "section.md")
	writeFile(t, includedPath, "See [bar](./bar.md) and ![img](./img.png).\n")
	writeFile(t, filepath.Join(dir, "root.tpl.md"),
		fmt.Sprintf("# Title\n<!-- @include: %s -->\n", includedPath))

	got, err := include.ProcessFile(
		filepath.Join(dir, "root.tpl.md"),
		filepath.Join(dir, "output.md"),
	)
	if err != nil {
		t.Fatal(err)
	}

	// links in section.md resolve to dir/docs/sub/{bar.md,img.png},
	// relative to dir/ (output dir) → "./docs/sub/bar.md", "./docs/sub/img.png"
	want := "# Title\nSee [bar](./docs/sub/bar.md) and ![img](./docs/sub/img.png).\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestProcessFile_LinkRewrite_AbsoluteAndExternalUnchanged tests that absolute URLs,
// absolute paths, and pure anchors are not rewritten.
func TestProcessFile_LinkRewrite_AbsoluteAndExternalUnchanged(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	content := "[ext](https://example.com) [abs](/abs/path.md) [anchor](#section) [relative](./keep.md)\n"
	writeFile(t, filepath.Join(dir, "sub", "file.md"), content)
	writeFile(t, filepath.Join(dir, "root.tpl.md"),
		"<!-- @include: sub/file.md -->\n")

	got, err := include.ProcessFile(
		filepath.Join(dir, "root.tpl.md"),
		filepath.Join(dir, "output.md"),
	)
	if err != nil {
		t.Fatal(err)
	}

	// only the relative link changes; others are untouched
	want := "[ext](https://example.com) [abs](/abs/path.md) [anchor](#section) [relative](./sub/keep.md)\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestProcessFile_LinkRewrite_WithFragment tests that URL fragments are preserved after rewriting.
func TestProcessFile_LinkRewrite_WithFragment(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "docs", "guide.md"), "[section](./other.md#heading)\n")
	writeFile(t, filepath.Join(dir, "root.tpl.md"),
		"<!-- @include: docs/guide.md -->\n")

	got, err := include.ProcessFile(
		filepath.Join(dir, "root.tpl.md"),
		filepath.Join(dir, "output.md"),
	)
	if err != nil {
		t.Fatal(err)
	}

	want := "[section](./docs/other.md#heading)\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestProcessFile_LinkRewrite_InsideCodeFenceUnchanged tests that links inside code fences
// are not rewritten.
func TestProcessFile_LinkRewrite_InsideCodeFenceUnchanged(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "docs", "guide.md"),
		"```\n[link](./example.md)\n```\n")
	writeFile(t, filepath.Join(dir, "root.tpl.md"),
		"<!-- @include: docs/guide.md -->\n")

	got, err := include.ProcessFile(
		filepath.Join(dir, "root.tpl.md"),
		filepath.Join(dir, "output.md"),
	)
	if err != nil {
		t.Fatal(err)
	}

	// link inside code fence must NOT be rewritten
	want := "```\n[link](./example.md)\n```\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_LevelAdjustment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
		param   string
		want    string
	}{
		{
			name:    "plus1",
			content: "## Title\n### Sub1\n### Sub2\n",
			param:   "level=+1",
			want:    "### Title\n#### Sub1\n#### Sub2\n",
		},
		{
			name:    "minus1",
			content: "## Title\n### Sub1\n",
			param:   "level=-1",
			want:    "# Title\n## Sub1\n",
		},
		{
			name:    "plus2",
			content: "# Title\n## Sub\n",
			param:   "level=+2",
			want:    "### Title\n#### Sub\n",
		},
		{
			name:    "zero_explicit",
			content: "## Title\n### Sub\n",
			param:   "level=0",
			want:    "## Title\n### Sub\n",
		},
		{
			name:    "no_param",
			content: "## Title\n### Sub\n",
			param:   "",
			want:    "## Title\n### Sub\n",
		},
		{
			name:    "clamp_at_h6",
			content: "###### Deep\n",
			param:   "level=+1",
			want:    "###### Deep\n",
		},
		{
			name:    "clamp_at_h1",
			content: "# Top\n",
			param:   "level=-1",
			want:    "# Top\n",
		},
		{
			name:    "non_heading_unchanged",
			content: "just text\n- list item\n",
			param:   "level=+1",
			want:    "just text\n- list item\n",
		},
		{
			name:    "heading_in_code_fence_unchanged",
			content: "# Outside\n```\n# Inside fence\n```\n",
			param:   "level=+1",
			want:    "## Outside\n```\n# Inside fence\n```\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, "child.md"), tc.content)
			directive := "<!-- @include: child.md"
			if tc.param != "" {
				directive += " " + tc.param
			}
			directive += " -->"
			writeFile(t, filepath.Join(dir, "root.md"), directive+"\n")

			got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tc.want)
			}
		})
	}
}

func TestProcessFile_DirectoryInclude(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	subDir := filepath.Join(dir, "docs")
	writeFile(t, filepath.Join(subDir, "01_first.md"), "first\n")
	writeFile(t, filepath.Join(subDir, "02_second.md"), "second\n")
	writeFile(t, filepath.Join(subDir, "03_third.md"), "third\n")
	// a non-.md file that should be ignored
	if err := os.WriteFile(filepath.Join(subDir, "skip.txt"), []byte("skip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "root.md"), "before\n<!-- @include: docs/ -->\nafter\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "before\nfirst\nsecond\nthird\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_DirectoryInclude_SortedOrder(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	subDir := filepath.Join(dir, "docs")
	// write in reverse order to verify sorting is by filename, not creation order
	writeFile(t, filepath.Join(subDir, "b.md"), "b\n")
	writeFile(t, filepath.Join(subDir, "a.md"), "a\n")
	writeFile(t, filepath.Join(subDir, "c.md"), "c\n")
	writeFile(t, filepath.Join(dir, "root.md"), "<!-- @include: docs/ -->\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "a\nb\nc\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_DirectoryInclude_MissingDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "root.md"), "<!-- @include: nonexistent/ -->\n")

	_, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
}

func TestProcessFile_DirectoryInclude_EmptyDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "root.md"), "before\n<!-- @include: empty/ -->\nafter\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "before\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_DirectoryInclude_WithLevelAdjustment(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	subDir := filepath.Join(dir, "docs")
	writeFile(t, filepath.Join(subDir, "a.md"), "## Section A\n")
	writeFile(t, filepath.Join(subDir, "b.md"), "## Section B\n")
	writeFile(t, filepath.Join(dir, "root.md"), "<!-- @include: docs/ level=+1 -->\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "### Section A\n### Section B\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_GlobInclude(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	subDir := filepath.Join(dir, "docs")
	writeFile(t, filepath.Join(subDir, "01_first.md"), "first\n")
	writeFile(t, filepath.Join(subDir, "02_second.md"), "second\n")
	writeFile(t, filepath.Join(subDir, "03_third.md"), "third\n")
	// a non-.md file that should not be matched
	if err := os.WriteFile(filepath.Join(subDir, "skip.txt"), []byte("skip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "root.md"), "before\n<!-- @include: docs/*.md -->\nafter\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "before\nfirst\nsecond\nthird\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_GlobInclude_SortedOrder(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	subDir := filepath.Join(dir, "docs")
	// write in reverse order to verify sorting is lexical, not creation order
	writeFile(t, filepath.Join(subDir, "c.md"), "c\n")
	writeFile(t, filepath.Join(subDir, "a.md"), "a\n")
	writeFile(t, filepath.Join(subDir, "b.md"), "b\n")
	writeFile(t, filepath.Join(dir, "root.md"), "<!-- @include: docs/*.md -->\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "a\nb\nc\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_GlobInclude_NoMatches(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "root.md"), "before\n<!-- @include: docs/*.md -->\nafter\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	// no matches → empty expansion, surrounding text preserved
	want := "before\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_GlobInclude_WithLevelAdjustment(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	subDir := filepath.Join(dir, "docs")
	writeFile(t, filepath.Join(subDir, "a.md"), "## Section A\n")
	writeFile(t, filepath.Join(subDir, "b.md"), "## Section B\n")
	writeFile(t, filepath.Join(dir, "root.md"), "<!-- @include: docs/*.md level=+1 -->\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"), filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "### Section A\n### Section B\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestProcessFile_LinkRewrite_WithTitle tests that optional link titles are preserved.
func TestProcessFile_LinkRewrite_WithTitle(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "docs", "guide.md"), `[foo](./bar.md "My Title")`+"\n")
	writeFile(t, filepath.Join(dir, "root.tpl.md"),
		"<!-- @include: docs/guide.md -->\n")

	got, err := include.ProcessFile(
		filepath.Join(dir, "root.tpl.md"),
		filepath.Join(dir, "output.md"),
	)
	if err != nil {
		t.Fatal(err)
	}

	want := `[foo](./docs/bar.md "My Title")` + "\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}
