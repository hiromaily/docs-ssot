## Cross-Tool Mapping

### Concept Mapping

| Concept | Claude Code | Cursor | Codex | Copilot |
|---------|-------------|--------|-------|---------|
| Common instructions | `AGENTS.md` | `AGENTS.md` | `AGENTS.md` | `AGENTS.md` |
| Primary instructions | `CLAUDE.md` | `.cursor/rules/` | `AGENTS.md` | `.github/copilot-instructions.md` |
| Scoped rules | `.claude/rules/*.md` | `.cursor/rules/*.mdc` | Nested `AGENTS.md` | `.github/instructions/*.instructions.md` |
| Skills | `.claude/skills/` | `.cursor/skills/` | `.agents/skills/` | `.github/skills/` |
| Commands | `.claude/commands/` (legacy) | Slash commands (migrating) | `~/.codex/prompts/` (deprecated) | `.github/prompts/*.prompt.md` |
| Subagents | `.claude/agents/*.md` | `.cursor/agents/*.md` | `.codex/agents/*.toml` | `.github/agents/*.agent.md` |
| Settings | `.claude/settings.json` | Cursor settings / `.cursor/cli.json` | `.codex/config.toml` | VS Code / GitHub settings |

### Skill Frontmatter Comparison

All four tools use `SKILL.md` with YAML frontmatter, but the supported fields differ:

| Field | Claude | Cursor | Codex | Copilot |
|-------|--------|--------|-------|---------|
| `name` | Optional | Yes | **Required** | **Required** |
| `description` | Recommended | Yes | **Required** | **Required** |
| `argument-hint` | Yes | — | — | — |
| `disable-model-invocation` | Yes | Yes | — | — |
| `user-invocable` | Yes | — | — | — |
| `allowed-tools` | Yes | — | — | Yes |
| `model` | Yes | — | — | — |
| `effort` | Yes | — | — | — |
| `context` | Yes | — | — | — |
| `license` | — | — | — | Yes |

**Claude** has the richest frontmatter. **Codex** has the most minimal (name + description only).

### Rules Frontmatter Comparison

| Field | Cursor `.mdc` | Copilot `.instructions.md` | Claude `.claude/rules/*.md` | Codex |
|-------|---------------|---------------------------|----------------------------|-------|
| `description` | Yes | — | — (no frontmatter) | — (uses `AGENTS.md`) |
| `globs` | Yes | — | — | — |
| `alwaysApply` | Yes | — | — | — |
| `applyTo` | — | Yes | — | — |
| `excludeAgent` | — | Yes | — | — |

### Functional Categories

Understanding what goes where:

| What you want | Mechanism | Tools that support it |
|--------------|-----------|----------------------|
| Always-active project rules | Instructions file | All four |
| Path-specific rules | Scoped rules | Cursor (`globs`), Copilot (`applyTo`), Codex (nested `AGENTS.md`) |
| Reusable multi-step workflows | Skills | All four |
| Manual-only slash commands | Skills + `disable-model-invocation: true` | Claude, Cursor |
| Runtime config (model, permissions) | Settings file | Claude (`settings.json`), Codex (`config.toml`) |
| Specialized agent roles | Subagent definitions | All four |
