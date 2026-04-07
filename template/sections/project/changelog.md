## Changelog

### v0.3 (Latest Release)

- `validate` command — dry-run over all templates; reports unresolvable includes without writing output files; exits non-zero on failure
- `include` command — expand includes in a given file and print to stdout (debugging tool)
- `version` command — print the build version

### v0.2

- Heading level adjustment — optional `level=+N` parameter on include directives shifts the heading depth of included content
- Directory include (`<!-- @include: docs/dir/ -->`) — include all `.md` files in a directory (sorted by filename)
- Glob include (`<!-- @include: docs/*.md -->`) — include files matching a glob pattern
- Recursive glob include (`<!-- @include: docs/**/*.md -->`) — include files matching a recursive glob

### v0.1

- Single file include directive (`<!-- @include: path -->`)
- Recursive include resolution
- Circular include detection
- Code fence passthrough
- Multiple output targets via `docsgen.yaml`
- README, CLAUDE.md, AGENTS.md generation
- Link path rewriting
