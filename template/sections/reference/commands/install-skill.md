## docs-ssot install-skill

Install a `SKILL.md` file for AI coding agents so they know how to migrate existing documentation to the docs-ssot SSOT structure, build output files, and validate the result.

```
docs-ssot install-skill [flags]
```

### What it does

- Writes a `SKILL.md` file into each target tool's skills directory
- The installed skill guides the agent through the full docs-ssot workflow: migrate → build → validate
- Prompts before overwriting an existing skill file

### Skill install locations

| Tool | Path |
|------|------|
| Claude Code | `.claude/skills/docs-ssot/SKILL.md` |
| Cursor | `.cursor/skills/docs-ssot/SKILL.md` |
| GitHub Copilot | `.github/skills/docs-ssot/SKILL.md` |
| Codex | `.agents/skills/docs-ssot/SKILL.md` |

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--tool` | all | Target tool(s): `claude`, `cursor`, `copilot`, `codex` (comma-separated) |

### Examples

Install for all tools:

```
docs-ssot install-skill
```

Install for Claude Code only:

```
docs-ssot install-skill --tool claude
```

Install for multiple specific tools:

```
docs-ssot install-skill --tool claude,cursor
```
