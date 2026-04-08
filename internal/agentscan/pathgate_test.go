package agentscan_test

import (
	"testing"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
)

func TestInferGlobs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		slug   string
		want   string
		wantOK bool
	}{
		{name: "exact_go", slug: "go", want: "**/*.go", wantOK: true},
		{name: "exact_typescript", slug: "typescript", want: "**/*.{ts,tsx}", wantOK: true},
		{name: "exact_python", slug: "python", want: "**/*.py", wantOK: true},
		{name: "exact_proto", slug: "proto", want: "**/*.proto", wantOK: true},
		{name: "contains_frontend", slug: "app-web-architecture", want: "frontend/app-web/**", wantOK: true},
		{name: "contains_backend", slug: "backend-rules", want: "backend/**", wantOK: true},
		{name: "exact_test", slug: "testing", want: "**/*_test.*", wantOK: true},
		{name: "unknown", slug: "github-flow", want: "", wantOK: false},
		{name: "unknown_generic", slug: "architecture", want: "", wantOK: false},
		{name: "case_insensitive", slug: "Go", want: "**/*.go", wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := agentscan.InferGlobs(tt.slug)
			if ok != tt.wantOK {
				t.Errorf("InferGlobs(%q) ok = %v, want %v", tt.slug, ok, tt.wantOK)
			}
			if got != tt.want {
				t.Errorf("InferGlobs(%q) = %q, want %q", tt.slug, got, tt.want)
			}
		})
	}
}
