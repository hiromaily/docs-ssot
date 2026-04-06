## Cursor

Cursor is transitioning from an IDE-centric model to an agent-first architecture. Rules are the primary configuration mechanism, with skills and subagents as newer additions.

### File Overview

| Category | File | Scope |
|----------|------|-------|
| Rules | `.cursor/rules/*.md` | Project rules (plain Markdown) |
| Rules | `.cursor/rules/*.mdc` | Project rules (MDC: Markdown + frontmatter) |
| Rules | `.cursorrules` | Legacy single-file rules (deprecated) |
| Skills | `.cursor/skills/<name>/SKILL.md` | Project-scoped skills |
| Subagents | `.cursor/agents/*.md` | Project-scoped subagents |
| Subagents | `.claude/agents/*.md` | Claude compatibility (also read by Cursor) |
| Subagents | `.codex/agents/*.toml` | Codex compatibility (also read by Cursor) |
| Settings | `~/.cursor/cli-config.json` | CLI permissions (user-level) |
| Settings | `.cursor/cli.json` | CLI permissions (project-level) |
| Ignore | `.cursorignore` | Files excluded from context |
| Compat | `AGENTS.md` | Cross-tool instructions (read by Cursor) |

### Rules

Rules are the core configuration for Cursor. Two formats are supported:

- **`.md`** — plain Markdown, always applied
- **`.mdc`** — MDC format with YAML frontmatter for conditional application

#### Rule Frontmatter (MDC format)

```yaml
---
description: Go layered architecture rules   # What the rule covers
globs:                                        # File patterns that activate this rule
  - "internal/**/*.go"
  - "pkg/**/*.go"
alwaysApply: true                             # Apply regardless of file context
---
```

Key fields:

- `description` — tells Cursor when the rule is relevant
- `globs` — file patterns that trigger the rule
- `alwaysApply` — force the rule to always be active

### Skills

Cursor supports Agent Skills (the emerging open standard). Skills can function as:

- **Implicit skills** — automatically invoked when the description matches the task
- **Explicit commands** — invoked only via `/skill-name` when `disable-model-invocation: true` is set

#### Skill Frontmatter

```yaml
---
name: add-endpoint                    # Skill identifier and slash command name
description: Add a new Go endpoint    # Used for auto-invocation matching
disable-model-invocation: true        # Optional. Makes it a manual-only command
---
```

Cursor can migrate existing rules and slash commands into skills. Migrated commands get `disable-model-invocation: true` to preserve explicit-invocation behavior.

### Subagents

Cursor reads subagent definitions from multiple locations for cross-tool compatibility:

1. `.cursor/agents/` (native)
2. `.claude/agents/` (Claude compatibility)
3. `.codex/agents/` (Codex compatibility)

### Migration Path

The trend in Cursor is:

- `.cursorrules` (single file) → `.cursor/rules/` (multiple files)
- Slash commands → Skills with `disable-model-invocation: true`
- Rules with reusable workflows → Skills
