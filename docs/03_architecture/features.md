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
| Front matter support | Planned | Parse and strip/merge YAML front matter from included files |

### CLI and Workflow Features

| Feature | Status | Notes |
|---------|--------|-------|
| `build` command | Implemented | Generates all output targets defined in `docsgen.yaml` |
| `include` command | Implemented | Expands includes in a file and prints the result to stdout; useful for debugging |
| `validate` command | Implemented | Dry-run over all templates; reports unresolvable includes without writing any output files |
| `version` command | Implemented | Prints the build version |
| Watch mode | Planned | Automatically rebuild on source file changes |
| Dry-run mode | Planned | Preview changes without writing output files |
| Diff / up-to-date check | Planned | Exit non-zero if generated files differ from committed versions (useful for CI) |
| Custom config file path | Planned | Allow specifying a non-default config file via CLI flag |

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
