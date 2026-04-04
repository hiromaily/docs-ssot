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
