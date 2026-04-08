# docs-ssot README Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign the docs-ssot README for global OSS appeal, restructure sections for fine-grained reuse, and update all templates to match.

**Architecture:** Split existing large section files into smaller reusable units, create new README-specific sections (hero, why, before-after, etc.), rewrite `README.tpl.md` using the new sections, and update `CLAUDE.tpl.md` / `AGENTS.tpl.md` to use directory includes where files were split.

**Tech Stack:** Markdown, docs-ssot include directives, shields.io badges

**Spec:** `poc/docs-ssot/spec/readme-redesign.md`

---

## File Map

### Files to Split (delete originals after splitting)

| Original | Split Into |
|----------|-----------|
| `template/sections/project/overview.md` (95 lines) | `project/overview.md`, `project/background.md`, `project/problem.md`, `project/solution.md` |
| `template/sections/reference/commands.md` (396 lines) | `reference/commands-summary.md`, `reference/commands/build.md`, `reference/commands/check.md`, `reference/commands/migrate.md`, `reference/commands/migrate-from.md`, `reference/commands/validate.md`, `reference/commands/include.md`, `reference/commands/version.md`, `reference/commands/workflow.md` |

### Files to Create (new README sections)

| File | Purpose |
|------|---------|
| `template/sections/product/hero.md` | Tagline + badges + hero diagram |
| `template/sections/product/why.md` | 3 audiences + pain points + silent killer message |
| `template/sections/product/before-after.md` | Testing instruction divergence example |
| `template/sections/product/supported-targets.md` | Human + AI target tables with badges |
| `template/sections/product/comparison.md` | Why not Hugo/MkDocs table + upstream diagram |
| `template/sections/product/self-hosting.md` | This repo uses docs-ssot |
| `template/sections/product/quick-start.md` | Install + migrate + docsgen.yaml + build |
| `template/sections/architecture/includes-syntax.md` | Include syntax quick reference (README-sized) |

### Templates to Rewrite/Update

| File | Change |
|------|--------|
| `template/pages/README.tpl.md` | Full rewrite using new sections |
| `template/pages/CLAUDE.tpl.md` | Update split-file references |
| `template/pages/AGENTS.tpl.md` | Update split-file references |

---

## Task 1: Split `project/overview.md` into 4 files

**Files:**

- Delete: `template/sections/project/overview.md`
- Create: `template/sections/project/overview.md` (new, smaller)
- Create: `template/sections/project/background.md`
- Create: `template/sections/project/problem.md`
- Create: `template/sections/project/solution.md`

All paths below are relative to `poc/docs-ssot/`.

- [ ] **Step 1: Create `template/sections/project/overview.md` (overview only)**

```markdown
## Overview

`docs-ssot` is a documentation Single Source of Truth (SSOT) generator.

It composes files such as README.md, CLAUDE.md, AGENTS.md, and other AI agent instruction files from small modular Markdown files.
```

- [ ] **Step 2: Create `template/sections/project/background.md`**

```markdown
## Background

AI-assisted development and AI agents are becoming a standard part of software development workflows.
Different AI tools and agents require different instruction and context files, for example:

- README.md
- AGENTS.md
- CLAUDE.md
- Agent-specific rule files like `.claude/rules`, `.cursor/rules`
- Development guidelines
- Architecture documentation

As the number of AI tools increases (Claude, Codex, Cursor, etc.), maintaining these files becomes difficult.

Common problems include:

- Documentation duplication
- Inconsistent information across files
- Outdated documentation
- Manual copy & paste maintenance
- Documentation drift over time

Maintaining multiple documentation files without duplication becomes increasingly difficult.
```

- [ ] **Step 3: Create `template/sections/project/problem.md`**

```markdown
## Problem

Documentation should follow the Single Source of Truth (SSOT) principle, but Markdown alone has limited reuse and composition capabilities.

Markdown is easy to write but lacks:

- File composition
- Reusable documentation modules
- Document templating
- Shared sections across multiple documents
- Structured documentation assembly

As a result, teams often duplicate content across multiple Markdown files.
```

- [ ] **Step 4: Create `template/sections/project/solution.md`**

```markdown
## Solution

`docs-ssot` solves this problem by introducing:

- Modular Markdown documentation
- Template-based document structure
- Include directives for Markdown files
- Generated documentation files
- Single Source of Truth documentation architecture

Instead of writing large README files directly, documentation is split into small reusable Markdown modules and assembled into final documents using templates.
```

- [ ] **Step 5: Delete the original `template/sections/project/overview.md`**

The original 95-line file has been replaced by the 4 new files above. The existing `template/sections/project/concept.md` already contains the Concept section, so no additional file is needed.

- [ ] **Step 6: Verify the split covers the original**

Run:

```sh
cat template/sections/project/overview.md template/sections/project/background.md template/sections/project/problem.md template/sections/project/solution.md
```

Confirm all content from the original overview.md is present across the 4 files. The `Concept` section (lines 65-95 of the original) already exists in `template/sections/project/concept.md`.

- [ ] **Step 7: Commit**

```sh
git add template/sections/project/overview.md template/sections/project/background.md template/sections/project/problem.md template/sections/project/solution.md
git commit -m "refactor(sections): split project/overview.md into overview, background, problem, solution"
```

---

## Task 2: Split `reference/commands.md` into individual command files

**Files:**

- Delete: `template/sections/reference/commands.md`
- Create: `template/sections/reference/commands-summary.md`
- Create: `template/sections/reference/commands/build.md`
- Create: `template/sections/reference/commands/check.md`
- Create: `template/sections/reference/commands/migrate.md`
- Create: `template/sections/reference/commands/migrate-from.md`
- Create: `template/sections/reference/commands/include.md`
- Create: `template/sections/reference/commands/validate.md`
- Create: `template/sections/reference/commands/version.md`
- Create: `template/sections/reference/commands/workflow.md`

- [ ] **Step 1: Create `template/sections/reference/commands-summary.md`**

Extract lines 1-17 of the original (the overview table only):

```markdown
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
```

- [ ] **Step 2: Create `template/sections/reference/commands/build.md`**

Extract lines 21-35 of the original:

```markdown
## docs-ssot build

Generate final documents (e.g., README.md, CLAUDE.md) from templates.

```text
docs-ssot build
```

### What it does

- Reads template files
- Resolves `@include` directives
- Expands included Markdown files
- Writes final generated documents
```

- [ ] **Step 3: Create `template/sections/reference/commands/check.md`**

Extract lines 38-98 of the original (the full `docs check` section including Flags, Examples, Output, Exit behaviour).

- [ ] **Step 4: Create `template/sections/reference/commands/migrate.md`**

Extract lines 102-217 of the original (the full `docs migrate` section including What it does, Section categorisation, Duplicate handling, Flags, Examples, Output, Post-migration workflow).

- [ ] **Step 5: Create `template/sections/reference/commands/migrate-from.md`**

Extract lines 221-302 of the original (the full agent-aware migration section including What it does, Flags, Examples, Output).

- [ ] **Step 6: Create `template/sections/reference/commands/include.md`**

Extract lines 306-320 of the original:

```markdown
## docs-ssot include

Resolve include directives and print the expanded result to stdout.

```text
docs-ssot include <file>
```

Example:

```text
docs-ssot include template/README.tpl.md
```

Useful for debugging template expansion without writing any output files.
```

- [ ] **Step 7: Create `template/sections/reference/commands/validate.md`**

Extract lines 324-356 of the original (the full validate section including Validation checks and Output).

- [ ] **Step 8: Create `template/sections/reference/commands/version.md`**

Extract lines 358-364 of the original:

```markdown
## docs-ssot version

Print the build version.

```text
docs-ssot version
```
```

- [ ] **Step 9: Create `template/sections/reference/commands/workflow.md`**

Extract lines 366-396 of the original (Typical Workflow and Recommended Makefile Shortcuts).

- [ ] **Step 10: Delete the original `template/sections/reference/commands.md`**

- [ ] **Step 11: Commit**

```sh
git add template/sections/reference/
git commit -m "refactor(sections): split reference/commands.md into individual command files"
```

---

## Task 3: Create `architecture/includes-syntax.md` (README-sized quick reference)

**Files:**

- Create: `template/sections/architecture/includes-syntax.md`

- [ ] **Step 1: Create `template/sections/architecture/includes-syntax.md`**

```markdown
## Include Directive

Compatible with [VitePress](https://vitepress.dev/) syntax:

```markdown
<!-- @include: path/to/file.md -->           Single file
<!-- @include: path/to/dir/ -->              All .md files in directory
<!-- @include: path/**/*.md -->              Recursive glob
<!-- @include: path/to/file.md level=+1 -->  Shift heading depth
```

Includes are resolved recursively. Circular includes are detected and cause a build error.
```

- [ ] **Step 2: Commit**

```sh
git add template/sections/architecture/includes-syntax.md
git commit -m "feat(sections): add includes-syntax.md for README quick reference"
```

---

## Task 4: Create new README sections (hero, why, before-after, supported-targets, quick-start, comparison, self-hosting)

**Files:**

- Create: `template/sections/product/hero.md`
- Create: `template/sections/product/why.md`
- Create: `template/sections/product/before-after.md`
- Create: `template/sections/product/supported-targets.md`
- Create: `template/sections/product/quick-start.md`
- Create: `template/sections/product/comparison.md`
- Create: `template/sections/product/self-hosting.md`

- [ ] **Step 1: Create `template/sections/product/hero.md`**

```markdown
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
```

- [ ] **Step 2: Create `template/sections/product/why.md`**

```markdown
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
```

- [ ] **Step 3: Create `template/sections/product/before-after.md`**

```markdown
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
```

- [ ] **Step 4: Create `template/sections/product/quick-start.md`**

```markdown
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
```

- [ ] **Step 5: Create `template/sections/product/supported-targets.md`**

```markdown
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
```

- [ ] **Step 6: Create `template/sections/product/comparison.md`**

```markdown
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
```

- [ ] **Step 7: Create `template/sections/product/self-hosting.md`**

```markdown
## This Repo Uses docs-ssot

The README you're reading, CLAUDE.md, AGENTS.md,
`.claude/rules/`, `.cursor/rules/`, `.github/instructions/` —
all generated from `template/sections/`.

See [docsgen.yaml](./docsgen.yaml) for the full target list.
```

- [ ] **Step 8: Commit**

```sh
git add template/sections/product/hero.md template/sections/product/why.md template/sections/product/before-after.md template/sections/product/quick-start.md template/sections/product/supported-targets.md template/sections/product/comparison.md template/sections/product/self-hosting.md
git commit -m "feat(sections): add README redesign sections (hero, why, before-after, quick-start, targets, comparison, self-hosting)"
```

---

## Task 5: Rewrite `README.tpl.md`

**Files:**

- Modify: `template/pages/README.tpl.md`

- [ ] **Step 1: Rewrite `template/pages/README.tpl.md`**

```markdown
<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT - template/pages/README.tpl.md
-->
# docs-ssot

<!-- @include: ../sections/product/hero.md -->

---

<!-- @include: ../sections/product/why.md -->

---

<!-- @include: ../sections/product/before-after.md -->

---

<!-- @include: ../sections/product/quick-start.md -->

---

<!-- @include: ../sections/product/supported-targets.md -->

---

## How It Works

docs-ssot adds one missing feature to Markdown: **`#include`**.

<!-- @include: ../sections/architecture/includes-syntax.md level=+1 -->

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

<!-- @include: ../sections/reference/commands-summary.md -->

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

<!-- @include: ../sections/product/comparison.md -->

---

<!-- @include: ../sections/product/self-hosting.md -->

---

<!-- @include: ../sections/development/contributing.md -->

---

## License

MIT
```

- [ ] **Step 2: Commit**

```sh
git add template/pages/README.tpl.md
git commit -m "feat(readme): rewrite README template with new design"
```

---

## Task 6: Update `CLAUDE.tpl.md` to use split sections

**Files:**

- Modify: `template/pages/CLAUDE.tpl.md`

- [ ] **Step 1: Update `template/pages/CLAUDE.tpl.md`**

Replace the single `project/overview.md` include with the split files, and replace the single `reference/commands.md` include with a directory include:

```markdown
<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT - template/CLAUDE.tpl.md
-->
# Project Context

<!-- @include: ../sections/project/overview.md -->

<!-- @include: ../sections/project/background.md -->

<!-- @include: ../sections/project/problem.md -->

<!-- @include: ../sections/project/solution.md -->

<!-- @include: ../sections/project/concept.md -->

---

<!-- @include: ../sections/ai/agents.md -->

---

# Repository Structure

<!-- @include: ../sections/reference/directory.md -->

---

# Architecture

<!-- @include: ../sections/architecture/overview.md -->

<!-- @include: ../sections/architecture/system.md -->

<!-- @include: ../sections/architecture/pipeline.md -->

---

<!-- @include: ../sections/architecture/includes.md level=-1 -->

---

# Development Guide

<!-- @include: ../sections/development/setup.md -->

<!-- @include: ../sections/development/test.md -->

<!-- @include: ../sections/development/lint.md -->

---

# Commands Reference

<!-- @include: ../sections/reference/commands-summary.md level=-1 -->

<!-- @include: ../sections/reference/commands/ level=-1 -->

---

# AI Agent Configuration

<!-- @include: ../sections/ai/overview.md -->

<!-- @include: ../sections/ai/claude.md -->

<!-- @include: ../sections/ai/hooks.md -->

<!-- @include: ../sections/ai/cross-tool-mapping.md -->

<!-- @include: ../sections/ai/best-practices.md -->

---

<!-- @include: ../sections/architecture/features.md level=-1 -->

---

<!-- @include: ../sections/ai/glossary.md level=-1 -->
```

- [ ] **Step 2: Commit**

```sh
git add template/pages/CLAUDE.tpl.md
git commit -m "refactor(claude): update CLAUDE template for split sections"
```

---

## Task 7: Update `AGENTS.tpl.md` to use split sections

**Files:**

- Modify: `template/pages/AGENTS.tpl.md`

- [ ] **Step 1: Update `template/pages/AGENTS.tpl.md`**

Apply same changes as CLAUDE.tpl.md — split project/overview and directory include for commands:

```markdown
<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT - template/AGENTS.tpl.md
-->
# Project Context

<!-- @include: ../sections/project/overview.md -->

<!-- @include: ../sections/project/background.md -->

<!-- @include: ../sections/project/problem.md -->

<!-- @include: ../sections/project/solution.md -->

<!-- @include: ../sections/project/concept.md -->

---

<!-- @include: ../sections/ai/agents.md -->

---

# Repository Structure

<!-- @include: ../sections/reference/directory.md -->

---

# Architecture

<!-- @include: ../sections/architecture/overview.md -->

<!-- @include: ../sections/architecture/system.md -->

---

<!-- @include: ../sections/architecture/includes.md level=-1 -->

---

# Development Guide

<!-- @include: ../sections/development/setup.md -->

<!-- @include: ../sections/development/test.md -->

<!-- @include: ../sections/development/lint.md -->

---

# Commands Reference

<!-- @include: ../sections/reference/commands-summary.md level=-1 -->

<!-- @include: ../sections/reference/commands/ level=-1 -->

---

# AI Agent Configuration

<!-- @include: ../sections/ai/overview.md -->

<!-- @include: ../sections/ai/claude.md -->

<!-- @include: ../sections/ai/codex.md -->

<!-- @include: ../sections/ai/cursor.md -->

<!-- @include: ../sections/ai/github-copilot.md -->

<!-- @include: ../sections/ai/cross-tool-mapping.md -->

<!-- @include: ../sections/ai/best-practices.md -->

---

<!-- @include: ../sections/ai/glossary.md level=-1 -->
```

- [ ] **Step 2: Commit**

```sh
git add template/pages/AGENTS.tpl.md
git commit -m "refactor(agents): update AGENTS template for split sections"
```

---

## Task 8: Build, verify, and fix

**Files:**

- All generated outputs (README.md, CLAUDE.md, AGENTS.md, etc.)

- [ ] **Step 1: Run `docs-ssot validate`**

```sh
cd poc/docs-ssot && make docs-validate
```

Expected: `OK` — all includes resolve. If errors appear, fix the paths in the template files.

- [ ] **Step 2: Run `docs-ssot build`**

```sh
cd poc/docs-ssot && make docs
```

Expected: All output files regenerated.

- [ ] **Step 3: Review the generated README.md**

```sh
head -80 poc/docs-ssot/README.md
```

Verify:

- Hero section with tagline, badges, and diagram appears first
- Why section with 3 audiences and pain points
- Before/After section with testing example
- Quick Start with docsgen.yaml examples
- No broken include directives (no raw `<!-- @include:` in output)

- [ ] **Step 4: Review generated CLAUDE.md and AGENTS.md**

```sh
head -40 poc/docs-ssot/CLAUDE.md
head -40 poc/docs-ssot/AGENTS.md
```

Verify:

- Project Context section includes overview + background + problem + solution + concept
- Commands Reference section includes summary table + all individual command details
- No missing content compared to previous versions

- [ ] **Step 5: Check for regressions in other generated files**

```sh
cd poc/docs-ssot && git diff --stat
```

Review the diff. Only README.md should have major changes. CLAUDE.md and AGENTS.md should have content intact (possibly reordered due to directory includes).

- [ ] **Step 6: Fix any issues found**

If any include paths are wrong, fix them in the template files and re-run `make docs`.

- [ ] **Step 7: Commit all generated files**

```sh
cd poc/docs-ssot && git add README.md CLAUDE.md AGENTS.md AGENTS-codex.md .claude/ .cursor/ .github/ .agents/
git commit -m "docs: regenerate all outputs with redesigned README and split sections"
```

---

## Task 9: Final review and cleanup

- [ ] **Step 1: Read through the full generated README.md end-to-end**

Verify the flow:

1. Hero (tagline + badges + diagram)
2. Why docs-ssot? (3 audiences + Markdown limitation + AI silent killer)
3. Before / After (testing instructions example)
4. Quick Start (install + migrate + docsgen.yaml + build)
5. Supported Targets (human + AI tables with badges)
6. How It Works (include syntax + pipeline + template example)
7. Commands (summary table)
8. SSOT Duplicate Detection (check command highlight)
9. Why not Hugo/MkDocs? (comparison table + upstream diagram)
10. This Repo Uses docs-ssot (self-hosting)
11. Contributing
12. License

- [ ] **Step 2: Verify shields.io badges render correctly**

Open the README.md in a Markdown previewer or push to a branch and check GitHub rendering. Confirm all 6 badges (Go, MIT, Claude, Cursor, Codex, Copilot) render with logos.

- [ ] **Step 3: Check for stale references**

Search for any references to the old file structure:

```sh
cd poc/docs-ssot && grep -r "project/overview.md" template/pages/
cd poc/docs-ssot && grep -r "reference/commands.md" template/pages/
```

Expected: No matches (all references should be updated).

- [ ] **Step 4: Squash or tidy commits if desired, then push**

```sh
git push origin HEAD
```
