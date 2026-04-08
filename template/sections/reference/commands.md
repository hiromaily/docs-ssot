## Commands Reference

This document describes the available CLI commands for docs-ssot.

### Overview

The CLI provides commands for generating documents from templates and managing documentation sources.

| Command | Description |
|---------|-------------|
| `docs-ssot build` | Generate final documents from templates |
| `docs-ssot check` | Check docs for SSOT violations by detecting near-duplicate sections |
| `docs-ssot include <file>` | Resolve includes and print expanded result to stdout |
| `docs-ssot migrate [files...]` | Decompose existing Markdown files into SSOT section structure |
| `docs-ssot migrate --agents` | Migrate AI tool configurations into SSOT multi-tool structure |
| `docs-ssot validate` | Validate documentation structure without generating output |
| `docs-ssot version` | Print the build version |

---

### docs build

Generate final documents (e.g., README.md, CLAUDE.md) from templates.

```
docs-ssot build
```

#### What it does

- Reads template files
- Resolves `@include` directives
- Expands included Markdown files
- Writes final generated documents

---

### docs check

Check docs for SSOT violations by detecting near-duplicate sections across Markdown files.

```
docs-ssot check [flags]
```

Uses TF-IDF cosine similarity to compare sections at the specified heading level. Sections scoring above the threshold are reported as potential SSOT violations — places where the same information exists in multiple source files.

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--root` | `docs` | Root directory to scan for Markdown files |
| `--threshold` | `0.82` | Similarity cutoff (0.0–1.0); pairs above this score are reported |
| `--min-chars` | `120` | Minimum character count for a section to be included in comparison |
| `--section-level` | `2` | Heading level used as section boundary (1–6) |
| `--format` | `text` | Output format: `text` or `json` |
| `--exclude` | — | Exclude path pattern (repeatable) |

#### Examples

Basic check with default settings:

```
docs-ssot check
```

Lower threshold to catch more candidates:

```
docs-ssot check --threshold 0.75
```

Compare at H3 level, exclude changelogs, output JSON:

```
docs-ssot check --section-level 3 --exclude docs/changelog/** --format json
```

#### Output

Text output (one block per similar pair):

```
score=0.891
A: docs/auth/overview.md [API > Authentication]
B: docs/setup/login.md [Setup > Authentication]
A title: Authentication
B title: Authentication
A snippet: Authentication tokens must be refreshed before they expire...
B snippet: Access tokens must be renewed prior to expiry...
----------------------------------------------------------------------------------------------------
```

A score of `1.0` means identical content; `0.82` (default threshold) catches near-duplicates while filtering loosely related content.

#### Exit behaviour

Exits `0` whether or not duplicates are found. Use `--format json` and inspect `result_count` in CI pipelines.

---

### docs migrate

Decompose existing monolithic Markdown files (e.g., README.md, CLAUDE.md) into the docs-ssot section structure.

```
docs-ssot migrate [files...] [flags]
```

This is the primary adoption command. It takes existing documentation files and converts them into modular, reusable sections with template files that reproduce the original document structure via `@include` directives.

#### What it does

1. **Splits** each input file by H2 headings into candidate sections
2. **Categorises** sections into directories (`project/`, `development/`, `architecture/`, `reference/`, `product/`, `misc/`) based on heading keyword heuristics
3. **Detects duplicates** across input files using TF-IDF cosine similarity (reuses the `check` command's engine)
4. **Creates section files** under `template/sections/<category>/<slug>.md`
5. **Creates template files** under `template/pages/<name>.tpl.md` with `@include` directives
6. **Creates `docsgen.yaml`** if it does not already exist
7. **Verifies round-trip** by running `build` and comparing output against originals

#### Section categorisation

Sections are assigned to categories based on heading keywords:

| Heading keywords | Category |
|-----------------|----------|
| Architecture, Design, System, Pipeline | `architecture/` |
| Overview, About, Introduction, Background | `project/` |
| Install, Setup, Getting Started, Prerequisites | `development/` |
| Test, Testing, CI | `development/` |
| Lint, Format, Code Quality | `development/` |
| Contributing, Contribute | `development/` |
| API, Commands, CLI, Reference | `reference/` |
| License, Changelog, Roadmap | `project/` |
| FAQ, Troubleshooting | `product/` |
| (fallback) | `misc/` |

#### Duplicate handling

When the same content appears in multiple input files:

1. TF-IDF cosine similarity is computed between all cross-file section pairs
2. Pairs scoring above the threshold are merged into a single section file
3. Both templates reference the shared section via `@include`

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--output-dir` | `template/sections` | Where to write section files |
| `--template-dir` | `template/pages` | Where to write template files |
| `--section-level` | `2` | Heading level used as section boundary (1–6) |
| `--threshold` | `0.82` | Similarity threshold for duplicate detection (0.0–1.0) |
| `--dry-run` | `false` | Print the migration plan without writing files |

#### Examples

Migrate existing README and CLAUDE.md:

```
docs-ssot migrate README.md CLAUDE.md
```

Preview migration plan without writing files:

```
docs-ssot migrate --dry-run README.md CLAUDE.md
```

Lower the duplicate detection threshold:

```
docs-ssot migrate --threshold 0.75 README.md
```

Split at H1 boundaries instead of H2:

```
docs-ssot migrate --section-level 1 README.md
```

#### Output

```
Parsed README.md: 8 sections
Parsed CLAUDE.md: 6 sections
Detected 3 duplicate sections (similarity > 0.82):
  "Architecture Overview" — merged into template/sections/architecture/overview.md (score=0.950)
  "Setup" — merged into template/sections/development/setup.md (score=1.000)
  "Testing" — merged into template/sections/development/testing.md (score=0.891)
Creating 11 unique section files in template/sections
  template/sections/project/overview.md
  template/sections/development/setup.md
  ...
Created template/pages/README.tpl.md (8 includes)
Created template/pages/CLAUDE.tpl.md (6 includes)
Created docsgen.yaml
Verifying round-trip...
Round-trip verification: OK
Migration complete.
```

#### Post-migration workflow

After `migrate`, the user's workflow becomes:

```sh
# Edit source sections
vim template/sections/development/setup.md

# Regenerate all outputs
docs-ssot build

# Verify
git diff README.md CLAUDE.md
```

---

#### Agent-aware migration (`--agents`)

With the `--agents` flag, `migrate` scans for AI tool configuration files (rules, skills, commands) and generates SSOT sections with per-tool templates for Claude, Cursor, Copilot, and Codex.

```
docs-ssot migrate --agents [flags]
```

##### What it does

1. **Scans** the repository for AI tool directories (`.claude/`, `.cursor/`, `.github/`, `.codex/`, `AGENTS.md`)
2. **Detects** which tool has the most configuration files (auto-selects as source of truth)
3. **Strips** frontmatter from source files and shifts H1→H2 headings
4. **Creates section files** under `template/sections/ai/<type>/<slug>.md`
5. **Generates templates** for each target tool with appropriate frontmatter
6. **Updates `docsgen.yaml`** with new build targets
7. **Verifies round-trip** by building and comparing against originals

##### Agent-specific flags

| Flag | Default | Description |
|------|---------|-------------|
| `--agents` | `false` | Enable agent-aware migration mode |
| `--source-tool` | `auto` | Source tool: `auto`, `claude`, `cursor`, `copilot` |
| `--target-tools` | `all` | Target tools: `all` or comma-separated (`claude,cursor,copilot,codex`) |
| `--dry-run` | `false` | Print the migration plan without writing files |

##### Examples

Migrate Claude rules to all tools:

```
docs-ssot migrate --agents
```

Preview migration plan:

```
docs-ssot migrate --agents --dry-run
```

Target specific tools only:

```
docs-ssot migrate --agents --target-tools cursor,copilot
```

Combine agent and file migration:

```
docs-ssot migrate --agents README.md CLAUDE.md
```

##### Output

```
Detected source tool: claude (3 files)
Target tools: claude, cursor, copilot, codex

Creating sections:
  template/sections/ai/rules/architecture.md
  template/sections/ai/rules/testing.md
  template/sections/ai/skills/deploy.md

Creating templates (4 tools × 3 files):
  claude: 3 templates
  cursor: 3 templates
  copilot: 3 templates
  codex: 2 templates

Updated docsgen.yaml (11 new targets)
Verifying round-trip...
Round-trip verification: OK
Agent migration complete.
```

---

### docs include

Resolve include directives and print the expanded result to stdout.

```
docs-ssot include <file>
```

Example:

```
docs-ssot include template/README.tpl.md
```

Useful for debugging template expansion without writing any output files.

---

### docs validate

Validate documentation structure without generating any output files.

```
docs-ssot validate
```

Performs a dry run over all templates in `docsgen.yaml`.

#### Validation checks

- Missing include files
- Circular includes
- Invalid paths

#### Output

Success:

```
OK
```

Failure (one line per failing template):

```
ERROR: include error (/path/to/file.md): open /path/to/file.md: no such file or directory
```

Exits with a non-zero status code when any error is found.

---

### docs version

Print the build version.

```
docs-ssot version
```

---

### Typical Workflow

```
docs-ssot validate
docs-ssot build
```

Or during development:

```
docs-ssot include template/README.tpl.md
```

---

### Recommended Makefile Shortcuts

```
make docs                                     # generate all output targets
make docs-validate                            # validate all templates
make docs-include FILE=template/README.tpl.md # expand and print a template
make docs-check                               # check docs for SSOT violations (default settings)
make docs-check ARGS="--threshold 0.75"       # check with custom flags
make docs-migrate FILES="README.md CLAUDE.md" # migrate existing docs to SSOT structure
make docs-migrate FILES="README.md" ARGS="--dry-run"  # preview migration plan
make docs-migrate-agents                      # migrate AI tool configs to SSOT structure
make docs-migrate-agents ARGS="--dry-run"     # preview agent migration plan
make docs-version                             # print the build version
```
