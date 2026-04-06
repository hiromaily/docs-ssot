## Best Practices for Multi-Tool Repositories

When a repository is used with multiple AI tools, follow these principles to minimize duplication and keep instructions consistent.

### Principle 1 — Single source of truth in `AGENTS.md`

Write shared project knowledge (architecture, coding rules, build commands, testing strategy) in `AGENTS.md` only. All four major tools read it.

### Principle 2 — Tool-specific files contain only deltas

`CLAUDE.md`, `.cursor/rules/`, `.github/copilot-instructions.md`, and `.codex/config.toml` should reference `AGENTS.md` and add only tool-specific behavior:

- Claude: subagent usage, skill invocation patterns
- Cursor: rule application granularity (`globs`)
- Codex: approval policy, sandbox mode, MCP connections
- Copilot: agent exclusions, prompt file conventions

### Principle 3 — Long procedures go into skills

Multi-step workflows (adding endpoints, running migrations, release checklists) belong in skills, not in instruction files. Skills are supported by all four tools.

### Principle 4 — Rules are constraints, not procedures

Rules should express what to do and what not to do. Procedures (how to do it step by step) belong in skills.

---

### Recommended Minimal Directory Structure

```
repo/
├── AGENTS.md                           # Cross-tool shared instructions
├── CLAUDE.md                           # Claude-specific delta
├── .claude/
│   ├── settings.json                   # Permissions, hooks, MCP
│   └── rules/
│       ├── architecture.md             # Topic-scoped rules
│       └── testing.md
├── .cursor/
│   └── rules/
│       ├── 00-core.mdc                 # Always-apply core rules
│       └── 10-go-architecture.mdc      # Path-gated Go rules
├── .codex/
│   └── config.toml                     # Runtime settings
├── .agents/
│   └── skills/
│       └── add-endpoint/
│          └── SKILL.md                 # Shared skill (Codex reads .agents/)
├── .github/
│   └── copilot-instructions.md         # Copilot-specific delta
└── README.md
```

### Recommended Full Directory Structure

```
repo/
├── AGENTS.md                           # Cross-tool shared instructions
├── CLAUDE.md                           # Claude delta
├── .claude/
│   ├── settings.json
│   ├── rules/
│   │   ├── architecture.md
│   │   ├── testing.md
│   │   └── db.md
│   ├── skills/
│   │   ├── add-endpoint/SKILL.md
│   │   ├── db-migration/SKILL.md
│   │   └── run-tests/SKILL.md
│   └── agents/
│       ├── reviewer.md
│       └── refactorer.md
├── .cursor/
│   ├── rules/
│   │   ├── 00-core.mdc
│   │   ├── 10-go-architecture.mdc
│   │   ├── 20-testing.mdc
│   │   └── 30-db.mdc
│   ├── skills/
│   │   ├── add-endpoint/SKILL.md
│   │   └── db-migration/SKILL.md
│   └── agents/
│       ├── reviewer.md
│       └── refactorer.md
├── .codex/
│   ├── config.toml
│   └── agents/
│       ├── reviewer.toml
│       └── refactorer.toml
├── .agents/
│   └── skills/
│       ├── add-endpoint/SKILL.md
│       ├── db-migration/SKILL.md
│       └── run-tests/SKILL.md
├── .github/
│   ├── copilot-instructions.md
│   ├── instructions/
│   │   ├── go.instructions.md
│   │   ├── testing.instructions.md
│   │   └── db.instructions.md
│   └── agents/
│       └── reviewer.agent.md
└── README.md
```

### What Goes Where — Decision Guide

| Question | Answer |
|----------|--------|
| Is it shared project knowledge? | `AGENTS.md` |
| Is it a runtime/execution setting? | Tool-specific config (`settings.json`, `config.toml`) |
| Is it a persistent constraint? | Rules (`.claude/rules/`, `.cursor/rules/`, `.github/instructions/`) |
| Is it a reusable multi-step procedure? | Skills (`.claude/skills/`, `.agents/skills/`, etc.) |
| Is it tool-specific behavior only? | Tool-specific instruction file (`CLAUDE.md`, `.github/copilot-instructions.md`) |
