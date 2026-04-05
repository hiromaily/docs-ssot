---
paths:
  - "**/*.md"
---

# Documentation Rules for docs-ssot

## SSOT Principle: Only Edit Source Files

This project generates documentation from source files. The pipeline is:

```
docs/**/*.md  +  template/*.tpl.md  +  docsgen.yaml
                          ↓
              docs-ssot build (make docs)
                          ↓
          README.md / AGENTS.md / CLAUDE.md
```

### Generated Files — Never Edit Directly

The following are **build artifacts** and must NOT be modified:

- `README.md`
- `AGENTS.md`
- `CLAUDE.md`

If you need to change content in these files, find the source file under `docs/` that contains that section and edit it there, then run `make docs`.

---

## Source Files — Where to Make Changes

All documentation content lives under `docs/`. Edit files here:

| Directory | Purpose |
|---|---|
| `docs/01_project/` | Project overview, vision, background |
| `docs/02_product/` | Product concept and feature descriptions |
| `docs/03_architecture/` | System architecture, pipeline, diagrams |
| `docs/04_development/` | Setup, testing, linting guides |
| `docs/05_ai/` | AI tool-specific instructions (Claude, Cursor, Codex, etc.) |
| `docs/06_reference/` | Commands reference, directory structure |

### Templates — Structure Only

Template files in `template/` define document structure using include directives. Modify templates only to change document structure (add/remove/reorder sections), not content.

```
template/README.tpl.md
template/AGENTS.tpl.md
template/CLAUDE.tpl.md
```

---

## Adding New Content

1. Create or edit the appropriate file under `docs/`.
2. If it's a new file, add an include directive in the relevant `template/*.tpl.md`:
   ```markdown
   <!-- @include: docs/XX_category/new-file.md -->
   ```
3. Run `make docs` to regenerate output files.
4. Commit both the source file and regenerated outputs together.

---

## Include Directive Format

Include directives follow the [VitePress](https://vitepress.dev/) style:

```markdown
<!-- @include: path/to/file.md -->
```

To include all `.md` files in a directory (sorted by filename), end the path with `/`:

```markdown
<!-- @include: docs/02_product/ -->
```

To include files matching a glob pattern (sorted lexically), use glob metacharacters (`*`, `?`, `[`):

```markdown
<!-- @include: docs/02_product/*.md -->
```

To include files matching a recursive glob (sorted lexically by full path), use `**`:

```markdown
<!-- @include: docs/**/*.md -->
```

`**` matches zero or more path segments, so `docs/**/*.md` matches both `docs/file.md` and `docs/sub/deep/file.md`.

An optional `level` parameter shifts the heading depth of the included content:

```markdown
<!-- @include: path/to/file.md level=+1 -->
<!-- @include: path/to/file.md level=-1 -->
<!-- @include: docs/02_product/ level=+1 -->
<!-- @include: docs/02_product/*.md level=+1 -->
<!-- @include: docs/**/*.md level=+1 -->
```

- `level=+N` deepens all headings by N levels (`##` → `###` for `+1`)
- `level=-N` shallows all headings by N levels (`###` → `##` for `-1`)
- Heading levels are clamped to `[1, 6]`; headings inside code fences are not adjusted

Other rules:

- Paths are resolved relative to the file containing the directive.
- Include directives inside code fences are treated as literal text (not expanded).
- Recursive includes are supported: included files may themselves contain include directives.
- Circular includes are detected and will cause a build error.

---

## Markdown Style Rules

- Each source file should cover one topic or section only (modular, single-responsibility).
- Do not duplicate content across multiple source files — the SSOT principle applies to docs too.
- Do not add front matter (`---`) to files under `docs/` — front matter is for templates if needed.
- Use ATX-style headings (`#`, `##`, `###`), not Setext-style (`===`, `---`).
- Use fenced code blocks with language identifiers (` ```go `, ` ```sh `, etc.).
- Prefer relative links when linking between docs source files.
- Do not hardcode generated file paths (e.g., `README.md`) in source docs — they are build artifacts.

---

## After Editing

Always regenerate and verify:

```sh
make docs
git diff README.md AGENTS.md CLAUDE.md
```

If generated files change unexpectedly, review which source file caused the change.
