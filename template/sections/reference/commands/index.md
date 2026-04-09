## docs-ssot index

Generate `INDEX.md` showing include relationships and orphan detection across all templates.

```
docs-ssot index [flags]
```

### What it does

- Scans all templates defined in `docsgen.yaml`
- Resolves include relationships between templates and section files
- Detects orphaned section files not referenced by any template
- Prints the index to stdout, or writes it to a file

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--output` | — | Write index to file instead of stdout |

### Examples

Print index to stdout:

```
docs-ssot index
```

Write index to a file:

```
docs-ssot index --output template/INDEX.md
```
