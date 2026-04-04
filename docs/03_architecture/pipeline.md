# Documentation Pipeline Architecture

This document describes how the documentation generation pipeline works.

## Overview

The docs-ssot system generates final documents (e.g., README.md, CLAUDE.md) from template files and modular Markdown sources.

The pipeline consists of the following stages:

1. Template Loading
2. Include Resolution
3. Recursive Expansion
4. Document Assembly
5. Output Generation

---

## Pipeline Flow

```
docs/ (source markdown)
↓
template/*.tpl.md
↓
Include Resolver
↓
Expanded Markdown
↓
Document Builder
↓
README.md / CLAUDE.md (generated)
```

---

## Step 1 — Template Loading

The system loads template files from the `template/` directory.

Example:

```
template/README.tpl.md
template/CLAUDE.tpl.md
```

Templates define the structure of the final documents.

---

## Step 2 — Include Resolution

Templates and Markdown files may contain include directives:

```
<!-- @include: docs/01_project/overview.md -->
```

The include resolver replaces this directive with the contents of the referenced file.

---

## Step 3 — Recursive Expansion

Included files may also contain include directives.

The system resolves includes recursively until all includes are expanded.

```
A.tpl.md
→ includes B.md
→ includes C.md
```

Final result:

```
A + B + C
```

---

## Step 4 — Document Assembly

After all includes are expanded, the document builder assembles the final Markdown document.

This includes:

- Merging expanded content
- Ensuring correct order
- Adding headers/footers if necessary

---

## Step 5 — Output Generation

The final document is written to the project root:

```
README.md
CLAUDE.md
```

These files are generated files and should not be edited directly.

---

## Pipeline Summary

```
Template
↓
Load
↓
Resolve Includes
↓
Recursive Expansion
↓
Assemble Document
↓
Write Output
```

---

## Design Goals

The pipeline is designed with the following goals:

- Single Source of Truth
- Modular documentation
- Reusable Markdown components
- Template-based document generation
- Deterministic document builds
- Simple and predictable behavior
