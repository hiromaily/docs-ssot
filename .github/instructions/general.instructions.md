---
applyTo: "**/*"
---

# General Development Rules for docs-ssot

## Project Purpose

`docs-ssot` is a CLI tool that generates documentation files (README.md, CLAUDE.md, AGENTS.md) from modular Markdown source files using a template-based composition system. It implements the Single Source of Truth (SSOT) principle for documentation.

---

## Critical Rule: Never Edit Generated Files

The following files are **build artifacts** — do NOT edit them directly:

- `README.md`
- `CLAUDE.md`
- `AGENTS.md`

Always edit the source files in `docs/` and regenerate with `make docs`.

---

## Repository Structure

```
docs-ssot/
├── cmd/docs-ssot/main.go       # CLI entry point (build command)
├── internal/
│   ├── config/config.go        # YAML config loader (Config, Target types)
│   ├── generator/generator.go  # Build orchestrator (Build func)
│   └── processor/              # Include resolver + transformer pipeline (ProcessFile func)
├── docs/                       # Source Markdown files (SSOT — edit here)
│   ├── 01_project/             # Project context and vision
│   ├── 02_product/             # Product concept and features
│   ├── 03_architecture/        # System architecture and pipeline
│   ├── 04_development/         # Setup, testing, linting guides
│   ├── 05_ai/                  # AI tool-specific docs (Claude, Cursor, Codex, etc.)
│   └── 06_reference/           # Commands and directory reference
├── template/                   # Document templates (define output structure)
│   ├── README.tpl.md
│   ├── CLAUDE.tpl.md
│   └── AGENTS.tpl.md
├── docsgen.yaml                # Build targets: template → output mapping
├── Makefile                    # Build, lint, test, docs commands
└── .golangci.yml               # Linting configuration (46+ linters)
```

---

## Go Package Responsibilities

| Package | File | Responsibility |
|---|---|---|
| `main` | `cmd/docs-ssot/main.go` | CLI arg parsing, dispatch to `generator.Build()` |
| `config` | `internal/config/config.go` | Load/save `docsgen.yaml` into `Config{Targets []Target}` |
| `generator` | `internal/generator/generator.go` | Iterate targets, call include resolver, write output |
| `processor` | `internal/processor/processor.go` | Include resolution, transformer pipeline (ProcessFile, Transformer, Apply) |
| `agentscan` | `internal/agentscan/agentscan.go` | Detect AI tools (.claude/, .cursor/, .github/, .codex/) and collect agent files |
| `frontmatter` | `internal/frontmatter/frontmatter.go` | Parse/strip/generate YAML frontmatter for different tool formats |
| `migrate` | `internal/migrate/agents.go` | Agent-aware migration pipeline (scan → section → template → config → verify) |

### Include Directive Pattern

```go
var includePattern = regexp.MustCompile(`^\s*<!--\s*@include:\s*(.*?)\s*-->\s*$`)
```

In Markdown, the directive looks like:

```markdown
<!-- @include: docs/01_project/overview.md -->
<!-- @include: docs/01_project/overview.md level=+1 -->
<!-- @include: docs/01_project/overview.md level=-1 -->
<!-- @include: docs/02_product/ -->
<!-- @include: docs/02_product/ level=+1 -->
<!-- @include: docs/02_product/*.md -->
<!-- @include: docs/02_product/*.md level=+1 -->
<!-- @include: docs/**/*.md -->
<!-- @include: docs/**/*.md level=+1 -->
```

The optional `level=±N` parameter shifts all ATX heading levels in the included content by N (clamped to `[1, 6]`). Headings inside code fences are not adjusted.

When the path ends with `/`, all `.md` files in that directory are included in sorted filename order (subdirectories are skipped). Combine with `level=±N` to adjust heading depths for the entire directory's content.

When the path contains `**`, all files matching the recursive glob pattern are included in sorted (lexical) path order. `**` matches zero or more path segments. If the root directory does not exist or no files match, no content is inserted (no error).

When the path contains glob metacharacters (`*`, `?`, `[`) but not `**`, all files matching the pattern are included in sorted (lexical) order. Directories matched by the pattern are skipped. If no files match, no content is inserted (no error).

`ProcessFile()` reads the template line-by-line, replaces include directives with file contents (applying heading adjustment if specified), and returns the assembled string.

---

## Build Pipeline

```
docsgen.yaml
  → config.Load()
  → for each target: include.ProcessFile(template)
  → os.WriteFile(output)
```

Recursive expansion is supported: included files may themselves contain include directives, resolved depth-first with circular reference detection.

---

## Common Commands

```sh
make build          # Compile: go build -o bin/docs-ssot ./cmd/docs-ssot
make docs           # Generate docs: go run ./cmd/docs-ssot build
make go-test        # Run tests: go test ./...
make go-lint        # Lint and auto-fix: golangci-lint run --fix
make go-lint-check  # Lint check only (no fix)
make go-fmt         # Format: golangci-lint fmt
make run            # go run ./cmd/docs-ssot build
make clean          # Remove bin/ and generated README.md, CLAUDE.md
```

---

## How to Add New Documentation

1. Create a new `.md` file in the appropriate `docs/` subdirectory.
2. Add an include directive where needed in the relevant `template/*.tpl.md` file:
   ```markdown
   <!-- @include: docs/04_development/new-guide.md -->
   ```
3. Run `make docs` to regenerate output files.
4. Commit both the source file and the regenerated output.

---

## How to Add a New Output Target

1. Create a new template file in `template/`, e.g., `template/MYFILE.tpl.md`.
2. Add a new entry in `docsgen.yaml`:
   ```yaml
   - input: template/MYFILE.tpl.md
     output: MYFILE.md
   ```
3. Run `make docs`.

---

## Current Limitations (Planned for Future)
- **No variable substitution**: No `{{ variable }}` placeholder support.

When implementing include-related features, the primary file to modify is `internal/processor/processor.go`.
To add a new content transformation, implement the `Transformer` interface in `internal/processor/` and register it in the relevant processing step.

---

## Code Quality Requirements

- **Go version**: 1.26.1
- **Linter**: `golangci-lint` via `go tool golangci-lint` (configured in `.golangci.yml`)
- **Active linters**: 46+ including `govet`, `staticcheck`, `gosec`, `errcheck`, `gofumpt`, `goimports`
- **Max line length**: 200 chars
- **Max cyclomatic complexity**: 16
- **Formatting**: `gofumpt` (stricter than `gofmt`)
- Always run `make go-lint` before committing Go changes.

---

## Testing Strategy

- Unit tests: cover include parsing, path resolution, circular detection, file loading.
- Integration tests: run generator on `testdata/`, compare output with `expected/` fixtures.
- Deterministic output required: same input must always produce same output.

After implementing features, verify with:

```sh
make go-test
make docs
git diff --exit-code README.md CLAUDE.md AGENTS.md
```

---

## Module Info

- **Module path**: `github.com/hiromaily/docs-ssot`
- **Key dependency**: `gopkg.in/yaml.v3` for YAML config parsing
- **Linting tools**: bundled as Go tool dependencies in `go.mod`
