<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT - template/pages/README.tpl.md
-->
# docs-ssot

**Single Source of Truth for the AI agent era.**

Generate README, CLAUDE.md, AGENTS.md, Cursor rules, Copilot instructions,
and VitePress docs — all from one source.

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Claude Code](https://img.shields.io/badge/Claude_Code-supported-blueviolet?logo=anthropic&logoColor=white)](https://docs.anthropic.com/en/docs/claude-code)
[![Cursor](https://img.shields.io/badge/Cursor-supported-00D1B2?logo=cursor&logoColor=white)](https://www.cursor.com/)
[![Codex](https://img.shields.io/badge/Codex-supported-412991?logo=openai&logoColor=white)](https://openai.com/codex)
[![GitHub Copilot](https://img.shields.io/badge/Copilot-supported-2088FF?logo=githubcopilot&logoColor=white)](https://github.com/features/copilot)

```text
                    ┌─── README.md
                    ├─── CLAUDE.md
                    ├─── AGENTS.md
  template/sections ├─── .claude/rules/*.md
  (single source) ──├─── .cursor/rules/*.mdc
                    ├─── .github/instructions/*.md
                    ├─── .agents/skills/
                    └─── VitePress docs site
```

---

## Why docs-ssot?

Software projects now maintain documentation for three different audiences:

📖 **For humans** — README, contributing guides, VitePress/Docusaurus docs sites
🤖 **For AI agents** — CLAUDE.md, AGENTS.md, .cursor/rules/, .github/instructions/
📋 **For both** — Architecture docs, coding rules, setup guides

And the list keeps growing:

| Audience | Files |
|----------|-------|
| Humans | README.md, VitePress docs, CONTRIBUTING.md |
| Claude Code | CLAUDE.md, .claude/rules/*.md, .claude/skills/ |
| Codex | AGENTS.md, .agents/skills/ |
| Cursor | .cursor/rules/*.mdc, .cursor/skills/ |
| Copilot | .github/copilot-instructions.md, .github/instructions/*.md |

Most of these files share the same underlying information —
architecture, coding rules, setup steps, testing strategy.

**The problem: Markdown has no `#include`.**

Every tool demands Markdown. Not YAML, not HTML — Markdown.
But Markdown has no way to share content across files.
So teams copy-paste. Then the copies drift. Information contradicts.

When humans read conflicting docs, they ask questions.
**When AI agents read conflicting docs, they silently act on the wrong one.**

An agent trusting stale architecture notes will refactor
the wrong module. An agent reading outdated rules will
bypass your testing strategy. And it won't ask — it will
just do it, confidently.

**Inconsistent documentation is the silent killer of AI-assisted development.**

docs-ssot solves this: write once, generate everywhere, always consistent.

---

## Before / After

### ❌ Before: Copy-paste chaos

```text
README.md              ← "Run tests with: make test"
CLAUDE.md              ← "Run tests with: go test ./..."
AGENTS.md              ← "Run tests with: make test-local"
.cursor/rules/go.mdc   ← "Always run go test before committing"
.github/instructions/   ← "Use make verify for pre-commit checks"
```

5 files. 5 different testing instructions.
An AI agent picks one — and skips your lint, coverage, and integration test pipeline.

### ✅ After: Single source of truth

```text
template/sections/development/testing.md    ← single source
           │
           ├──→ README.md
           ├──→ CLAUDE.md
           ├──→ AGENTS.md
           ├──→ .cursor/rules/go.mdc
           └──→ VitePress docs site

$ docs-ssot build
Generated 12 files from 1 source.
```

1 file. 1 version. Always consistent.

---

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

---

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

---

## How It Works

docs-ssot adds one missing feature to Markdown: **`#include`**.

### Include Directive

Compatible with [VitePress](https://vitepress.dev/) syntax:

```markdown
<!-- @include: path/to/file.md -->           Single file
<!-- @include: path/to/dir/ -->              All .md files in directory
<!-- @include: path/**/*.md -->              Recursive glob
<!-- @include: path/to/file.md level=+1 -->  Shift heading depth
```

Includes are resolved recursively. Circular includes are detected and cause a build error.

### The Pipeline

```text
template/sections/          template/pages/
(single-source docs)        (document structure)
        │                          │
        └──────────┬───────────────┘
                   ▼
          docs-ssot build
                   │
     ┌─────────────┼──────────────────┐
     ▼             ▼                  ▼
  README.md    CLAUDE.md    .cursor/rules/*.mdc
               AGENTS.md    .github/instructions/
                            .claude/rules/*.md
```

### Template Example

```markdown
<!-- template/pages/CLAUDE.tpl.md -->

# Project Context

<!-- @include: ../sections/project/overview.md -->

# Architecture

<!-- @include: ../sections/architecture/overview.md -->
<!-- @include: ../sections/architecture/pipeline.md -->

# Development Guide

<!-- @include: ../sections/development/ -->
```

One template defines the structure. Sections provide the content.
`docs-ssot build` assembles the final document.

---

## Commands

| Command | Description |
|---------|-------------|
| `docs-ssot build` | Generate final documents from templates |
| `docs-ssot check` | Check docs for SSOT violations by detecting near-duplicate sections |
| `docs-ssot include <file>` | Resolve includes and print expanded result to stdout |
| `docs-ssot migrate [files...]` | Decompose existing Markdown files into SSOT section structure |
| `docs-ssot migrate --from <tool>` | Migrate AI tool configs from one tool to others |
| `docs-ssot validate` | Validate documentation structure without generating output |
| `docs-ssot version` | Print the build version |

---

## SSOT Duplicate Detection

docs-ssot doesn't just generate — it helps you find existing duplication.

```sh
docs-ssot check --threshold 0.75
```

```text
score=0.891
A: docs/auth/overview.md [API > Authentication]
B: docs/setup/login.md [Setup > Authentication]
```

Uses TF-IDF cosine similarity to surface sections that say the same
thing in different places. Fix them before they confuse your AI agents.

---

## Why not Hugo / MkDocs / Docusaurus?

Static site generators build **websites**.
docs-ssot builds **any Markdown file from shared sources**.

|  | Hugo / MkDocs | docs-ssot |
|---|---|---|
| Output | HTML website | Any Markdown file |
| CLAUDE.md generation | ❌ | ✅ |
| .cursor/rules/ generation | ❌ | ✅ |
| AI agent config migration | ❌ | ✅ |
| Works alongside SSGs | — | ✅ (generates source .md for VitePress) |
| Duplicate detection | ❌ | ✅ |
| Markdown include syntax | Varies | VitePress-compatible |

docs-ssot is not a replacement for static site generators.
It sits **upstream** — generating the Markdown that SSGs then render.

```text
template/sections/ → docs-ssot build → docs/*.md → VitePress build → website
                                      → README.md
                                      → CLAUDE.md
                                      → .cursor/rules/
```

---

## This Repo Uses docs-ssot

The README you're reading, CLAUDE.md, AGENTS.md,
`.claude/rules/`, `.cursor/rules/`, `.github/instructions/` —
all generated from `template/sections/`.

See [docsgen.yaml](./template/sections/product/docsgen.yaml) for the full target list.

---

## Contributing

Thank you for your interest in contributing to `docs-ssot`.

---

### Development Workflow

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run checks before committing:

```sh
make go-fmt        # Format code
make go-lint       # Lint
make go-test       # Run tests
make docs          # Regenerate documentation
```

5. Commit using [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(processor): add URL include support
fix(generator): handle empty template gracefully
docs: update architecture overview
```

6. Push your branch and open a Pull Request

---

### Branch Naming

Use the prefix that matches the change type:

| Prefix | When to use |
|--------|-------------|
| `feature/` | New capability or behaviour |
| `fix/` | Bug fix |
| `refactor/` | Code restructuring without behaviour change |
| `chore/` | Maintenance, dependency updates, config |
| `docs/` | Documentation only |

---

### Code Quality Requirements

- Go 1.26+
- All code must pass `golangci-lint` (46+ linters enabled)
- Max line length: 200 characters
- Max cyclomatic complexity: 16
- Formatting: `gofumpt` (stricter than `gofmt`)

---

### Documentation Changes

This project uses SSOT (Single Source of Truth) for documentation.

- **Never edit** `README.md`, `CLAUDE.md`, or `AGENTS.md` directly
- **Edit source files** under `template/sections/`
- **Edit templates** under `template/pages/` to change document structure
- Run `make docs` to regenerate output files
- Commit both source and generated files together

---

### Testing

- Add unit tests for new functionality
- Integration tests should compare generated output with expected fixtures
- Ensure deterministic output: same input must always produce same output

```sh
make go-test                          # Run all tests
make docs && git diff --exit-code     # Verify generated files are up to date
```

---

### Pull Request Guidelines

- Keep PRs focused on a single concern
- Include a clear description of what and why
- Ensure all CI checks pass before requesting review
- Link related issues if applicable

---

## License

MIT
