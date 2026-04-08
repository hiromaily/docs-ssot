package cli_test

import (
	"reflect"
	"testing"

	"github.com/hiromaily/docs-ssot/internal/agentscan"
	"github.com/hiromaily/docs-ssot/internal/cli"
)

func TestResolveTargetTools(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		to      string
		from    string
		want    []agentscan.Tool
		wantErr bool
	}{
		{
			name: "empty_to_from_claude",
			to:   "",
			from: "claude",
			want: []agentscan.Tool{agentscan.ToolCursor, agentscan.ToolCopilot, agentscan.ToolCodex},
		},
		{
			name: "explicit_to",
			to:   "cursor,codex",
			from: "claude",
			want: []agentscan.Tool{agentscan.ToolCursor, agentscan.ToolCodex},
		},
		{
			name: "empty_to_empty_from",
			to:   "",
			from: "",
			want: agentscan.AllTools(),
		},
		{
			name: "empty_to_from_auto",
			to:   "",
			from: "auto",
			want: agentscan.AllTools(),
		},
		{
			name: "all_to_from_claude",
			to:   "all",
			from: "claude",
			want: []agentscan.Tool{agentscan.ToolCursor, agentscan.ToolCopilot, agentscan.ToolCodex},
		},
		{
			name: "to_includes_source_filtered",
			to:   "claude,cursor",
			from: "claude",
			want: []agentscan.Tool{agentscan.ToolCursor},
		},
		{
			name:    "invalid_from",
			to:      "",
			from:    "vim",
			wantErr: true,
		},
		{
			name:    "invalid_to",
			to:      "vim",
			from:    "claude",
			wantErr: true,
		},
		{
			name: "duplicate_to",
			to:   "cursor,cursor",
			from: "claude",
			want: []agentscan.Tool{agentscan.ToolCursor},
		},
		{
			name: "trailing_comma",
			to:   "cursor,",
			from: "claude",
			want: []agentscan.Tool{agentscan.ToolCursor},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := cli.ResolveTargetTools(tt.to, tt.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveTargetTools(%q, %q) error = %v, wantErr %v", tt.to, tt.from, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveTargetTools(%q, %q) = %v, want %v", tt.to, tt.from, got, tt.want)
			}
		})
	}
}

func TestParseToolList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []agentscan.Tool
		wantErr bool
	}{
		{name: "single", input: "claude", want: []agentscan.Tool{agentscan.ToolClaude}},
		{name: "multiple", input: "claude,cursor", want: []agentscan.Tool{agentscan.ToolClaude, agentscan.ToolCursor}},
		{name: "with_spaces", input: " claude , cursor ", want: []agentscan.Tool{agentscan.ToolClaude, agentscan.ToolCursor}},
		{name: "dedup", input: "claude,claude", want: []agentscan.Tool{agentscan.ToolClaude}},
		{name: "empty_segments", input: ",claude,,cursor,", want: []agentscan.Tool{agentscan.ToolClaude, agentscan.ToolCursor}},
		{name: "invalid", input: "vim", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := cli.ParseToolList(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToolList(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseToolList(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
