## Linting

This project uses `golangci-lint` (v2) with 46+ linters enabled.

The linter is pinned as a Go tool dependency and invoked via `go tool golangci-lint`.

### Commands

| Command | Description |
|---------|-------------|
| `make go-lint` | Lint and auto-fix |
| `make go-lint-check` | Lint check only (no fix) |
| `make go-lint-fast` | Fast linters only with auto-fix |
| `make go-fmt` | Format all Go files with gofumpt |

### Key Rules

- **Max line length**: 200 characters
- **Max cyclomatic complexity**: 16
- **Formatting**: gofumpt (stricter than gofmt)

### Git Hooks (lefthook)

| Hook | Command | Trigger |
|------|---------|---------|
| `pre-commit` | `make go-fmt` | `*.go` files staged |
| `pre-push` | `make go-lint` | `*.go` files pushed |

Run `make install-dev` to set up hooks.
