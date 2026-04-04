# docs-ssot

## Overview

docs-ssot is a documentation single source of truth generator.

It composes README.md and CLAUDE.md from small modular markdown files.
## Vision

Documentation should be treated as a system, not as static files.

In many projects, documentation becomes fragmented, duplicated, and inconsistent.
The same information is rewritten across README files, design documents, and internal notes.

docs-ssot aims to solve this by applying the **Single Source of Truth (SSOT)** principle to Markdown documentation.

### Core Ideas

- **Write once, reuse everywhere**
  - Each piece of information exists in exactly one place
  - Reused via includes across multiple documents

- **Modular documentation**
  - Split documentation into small, composable Markdown files
  - Treat each file as a reusable unit

- **Docs as Code**
  - Documentation follows the same principles as software:
    - modularity
    - composition
    - build pipelines

- **Generated outputs**
  - Final documents (README, CLAUDE.md, etc.) are build artifacts
  - Never edited manually

### Why this matters

Without SSOT:
- Documentation diverges
- Updates are error-prone
- Context is duplicated and inconsistent

With docs-ssot:
- Documentation stays consistent
- Changes propagate automatically
- Different audiences (users, developers, AI) get tailored outputs from the same source

### Goal

To build a lightweight documentation system where:

- Markdown is the source of truth
- Templates define structure
- A generator composes final documents

Turning documentation into a **maintainable, scalable system**.


---

## Product

### Concept

Instead of maintaining large README files, this project splits documents into
small reusable markdown modules and composes them into final documents.

### Features

docs-ssot provides a simple documentation generation system based on Markdown includes and templates.

#### 1. Markdown Include System
Split large documents into small reusable Markdown files and include them where needed.

Example:

```md

<!-- @include: docs/01_project/overview.md -->

```

This allows documentation to be modular and reusable across multiple documents.

#### 2. Template-Based Document Structure
Document structure is defined in template files.

For example:
- README.template.md
- CLAUDE.template.md

Each template defines how documents are composed, while the actual content lives in the docs directory.

#### 3. Multiple Output Documents
The same source Markdown files can generate multiple documents:

- README.md (for GitHub)
- CLAUDE.md (for AI context)
- Documentation files
- Internal docs

This enables different audiences to receive different document structures from the same source.

#### 4. Single Source of Truth (SSOT)
Each piece of information exists in only one Markdown file.
All final documents are generated from these source files.

This prevents:
- duplicated documentation
- inconsistent information
- outdated README sections

#### 5. Recursive Includes
Included Markdown files can themselves include other files, allowing hierarchical document composition.

This enables building large documents from small components.

#### 6. Docs as Code Workflow
Documentation becomes a build artifact:

```
docs/        → source
template/    → structure
generator    → build tool
README.md    → output
```

This makes documentation maintainable, scalable, and version-controlled like code.


---

## Architecture

### Architecture Overview

The system consists of:

- Markdown modules
- Template files
- Generator CLI

### System Architecture

docs-ssot is composed of three main parts:

1. Markdown source files
2. Template files
3. Generator CLI

The generator reads template files, resolves include directives, and produces final documents such as README.md and CLAUDE.md.

### Components

#### docs/
The docs directory contains the Single Source of Truth Markdown files.
Each file represents a small, reusable piece of documentation.

These files should:
- be small
- be reusable
- contain only one topic
- not depend on document structure

#### template/
Template files define document structure.
They do not contain actual documentation content, only structure and include directives.

Examples:
- README.template.md
- CLAUDE.template.md

Templates decide:
- document order
- document sections
- which content appears in which output

#### Generator (docs-ssot)
The generator is a CLI tool that:
1. Reads a template file
2. Resolves include directives
3. Recursively expands included Markdown files
4. Writes the final output file

### Document Build Flow

The document generation flow works like this:

```
template/README.template.md
↓
include resolver
↓
docs/*.md
↓
combine
↓
README.md
```

### Design Principles

The system is designed with the following principles:

- Single Source of Truth
- Modular documentation
- Template-based composition
- Generated outputs
- Documentation as code
- Simple implementation
- No heavy static site generator


---

## Development

### Setup

```sh
make build
make docs
```


---

## AI

### AI Context

This repository uses docs-ssot, a documentation single source of truth system.

All documentation is written as small modular Markdown files under the `docs/` directory.
Final documents such as README.md and CLAUDE.md are generated from template files.

### How Documentation Works

Documentation is built using three main parts:

1. docs/ (Markdown source files)
2. template/ (document structure)
3. generator (include resolver and builder)

The generator reads template files and expands include directives like:

```

<!-- @include: docs/01_project/overview.md -->

```

Included files may also include other files (recursive includes).

### Important Rules

When editing documentation:

- Do NOT edit README.md directly
- Do NOT edit CLAUDE.md directly
- Edit files under docs/ instead
- Templates define document structure
- docs directory contains the source of truth

### Directory Roles

```

docs/       → documentation source (SSOT)
template/   → document templates
internal/   → generator implementation
cmd/        → CLI entrypoint
README.md   → generated output
CLAUDE.md   → generated output for AI context

```

### Documentation Philosophy

This project follows these principles:

- Single Source of Truth
- Modular documentation
- Documentation as Code
- Generated documents
- Reusable Markdown modules
- Template-based composition


## Reference

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

