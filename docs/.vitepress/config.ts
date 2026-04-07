import { defineConfig } from "vitepress";

export default defineConfig({
  title: "docs-ssot",
  description: "Documentation SSOT Generator",
  themeConfig: {
    nav: [
      { text: "Guide", link: "/guide/getting-started" },
      { text: "Architecture", link: "/architecture/overview" },
      { text: "AI Agents", link: "/ai/overview" },
      { text: "Reference", link: "/reference/commands" },
    ],
    sidebar: {
      "/guide/": [
        {
          text: "Guide",
          items: [
            { text: "Getting Started", link: "/guide/getting-started" },
            { text: "Setup", link: "/guide/setup" },
            { text: "Testing", link: "/guide/testing" },
            { text: "Linting", link: "/guide/linting" },
          ],
        },
      ],
      "/architecture/": [
        {
          text: "Architecture",
          items: [
            { text: "Overview", link: "/architecture/overview" },
            { text: "System Design", link: "/architecture/system" },
            { text: "Build Pipeline", link: "/architecture/pipeline" },
            { text: "Include Directives", link: "/architecture/includes" },
            { text: "Feature Details", link: "/architecture/features" },
          ],
        },
      ],
      "/ai/": [
        {
          text: "AI Agent Configuration",
          items: [
            { text: "Overview", link: "/ai/overview" },
            { text: "Claude Code", link: "/ai/claude" },
            { text: "Hooks", link: "/ai/hooks" },
            { text: "Cursor", link: "/ai/cursor" },
            { text: "GitHub Copilot", link: "/ai/github-copilot" },
            { text: "Codex", link: "/ai/codex" },
            { text: "Cross-Tool Mapping", link: "/ai/cross-tool-mapping" },
            { text: "Best Practices", link: "/ai/best-practices" },
            { text: "Glossary", link: "/ai/glossary" },
          ],
        },
      ],
      "/reference/": [
        {
          text: "Reference",
          items: [
            { text: "Commands", link: "/reference/commands" },
            { text: "Directory Structure", link: "/reference/directory" },
            { text: "Roadmap", link: "/reference/roadmap" },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: "github", link: "https://github.com/hiromaily/docs-ssot" },
    ],
  },
});
