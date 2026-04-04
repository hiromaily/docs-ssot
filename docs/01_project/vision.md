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
