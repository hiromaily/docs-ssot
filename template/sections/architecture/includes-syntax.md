## Include Directive

Compatible with [VitePress](https://vitepress.dev/) syntax:

```markdown
<!-- @include: path/to/file.md -->           Single file
<!-- @include: path/to/dir/ -->              All .md files in directory
<!-- @include: path/**/*.md -->              Recursive glob
<!-- @include: path/to/file.md level=+1 -->  Shift heading depth
```

Includes are resolved recursively. Circular includes are detected and cause a build error.
