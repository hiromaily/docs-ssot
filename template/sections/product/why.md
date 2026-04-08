## Why docs-ssot?

Software projects now maintain documentation for three different audiences:

📖 **For humans** — README, contributing guides, VitePress/Docusaurus docs sites
🤖 **For AI agents** — CLAUDE.md, AGENTS.md, .cursor/rules/, .github/instructions/
📋 **For both** — Architecture docs, coding rules, setup guides

And the list keeps growing:

| Audience | Files |
|----------|-------|
| Humans | README.md, VitePress docs, CONTRIBUTING.md |
| Claude Code | CLAUDE.md, .claude/rules/*.md, .claude/skills/ |
| Codex | AGENTS.md, .agents/skills/ |
| Cursor | .cursor/rules/*.mdc, .cursor/skills/ |
| Copilot | .github/copilot-instructions.md, .github/instructions/*.md |

Most of these files share the same underlying information —
architecture, coding rules, setup steps, testing strategy.

**The problem: Markdown has no `#include`.**

Every tool demands Markdown. Not YAML, not HTML — Markdown.
But Markdown has no way to share content across files.
So teams copy-paste. Then the copies drift. Information contradicts.

When humans read conflicting docs, they ask questions.
**When AI agents read conflicting docs, they silently act on the wrong one.**

An agent trusting stale architecture notes will refactor
the wrong module. An agent reading outdated rules will
bypass your testing strategy. And it won't ask — it will
just do it, confidently.

**Inconsistent documentation is the silent killer of AI-assisted development.**

docs-ssot solves this: write once, generate everywhere, always consistent.
