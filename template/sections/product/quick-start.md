## Quick Start

### Install

```sh
go install github.com/hiromaily/docs-ssot/cmd/docs-ssot@latest
```

### Try it in 30 seconds

```sh
# 1. Migrate your existing docs into SSOT structure
docs-ssot migrate README.md CLAUDE.md AGENTS.md

# 2. Check the generated config
cat docsgen.yaml
```

```yaml
# docsgen.yaml — defines what gets generated from where
targets:
  - input: template/pages/README.tpl.md
    output: README.md
  - input: template/pages/CLAUDE.tpl.md
    output: CLAUDE.md
  - input: template/pages/AGENTS.tpl.md
    output: AGENTS.md
```

```sh
# 3. Edit the single source
vim template/sections/development/testing.md

# 4. Regenerate everything
docs-ssot build
# → README.md, CLAUDE.md, AGENTS.md all updated. One edit, done.
```

### Migrate AI agent configs too

```sh
# Migrate Claude rules to Cursor, Codex, and Copilot
docs-ssot migrate --from claude
```

```yaml
# docsgen.yaml — agent configs added automatically
targets:
  # Documents
  - input: template/pages/README.tpl.md
    output: README.md
  - input: template/pages/CLAUDE.tpl.md
    output: CLAUDE.md

  # Claude Code rules
  - input: template/pages/ai-agents/claude/rules/go.tpl.md
    output: .claude/rules/go.md

  # Cursor rules (generated from same source)
  - input: template/pages/ai-agents/cursor/rules/go.tpl.mdc
    output: .cursor/rules/go.mdc

  # Copilot instructions (generated from same source)
  - input: template/pages/ai-agents/copilot/instructions/go.tpl.md
    output: .github/instructions/go.instructions.md
```
