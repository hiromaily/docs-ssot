## Agent-aware migration (--from)

With `--from`, `migrate` scans AI tool configuration files (rules, skills, commands, subagents) from the specified tool and generates SSOT sections with per-tool templates for the target tools.

```
docs-ssot migrate --from <tool> [--to <tools>] [flags]
```

### What it does

1. **Scans** the source tool's configuration directory for rules, skills, commands, and subagents
2. **Strips** frontmatter from source files and shifts H1→H2 headings
3. **Creates section files** under `template/sections/ai/<type>/<slug>.md`
4. **Generates templates** for each target tool with appropriate frontmatter
5. **Updates `docsgen.yaml`** with new build targets
6. **Verifies round-trip** by building and comparing against originals

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--from` | (required) | Source AI tool to migrate from (`claude`, `cursor`, `copilot`) |
| `--to` | all except `--from` | Target tools, comma-separated (`cursor,copilot,codex`) |
| `--convert-commands` | `false` | Convert legacy commands to skills during migration |
| `--infer-globs` | `false` | Infer path-gated globs from rule slug names |
| `--dry-run` | `false` | Print the migration plan without writing files |

### Examples

Migrate Claude configs to all other tools:

```
docs-ssot migrate --from claude
```

Migrate to specific tools only:

```
docs-ssot migrate --from claude --to cursor,codex
```

Preview migration plan:

```
docs-ssot migrate --from claude --dry-run
```

Migrate with path inference and command conversion:

```
docs-ssot migrate --from claude --to cursor --infer-globs --convert-commands
```

Combine agent and file migration:

```
docs-ssot migrate --from claude --to cursor README.md CLAUDE.md
```

### Output

```
Detected source tool: claude (5 files)
Target tools: cursor, copilot, codex

Creating sections:
  template/sections/ai/rules/architecture.md
  template/sections/ai/rules/testing.md
  template/sections/ai/skills/deploy.md
  template/sections/ai/subagents/critic.md
  template/sections/ai/subagents/debugger.md

Creating templates (3 tools × 5 files):
  cursor: 5 templates
  copilot: 5 templates
  codex: 4 templates

Updated docsgen.yaml (14 new targets)
Verifying round-trip...
Round-trip verification: OK
Agent migration complete.
```
