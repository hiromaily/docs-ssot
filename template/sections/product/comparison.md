## Why not Hugo / MkDocs / Docusaurus?

Static site generators build **websites**.
docs-ssot builds **any Markdown file from shared sources**.

|  | Hugo / MkDocs | docs-ssot |
|---|---|---|
| Output | HTML website | Any Markdown file |
| CLAUDE.md generation | ❌ | ✅ |
| .cursor/rules/ generation | ❌ | ✅ |
| AI agent config migration | ❌ | ✅ |
| Works alongside SSGs | — | ✅ (generates source .md for VitePress) |
| Duplicate detection | ❌ | ✅ |
| Markdown include syntax | Varies | VitePress-compatible |

docs-ssot is not a replacement for static site generators.
It sits **upstream** — generating the Markdown that SSGs then render.

```text
template/sections/ → docs-ssot build → docs/*.md → VitePress build → website
                                      → README.md
                                      → CLAUDE.md
                                      → .cursor/rules/
```
