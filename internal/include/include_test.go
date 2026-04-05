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
	dir := t.TempDir()

	childPath := filepath.Join(dir, "child.md")
	writeFile(t, childPath, "child content\n")
	writeFile(t, filepath.Join(dir, "root.md"), fmt.Sprintf("before\n<!-- @include: %s -->\nafter\n", childPath))

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "before\nchild content\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_RecursiveInclude(t *testing.T) {
	dir := t.TempDir()

	cPath := filepath.Join(dir, "c.md")
	bPath := filepath.Join(dir, "b.md")
	writeFile(t, cPath, "level-c\n")
	writeFile(t, bPath, fmt.Sprintf("level-b\n<!-- @include: %s -->\n", cPath))
	writeFile(t, filepath.Join(dir, "a.md"), fmt.Sprintf("level-a\n<!-- @include: %s -->\n", bPath))

	got, err := include.ProcessFile(filepath.Join(dir, "a.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "level-a\nlevel-b\nlevel-c\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_CircularInclude(t *testing.T) {
	dir := t.TempDir()

	aPath := filepath.Join(dir, "a.md")
	bPath := filepath.Join(dir, "b.md")
	writeFile(t, aPath, fmt.Sprintf("<!-- @include: %s -->\n", bPath))
	writeFile(t, bPath, fmt.Sprintf("<!-- @include: %s -->\n", aPath))

	_, err := include.ProcessFile(aPath)
	if err == nil {
		t.Fatal("expected circular include error, got nil")
	}
}

func TestProcessFile_MissingFile(t *testing.T) {
	dir := t.TempDir()

	missingPath := filepath.Join(dir, "missing.md")
	writeFile(t, filepath.Join(dir, "root.md"), fmt.Sprintf("<!-- @include: %s -->\n", missingPath))

	_, err := include.ProcessFile(filepath.Join(dir, "root.md"))
	if err == nil {
		t.Fatal("expected error for missing include, got nil")
	}
}

func TestProcessFile_NoIncludes(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "plain.md"), "just text\nno includes\n")

	got, err := include.ProcessFile(filepath.Join(dir, "plain.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "just text\nno includes\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_IncludeInsideCodeFence(t *testing.T) {
	dir := t.TempDir()

	childPath := filepath.Join(dir, "child.md")
	// child.md does NOT exist — if the include inside the fence were expanded, it would error
	writeFile(t, filepath.Join(dir, "root.md"), fmt.Sprintf("```md\n<!-- @include: %s -->\n```\n", childPath))

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"))
	if err != nil {
		t.Fatalf("include inside code fence should not be expanded, got error: %v", err)
	}

	want := fmt.Sprintf("```md\n<!-- @include: %s -->\n```\n", childPath)
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_RelativeIncludePath(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "sub", "child.md"), "child content\n")
	// Use a relative path from the root file's directory
	writeFile(t, filepath.Join(dir, "root.md"), "before\n<!-- @include: sub/child.md -->\nafter\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := "before\nchild content\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_MixedFenceTypes(t *testing.T) {
	dir := t.TempDir()

	// A tilde fence should NOT be closed by backticks — the include inside remains literal.
	// child.md does NOT exist; if the include were expanded it would error.
	writeFile(t, filepath.Join(dir, "root.md"), "~~~\n<!-- @include: missing.md -->\n~~~\nafter\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"))
	if err != nil {
		t.Fatalf("include inside tilde fence should not be expanded, got error: %v", err)
	}

	want := "~~~\n<!-- @include: missing.md -->\n~~~\nafter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestProcessFile_IncludeInlineNotExpanded(t *testing.T) {
	dir := t.TempDir()

	// Directive embedded within other text should NOT be expanded (regex is anchored).
	writeFile(t, filepath.Join(dir, "root.md"), "Note: <!-- @include: missing.md --> end\n")

	got, err := include.ProcessFile(filepath.Join(dir, "root.md"))
	if err != nil {
		t.Fatalf("inline directive should not be expanded, got error: %v", err)
	}

	want := "Note: <!-- @include: missing.md --> end\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}
