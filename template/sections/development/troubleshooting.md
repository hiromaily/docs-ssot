## Troubleshooting

### Build Errors

#### `include error: open /path/to/file.md: no such file or directory`

An `@include` directive references a file that does not exist.

**Fix:** Check the path in the include directive. Paths are resolved relative to the file containing the directive, not the project root.

```markdown
<!-- In template/pages/README.tpl.md -->
<!-- @include: ../sections/project/overview.md -->
<!--            ^ relative to template/pages/ -->
```

---

#### `circular include detected`

Two or more files include each other, forming a loop.

**Fix:** Check the include chain. For example, if `A.md` includes `B.md` and `B.md` includes `A.md`, remove one of the references. Run `docs-ssot validate` to identify the cycle.

---

#### Generated files differ from committed versions

After running `make docs`, `git diff` shows unexpected changes in generated files.

**Fix:** This usually means the source files were edited but `make docs` was not run before committing. Always run:

```sh
make docs
git diff --exit-code README.md CLAUDE.md AGENTS.md
```

---

### Include Directives Not Expanding

#### Directive appears as literal text in the output

The include directive is inside a fenced code block and is intentionally treated as literal text.

**Fix:** If you want the directive to be expanded, move it outside the code fence. Include directives must be on their own line and not inside `` ``` `` blocks.

---

#### Directory include produces no output

A directory include (`<!-- @include: docs/somedir/ -->`) generates nothing.

**Fix:**
- Ensure the path ends with `/` (required for directory mode)
- Verify the directory contains `.md` files (subdirectories and non-`.md` files are skipped)
- Check that the path is correct relative to the file containing the directive

---

#### Glob include produces no output

A glob include (`<!-- @include: docs/*.md -->`) generates nothing.

**Fix:**
- Verify that files matching the pattern exist
- Glob includes silently produce no output if no files match (this is by design)
- Check path relativity — patterns are resolved from the file containing the directive

---

### Heading Level Issues

#### Headings are too deep or too shallow after inclusion

When using `level=+N` or `level=-N`, headings may end up at unexpected levels.

**Fix:** Check the source file's heading levels. Source files under `template/sections/` should start at `##` (H2). Combine with the `level` parameter to achieve the desired depth:

```markdown
<!-- level=0 (default): ## stays ## -->
<!-- level=+1: ## becomes ### -->
<!-- level=-1: ## becomes # -->
```

Heading levels are clamped to the valid range `[1, 6]`.

---

### Linting

#### `golangci-lint` fails with version error

**Fix:** This project uses `golangci-lint` v2 as a Go tool dependency. Do not install it globally — use:

```sh
go tool golangci-lint run
# or
make go-lint
```

---

### Getting Help

If your issue is not listed here:

1. Run `docs-ssot validate` to check for structural errors
2. Run `docs-ssot include <file>` to debug a specific template's expansion
3. Open an issue on [GitHub](https://github.com/hiromaily/docs-ssot/issues)
