# Design: Agent-Aware Migration (`migrate --agents`)

## Status

Draft

## Summary

Extend the `migrate` command with an `--agents` mode that converts existing AI tool configuration files (rules, skills, commands) into the docs-ssot SSOT structure, enabling **one-command multi-tool expansion** from a single tool's configuration.

## Motivation

### Problem: "I have Claude rules, now I need Cursor/Copilot/Codex too"

A repository that has invested in Claude Code configuration (`.claude/rules/*.md`, `.claude/skills/*/SKILL.md`, `.claude/commands/*.md`) often wants to support other tools. Today this requires:

1. Manually copying each rule/skill/command
2. Adapting format and frontmatter per tool (`.mdc` for Cursor, `applyTo` for Copilot, etc.)
3. Maintaining N copies of the same content going forward

This is the exact duplication problem docs-ssot was built to solve.

### Goal: "Expand my Claude rules to all tools"

```sh
docs-ssot migrate --agents
# Scans .claude/rules/, .claude/skills/, .claude/commands/
# Creates shared sections (1 file per rule/skill/command)
# Generates templates for Claude, Cursor, Copilot, Codex
# Adds all targets to docsgen.yaml
# From now on: edit section → build → all tools updated
```

### Non-goal: Deep content decomposition

Phase 1 does **not** split agent files into sub-sections. Each existing file becomes one section. The value comes from multi-tool template generation and frontmatter adaptation, not from restructuring content.

---

## Design

### Core principle: 1 file = 1 section, N templates

```
Existing:  .claude/rules/architecture.md
                ↓ migrate --agents
Section:   template/sections/ai/rules/architecture.md     ← content (no frontmatter)
Templates: template/pages/ai-agents/claude/rules/architecture.tpl.md    ← @include level=-1
           template/pages/ai-agents/cursor/rules/architecture.tpl.mdc   ← frontmatter + @include level=-1
           template/pages/ai-agents/copilot/instructions/architecture.tpl.md ← frontmatter + @include level=-1
Config:    docsgen.yaml targets for each template → output mapping
```

### Behaviour

1. **Detect AI tools** — scan for `.claude/`, `.cursor/`, `.github/copilot-instructions.md`, `.codex/`, `AGENTS.md`
2. **Collect agent files** — enumerate rules, skills, commands from the detected tool directories
3. **Determine source tool** — identify which tool's files are the "source of truth" (the most complete set)
4. **Extract frontmatter** — strip YAML frontmatter from source files; preserve for the source tool's template
5. **Create section files** — copy content (without frontmatter) to `template/sections/ai/<type>/<slug>.md`
6. **Generate templates for all target tools** — create tool-specific templates with appropriate frontmatter
7. **Update docsgen.yaml** — add targets for all generated templates
8. **Verify round-trip** — build and compare against originals

### Agent file type mapping

| Source type | Section location | Claude template | Cursor template | Copilot template | Codex |
|-------------|-----------------|-----------------|-----------------|------------------|-------|
| Rules | `sections/ai/rules/<slug>.md` | `claude/rules/<slug>.tpl.md` | `cursor/rules/<slug>.tpl.mdc` | `copilot/instructions/<slug>.tpl.md` | Embedded in `AGENTS.tpl.md` |
| Skills | `sections/ai/skills/<slug>.md` | `claude/skills/<slug>/SKILL.tpl.md` | `cursor/skills/<slug>/SKILL.tpl.md` | `copilot/skills/<slug>/SKILL.tpl.md` | `codex/skills/<slug>/SKILL.tpl.md` |
| Commands | `sections/ai/commands/<slug>.md` | `claude/commands/<slug>.tpl.md` | — (not supported) | — (not supported) | — (not supported) |

### Frontmatter adaptation

Each tool requires different frontmatter. The migrate command generates appropriate frontmatter based on the source file's content and metadata:

#### Rules

| Tool | Format | Generated frontmatter |
|------|--------|-----------------------|
| Claude Code | `.md` | None (plain markdown) |
| Cursor | `.mdc` | `description` (derived from first heading or filename), `alwaysApply: true` |
| Copilot | `.md` | `applyTo: "**/*"` (default, user can narrow later) |
| Codex | N/A | Content embedded in combined `AGENTS.md` via `@include` |

#### Skills

| Tool | Format | Generated frontmatter |
|------|--------|-----------------------|
| Claude Code | `SKILL.md` | `name`, `description` (preserved from source or derived) |
| Cursor | `SKILL.md` | `name`, `description` |
| Copilot | `SKILL.md` | `name`, `description` |
| Codex | `SKILL.md` | `name`, `description` |

### Template generation pattern

Each generated template follows this structure:

```markdown
<!-- For .md templates (Claude, Copilot) -->
---                              ← only if frontmatter is needed
applyTo: "**/*"                  ← tool-specific fields
---

<!-- @include: ../../../../sections/ai/rules/<slug>.md level=-1 -->
```

```markdown
<!-- For .mdc templates (Cursor) -->
---
description: <derived from heading or filename>
alwaysApply: true
---

<!-- @include: ../../../../sections/ai/rules/<slug>.md level=-1 -->
```

The `level=-1` is added because section files use H2+ headings (per docs-ssot convention), and standalone rule/skill files need H1 headings.

### Heading level handling

Section files created from agent files follow the same H2+ convention as all other sections. If the source file uses H1 headings, they are shifted to H2 during section creation. Templates use `level=-1` to restore them to H1 in the generated output.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--agents` | `false` | Enable agent-aware migration mode |
| `--source-tool` | `auto` | Which tool's files to use as source: `auto`, `claude`, `cursor`, `copilot` |
| `--target-tools` | `all` | Comma-separated list of target tools: `claude,cursor,copilot,codex` or `all` |
| `--dry-run` | `false` | Print migration plan without writing files |

### Auto-detection of source tool

When `--source-tool=auto`:

1. Count agent files per tool
2. Select the tool with the most files as the source
3. Report the selection to the user

### Example

```sh
# Repository has Claude rules and skills
$ ls .claude/rules/
architecture.md  general.md  testing.md
$ ls .claude/skills/
check-ssot/SKILL.md

# Migrate to support all tools
$ docs-ssot migrate --agents --dry-run

Detected source tool: claude (3 rules, 1 skill)
Target tools: claude, cursor, copilot, codex

Would create sections:
  template/sections/ai/rules/architecture.md
  template/sections/ai/rules/general.md
  template/sections/ai/rules/testing.md
  template/sections/ai/skills/check-ssot.md

Would create templates (4 tools × 4 files):
  Claude:  3 rule templates, 1 skill template
  Cursor:  3 rule templates (.mdc), 1 skill template
  Copilot: 3 instruction templates, 1 skill template
  Codex:   1 combined AGENTS.tpl.md (3 includes), 1 skill template

Would add 15 targets to docsgen.yaml
```

### Interaction with existing `migrate`

`--agents` can be combined with regular file arguments:

```sh
# Migrate docs AND agent files in one command
docs-ssot migrate --agents README.md CLAUDE.md
```

When used together:
- Regular files (README.md, CLAUDE.md) go through the existing section decomposition pipeline
- Agent files go through the agent-aware pipeline (no decomposition, 1:1 section mapping)
- Both share the same docsgen.yaml and round-trip verification

---

## Implementation

### New packages

| Package | Purpose |
|---------|---------|
| `agentscan` | Detect AI tools and collect agent files from a repository |
| `frontmatter` | Parse, strip, and generate YAML frontmatter for different tool formats |

### Modified packages

| Package | Changes |
|---------|---------|
| `migrate` | Add `--agents` mode, call `agentscan` and `frontmatter` packages |
| `cli` | Add `--agents`, `--source-tool`, `--target-tools` flags to `migrateCmd` |

### Reused packages

| Package | Reuse |
|---------|-------|
| `config` | Load/write docsgen.yaml |
| `generator` | Round-trip verification (build + diff) |

---

## Open Questions

1. **Codex AGENTS.md structure** — Should the combined AGENTS.md template include all rules inline, or should it be a single `@include` of a directory? The current project uses inline includes per section.

2. **Path-gated rules** — For Cursor, should `--agents` attempt to infer `globs` from filenames (e.g., `go.md` → `globs: ["**/*.go"]`)? Or always default to `alwaysApply: true`?

3. **Existing multi-tool files** — If the repository already has both `.claude/rules/` and `.cursor/rules/`, should migrate detect cross-tool duplicates and merge them? Or only use the source tool's files?

4. **Commands deprecation** — Claude commands (`.claude/commands/`) are being deprecated in favour of skills. Should migrate convert commands to skills during migration?

5. **Skill frontmatter preservation** — Skills have tool-specific frontmatter fields (e.g., Claude's `allowed-tools`, `model`, `effort`). Should these be preserved verbatim in the source tool's template and omitted from other tools' templates?
