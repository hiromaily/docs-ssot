## FAQ

### What problem does docs-ssot solve?

When using multiple AI coding tools (Claude Code, Codex, Cursor, Copilot), each requires its own instruction file (`CLAUDE.md`, `AGENTS.md`, etc.). Maintaining these files manually leads to duplication, inconsistency, and drift. `docs-ssot` lets you write each piece of documentation once and generate all output files from that single source.

---

### How is this different from a static site generator?

Static site generators (Hugo, Jekyll, VitePress) convert Markdown to HTML for web viewing. `docs-ssot` generates Markdown-to-Markdown — it composes small Markdown modules into larger Markdown files like `README.md` or `CLAUDE.md`. They serve different purposes and can be used together (this project uses both `docs-ssot` and VitePress).

---

### Can I use docs-ssot with VitePress?

Yes. The `@include` directive syntax is compatible with [VitePress markdown includes](https://vitepress.dev/guide/markdown#markdown-file-inclusion). You can share the same source files between `docs-ssot` and VitePress.

---

### What happens if an included file is missing?

The build fails with an error. Includes never fail silently — this is by design to prevent broken documentation from being generated.

---

### Does docs-ssot support variables or conditional includes?

Not yet. Variable substitution (`{{ variable }}`) and conditional includes are planned for a future release. See the [Roadmap](/reference/roadmap) for details.

---

### Can included files include other files?

Yes. Includes are resolved recursively. If `A.md` includes `B.md` and `B.md` includes `C.md`, the final output contains all three. Circular includes are detected and produce a build error.

---

### How do I add a new output target?

1. Create a template file in `template/pages/` (e.g., `MYFILE.tpl.md`)
2. Add a new entry in `docsgen.yaml`:
   ```yaml
   - input: template/pages/MYFILE.tpl.md
     output: MYFILE.md
   ```
3. Run `make docs`

---

### Why Markdown-to-Markdown instead of Markdown-to-HTML?

AI tools consume Markdown directly — they don't read HTML. The primary consumers of `CLAUDE.md`, `AGENTS.md`, and `README.md` are AI agents and GitHub's Markdown renderer, not web browsers. Generating HTML would add complexity without serving the core use case.
