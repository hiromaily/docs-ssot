## AI Agent Configuration Landscape (April 2026)

AI coding agents (Claude Code, Codex, Cursor, GitHub Copilot) each require configuration files to understand project context. These files fall into four layers:

### Layer 1 — Persistent Instructions

Files the agent reads every session to understand project rules and architecture.

| Tool | Primary file | Scoped files |
|------|-------------|-------------|
| Claude Code | `CLAUDE.md` | Subdirectory `CLAUDE.md`, `.claude/CLAUDE.md`, `CLAUDE.local.md` |
| Codex | `AGENTS.md` | Nested `AGENTS.md` per directory, `AGENTS.override.md` |
| Cursor | `.cursor/rules/*.mdc` | Per-file via `globs` frontmatter |
| Copilot | `.github/copilot-instructions.md` | `.github/instructions/*.instructions.md`, `AGENTS.md` |

### Layer 2 — Scoped Rules

Topic-specific or path-gated rules that supplement the primary instruction file.

| Tool | Location | Format |
|------|----------|--------|
| Claude Code | `.claude/rules/*.md` | Markdown, optionally path-gated |
| Codex | Nested `AGENTS.md` hierarchy | Markdown, directory-scoped |
| Cursor | `.cursor/rules/*.mdc` | MDC (Markdown + YAML frontmatter) |
| Copilot | `.github/instructions/*.instructions.md` | Markdown + `applyTo` frontmatter |

### Layer 3 — Reusable Workflows (Skills / Commands)

Packaged multi-step procedures the agent invokes on demand or automatically.

| Tool | Skills location | Command location | Trend |
|------|----------------|-----------------|-------|
| Claude Code | `.claude/skills/<name>/SKILL.md` | `.claude/commands/*.md` (legacy) | Commands integrated into skills |
| Codex | `.agents/skills/<name>/SKILL.md` | `~/.codex/prompts/*.md` (deprecated) | Custom prompts deprecated, skills preferred |
| Cursor | `.cursor/skills/<name>/SKILL.md` | Slash commands | Commands migrating to skills |
| Copilot | `.github/skills/<name>/SKILL.md` | `.github/prompts/*.prompt.md` | Prompt files for explicit invocation |

### Layer 4 — Agent Execution Settings

Runtime configuration controlling model selection, permissions, subagents, and external connections.

| Tool | Settings file | Subagents | Hooks |
|------|--------------|-----------|-------|
| Claude Code | `.claude/settings.json` | `.claude/agents/*.md` | Hooks in `settings.json` |
| Codex | `.codex/config.toml` | `.codex/agents/*.toml` | `.codex/hooks.json` |
| Cursor | `.cursor/cli.json` | `.cursor/agents/*.md` | — |
| Copilot | VS Code / GitHub settings | `.github/agents/*.agent.md` | — |

---

### Key Trends

1. **`AGENTS.md` is the de facto cross-tool standard** — supported by Claude, Codex, Cursor, and Copilot
2. **Claude Code has the richest configuration** — skills, agents, commands, memory, hooks, settings
3. **Cursor is evolving from IDE to Agent-first** — rules and commands are migrating to skills
4. **Copilot is GitHub-native** — deeply integrated with issues, PRs, and the `.github/` directory
5. **Skills are converging** — all four tools support `SKILL.md`-based skills with YAML frontmatter
6. **Commands are being deprecated or merged into skills** across all platforms
