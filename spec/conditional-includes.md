# Design: Conditional Includes

## Status

Draft

## Summary

Add conditional include directives to `docs-ssot`, enabling a single template to produce different output depending on build-time variables defined in `docsgen.yaml`.

## Motivation

Currently, generating tool-specific documents (e.g., AGENTS.md for Codex with rules included, vs. AGENTS.md for generic use without rules) requires separate template files with near-identical structure. Conditional includes eliminate this duplication by allowing a single template to express "include this section only when variable X equals Y."

### Concrete use case

Codex uses `AGENTS.md` as its primary instruction mechanism and has no separate rules directory. Other tools (Claude, Cursor, Copilot) have dedicated rule file outputs. To include shared rule content in the Codex-targeted `AGENTS.md` only, a conditional include is needed:

```markdown
<!-- @include: docs/07_rules/ if=tool:codex -->
```

Without this feature, the workaround is to maintain a separate `AGENTS-codex.tpl.md` template — which defeats the SSOT goal of this tool.

---

## Design

### 1. docsgen.yaml — Variable Declaration

Each target can optionally declare `vars`, a flat key-value map of strings:

```yaml
targets:
  - input: template/AGENTS.tpl.md
    output: AGENTS.md
    vars:
      tool: generic

  - input: template/AGENTS.tpl.md
    output: .codex/AGENTS.md
    vars:
      tool: codex
```

When `vars` is omitted or empty, the target behaves exactly as today (no variables, all conditional directives are skipped).

### 2. Include Directive — Conditional Syntax

A new optional `if=` parameter is added to the existing include directive:

```
<!-- @include: <path> [level=<delta>] [if=<condition>] -->
```

#### Condition expressions

| Expression | Meaning |
|---|---|
| `if=key:value` | Include only when `vars[key] == value` |
| `if=key` | Include only when `vars[key]` is defined (any non-empty value) |
| `if=!key` | Include only when `vars[key]` is NOT defined |
| `if=!key:value` | Include only when `vars[key] != value` |

Only one `if=` parameter is allowed per directive. For multiple conditions, use nested includes (a wrapper file that contains further conditionals).

The condition value is split on the **first** colon only. `if=key:value:with:colons` parses as `key="key"`, `value="value:with:colons"`.

#### Examples

```markdown
<!-- @include: docs/07_rules/ if=tool:codex -->
<!-- @include: docs/07_rules/ if=tool:codex level=+1 -->
<!-- @include: docs/05_ai/codex.md if=tool -->
<!-- @include: docs/05_ai/cursor.md if=!tool:codex -->
```

#### Parameter order

`if=` and `level=` can appear in any order after the path:

```markdown
<!-- @include: path level=+1 if=tool:codex -->
<!-- @include: path if=tool:codex level=+1 -->
```

Both are equivalent.

### 3. Conditional Block — Non-Include Content

For conditional inclusion of inline content (not from a file), a block syntax is provided:

```markdown
<!-- @if: tool:codex -->
This content only appears when tool=codex.
<!-- @endif -->
```

Block conditionals:
- Support the same condition expressions as `if=` on include directives
- Can nest (inner blocks are only evaluated if the outer block is active)
- Are not expanded inside fenced code blocks (treated as literal text)
- The `@if:` and `@endif` lines are consumed and do not appear in output
- **Must be balanced within a single file.** An `@if` opened in a parent file cannot be closed in an included file, and vice versa. Unbalanced blocks at end-of-file produce a build error.

### 4. Processing Pipeline Changes

#### 4.1. Config changes

```go
type Target struct {
    Input  string            `yaml:"input"`
    Output string            `yaml:"output"`
    Vars   map[string]string `yaml:"vars"`
}
```

#### 4.2. ProcessFile signature change

`ProcessFile` receives the variable map so the processor can evaluate conditions:

```go
func ProcessFile(path, outputPath string, vars map[string]string) (string, error)
```

For backward compatibility, `vars` may be `nil` (treated as empty map — all conditional includes are skipped, all `@if` blocks are skipped).

The internal recursive function `processFile` also receives `vars` so that conditions are evaluated consistently within included files:

```go
func processFile(absPath string, ancestors []string, absOutputPath string, vars map[string]string) (string, error)
```

#### 4.3. parseIncludeArgs extension

The current `parseIncludeArgs` function:

```go
func parseIncludeArgs(args string) (string, int)
```

Is extended to return a condition:

```go
type IncludeArgs struct {
    Path       string
    LevelDelta int
    Condition  *Condition // nil when no if= parameter
}

type Condition struct {
    Key    string
    Value  string // empty when condition is just "if=key" or "if=!key"
    Negate bool   // true for "if=!key" or "if=!key:value"
}

func parseIncludeArgs(args string) IncludeArgs
```

#### 4.4. Condition evaluation

```go
func (c *Condition) Evaluate(vars map[string]string) bool
```

| Condition | vars | Result |
|---|---|---|
| `if=tool:codex` | `{"tool": "codex"}` | `true` |
| `if=tool:codex` | `{"tool": "cursor"}` | `false` |
| `if=tool:codex` | `{}` | `false` |
| `if=tool` | `{"tool": "codex"}` | `true` |
| `if=tool` | `{}` | `false` |
| `if=!tool` | `{}` | `true` |
| `if=!tool:codex` | `{"tool": "cursor"}` | `true` |
| `if=!tool:codex` | `{}` | `true` |

#### 4.5. Block conditional processing

The processor's line-by-line loop gains a condition stack:

```go
type condFrame struct {
    active bool // whether this frame's condition evaluated to true
}
```

- On `<!-- @if: condition -->`: push a frame. If the outer frame is inactive, the inner frame is also inactive regardless of its own condition.
- On `<!-- @endif -->`: pop a frame. Error if stack is empty.
- Lines between `@if` and `@endif` are only processed (include resolution, link rewriting, output) when the top frame is active.
- At end of file, if the stack is non-empty, return an error (unclosed `@if`).

#### 4.6. processFile flow (updated)

```
for each line:
  1. detect fence open/close
  2. if inside fence → write as-is, continue
  3. if line is @if directive → push condition frame, continue
  4. if line is @endif directive → pop condition frame, continue
  5. if top condition frame is inactive → skip line, continue
  6. if line is @include directive:
     a. parse args (path, level, condition)
     b. if condition present and evaluates to false → skip, continue
     c. resolve and expand include
  7. write line to output
```

### 5. Regex Pattern Changes

Current:

```go
var includePattern = regexp.MustCompile(`^\s*<!--\s*@include:\s*(.*?)\s*-->\s*$`)
```

Additional patterns:

```go
var ifPattern    = regexp.MustCompile(`^\s*<!--\s*@if:\s*(.*?)\s*-->\s*$`)
var endifPattern = regexp.MustCompile(`^\s*<!--\s*@endif\s*-->\s*$`)
```

The `includePattern` does not change — the `if=` parameter is parsed inside `parseIncludeArgs` from the captured group.

### 6. Validate Command

The `validate` command performs a dry run. For conditional includes, validation must:
- Parse all `if=` conditions (syntax validation)
- Expand includes regardless of condition evaluation (to verify all paths are resolvable)
- Verify `@if`/`@endif` blocks are balanced

Validation runs in **"ignore conditions" mode**: all `if=` conditions on include directives are treated as `true`, and all `@if`/`@endif` blocks are expanded unconditionally. This avoids the contradiction where `if=tool:codex` and `if=!tool:codex` cannot both be true simultaneously. The implementation passes a boolean flag (e.g., `validateMode bool`) through the processing pipeline rather than manipulating `vars`.

### 7. docsgen.yaml — Full Example

```yaml
targets:
  # Documents
  - input: template/AGENTS.tpl.md
    output: AGENTS.md
    vars:
      tool: generic

  - input: template/AGENTS.tpl.md
    output: .codex/AGENTS.md
    vars:
      tool: codex

  # Claude rules (no vars needed — unconditional)
  - input: template/claude/rules/general.tpl.md
    output: .claude/rules/general.md

  # Conditional content inside a shared template
  - input: template/CLAUDE.tpl.md
    output: CLAUDE.md
    vars:
      include_glossary: "true"
```

### 8. Template Example

```markdown
# Project Context

<!-- @include: ./docs/01_project/overview.md -->

---

# Architecture

<!-- @include: ./docs/03_architecture/overview.md -->

---

<!-- @include: ./docs/07_rules/ if=tool:codex level=-1 -->

<!-- @if: tool:codex -->
# Development Rules

The following rules are included because this file targets Codex,
which uses AGENTS.md as its primary instruction source.
<!-- @endif -->

---

# Glossary

<!-- @include: ./docs/05_ai/glossary.md if=include_glossary level=-1 -->
```

---

## Backward Compatibility

- Targets without `vars` behave identically to today
- Include directives without `if=` behave identically to today
- No `@if`/`@endif` blocks means no change in behavior
- The `ProcessFile` signature changes; callers are `generator.Build`, `generator.Validate`, and `cli/include.go`. The `include` command passes `nil` for vars (no conditional behavior).

## Error Handling

| Error | Behavior |
|---|---|
| Unknown `if=` syntax (e.g., `if=:value`) | Build error with line number |
| Unclosed `@if` block | Build error: "unclosed @if at line N" |
| `@endif` without matching `@if` | Build error: "unexpected @endif at line N" |
| `if=` references undefined variable | Condition evaluates to `false` (not an error) |
| `@if` / `@endif` inside code fence | Treated as literal text (not processed) |

## Scope Exclusions

The following are explicitly out of scope for this feature:

- **Boolean operators** (`if=a:1 AND b:2`, `if=a:1 OR b:2`) — use nested includes or `@if` blocks instead
- **`@else` / `@elif` blocks** — use paired `if=key:value` and `if=!key:value` directives instead. May be added in a future iteration if the pattern proves too verbose.
- **Variable substitution in content** (`{{ tool }}` placeholders) — separate feature, tracked in roadmap v0.5
- **Computed variables** (deriving vars from env or other sources) — vars are static per target
- **Cross-file `@if`/`@endif` blocks** — each file's condition stack is independent. An `@if` opened in a parent cannot be closed by an `@endif` in an included file.

## Test Plan

### Unit tests (processor)

| Test case | Description |
|---|---|
| `if=key:value` match | Include is expanded when vars match |
| `if=key:value` no match | Include is skipped when vars don't match |
| `if=key` defined | Include is expanded when key exists |
| `if=key` undefined | Include is skipped when key is absent |
| `if=!key` negation | Include is expanded when key is absent |
| `if=!key:value` negation | Include is expanded when key has different value |
| `if=` with `level=` combined | Both parameters apply correctly |
| `if=` parameter order | `if=` before and after `level=` both work |
| `if=` inside code fence | Directive treated as literal text |
| `@if`/`@endif` block active | Content between markers is included |
| `@if`/`@endif` block inactive | Content between markers is skipped |
| `@if` nested | Inner block respects outer block state |
| `@if` unclosed | Returns error |
| `@endif` unmatched | Returns error |
| `@if` inside code fence | Treated as literal text |
| `nil` vars | All conditionals skipped, no panic |
| Empty vars map | All conditionals skipped |

### Integration tests (generator)

| Test case | Description |
|---|---|
| Same template, different vars | Two targets produce different output |
| Validate with conditional includes | All paths checked regardless of condition |

### Regression tests

| Test case | Description |
|---|---|
| Existing templates without vars | Unchanged behavior |
| Existing include directives | No `if=` parameter, unchanged behavior |

## Implementation Order

1. **Config**: Add `Vars` field to `Target` struct
2. **Processor — parseIncludeArgs**: Parse `if=` parameter, return `IncludeArgs` struct
3. **Processor — Condition**: Implement `Evaluate(vars)`
4. **Processor — processFile**: Thread `vars` through, evaluate `if=` on includes
5. **Processor — block conditionals**: Implement `@if`/`@endif` with condition stack
6. **Generator**: Pass `target.Vars` to `ProcessFile`
7. **Validate**: Implement "expand all" mode for validation
8. **Tests**: Unit + integration + regression
9. **Documentation**: Update `template/docs/03_architecture/includes.md` and `template/docs/06_reference/commands.md`
