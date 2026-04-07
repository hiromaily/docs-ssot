## GitHub Copilot

Copilot integrates natively with GitHub's directory conventions. It supports multiple instruction mechanisms, including the cross-tool `AGENTS.md` standard.

### File Overview

| Category | File | Scope |
|----------|------|-------|
| Instructions | `.github/copilot-instructions.md` | Repository-wide instructions |
| Instructions | `.github/instructions/*.instructions.md` | Path-specific instructions |
| Instructions | `AGENTS.md` | Cross-tool instructions (supported by coding agent) |
| Prompts | `.github/prompts/*.prompt.md` | Reusable prompt templates (public preview) |
| Skills | `.github/skills/<name>/SKILL.md` | Repository-scoped agent skills |
| Agents | `.github/agents/*.agent.md` | Custom agent definitions |
| CLI | `~/.copilot/copilot-instructions.md` | Local CLI instructions |

### Repository-Wide Instructions

`.github/copilot-instructions.md` is the primary instruction file. It applies to all Copilot interactions in the repository.

### Path-Specific Instructions

`.github/instructions/*.instructions.md` files provide scoped rules using YAML frontmatter:

#### Instruction Frontmatter

```yaml
---
applyTo: "internal/**/*.go"         # Required. File pattern this instruction applies to
excludeAgent: "code-review"          # Optional. Exclude specific agents
---
```

Key fields:

- `applyTo` — determines which files trigger this instruction
- `excludeAgent` — prevents specific agents (e.g., `code-review`, `coding-agent`) from seeing this instruction

### AGENTS.md Support

The Copilot coding agent reads `AGENTS.md` files, including nested ones in subdirectories. This provides cross-tool compatibility with Codex and other tools that use the `AGENTS.md` convention.

### Prompt Files

`.github/prompts/*.prompt.md` are reusable prompt templates. They serve as the closest equivalent to slash commands in Copilot — explicit-invocation templates for common tasks.

Prompt files are primarily Markdown body content with minimal frontmatter.

### Skills

Copilot supports agent skills via `SKILL.md` files. It also reads skill directories from other tool locations for compatibility.

#### Skill Frontmatter

```yaml
---
name: image-convert                              # Required. Skill identifier
description: Converts SVG images to PNG format   # Required. Trigger description
license: MIT                                     # Optional. License declaration
allowed-tools: Bash(convert-svg-to-png.sh)       # Optional. Permitted tools
---
```

#### Skill Discovery Paths

Copilot looks for skills in multiple locations:

1. `.github/skills/<name>/SKILL.md`
2. `.claude/skills/<name>/SKILL.md`
3. `.agents/skills/<name>/SKILL.md`

### CLI-Specific Configuration

For Copilot CLI usage:

- `~/.copilot/copilot-instructions.md` — local instructions
- `COPILOT_CUSTOM_INSTRUCTIONS_DIRS` — environment variable pointing to additional instruction directories
