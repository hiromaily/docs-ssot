## Template Design Rationale

### Shared base: `sections/ai/agents-base.md`

`agents-base.md` is a shared template fragment included by both `AGENTS.tpl.md` and `CLAUDE.tpl.md`. It covers the full project context, repository structure, architecture, development guide, and commands reference.

These two templates generate **comprehensive documentation files** — their readers need the full project context to work effectively.

### Independent template: `AGENTS-codex.tpl.md`

`AGENTS-codex.tpl.md` does **not** use `agents-base.md`. This is intentional.

The Codex AGENTS file is a **focused AI instruction file**, not a comprehensive documentation reference. Its structural differences from the shared base are deliberate:

| Aspect | `AGENTS.tpl.md` / `CLAUDE.tpl.md` (via `agents-base.md`) | `AGENTS-codex.tpl.md` |
|---|---|---|
| Project Context | Full — overview, background, problem, solution, concept | Minimal — overview only |
| Architecture | Includes pipeline documentation | No pipeline section |
| Commands | `# Commands Reference` section (H1) | Headingless block at H2 (`level=-1`) |
| After Development Guide | — | `# Development Rules` with per-topic rule files |

Forcing `AGENTS-codex.tpl.md` through `agents-base.md` would silently expand its scope — adding background, problem, solution, and pipeline sections that the Codex file intentionally omits — or require parameterised includes that the system does not support.

**Rule of thumb:** Use `agents-base.md` for templates that generate comprehensive documentation files. Keep `AGENTS-codex.tpl.md` independent because its purpose is targeted instructions, not exhaustive project context.
