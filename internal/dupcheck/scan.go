package dupcheck

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// markdownFiles walks root recursively and returns all .md files, skipping
// directories excluded by the given exclude patterns.
func markdownFiles(root string, excludes []string) ([]string, error) {
	var files []string

	normalizedRoot := normalizePath(root)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		normalized := normalizePath(path)

		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", "vendor":
				return filepath.SkipDir
			}
			if normalized != normalizedRoot && isExcluded(normalized, excludes) {
				return filepath.SkipDir
			}
			return nil
		}

		if isExcluded(normalized, excludes) {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(path), ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func isExcluded(path string, excludes []string) bool {
	for _, ex := range excludes {
		if matchExclude(path, ex) {
			return true
		}
	}
	return false
}

func matchExclude(path, pattern string) bool {
	normalized := normalizePath(pattern)

	if prefix, ok := strings.CutSuffix(normalized, "/**"); ok {
		return path == prefix || strings.HasPrefix(path, prefix+"/")
	}

	if ok, _ := filepath.Match(normalized, path); ok {
		return true
	}

	if ok, _ := filepath.Match(normalized, filepath.Base(path)); ok {
		return true
	}

	trimmed := strings.TrimSuffix(normalized, "/")
	return path == trimmed || strings.HasPrefix(path, trimmed+"/")
}

func normalizePath(p string) string {
	return filepath.ToSlash(filepath.Clean(p))
}
