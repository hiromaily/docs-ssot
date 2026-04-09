## Feature Status

This document is the single source of truth for the feature roadmap and implementation status of `docs-ssot`.
Other architecture documents should reference this file rather than duplicating status information.

### Include Resolver Features

| Feature | Status | Notes |
|---------|--------|-------|
| Single file include | Implemented | `<!-- @include: path/to/file.md -->` |
| Recursive include | Implemented | Included files may themselves contain include directives |
| Circular include detection | Implemented | Circular references produce a build error |
| Missing file error | Implemented | Missing included file stops the build with an error |
| Code fence passthrough | Implemented | Include directives inside fenced code blocks are treated as literal text |
| Directory include | Implemented | Include all `.md` files in a directory (sorted by filename); trailing `/` in path triggers directory mode |
| Glob include | Implemented | Include files matching a glob pattern (e.g. `*.md`); glob metacharacters (`*`, `?`, `[`) in path trigger glob mode |
| Recursive glob include | Implemented | Include files matching a recursive glob (e.g. `**/*.md`); `**` matches zero or more path segments |
| Link path rewriting | Implemented | Relative links and image URLs in all files are rewritten to be correct relative to the output file location |
| Heading level adjustment | Implemented | Optional `level=±N` parameter on include directives shifts heading depth of included content |
| Include from URL | Planned | Fetch and include a remote Markdown file |

### Generator Features

| Feature | Status | Notes |
|---------|--------|-------|
| Multiple output targets | Implemented | One `docsgen.yaml` can define many template → output pairs |
| Template-based generation | Implemented | Templates in `template/` define output structure |
| Deterministic output | Implemented | Same input always produces identical output |
| Variable substitution | Planned | Allow `{{ variable }}` placeholders expanded at build time |
| Conditional includes | Planned | Include or exclude sections based on build-time flags |
| Front matter support | Partial | Parse and strip YAML front matter implemented in `frontmatter` package; merge/pass-through not yet supported |

### CLI and Workflow Features

| Feature | Status | Notes |
|---------|--------|-------|
| `build` command | Implemented | Generates all output targets defined in `docsgen.yaml` |
| `check` command | Implemented | Scans docs for near-duplicate sections using TF-IDF cosine similarity; reports potential SSOT violations |
| `include` command | Implemented | Expands includes in a file and prints the result to stdout; useful for debugging |
| `validate` command | Implemented | Dry-run over all templates; reports unresolvable includes without writing any output files |
| `version` command | Implemented | Prints the build version |
| `migrate` command | Implemented | Decomposes existing Markdown files into SSOT section structure with duplicate detection and round-trip verification |
| Watch mode | Planned | Automatically rebuild on source file changes |
| Dry-run mode | Planned | Preview changes without writing output files |
| Diff / up-to-date check | Planned | Exit non-zero if generated files differ from committed versions (useful for CI) |
| Custom config file path | Planned | Allow specifying a non-default config file via CLI flag |

### Agent Migration Features (`migrate --from`)

| Feature | Status | Notes |
|---------|--------|-------|
| AI tool detection | Implemented | Scans `.claude/`, `.cursor/`, `.github/`, `.codex/`, `AGENTS.md` to detect configured tools |
| `--from` / `--to` flags | Implemented | Specify source and target tools; `--to` defaults to all tools except source |
| Rules migration | Implemented | Rules converted with tool-specific frontmatter (`.mdc` for Cursor, `applyTo` for Copilot, `@include` for Codex) |
| Skills migration | Implemented | Skills generated for all target tools with `name` + `description`; Claude preserves extra fields (`model`, `effort`, `allowed-tools`) |
| Subagent migration | Implemented | `.claude/agents/*.md` scanned and migrated to `.cursor/agents/`, `.github/agents/`, `.codex/agents/` (TOML format) |
| Command migration | Implemented | Claude commands migrated (Claude-only by default); convertible to skills via `--convert-commands` |
| `--convert-commands` | Implemented | Converts legacy `.claude/commands/` to cross-tool skills during migration |
| `--infer-globs` | Implemented | Infers path-gated rules from slug names (e.g., `go` → `**/*.go`, `frontend-*` → `frontend/**`) |
| Frontmatter parsing | Implemented | Full YAML parsing via `yaml.Unmarshal`; handles multi-line values (lists, maps) |
| CRLF handling | Implemented | Normalizes `\r\n` to `\n` before parsing |
| Config deduplication | Implemented | `docsgen.yaml` targets deduplicated on re-run (idempotent) |
| Round-trip verification | Implemented | Builds output and compares against source content after migration |
| Codex combined AGENTS.md | Implemented | All rules aggregated into single `AGENTS.tpl.md` via `@include` directives |

### Output Header Features

| Feature | Status | Notes |
|---------|--------|-------|
| Auto-generated file header | Planned | Prepend a `<!-- ⚠️ AUTO-GENERATED FILE — DO NOT EDIT -->` banner to all generated files |

### Output Format Features

| Feature | Status | Notes |
|---------|--------|-------|
| Markdown output | Implemented | Generated files are standard Markdown |
| HTML output | Planned | Convert generated Markdown to HTML |
| PDF output | Planned | Convert generated Markdown to PDF |
