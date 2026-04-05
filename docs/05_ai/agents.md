<!--
Used for AGENTS.md / CLAUDE.md
-->

This repository uses `docs-ssot`, a documentation single source of truth system.

All documentation is written as small modular Markdown files under the `docs/` directory.
Final documents such as README.md and CLAUDE.md are generated from template files.

## How Documentation Works

Documentation is built using three main parts:

1. `docs/` (Markdown source files)
2. `template/` (document structure)
3. generator (include resolver and builder)

The generator reads template files and expands include directives like:

```markdown
<!-- @include: docs/01_project/overview.md -->
```

Included files may also include other files (recursive includes).

## Important Rules

When editing documentation:

- Do NOT edit `README.md` directly
- Do NOT edit `CLAUDE.md` directly
- Edit files under `docs/` instead
- Templates define document structure
- docs directory contains the source of truth

## Directory Roles

```
docs/       → documentation source (SSOT)
template/   → document templates
internal/   → generator implementation
cmd/        → CLI entrypoint
README.md   → generated output
CLAUDE.md   → generated output for AI context
```

## Documentation Philosophy

This project follows these principles:

- Single Source of Truth
- Modular documentation
- Documentation as Code
- Generated documents
- Reusable Markdown modules
- Template-based composition
