## OpenAI Codex

Codex uses `AGENTS.md` as its primary instruction mechanism, with a hierarchical directory-based scoping model.

### File Overview

| Category | File | Scope |
|----------|------|-------|
| Instructions | `AGENTS.md` | Project-wide and per-directory guidance |
| Instructions | `AGENTS.override.md` | Override file (takes precedence over `AGENTS.md`) |
| Instructions | `~/.codex/AGENTS.md` | User-level global instructions |
| Skills | `.agents/skills/<name>/SKILL.md` | Repository-scoped skills |
| Skills | `~/.agents/skills/<name>/SKILL.md` | User-level skills |
| Skills | `/etc/codex/skills/<name>/SKILL.md` | Admin-level skills |
| Subagents | `.codex/agents/*.toml` | Project-scoped subagents |
| Subagents | `~/.codex/agents/*.toml` | User-level subagents |
| Settings | `.codex/config.toml` | Project-scoped runtime config |
| Settings | `~/.codex/config.toml` | User-level runtime config |
| Hooks | `.codex/hooks.json` | Project-scoped hooks (experimental) |
| Hooks | `~/.codex/hooks.json` | User-level hooks |
| Commands | `~/.codex/prompts/*.md` | Custom prompts (deprecated, use skills) |

### Instruction Hierarchy

Codex builds an instruction chain by walking from the repository root to the current working directory. At each level, it looks for (in priority order):

1. `AGENTS.override.md`
2. `AGENTS.md`
3. Fallback filenames (configurable in `config.toml`)

Instructions accumulate — deeper directories add specificity to root-level rules. This makes `AGENTS.md` placement a structural design decision.

### Skills

Codex skills use the `.agents/skills/` directory (not `.codex/skills/`). Each skill requires a `SKILL.md` with YAML frontmatter.

#### Skill Frontmatter

```yaml
---
name: add-endpoint          # Required. Skill identifier
description: Add a new ...  # Required. Used for progressive disclosure matching
---
```

Codex has the most minimal skill frontmatter. The `description` is critical because Codex reads metadata first and loads full content only when needed (progressive disclosure).

Additional configuration can go in `agents/openai.yaml` within the skill directory.

### Settings (`config.toml`)

`config.toml` separates runtime concerns from project guidance:

| Setting | Purpose |
|---------|---------|
| `model` | Default model selection |
| `approval_policy` | `untrusted` / `on-request` / `never` |
| `sandbox_mode` | File access scope |
| `project_root_markers` | How Codex finds project boundaries |
| `project_doc_fallback_filenames` | Additional instruction file names |
| `[mcp_servers.*]` | External tool connections |
| `[agents.*]` | Subagent definitions |
| `[tools]` | Tool enablement (e.g., `web_search`) |

**Key distinction**: `AGENTS.md` = how to think about the code; `config.toml` = how to run the agent.

### Subagents

Codex subagents are defined as TOML files in `.codex/agents/`. They specify role, model, thread limits, and other execution parameters.
