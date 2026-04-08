## Supported Targets

### 📖 Human Documentation

| Target | Output |
|--------|--------|
| README | `README.md` |
| VitePress / Docusaurus | `docs/**/*.md` |
| Contributing guide | `CONTRIBUTING.md` |
| Any Markdown file | Configurable in `docsgen.yaml` |

### 🤖 AI Agent Instructions

| Agent | Output files |
|-------|-------------|
| [![Claude Code](https://img.shields.io/badge/Claude_Code-blueviolet?logo=anthropic&logoColor=white)](https://docs.anthropic.com/en/docs/claude-code) | `CLAUDE.md`, `.claude/rules/*.md`, `.claude/skills/`, `.claude/commands/` |
| [![Codex](https://img.shields.io/badge/Codex-412991?logo=openai&logoColor=white)](https://openai.com/codex) | `AGENTS.md`, `.agents/skills/` |
| [![Cursor](https://img.shields.io/badge/Cursor-00D1B2?logo=cursor&logoColor=white)](https://www.cursor.com/) | `.cursor/rules/*.mdc`, `.cursor/skills/` |
| [![GitHub Copilot](https://img.shields.io/badge/Copilot-2088FF?logo=githubcopilot&logoColor=white)](https://github.com/features/copilot) | `.github/copilot-instructions.md`, `.github/instructions/*.md`, `.github/skills/` |

> **All generated from the same `template/sections/` directory.**
> Change one source file → every target stays in sync.
