## Setup

### Prerequisites

- Go 1.26+
- make

### Install

#### Homebrew (macOS/Linux)

```sh
brew tap hiromaily/tap
brew install docs-ssot
```

#### Go install

```sh
go install github.com/hiromaily/docs-ssot/cmd/docs-ssot@latest
```

#### Build from source

```sh
git clone https://github.com/hiromaily/docs-ssot.git
cd docs-ssot
make build
```

The binary is output to `bin/docs-ssot`.

### Quick Start

1. Create source Markdown files under `template/docs/`
2. Create template files under `template/` (e.g., `README.tpl.md`)
3. Define build targets in `docsgen.yaml`:

```yaml
targets:
  - input: template/README.tpl.md
    output: README.md
```

4. Generate documents:

```sh
docs-ssot build
```

### Development Setup

```sh
make install-dev   # Install lefthook and golangci-lint
make build         # Build the binary
make docs          # Generate documentation
make test          # Run tests
```
