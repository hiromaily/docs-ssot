# Design: Template Directory Restructure & INDEX.md Auto-Generation

## Status

Draft

## Summary

Restructure `template/` to separate concerns by role (pages / ai-agents / docs), reorganize `docs/` into sections / rules / commands, and add automatic INDEX.md generation with include relationship mapping and orphan detection.

## Motivation

### Problem 1: Role confusion in `template/`

The current `template/` directory mixes three distinct roles at the same level:

| Current location                     | Actual role                                              |
| ------------------------------------ | -------------------------------------------------------- |
| `template/*.tpl.md`                  | Output structure templates (README, CLAUDE, AGENTS)      |
| `template/claude/`, `cursor/`, etc.  | AI agent config templates (rules, skills, commands)      |
| `template/docs/`                     | SSOT source content                                      |

A new contributor cannot tell which files define structure and which hold content.

### Problem 2: Numbered directories hide role differences

Directories `01_project/` through `06_reference/` are content sections included by page templates. Directories `07_rules/` and `08_commands/` are snippets included exclusively by AI agent templates. The numbering scheme treats them as equals, but their roles in the include graph are fundamentally different.

### Problem 3: No discoverability

There is no catalog of what exists, what includes what, or what is unused. This is the root cause of duplication — people create new files because they cannot find existing ones (docs-index2 section 4).

---

## Design

### Part 1: Directory Restructure

#### Proposed structure

```text
template/
├── INDEX.md                          # Auto-generated catalog
│
├── pages/                            # Output templates (define final document structure)
│   ├── README.tpl.md
│   ├── CLAUDE.tpl.md
│   └── AGENTS.tpl.md
│
├── ai-agents/                        # AI agent config templates
│   ├── claude/
│   │   ├── rules/
│   │   │   ├── general.tpl.md
│   │   │   ├── docs.tpl.md
│   │   │   ├── git.tpl.md
│   │   │   ├── go.tpl.md
│   │   │   └── go-test.tpl.md
│   │   └── commands/
│   │       └── fix-pr-reviews.tpl.md
│   ├── codex/
│   │   ├── AGENTS.tpl.md
│   │   └── skills/
│   │       └── fix-pr-reviews/SKILL.tpl.md
│   ├── cursor/
│   │   ├── rules/
│   │   │   ├── general.tpl.mdc
│   │   │   ├── docs.tpl.mdc
│   │   │   ├── git.tpl.mdc
│   │   │   ├── go.tpl.mdc
│   │   │   └── go-test.tpl.mdc
│   │   └── skills/
│   │       └── fix-pr-reviews/SKILL.tpl.md
│   └── copilot/
│       ├── instructions/
│       │   ├── general.tpl.md
│       │   ├── docs.tpl.md
│       │   ├── git.tpl.md
│       │   ├── go.tpl.md
│       │   └── go-test.tpl.md
│       └── skills/
│           └── fix-pr-reviews/SKILL.tpl.md
│
└── docs/                             # SSOT source content
    ├── sections/                     # Content sections (included by pages)
    │   ├── project/
    │   │   ├── overview.md
    │   │   ├── vision.md
    │   │   └── roadmap.md
    │   ├── product/
    │   │   ├── concept.md
    │   │   └── features.md
    │   ├── architecture/
    │   │   ├── overview.md
    │   │   ├── system.md
    │   │   ├── pipeline.md
    │   │   ├── includes.md
    │   │   ├── features.md
    │   │   └── diagrams/
    │   │       ├── include-resolution.md
    │   │       └── pipeline-flow.md
    │   ├── development/
    │   │   ├── setup.md
    │   │   ├── test.md
    │   │   └── lint.md
    │   ├── ai/
    │   │   ├── overview.md
    │   │   ├── agents.md
    │   │   ├── claude.md
    │   │   ├── codex.md
    │   │   ├── cursor.md
    │   │   ├── github-copilot.md
    │   │   ├── cross-tool-mapping.md
    │   │   ├── best-practices.md
    │   │   └── glossary.md
    │   └── reference/
    │       ├── commands.md
    │       └── directory.md
    │
    ├── rules/                        # Rule definitions (included by ai-agents)
    │   ├── general.md
    │   ├── docs.md
    │   ├── git.md
    │   ├── go.md
    │   └── go-test.md
    │
    └── commands/                     # Command/skill procedures (included by ai-agents)
        └── fix-pr-reviews.md
```

#### Design principles

1. **Role-based separation**: `pages/` = output structure, `ai-agents/` = AI tool configs, `docs/` = content
2. **Content sub-roles**: `sections/` (page content), `rules/` (constraints), `commands/` (procedures)
3. **No number prefixes**: Include directives control order explicitly; names should be semantic
4. **Include direction is one-way and shallow**:
   - `pages/*.tpl.md` → `docs/sections/**/*.md`
   - `ai-agents/**/*.tpl.*` → `docs/rules/*.md`, `docs/commands/*.md`
   - `docs/` files do NOT include other `docs/` files as a general rule
   - **Exception**: `docs/sections/architecture/` files may include from `diagrams/` subdirectory. This is the only permitted intra-docs include, kept because diagrams are tightly coupled to the architecture sections that reference them.

#### Migration mapping

| Current                          | Proposed                                 |
| -------------------------------- | ---------------------------------------- |
| `template/*.tpl.md`              | `template/pages/*.tpl.md`                |
| `template/claude/`               | `template/ai-agents/claude/`             |
| `template/cursor/`               | `template/ai-agents/cursor/`             |
| `template/codex/`                | `template/ai-agents/codex/`              |
| `template/copilot/`              | `template/ai-agents/copilot/`            |
| `template/docs/01_project/`      | `template/docs/sections/project/`        |
| `template/docs/02_product/`      | `template/docs/sections/product/`        |
| `template/docs/03_architecture/` | `template/docs/sections/architecture/`   |
| `template/docs/04_development/`  | `template/docs/sections/development/`    |
| `template/docs/05_ai/`           | `template/docs/sections/ai/`             |
| `template/docs/06_reference/`    | `template/docs/sections/reference/`      |
| `template/docs/07_rules/`        | `template/docs/rules/`                   |
| `template/docs/08_commands/`     | `template/docs/commands/`                |

#### Include path updates

All `<!-- @include: ... -->` paths in template files must be updated to match the new directory structure.

Page templates (depth = `pages/`):

```markdown
<!-- before: template/README.tpl.md -->
<!-- @include: ./docs/01_project/overview.md -->

<!-- after: template/pages/README.tpl.md -->
<!-- @include: ../docs/sections/project/overview.md -->
```

AI agent templates (depth = `ai-agents/<tool>/rules/`):

```markdown
<!-- before: template/claude/rules/general.tpl.md (2 levels up to template/) -->
<!-- @include: ../../docs/07_rules/general.md -->

<!-- after: template/ai-agents/claude/rules/general.tpl.md (3 levels up to template/) -->
<!-- @include: ../../../docs/rules/general.md -->
```

#### docsgen.yaml updates

```yaml
targets:
  # Pages
  - input: template/pages/README.tpl.md
    output: README.md
  - input: template/pages/CLAUDE.tpl.md
    output: CLAUDE.md
  - input: template/pages/AGENTS.tpl.md
    output: AGENTS.md

  # AI agent configs
  - input: template/ai-agents/claude/rules/general.tpl.md
    output: .claude/rules/general.md
  # ... (all other ai-agent targets)

index:
  output: template/INDEX.md
```

---

### Part 2: INDEX.md Auto-Generation

#### Data sources

The index is generated by scanning the file system and parsing include directives:

| Data                               | Source                                                                          |
| ---------------------------------- | ------------------------------------------------------------------------------- |
| Page templates and their outputs   | `docsgen.yaml` targets                                                          |
| Section/rule/command files         | Recursive scan of `template/docs/`                                              |
| AI agent template files            | Recursive scan of `template/ai-agents/`                                         |
| Include relationships              | Parse `<!-- @include: ... -->` directives from all `.tpl.md` and `.tpl.mdc` files |

#### Include relationship analysis

For every file in `docs/`, build a reverse index showing which templates reference it:

```text
docs/sections/project/overview.md
  <- pages/README.tpl.md
  <- pages/CLAUDE.tpl.md
  <- pages/AGENTS.tpl.md

docs/rules/general.md
  <- ai-agents/claude/rules/general.tpl.md
  <- ai-agents/cursor/rules/general.tpl.mdc
  <- ai-agents/copilot/instructions/general.tpl.md
```

Glob includes (e.g., `docs/sections/ai/*.md`) and directory includes (e.g., `docs/sections/ai/`) must be expanded to resolve all matched files. Reuse the existing `processor.go` glob/directory expansion logic for consistency.

#### Orphan detection

A file in `docs/` that is not referenced by any include directive (directly or via glob/directory expansion) is an orphan. Report these in a dedicated section so they can be removed or wired in.

#### Generated output format

```markdown
<!-- AUTO-GENERATED FILE — DO NOT EDIT -->
<!-- Regenerate with: docs-ssot index -->
# Template Index

## Pages

| Template             | Output     | Sections included |
| -------------------- | ---------- | ----------------- |
| pages/README.tpl.md  | README.md  | 8                 |
| pages/CLAUDE.tpl.md  | CLAUDE.md  | 15                |
| pages/AGENTS.tpl.md  | AGENTS.md  | 14                |

## Sections

| File                                   | Referenced by          |
| -------------------------------------- | ---------------------- |
| docs/sections/project/overview.md      | README, CLAUDE, AGENTS |
| docs/sections/architecture/system.md   | README, CLAUDE, AGENTS |
| docs/sections/ai/codex.md              | AGENTS                 |
| ...                                    | ...                    |

## Rules

| File                  | Referenced by          |
| --------------------- | ---------------------- |
| docs/rules/general.md | claude, cursor, copilot |
| docs/rules/git.md     | claude, cursor, copilot |
| ...                   | ...                    |

## Commands

| File                              | Referenced by                  |
| --------------------------------- | ------------------------------ |
| docs/commands/fix-pr-reviews.md   | claude, codex, cursor, copilot |

## Orphans

| File   | Note                     |
| ------ | ------------------------ |
| (none) | All files are referenced |
```

#### CLI interface

New subcommand:

```sh
docs-ssot index                       # Print to stdout
docs-ssot index --output INDEX.md     # Write to file
```

If `index.output` is configured in `docsgen.yaml`, `docs-ssot build` also regenerates the index alongside other targets.

#### CI integration

```sh
docs-ssot build
git diff --exit-code template/INDEX.md
```

Ensures the committed INDEX.md matches the actual file tree and include graph.

---

## Implementation Order

### Phase 1: Directory restructure (no code changes)

1. Create new directory structure (`pages/`, `ai-agents/`, `docs/sections/`, `docs/rules/`, `docs/commands/`)
2. Move page templates: `template/*.tpl.md` → `template/pages/`
3. Move AI agent templates: `template/{claude,cursor,codex,copilot}/` → `template/ai-agents/`
4. Move content sections: `template/docs/01_project/` → `template/docs/sections/project/` (repeat for all `01`–`06`)
5. Move rules and commands: `template/docs/07_rules/` → `template/docs/rules/`, `template/docs/08_commands/` → `template/docs/commands/`
6. Update all `<!-- @include: ... -->` paths in page templates and AI agent templates
7. Update `docsgen.yaml` input paths
8. Run `docs-ssot validate` — all paths must resolve
9. Run `docs-ssot build && git diff --exit-code README.md CLAUDE.md AGENTS.md` — output must be bit-for-bit identical
10. Update documentation source files (`docs/sections/reference/directory.md` etc.) to reflect new structure

### Phase 2: INDEX.md auto-generation (code changes)

1. Add `internal/index/` package — file scanner and include directive parser
2. Build reverse include map (file → list of referencing templates)
3. Detect orphans (files with zero references)
4. Render INDEX.md in Markdown table format from collected data
5. Add `index` CLI subcommand (`internal/cli/index.go`)
6. Add `index` config section to `docsgen.yaml` schema (`internal/config/config.go`)
7. Integrate into `docs-ssot build` — generate index when `index.output` is configured
8. Add tests for index generation, orphan detection, glob/directory expansion, and deterministic output

---

## Test Plan

### Directory restructure

| Test                                              | Method                                                                    |
| ------------------------------------------------- | ------------------------------------------------------------------------- |
| Output files are identical after restructure      | `docs-ssot build && git diff --exit-code README.md CLAUDE.md AGENTS.md`   |
| All include paths resolve                         | `docs-ssot validate` exits 0                                              |
| No broken relative links in generated output      | Manual inspection or link checker                                         |

### INDEX.md generation

| #  | Test case                                                   | Expected                                                   |
| -- | ----------------------------------------------------------- | ---------------------------------------------------------- |
| 1  | Generate index for current template tree                    | All sections, rules, commands listed with correct referents |
| 2  | Add a new file to `docs/sections/` without including it     | Appears in Orphans section                                 |
| 3  | Remove a file that is still included                        | `docs-ssot validate` fails (separate from index)           |
| 4  | File included by multiple pages                             | All referencing templates listed                           |
| 5  | File included by zero templates                             | Listed as orphan                                           |
| 6  | Glob include (`docs/sections/ai/*.md`)                      | All matched files shown as included                        |
| 7  | Directory include (`docs/sections/ai/`)                     | All `.md` files in directory shown as included             |
| 8  | Recursive glob include (`docs/**/*.md`)                     | All matched files shown as included                        |
| 9  | Index output is deterministic                               | Two consecutive runs produce identical output              |
| 10 | `--output` flag writes to specified path                    | File created with expected content                         |

---

## Risks and Mitigations

| Risk                                                  | Impact                        | Mitigation                                                                                       |
| ----------------------------------------------------- | ----------------------------- | ------------------------------------------------------------------------------------------------ |
| Include path breakage after move                      | Build fails                   | Phase 1 ends with `docs-ssot validate` + `git diff --exit-code` verification                    |
| Large diff from rename                                | Hard to review                | Split into one commit per move category (pages, ai-agents, sections, rules)                      |
| INDEX.md generation performance                       | Slow build on large trees     | Scan is O(files x templates); for this project (<100 files) this is negligible                   |
| Glob/directory includes make reverse mapping complex  | Incorrect "referenced by" data | Reuse existing `processor.go` glob/directory expansion logic for consistency                     |
