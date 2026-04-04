# Commands Reference

This document describes the available CLI commands for docs-ssot.

## Overview

The CLI provides commands for generating documents from templates and managing documentation sources.

---

## docs build

Generate final documents (e.g., README.md, CLAUDE.md) from templates.

```
docs-ssot build
```

### What it does

- Reads template files
- Resolves `@include` directives
- Expands included Markdown files
- Writes final generated documents

---

## docs include

Resolve include directives and print the expanded result.

```
docs-ssot include template/README.tpl.md
```

Useful for debugging template expansion.

---

## docs validate

Validate documentation structure.

```
docs-ssot validate
```

### Validation includes

- Missing include files
- Circular includes
- Invalid paths
- Broken documentation structure

---

## docs clean

Remove generated files.

```
docs-ssot clean
```

Example files removed:

- README.md
- CLAUDE.md
- generated docs

---

## Typical Workflow

```
docs-ssot validate
docs-ssot build
```

Or during development:

```

docs-ssot include template/README.tpl.md

```

---

## Recommended Makefile Shortcuts

```
make docs
make docs-build
make docs-validate
make docs-clean
```
