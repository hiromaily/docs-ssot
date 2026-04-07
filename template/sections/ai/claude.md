## Claude Code

Claude Code has the most comprehensive configuration system among AI coding agents.

### File Overview

| Category | File | Scope |
|----------|------|-------|
| Instructions | `CLAUDE.md` | Project-wide context and rules |
| Instructions | `CLAUDE.local.md` | Personal overrides (gitignored) |
| Instructions | `.claude/CLAUDE.md` | Alternative project instructions |
| Rules | `.claude/rules/*.md` | Topic-scoped or path-gated rules |
| Rules | `~/.claude/rules/*.md` | User-level global rules |
| Skills | `.claude/skills/<name>/SKILL.md` | Reusable workflows |
| Skills | `~/.claude/skills/<name>/SKILL.md` | User-level global skills |
| Commands | `.claude/commands/*.md` | Legacy custom commands (integrated into skills) |
| Subagents | `.claude/agents/*.md` | Project-scoped custom subagents |
| Subagents | `~/.claude/agents/*.md` | User-level custom subagents |
| Settings | `.claude/settings.json` | Permissions, hooks, MCP servers |
| Settings | `.claude/settings.local.json` | Personal settings overrides |
| Settings | `~/.claude/settings.json` | User-level global settings |

### Instruction Hierarchy

Claude Code merges instructions from multiple scopes (global > project > subdirectory):

1. `~/.claude/CLAUDE.md` (user-level)
2. `CLAUDE.md` (project root)
3. Subdirectory `CLAUDE.md` files (deeper = more specific)

### Rules

`.claude/rules/*.md` files provide topic-specific or path-gated instructions. They supplement `CLAUDE.md` without duplicating its content.

Rules have no required frontmatter. They are plain Markdown files that Claude reads automatically.

### Skills

Skills are the primary mechanism for reusable workflows. Custom commands (`.claude/commands/*.md`) are integrated into skills — both create slash commands, but skills offer richer configuration.

#### Skill Frontmatter

```yaml
---
name: deploy                        # Optional. Creates /deploy slash command
description: Deploy to production   # Recommended. Used for auto-invocation matching
argument-hint: [env]                # Optional. Shown in completion UI
disable-model-invocation: true      # Optional. Prevents auto-invocation (manual /name only)
user-invocable: true                # Optional. Whether it appears in slash command menu
allowed-tools:                      # Optional. Restricts available tools during skill execution
  - Read
  - Edit
  - Bash(make *)
model: opus                         # Optional. Override model for this skill
effort: high                        # Optional. Reasoning effort level
context: fork                       # Optional. Run in forked subagent context
---
```

Claude has the most feature-rich skill frontmatter of all tools.

### Subagents

Custom subagents are defined as Markdown files in `.claude/agents/`. Each file defines a specialized agent with its own system prompt, available tools, and model configuration.

### Settings

`.claude/settings.json` controls:

- Permissions and approval policies
- Hook definitions (pre/post tool execution)
- MCP server connections
- Environment variables
- Model defaults
