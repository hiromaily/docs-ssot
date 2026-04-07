## Contributing

Thank you for your interest in contributing to `docs-ssot`.

---

### Development Workflow

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run checks before committing:

```sh
make go-fmt        # Format code
make go-lint       # Lint
make go-test       # Run tests
make docs          # Regenerate documentation
```

5. Commit using [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(processor): add URL include support
fix(generator): handle empty template gracefully
docs: update architecture overview
```

6. Push your branch and open a Pull Request

---

### Branch Naming

Use the prefix that matches the change type:

| Prefix | When to use |
|--------|-------------|
| `feature/` | New capability or behaviour |
| `fix/` | Bug fix |
| `refactor/` | Code restructuring without behaviour change |
| `chore/` | Maintenance, dependency updates, config |
| `docs/` | Documentation only |

---

### Code Quality Requirements

- Go 1.26+
- All code must pass `golangci-lint` (46+ linters enabled)
- Max line length: 200 characters
- Max cyclomatic complexity: 16
- Formatting: `gofumpt` (stricter than `gofmt`)

---

### Documentation Changes

This project uses SSOT (Single Source of Truth) for documentation.

- **Never edit** `README.md`, `CLAUDE.md`, or `AGENTS.md` directly
- **Edit source files** under `template/sections/`
- **Edit templates** under `template/pages/` to change document structure
- Run `make docs` to regenerate output files
- Commit both source and generated files together

---

### Testing

- Add unit tests for new functionality
- Integration tests should compare generated output with expected fixtures
- Ensure deterministic output: same input must always produce same output

```sh
make go-test                          # Run all tests
make docs && git diff --exit-code     # Verify generated files are up to date
```

---

### Pull Request Guidelines

- Keep PRs focused on a single concern
- Include a clear description of what and why
- Ensure all CI checks pass before requesting review
- Link related issues if applicable
