---
name: check-ssot
description: Verify that generated documentation files (README.md, CLAUDE.md, AGENTS.md, etc.) are in sync with their source templates. Run after editing template/sections/ files.
argument-hint:
---

# Check SSOT Compliance

Verify that all generated documentation files are in sync with their source templates.

This project uses SSOT (Single Source of Truth) documentation generation. Output files like `README.md`, `CLAUDE.md`, and `AGENTS.md` are **build artifacts** generated from source files under `template/` and `template/sections/`.

## When to Run

- Before committing changes that touch files under `template/` or `template/sections/`
- When reviewing a PR that modifies documentation source files
- When you suspect generated files are out of date

## Procedure

### Step 1: Regenerate All Outputs

```bash
make docs
```

This runs `docs-ssot build` which processes all targets defined in `docsgen.yaml`.

### Step 2: Check for Drift

```bash
git diff --stat
```

If there is no diff, all generated files are in sync. Report success and stop.

If there is a diff, continue to Step 3.

### Step 3: Report Stale Files

For each file with changes, show:

1. **File path** of the stale generated file
2. **Diff summary** — lines added/removed
3. **Root cause** — which source file under `template/sections/` likely caused the drift (check `git log --oneline -3 <generated-file>` vs `git log --oneline -3 template/`)

### Step 4: Suggest Fix

Provide a remediation command:

```bash
make docs && git add <stale-files> && git status
```

If the generated files were manually edited (violating SSOT), warn the user that manual edits will be overwritten and they should move the changes to the source files instead.

## Output Format

### All in sync

```
SSOT Check: PASS
All N generated files match their source templates.
```

### Drift detected

```
SSOT Check: DRIFT DETECTED
N of M generated files are out of sync.

Stale files:
- README.md (3 lines changed)
- CLAUDE.md (12 lines changed)

Run `make docs` to regenerate, then commit the updated files.
```
