## Documentation Pipeline Architecture

This document describes how the documentation generation pipeline works.

### Pipeline Overview

The `docs-ssot` system generates final documents (e.g., README.md, CLAUDE.md) from template files and modular Markdown sources.

The pipeline consists of the following stages:

1. Template Loading
2. Include Resolution
3. Recursive Expansion
4. Document Assembly
5. Output Generation

---

### Pipeline Flow

<!-- @include: docs/03_architecture/diagrams/pipeline-flow.md -->

---

### Step 1 — Template Loading

The system loads template files from the `template/` directory.

Example:

```
template/README.tpl.md
template/AGENTS.tpl.md
template/CLAUDE.tpl.md
```

Templates define the structure of the final documents.

---

### Step 2 — Include Resolution

Templates and Markdown files may contain include directives:

The following style is compatible with [VitePress](https://vitepress.dev/).

```markdown
<!-- @include: docs/01_project/overview.md -->
```

The include resolver replaces this directive with the contents of the referenced file.

---

### Step 3 — Recursive Expansion

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

### Step 4 — Resolving the link path (Planning)

When each component of a markdown file is expanded into a template, the link paths included in the content are adjusted according to the template's expansion destination.

```markdown
[docsgen.yaml](../../docsgen.yaml)
```

---

### Step 5 — Document Assembly

After all includes are expanded, the document builder assembles the final Markdown document.

This includes:

- Merging expanded content
- Ensuring correct order

---

### Step 6 — Output Generation

The final document is written to where defined at [docsgen.yaml](../../docsgen.yaml):

```
README.md
AGENTS.md
CLAUDE.md
```

These files are generated files and should not be edited directly.

---

### Include Resolution Detail

The include resolver processes directives recursively. The following diagram shows the exact resolution algorithm:

<!-- @include: docs/03_architecture/diagrams/include-resolution.md -->

---

### Design Goals

The pipeline is designed with the following goals:

- Single Source of Truth (SSOT)
- Modular documentation
- Reusable Markdown components
- Template-based document generation
- Deterministic document builds
- Simple and predictable behavior
