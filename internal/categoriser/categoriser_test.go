package categoriser_test

import (
	"testing"

	"github.com/hiromaily/docs-ssot/internal/categoriser"
)

func TestCategorise(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{name: "overview", title: "Overview", expected: "project"},
		{name: "project_overview", title: "Project Overview", expected: "project"},
		{name: "about", title: "About This Project", expected: "project"},
		{name: "introduction", title: "Introduction", expected: "project"},
		{name: "background", title: "Background", expected: "project"},
		{name: "installation", title: "Installation", expected: "development"},
		{name: "setup_guide", title: "Setup Guide", expected: "development"},
		{name: "getting_started", title: "Getting Started", expected: "development"},
		{name: "prerequisites", title: "Prerequisites", expected: "development"},
		{name: "architecture_overview", title: "Architecture Overview", expected: "architecture"},
		{name: "system_design", title: "System Design", expected: "architecture"},
		{name: "pipeline_architecture", title: "Pipeline Architecture", expected: "architecture"},
		{name: "testing", title: "Testing", expected: "development"},
		{name: "unit_tests", title: "Unit Tests", expected: "development"},
		{name: "ci_configuration", title: "CI Configuration", expected: "development"},
		{name: "linting", title: "Linting", expected: "development"},
		{name: "code_quality", title: "Code Quality", expected: "development"},
		{name: "api_reference", title: "API Reference", expected: "reference"},
		{name: "cli_commands", title: "CLI Commands", expected: "reference"},
		{name: "commands_reference", title: "Commands Reference", expected: "reference"},
		{name: "contributing", title: "Contributing", expected: "development"},
		{name: "how_to_contribute", title: "How to Contribute", expected: "development"},
		{name: "license", title: "License", expected: "project"},
		{name: "changelog", title: "Changelog", expected: "project"},
		{name: "roadmap", title: "Roadmap", expected: "project"},
		{name: "faq", title: "FAQ", expected: "product"},
		{name: "troubleshooting", title: "Troubleshooting", expected: "product"},
		{name: "fallback_misc", title: "Random Section", expected: "misc"},
		{name: "fallback_misc_2", title: "Something Else", expected: "misc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := categoriser.Categorise(tt.title)
			if got != tt.expected {
				t.Errorf("Categorise(%q) = %q, want %q", tt.title, got, tt.expected)
			}
		})
	}
}
