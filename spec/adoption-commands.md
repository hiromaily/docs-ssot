# Design: Adoption Commands — init, migrate, drift

## Status

Draft

## Summary

Add three new CLI commands (`init`, `migrate`, `drift`) that lower the barrier to adopting docs-ssot in any repository. Together they cover the full adoption lifecycle: bootstrapping a new SSOT structure, converting existing documentation, and continuously enforcing SSOT integrity.

## Motivation

### Problem: High adoption barrier

Today, using docs-ssot in a new repository requires the user to:

1. Understand the SSOT concept and directory conventions
2. Manually create `template/sections/`, `template/pages/`, and `docsgen.yaml`
3. Decompose existing documentation (README.md, CLAUDE.md, etc.) into modular sections by hand
4. Set up CI or git hooks to keep generated files in sync

This is too much upfront work for most projects. The result is that only projects that already deeply understand SSOT will adopt the tool — the exact projects that need it least.

### Goal: "Just hand me your README"

The ideal adoption experience should be:

```
docs-ssot migrate README.md CLAUDE.md   # decompose existing docs
docs-ssot build                          # verify round-trip
```

Or for a brand-new project:

```
docs-ssot init                           # scaffold everything
```

And for ongoing maintenance:

```
docs-ssot drift                          # CI check: are generated files in sync?
```

---

## Proposed Commands

### 1. `docs-ssot init` — Scaffold SSOT structure

**Purpose:** Bootstrap a complete docs-ssot setup in an existing or new repository.

#### Behaviour

1. **Detect AI tools in use** — scan for `.claude/`, `.cursor/`, `.github/copilot-instructions.md`, `.codex/`, `AGENTS.md`, `CLAUDE.md`
2. **Generate directory structure:**
   ```
   template/
   ├── sections/          # empty section stubs or populated from migrate
   │   ├── project/
   │   ├── development/
   │   └── architecture/
   └── pages/
       ├── README.tpl.md
       ├── CLAUDE.tpl.md  # only if Claude detected
       └── AGENTS.tpl.md  # only if AGENTS.md detected or multi-tool
   docsgen.yaml
   ```
3. **Generate `docsgen.yaml`** — map each detected template to its output path
4. **Optionally generate CI workflow** — GitHub Actions snippet for `docs-ssot drift`
5. **Optionally generate git hook** — lefthook or husky config for pre-commit/pre-push build check

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--preset` | `auto` | Project type preset: `auto`, `go`, `node`, `python`, `generic` |
| `--ci` | `false` | Generate GitHub Actions workflow for drift detection |
| `--hooks` | `false` | Generate git hook configuration (lefthook/husky) |
| `--dry-run` | `false` | Print what would be created without writing files |

#### Presets

Presets populate section stubs with language/framework-appropriate placeholders:

| Preset | Sections generated |
|--------|-------------------|
| `go` | project overview, setup (Go install, make commands), testing (go test), linting (golangci-lint) |
| `node` | project overview, setup (npm/pnpm/bun), testing (vitest/jest), linting (eslint/biome) |
| `python` | project overview, setup (pip/poetry/uv), testing (pytest), linting (ruff) |
| `generic` | project overview, setup, testing, architecture |
| `auto` | Detect from go.mod / package.json / pyproject.toml / Cargo.toml and select preset |

#### Example

```sh
cd my-go-project
docs-ssot init --preset go --ci --hooks

# Created:
#   template/sections/project/overview.md
#   template/sections/development/setup.md
#   template/sections/development/testing.md
#   template/sections/development/linting.md
#   template/sections/architecture/overview.md
#   template/pages/README.tpl.md
#   template/pages/CLAUDE.tpl.md
#   docsgen.yaml
#   .github/workflows/docs-drift.yml
#   lefthook.yml (updated)
```

---

### 2. `docs-ssot migrate` — Decompose existing documents

**Purpose:** Convert existing monolithic Markdown files into the SSOT section structure.

This is the highest-impact feature. Nearly every project already has a README.md — being able to say "just give me your README and I'll SSOT-ify it" is the killer adoption story.

#### Behaviour

1. **Parse input files** — split each file by H2 headings into candidate sections
2. **Detect duplicates across files** — reuse the existing TF-IDF cosine similarity engine (`dupcheck` package) to identify overlapping sections between, e.g., README.md and CLAUDE.md
3. **Create section files** — write each unique section to `template/sections/<category>/<slug>.md`
4. **Create template files** — generate `template/pages/<name>.tpl.md` with `@include` directives that reproduce the original document structure
5. **Create `docsgen.yaml`** — if it does not already exist
6. **Verify round-trip** — run `build` internally and diff the output against the original files; warn if there are differences

#### Section categorisation

The command assigns sections to categories using heading-based heuristics:

| Heading keywords | Category |
|-----------------|----------|
| Overview, About, Introduction, Background | `project/` |
| Install, Setup, Getting Started, Prerequisites | `development/` |
| Architecture, Design, System, Pipeline | `architecture/` |
| Test, Testing, CI | `development/` |
| Lint, Format, Code Quality | `development/` |
| API, Commands, CLI, Reference | `reference/` |
| Contributing, Contribute | `development/` |
| License, Changelog, Roadmap | `project/` |
| FAQ, Troubleshooting | `product/` |
| (fallback) | `misc/` |

Users can override categorisation with `--section-map` or by editing after migration.

#### Duplicate handling

When the same content appears in multiple input files (e.g., an "Architecture" section in both README.md and CLAUDE.md):

1. Compute TF-IDF cosine similarity between candidate sections
2. If similarity > threshold (default 0.82, same as `check` command), merge into a single section file
3. Both templates reference the shared section via `@include`
4. Report merged sections to the user

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--output-dir` | `template/sections` | Where to write section files |
| `--template-dir` | `template/pages` | Where to write template files |
| `--section-level` | `2` | Heading level used as section boundary |
| `--threshold` | `0.82` | Similarity threshold for duplicate detection |
| `--dry-run` | `false` | Print the migration plan without writing files |
| `--section-map` | — | YAML file mapping heading patterns to categories |

#### Example

```sh
# Migrate existing README and CLAUDE.md
docs-ssot migrate README.md CLAUDE.md

# Output:
# Parsed README.md: 8 sections
# Parsed CLAUDE.md: 6 sections
# Detected 3 duplicate sections (similarity > 0.82):
#   "Architecture Overview" — merged into template/sections/architecture/overview.md
#   "Setup" — merged into template/sections/development/setup.md
#   "Testing" — merged into template/sections/development/testing.md
# Created 11 unique section files in template/sections/
# Created template/pages/README.tpl.md (8 includes)
# Created template/pages/CLAUDE.tpl.md (6 includes)
# Created docsgen.yaml
# Round-trip verification: OK (generated files match originals)
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

### 3. `docs-ssot drift` — Detect generated file divergence

**Purpose:** Exit non-zero if generated files differ from their current state on disk. Designed for CI pipelines and git hooks.

#### Behaviour

1. Run `build` to a temporary directory
2. Compare each generated file against the current on-disk version
3. If any file differs, print a diff summary and exit with code 1
4. If all files match, print "OK" and exit with code 0

#### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `text` | Output format: `text`, `json`, or `quiet` |
| `--fail-on-missing` | `true` | Also fail if a generated file does not exist on disk |

#### Example

```sh
# In CI
docs-ssot drift
# Exit 0: all generated files are in sync
# Exit 1: README.md is out of date (diff shown)
```

#### GitHub Actions integration

```yaml
- name: Check documentation drift
  run: docs-ssot drift
```

#### Relationship to existing commands

`drift` is essentially `build` + `diff --exit-code` packaged as a single command. The value is:

- No temporary file management needed by the user
- Clear exit codes for CI
- Structured output formats (JSON for programmatic use)
- `--fail-on-missing` catches the case where a new target was added to `docsgen.yaml` but never built

---

## Implementation Priority

| Command | Impact | Cost | Priority | Rationale |
|---------|--------|------|----------|-----------|
| `migrate` | Very High | Medium–High | **1** | Unlocks adoption for every project that already has a README. The TF-IDF engine already exists in `dupcheck`. |
| `init` | High | Low–Medium | **2** | Removes the cold-start problem for new projects. Mostly file generation and template logic. |
| `drift` | Medium | Low | **3** | Simple to implement (build to tmpdir + diff). Existing `validate` and `build && git diff` cover most of this already. |

### Suggested release plan

| Release | Commands |
|---------|----------|
| v0.4 | `drift` (low cost, high CI value, quick win) |
| v0.5 | `init` (scaffolding, presets) |
| v0.6 | `migrate` (section decomposition, duplicate detection, round-trip verification) |

Note: `migrate` has the highest impact but also the highest implementation complexity. Shipping `drift` and `init` first provides immediate value while `migrate` is being built.

---

## Architecture Considerations

### Reuse of existing packages

| Package | Reuse in new commands |
|---------|----------------------|
| `dupcheck` | `migrate` reuses TF-IDF cosine similarity for cross-file duplicate detection |
| `processor` | `drift` reuses `ProcessFile` for building to a temporary directory |
| `config` | All commands reuse `config.Load()` for `docsgen.yaml` |
| `generator` | `drift` reuses `generator.Build()` with an alternative output directory |

### New packages needed

| Package | Purpose |
|---------|---------|
| `splitter` | Parse Markdown into sections by heading level (used by `migrate`) |
| `categoriser` | Assign sections to categories based on heading heuristics (used by `migrate`) |
| `scaffold` | Generate directory structure, templates, and config files (used by `init`) |
| `differ` | Compare generated output against on-disk files and produce diffs (used by `drift`) |

### CLI integration

All new commands are added as Cobra subcommands under the existing CLI structure in `internal/cli/`.

---

## Open Questions

1. **`migrate` — should it back up original files?** If the user runs `migrate` on their README.md and then `build` overwrites it, the original is lost. Options: (a) always back up to `.docs-ssot-backup/`, (b) require `--force` to overwrite, (c) warn and require confirmation.

2. **`init` — interactive mode?** Should `init` support an interactive questionnaire (project name, description, which tools to target) or stay non-interactive with flags only?

3. **`migrate` — heading level flexibility.** Some projects use H1 for major sections and H2 for subsections. Should `--section-level` support splitting at H1 as well?

4. **`drift` — should it auto-fix?** Add a `--fix` flag that runs `build` and writes the updated files? Or keep `drift` as read-only and let users run `build` explicitly?

5. **Preset extensibility.** Should users be able to define custom presets via a YAML file, or are built-in presets sufficient for v1?
