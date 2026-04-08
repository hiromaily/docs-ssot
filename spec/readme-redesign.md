# docs-ssot README Redesign Spec

## Goal

docs-ssot の README を全面リデザインし、海外のソフトウェアエンジニア（AI agent ユーザー、テックリード、OSS 発見者）に向けて魅力的に訴求する。

## Target Audience

- AI agent（Claude, Codex, Cursor, Copilot）を日常的に使う開発者
- チーム内でドキュメント管理に課題を感じているテックリード
- GitHub で OSS ツールを探している一般的な開発者
- 特に海外のソフトウェアエンジニア

## Tone & Style

- フレンドリー＆実践的（Vite, pnpm のような親しみやすさ）
- 具体例豊富、コピペで試せる

## Core Message

> Single Source of Truth for the AI agent era.

最大の強み: README, CLAUDE.md, AGENTS.md, VitePress docs, Claude/Cursor/Codex/Copilot の agent 向けリソースまで、すべてを 1 つのソースから SSOT で管理・生成できる。

## Background & Problem Statement

- AI Agent による開発が当たり前になり、ドキュメント整備がより重要になった
- 人間用ドキュメント（README, VitePress docs）、AI 用（CLAUDE.md, .cursor/rules/）、共通（アーキテクチャ, コーディングルール）と 3 つのオーディエンスが存在
- Markdown が求められているが、Markdown には include の仕組みがない
- 情報が陳腐化し、異なるファイル間で矛盾が生じる
- 人間は矛盾に気づいて質問するが、AI Agent は疑わず間違った情報で自信を持って行動する — これが最も致命的

---

## Section Structure

### 1. Hero

```markdown
# docs-ssot

**Single Source of Truth for the AI agent era.**

Generate README, CLAUDE.md, AGENTS.md, Cursor rules, Copilot instructions,
and VitePress docs — all from one source.

<!-- shields.io badges for Claude Code, Cursor, Codex, GitHub Copilot -->
```

Hero diagram:

```text
                    ┌─── README.md
                    ├─── CLAUDE.md
                    ├─── AGENTS.md
  template/sections ├─── .claude/rules/*.md
  (single source) ──├─── .cursor/rules/*.mdc
                    ├─── .github/instructions/*.md
                    ├─── .agents/skills/
                    └─── VitePress docs site
```

### 2. Why docs-ssot?

Three audiences, one problem:

- 📖 **For humans** — README, contributing guides, VitePress/Docusaurus docs sites
- 🤖 **For AI agents** — CLAUDE.md, AGENTS.md, .cursor/rules/, .github/instructions/
- 📋 **For both** — Architecture docs, coding rules, setup guides

Table showing the file explosion per AI tool.

Key messages:

- Markdown has no `#include` — but every tool demands Markdown
- Teams copy-paste → copies drift → information contradicts
- When humans read conflicting docs, they ask questions
- When AI agents read conflicting docs, they silently act on the wrong one
- Inconsistent documentation is the silent killer of AI-assisted development

### 3. Before / After

Before: 5 files with different testing instructions

```text
README.md              ← "Run tests with: make test"
CLAUDE.md              ← "Run tests with: go test ./..."
AGENTS.md              ← "Run tests with: make test-local"
.cursor/rules/go.mdc   ← "Always run go test before committing"
.github/instructions/   ← "Use make verify for pre-commit checks"
```

5 files. 5 different testing instructions. An AI agent picks one — and skips your lint, coverage, and integration test pipeline.

After: 1 source → all targets

```text
template/sections/development/testing.md    ← single source
           │
           ├──→ README.md
           ├──→ CLAUDE.md
           ├──→ AGENTS.md
           ├──→ .cursor/rules/go.mdc
           └──→ VitePress docs site

$ docs-ssot build
Generated 12 files from 1 source.
```

1 file. 1 version. Always consistent.

### 4. Quick Start

1. Install: `go install github.com/hiromaily/docs-ssot/cmd/docs-ssot@latest`
2. Migrate existing docs: `docs-ssot migrate README.md CLAUDE.md AGENTS.md`
3. Show generated `docsgen.yaml` (input/output mapping)
4. Edit single source → `docs-ssot build`
5. Bonus: `docs-ssot migrate --from claude` with expanded `docsgen.yaml` showing cross-tool generation

### 5. Supported Targets

Two tables:

📖 Human Documentation: README, VitePress, CONTRIBUTING, any Markdown file

🤖 AI Agent Instructions: Claude Code, Codex, Cursor, GitHub Copilot with specific output file paths

shields.io badges for each AI tool.

Closing quote: All generated from the same `template/sections/` directory.

### 6. How It Works

Opening: docs-ssot adds one missing feature to Markdown — `#include`.

- Pipeline diagram: sections + templates → `docs-ssot build` → all outputs
- Include directive syntax (4 formats: file, directory, glob, recursive glob + level parameter)
- VitePress-compatible syntax
- Template example showing how CLAUDE.tpl.md composes sections

### 7. Commands

Table of all CLI commands: build, migrate, migrate --from, check, validate, include, version

### 8. SSOT Duplicate Detection

`docs-ssot check` feature highlight with example output (TF-IDF cosine similarity).

### 9. Why not Hugo / MkDocs?

Comparison table:

- SSGs build HTML websites; docs-ssot builds any Markdown file
- docs-ssot sits upstream of SSGs
- Pipeline diagram: sections → docs-ssot → .md files → VitePress → website

### 10. Self-Hosting Example

This repository itself uses docs-ssot. README, CLAUDE.md, AGENTS.md, all agent configs are generated from `template/sections/`.

### 11. Contributing

Minimal setup instructions.

### 12. License

MIT.

---

## Section Restructuring

README リデザインに伴い、既存の `template/sections/` を細粒度に再分割する。
現状はファイルが大きすぎて「README にはサマリーだけ、CLAUDE.md には全詳細」という使い分けができない。

### 分割対象と理由

| 既存ファイル                 | 行数 | 問題                                                          |
|-----------------------------|------|---------------------------------------------------------------|
| `project/overview.md`       | 95   | Overview + Background + Problem + Solution + Concept が 1 ファイル |
| `reference/commands.md`     | 396  | 全コマンド詳細が 1 ファイル                                      |
| `architecture/system.md`    | 157  | Core Components 4 つが全部入り                                   |
| `architecture/includes.md`  | 188  | Include 仕様の全詳細                                            |

### 分割後の構造

```text
template/sections/
├── project/
│   ├── overview.md            ← 簡潔な 1 段落の概要（既存から分離）
│   ├── background.md          ← AI 時代の背景（既存から分離）
│   ├── problem.md             ← Markdown の限界（既存から分離）
│   ├── solution.md            ← docs-ssot のアプローチ（既存から分離）
│   ├── concept.md             ← (既存のまま)
│   ├── vision.md              ← (既存のまま)
│   └── roadmap.md             ← (既存のまま)
├── product/
│   ├── hero.md                ← ★ 新規: キャッチコピー + バッジ + 放射図
│   ├── why.md                 ← ★ 新規: 3 audiences + pain points
│   ├── before-after.md        ← ★ 新規: testing instruction の例
│   ├── supported-targets.md   ← ★ 新規: テーブル + バッジ
│   ├── comparison.md          ← ★ 新規: Why not Hugo/MkDocs
│   ├── self-hosting.md        ← ★ 新規: この repo 自体が例
│   ├── features.md            ← (既存のまま)
│   ├── concept.md             ← (既存のまま)
│   └── faq.md                 ← (既存のまま)
├── reference/
│   ├── commands-summary.md    ← ★ 新規: コマンド一覧テーブルのみ
│   ├── commands/
│   │   ├── build.md           ← commands.md から分離
│   │   ├── check.md
│   │   ├── migrate.md
│   │   ├── migrate-from.md
│   │   ├── validate.md
│   │   └── include.md
│   └── directory.md           ← (既存のまま)
├── architecture/
│   ├── overview.md            ← (既存のまま)
│   ├── system.md              ← CLI Core Components のみに縮小
│   ├── pipeline.md            ← (既存のまま)
│   ├── includes-syntax.md     ← ★ 新規: 構文例のみ（README 向け）
│   ├── includes.md            ← (既存のまま、フル仕様、CLAUDE.md 向け)
│   ├── features.md            ← (既存のまま)
│   └── diagrams/              ← (既存のまま)
├── development/               ← (変更なし)
└── ai/                        ← (変更なし)
```

### テンプレートへの影響

- `README.tpl.md`: 全面書き換え。新規 sections（hero, why, before-after 等）を使用
- `CLAUDE.tpl.md`: `project/overview.md` → `project/` ディレクトリ include に変更。`reference/commands.md` → `reference/commands/` ディレクトリ include に変更
- `AGENTS.tpl.md`: CLAUDE.tpl.md と同様の変更

### 既存テンプレートの互換性

- 分割元ファイル（`project/overview.md` 等）は分割後も残さない（削除）
- テンプレート側を同時に更新して整合性を保つ
- `docs-ssot build` + `git diff` で出力が意図通りか検証

---

## Design Decisions

- **shields.io badges** over logo images — no external hosting needed, easy to add new tools
- **`migrate` command in Quick Start** — shows instant value for existing projects, not just greenfield
- **`docsgen.yaml` shown twice** in Quick Start — before and after agent migration, to demonstrate incremental adoption
- **"Why not X?" positioned late** — by this point the reader already understands the value; comparison reinforces rather than leads
- **Self-hosting example** — the repo itself is proof the tool works at scale
