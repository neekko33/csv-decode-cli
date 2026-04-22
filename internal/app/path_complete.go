package app

import (
	"os"
	"path/filepath"
	"strings"
)

func completePath(input string) (string, bool) {
	value := strings.TrimSpace(input)
	if value == "" {
		return input, false
	}

	expandedValue, expandOK := expandHome(value)
	if !expandOK {
		return input, false
	}

	dirPart := ""
	namePart := expandedValue
	searchDir := "."

	if idx := strings.LastIndex(expandedValue, string(filepath.Separator)); idx >= 0 {
		dirPart = expandedValue[:idx+1]
		namePart = expandedValue[idx+1:]
		searchDir = strings.TrimSuffix(dirPart, string(filepath.Separator))
		if searchDir == "" {
			searchDir = string(filepath.Separator)
		}
	}

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return input, false
	}

	matches := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), namePart) {
			matches = append(matches, entry)
		}
	}
	if len(matches) == 0 {
		return input, false
	}

	prefix := matches[0].Name()
	for i := 1; i < len(matches); i++ {
		prefix = commonPrefix(prefix, matches[i].Name())
	}
	if prefix == "" || prefix == namePart {
		if len(matches) == 1 {
			completed := dirPart + matches[0].Name()
			if matches[0].IsDir() {
				completed += string(filepath.Separator)
			}
			return completed, true
		}
		return input, false
	}

	return dirPart + prefix, true
}

func expandHome(path string) (string, bool) {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", false
		}
		return home, true
	}

	prefix := "~" + string(filepath.Separator)
	if strings.HasPrefix(path, prefix) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", false
		}
		return filepath.Join(home, strings.TrimPrefix(path, prefix)), true
	}

	return path, true
}

func commonPrefix(a, b string) string {
	max := len(a)
	if len(b) < max {
		max = len(b)
	}

	i := 0
	for i < max && a[i] == b[i] {
		i++
	}
	return a[:i]
}
