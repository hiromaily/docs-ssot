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
