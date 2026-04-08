<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT - template/pages/README.tpl.md
-->
# docs-ssot

<!-- @include: ../sections/product/hero.md -->

---

<!-- @include: ../sections/product/why.md -->

---

<!-- @include: ../sections/product/before-after.md -->

---

<!-- @include: ../sections/product/quick-start.md -->

---

<!-- @include: ../sections/product/supported-targets.md -->

---

## How It Works

docs-ssot adds one missing feature to Markdown: **`#include`**.

<!-- @include: ../sections/architecture/includes-syntax.md level=+1 -->

### The Pipeline

```text
template/sections/          template/pages/
(single-source docs)        (document structure)
        │                          │
        └──────────┬───────────────┘
                   ▼
          docs-ssot build
                   │
     ┌─────────────┼──────────────────┐
     ▼             ▼                  ▼
  README.md    CLAUDE.md    .cursor/rules/*.mdc
               AGENTS.md    .github/instructions/
                            .claude/rules/*.md
```

### Template Example

```markdown
<!-- template/pages/CLAUDE.tpl.md -->

# Project Context

<!-- @include: ../sections/project/overview.md -->

# Architecture

<!-- @include: ../sections/architecture/overview.md -->
<!-- @include: ../sections/architecture/pipeline.md -->

# Development Guide

<!-- @include: ../sections/development/ -->
```

One template defines the structure. Sections provide the content.
`docs-ssot build` assembles the final document.

---

<!-- @include: ../sections/reference/commands-summary.md -->

---

## SSOT Duplicate Detection

docs-ssot doesn't just generate — it helps you find existing duplication.

```sh
docs-ssot check --threshold 0.75
```

```text
score=0.891
A: docs/auth/overview.md [API > Authentication]
B: docs/setup/login.md [Setup > Authentication]
```

Uses TF-IDF cosine similarity to surface sections that say the same
thing in different places. Fix them before they confuse your AI agents.

---

<!-- @include: ../sections/product/comparison.md -->

---

<!-- @include: ../sections/product/self-hosting.md -->

> **Documentation Site:** <https://hiromaily.github.io/docs-ssot/>

---

## Contributing

Contributions are welcome!

```sh
git clone https://github.com/hiromaily/docs-ssot.git
cd docs-ssot
make install-dev  # Install hooks and tools
make build        # Build the binary
make test         # Run tests
make docs         # Regenerate documentation
```

**Important:** Never edit `README.md`, `CLAUDE.md`, or `AGENTS.md` directly — edit source files under `template/sections/` and run `make docs`.

---

## License

MIT
