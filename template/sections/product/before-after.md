## Before / After

### ❌ Before: Copy-paste chaos

```text
README.md              ← "Run tests with: make test"
CLAUDE.md              ← "Run tests with: go test ./..."
AGENTS.md              ← "Run tests with: make test-local"
.cursor/rules/go.mdc   ← "Always run go test before committing"
.github/instructions/   ← "Use make verify for pre-commit checks"
```

5 files. 5 different testing instructions.
An AI agent picks one — and skips your lint, coverage, and integration test pipeline.

### ✅ After: Single source of truth

```text
template/sections/development/testing.md    ← single source
           │
           ├──→ README.md
           ├──→ CLAUDE.md
           ├──→ AGENTS.md
           ├──→ .cursor/rules/go.mdc
           └──→ VitePress docs site

$ docs-ssot build
Generated 12 files from 1 source.
```

1 file. 1 version. Always consistent.
