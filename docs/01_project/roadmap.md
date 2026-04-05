<!-- Status legend:
- (Released): Tagged and released.
- (Ready for Release): Implemented but not yet tagged; planned for release once all WIP items are complete.
- No status: Planned for implementation.
-->

## Roadmap

### v0.1 (Released)

- Single file include directive (`<!-- @include: path -->`)
- Recursive include resolution (included files may themselves contain include directives)
- Circular include detection (circular references produce a build error)
- Code fence passthrough (include directives inside fenced code blocks are treated as literal text)
- Multiple output targets via `docsgen.yaml`
- README, CLAUDE.md, AGENTS.md generation
- Link path rewriting — relative links and image URLs in all processed files are rewritten to be correct relative to the output file location

### v0.2 (Ready for Release)

- Heading level adjustment — optional `level=±N` parameter on include directives shifts the heading depth of included content (e.g. `<!-- @include: file.md level=+1 -->`)
- Directory include (`<!-- @include: docs/dir/ -->`) — include all `.md` files in a directory (sorted by filename)
- Glob include (`<!-- @include: docs/*.md -->`) — include files matching a glob pattern

### v0.3

- Recursive glob include (`<!-- @include: docs/**/*.md -->`) — include files matching a recursive glob
- `validate` command — check include paths and detect missing files without generating output
- Diff / up-to-date check — exit non-zero if generated files differ from committed versions (CI use)
- Dry-run mode — preview changes without writing output files
- ~~Watch mode — automatically rebuild on source file changes~~

### v0.4

- Variable substitution — allow `{{ variable }}` placeholders expanded at build time
- Front matter support — parse and strip/merge YAML front matter from included files
- Conditional includes — include or exclude sections based on build-time flags

### v0.5

- HTML output — convert generated Markdown to HTML
- PDF output — convert generated Markdown to PDF
- TOC generation — automatically insert a table of contents
