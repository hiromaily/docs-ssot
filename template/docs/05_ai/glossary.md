# Glossary

This glossary defines important terms used in this project so that AI agents and contributors use consistent terminology.

## Documentation System Terms

### SSOT (Single Source of Truth)

A design principle where documentation content exists in only one place.
All generated documents (e.g., README.md, AGENTS.md, CLAUDE.md) are built from the docs/ directory, which is the single source of truth.

### Docs Directory

The `docs/` directory contains all documentation source files.
These files are modular Markdown documents and should be edited instead of generated files.

### Template

Template files define the structure of generated documents.
Templates usually live in the `template/` directory and include documentation files using include directives.

Example:

```
<!-- @include: docs/01_project/overview.md -->
```

### Include Directive

A special comment directive used to include another Markdown file into a template or document.

Format:

```
<!-- @include: path/to/file.md -->
```

The include resolver replaces this directive with the contents of the referenced file.

### Include Resolver

A component that processes include directives and expands them into actual content.
It also handles recursive includes and circular include detection.

### Generator

The generator is the main program that builds final documents from templates and docs sources.

Responsibilities:

- Load templates
- Resolve includes
- Assemble documents
- Write generated files

### Pipeline

The documentation generation process consisting of multiple stages:

1. Template Loading
2. Include Resolution
3. Recursive Expansion
4. Document Assembly
5. Output Generation

### Generated Files

Files produced by the generator, such as:

- README.md
- CLAUDE.md

These files should not be edited manually.

### Template Expansion

The process of resolving include directives inside templates and Markdown files to produce a final document.

### Recursive Include

When an included file itself contains include directives that must also be resolved.

Example:

```
A.md includes B.md
B.md includes C.md
```

Final document becomes:

```
A + B + C
```

### Circular Include

A circular reference between included files.

Example:

```
A.md includes B.md
B.md includes A.md
```

The system must detect and prevent circular includes.

---

## Project Structure Terms

### Modular Documentation

Documentation written as small reusable Markdown files instead of one large document.

### Documentation as Code

Treating documentation like source code:

- Version controlled
- Modular
- Reviewed
- Generated
- Tested

### Template-Based Documentation

Final documents are not written directly.
Instead, templates define structure and content is included from source files.

---

## AI Documentation Terms

### CLAUDE.md

A generated document intended to provide context and instructions for AI agents working in this repository.

### AI Context

Information provided to AI tools so they understand:

- Project structure
- Documentation rules
- Architecture
- Terminology
- Development workflow
