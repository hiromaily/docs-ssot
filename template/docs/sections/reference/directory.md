## Directory Structure

```
docs-ssot/
├── cmd/docs-ssot/main.go       # CLI entry point
├── internal/
│   ├── cli/                    # Cobra subcommands (build, check, include, validate, version)
│   ├── config/config.go        # YAML config loader (docsgen.yaml)
│   ├── dupcheck/               # Near-duplicate section detector (TF-IDF cosine similarity)
│   ├── generator/generator.go  # Build orchestrator (Build, Validate)
│   └── processor/              # Include resolver + content transformers
│       ├── processor.go        # Core: ProcessFile, include resolution, glob/directory support
│       ├── heading.go          # HeadingTransformer: level=±N adjustment
│       ├── link.go             # LinkTransformer: relative path rewriting
│       └── transformer.go      # Transformer interface and Apply function
├── template/
│   ├── docs/                   # Source Markdown files (SSOT — edit here)
│   │   ├── 01_project/         # Project overview, vision, roadmap
│   │   ├── 02_product/         # Product concept and features
│   │   ├── 03_architecture/    # System architecture, pipeline, diagrams
│   │   ├── 04_development/     # Setup, testing, linting guides
│   │   ├── 05_ai/              # AI agent context and rules
│   │   └── 06_reference/       # Commands and directory reference
│   ├── README.tpl.md           # Template for README.md
│   ├── CLAUDE.tpl.md           # Template for CLAUDE.md
│   └── AGENTS.tpl.md           # Template for AGENTS.md
├── docsgen.yaml                # Build targets: template → output mapping
├── Makefile                    # Build, lint, test, docs commands
├── .golangci.yml               # Linting configuration (46+ linters)
├── .goreleaser.yaml            # Release automation
└── lefthook.yml                # Git hooks (pre-commit, pre-push)
```

### Generated Files (do not edit)

- `README.md` — generated from `template/README.tpl.md`
- `CLAUDE.md` — generated from `template/CLAUDE.tpl.md`
- `AGENTS.md` — generated from `template/AGENTS.tpl.md`
