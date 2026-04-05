package dupcheck

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// markdownFiles walks root recursively and returns all .md files, skipping
// directories excluded by the given exclude patterns.
func markdownFiles(root string, excludes []string) ([]string, error) {
	normalizedExcludes, err := normalizeExcludes(excludes)
	if err != nil {
		return nil, err
	}

	var files []string
	normalizedRoot := normalizePath(root)

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		normalized := normalizePath(path)

		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", "vendor":
				return filepath.SkipDir
			}
			excluded, exErr := isExcluded(normalized, normalizedExcludes)
			if exErr != nil {
				return exErr
			}
			if normalized != normalizedRoot && excluded {
				return filepath.SkipDir
			}
			return nil
		}

		excluded, exErr := isExcluded(normalized, normalizedExcludes)
		if exErr != nil {
			return exErr
		}
		if excluded {
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

// normalizeExcludes normalizes and validates all exclude patterns up front,
// returning an error for any invalid glob syntax.
func normalizeExcludes(excludes []string) ([]string, error) {
	out := make([]string, len(excludes))
	for i, ex := range excludes {
		n := normalizePath(ex)
		// Validate glob syntax for patterns that will be passed to filepath.Match.
		if !strings.HasSuffix(n, "/**") {
			if _, err := filepath.Match(n, ""); err != nil {
				return nil, fmt.Errorf("invalid exclude pattern %q: %w", ex, err)
			}
		}
		out[i] = n
	}
	return out, nil
}

func isExcluded(path string, excludes []string) (bool, error) {
	for _, ex := range excludes {
		ok, err := matchExclude(path, ex)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// matchExclude reports whether path matches the pre-normalized pattern.
func matchExclude(path, pattern string) (bool, error) {
	if prefix, ok := strings.CutSuffix(pattern, "/**"); ok {
		return path == prefix || strings.HasPrefix(path, prefix+"/"), nil
	}

	if ok, err := filepath.Match(pattern, path); err != nil {
		return false, fmt.Errorf("exclude pattern %q: %w", pattern, err)
	} else if ok {
		return true, nil
	}

	if ok, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
		return false, fmt.Errorf("exclude pattern %q: %w", pattern, err)
	} else if ok {
		return true, nil
	}

	trimmed := strings.TrimSuffix(pattern, "/")
	return path == trimmed || strings.HasPrefix(path, trimmed+"/"), nil
}

func normalizePath(p string) string {
	return filepath.ToSlash(filepath.Clean(p))
}
