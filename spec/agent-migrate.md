# Design: Agent-Aware Migration (`migrate --from`)

## Status

Implemented (Phase 1 + Phase 2)

## Summary

Extend the `migrate` command with `--from` / `--to` flags that convert existing AI tool configuration files (rules, skills, commands, subagents) into the docs-ssot SSOT structure, enabling **one-command multi-tool expansion** from a single tool's configuration.

## Motivation

### Problem: "I have Claude rules, now I need Cursor/Copilot/Codex too"

A repository that has invested in Claude Code configuration (`.claude/rules/*.md`, `.claude/skills/*/SKILL.md`, `.claude/commands/*.md`, `.claude/agents/*.md`) often wants to support other tools. Today this requires:

1. Manually copying each rule/skill/command/subagent
2. Adapting format and frontmatter per tool (`.mdc` for Cursor, `applyTo` for Copilot, etc.)
3. Maintaining N copies of the same content going forward

This is the exact duplication problem docs-ssot was built to solve.

### Goal: "Expand my Claude rules to all tools"

```sh
docs-ssot migrate --from claude
# Scans .claude/rules/, .claude/skills/, .claude/commands/, .claude/agents/
# Creates shared sections (1 file per rule/skill/command/subagent)
# Generates templates for Cursor, Copilot, Codex (all except source)
# Adds all targets to docsgen.yaml
# From now on: edit section → build → all tools updated
```

### Non-goal: Deep content decomposition

Each existing file becomes one section. The value comes from multi-tool template generation and frontmatter adaptation, not from restructuring content.

---

## Design

### Core principle: 1 file = 1 section, N templates

```
Existing:  .claude/rules/architecture.md
                ↓ migrate --from claude
Section:   template/sections/ai/rules/architecture.md     ← content (no frontmatter)
Templates: template/pages/ai-agents/cursor/rules/architecture.tpl.mdc   ← frontmatter + @include level=-1
           template/pages/ai-agents/copilot/instructions/architecture.tpl.md ← frontmatter + @include level=-1
Config:    docsgen.yaml targets for each template → output mapping
```

### Behaviour

1. **Scan source tool** — collect rules, skills, commands, subagents from `--from` tool's directories
2. **Extract frontmatter** — parse YAML frontmatter using `yaml.Unmarshal` (handles multi-line values); preserve for source tool's template
3. **Create section files** — strip frontmatter, shift H1→H2 headings, write to `template/sections/ai/<type>s/<slug>.md`
4. **Generate templates for target tools** — create tool-specific templates with appropriate frontmatter
5. **Update docsgen.yaml** — append targets using structured `config.Save()` with deduplication
6. **Verify round-trip** — build and compare section content against originals

### Agent file type mapping

| Source type | Section location | Claude template | Cursor template | Copilot template | Codex |
|-------------|-----------------|-----------------|-----------------|------------------|-------|
| Rules | `sections/ai/rules/<slug>.md` | `claude/rules/<slug>.tpl.md` | `cursor/rules/<slug>.tpl.mdc` | `copilot/instructions/<slug>.tpl.md` | Embedded in `AGENTS.tpl.md` |
| Skills | `sections/ai/skills/<slug>.md` | `claude/skills/<slug>/SKILL.tpl.md` | `cursor/skills/<slug>/SKILL.tpl.md` | `copilot/skills/<slug>/SKILL.tpl.md` | `codex/skills/<slug>/SKILL.tpl.md` |
| Commands | `sections/ai/commands/<slug>.md` | `claude/commands/<slug>.tpl.md` | — (not supported) | — (not supported) | — (not supported) |
| Subagents | `sections/ai/subagents/<slug>.md` | `claude/agents/<slug>.tpl.md` | `cursor/agents/<slug>.tpl.md` | `copilot/agents/<slug>.tpl.md` | `codex/agents/<slug>.md` |

### Frontmatter adaptation

Each tool requires different frontmatter. The migrate command generates appropriate frontmatter based on the source file's content and metadata:

#### Rules

| Tool | Format | Generated frontmatter |
|------|--------|-----------------------|
| Claude Code | `.md` | None (plain markdown) |
| Cursor | `.mdc` | `description`, `alwaysApply: true` (or `globs` if `--infer-globs`) |
| Copilot | `.md` | `applyTo: "**/*"` (or inferred pattern if `--infer-globs`) |
| Codex | N/A | Content embedded in combined `AGENTS.md` via `@include` |

#### Skills

| Tool | Format | Generated frontmatter |
|------|--------|-----------------------|
| Claude Code | `SKILL.md` | All original fields preserved (`name`, `description`, `model`, `effort`, `allowed-tools`, etc.) |
| Others | `SKILL.md` | `name`, `description` only |

#### Subagents

| Tool | Format | Generated frontmatter |
|------|--------|-----------------------|
| Claude Code | `.md` | All original fields preserved (`name`, `description`, `disallowedTools`, etc.) |
| Copilot | `.agent.md` | `name`, `description` only |
| Others | `.md` | `name`, `description` only |

### Heading level handling

Section files follow H2+ convention. H1 headings in source files are shifted to H2 during section creation. Templates use `level=-1` to restore them to H1 in generated output.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--from` | — | Source AI tool to migrate from (`claude`, `cursor`, `copilot`) |
| `--to` | all except `--from` | Target tools, comma-separated (`cursor,copilot,codex`) |
| `--convert-commands` | `false` | Convert legacy commands to skills during migration |
| `--infer-globs` | `false` | Infer path-gated globs from rule slug names |
| `--dry-run` | `false` | Print migration plan without writing files |

Legacy flags `--agents`, `--source-tool`, `--target-tools` are hidden but functional for backward compatibility.

### `--to` default behaviour

When `--to` is omitted, all tools **except** the source tool are used as targets. This matches the most common use case: "I have Claude, generate everything else."

### Path-gated rules inference (`--infer-globs`)

When enabled, slug names are matched against known patterns:

| Slug | Inferred pattern |
|------|-----------------|
| `go` | `**/*.go` |
| `typescript` | `**/*.{ts,tsx}` |
| `frontend-*` | `frontend/**` |
| `backend-*` | `backend/**` |
| `testing` | `**/*_test.*` |

Matching uses deterministic sorted keys with longest-match preference. Unknown slugs default to `alwaysApply: true` (Cursor) / `applyTo: "**/*"` (Copilot).

### Command conversion (`--convert-commands`)

When enabled, legacy `.claude/commands/*.md` files are re-typed as skills during migration. This allows them to be generated for all target tools (commands are Claude-only, skills are cross-tool).

### Examples

```sh
# Migrate Claude configs to all other tools
docs-ssot migrate --from claude

# Migrate to specific tools only
docs-ssot migrate --from claude --to cursor,codex

# Preview migration plan
docs-ssot migrate --from claude --dry-run

# With path inference and command conversion
docs-ssot migrate --from claude --to cursor --infer-globs --convert-commands

# Combined with file migration
docs-ssot migrate --from claude --to cursor README.md CLAUDE.md
```

### Interaction with existing `migrate`

`--from` can be combined with regular file arguments. When used together:
- Regular files go through the existing section decomposition pipeline
- Agent files go through the agent-aware pipeline (no decomposition, 1:1 section mapping)
- Both share the same docsgen.yaml and round-trip verification

---

## Implementation

### New packages

| Package | Purpose |
|---------|---------|
| `agentscan` | Detect AI tools, collect agent files, infer path-gated globs |
| `frontmatter` | Parse (via `yaml.Unmarshal`), strip, and generate YAML frontmatter for different tool formats |

### Modified packages

| Package | Changes |
|---------|---------|
| `migrate` | Add `RunAgents()` for agent-aware migration pipeline |
| `cli` | Add `--from`, `--to`, `--convert-commands`, `--infer-globs` flags |
| `config` | Add `Save()` for structured YAML serialization |

### Reused packages

| Package | Reuse |
|---------|-------|
| `config` | Load/save docsgen.yaml |
| `generator` | Round-trip verification (build + diff) |

---

## Resolved Open Questions

1. **Codex AGENTS.md structure** — Resolved: inline includes per section. Each rule gets its own `@include` in the combined `AGENTS.tpl.md`.

2. **Path-gated rules** — Resolved: opt-in via `--infer-globs`. When disabled, defaults to `alwaysApply: true` / `applyTo: "**/*"`. When enabled, uses deterministic longest-match against known extension and path patterns.

3. **Existing multi-tool files** — Resolved: only the `--from` tool's files are used as source. Existing files for other tools are not scanned or merged.

4. **Commands deprecation** — Resolved: opt-in via `--convert-commands`. When enabled, commands are re-typed as skills and generated for all target tools. When disabled, commands remain Claude-only.

5. **Skill frontmatter preservation** — Resolved: Claude-specific fields (`allowed-tools`, `model`, `effort`, etc.) are preserved in the source tool's template. Other tools receive only `name` and `description`. Multi-line YAML values (lists) are properly parsed via `yaml.Unmarshal` and serialised to compact YAML.
