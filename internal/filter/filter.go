package filter

import (
	"path/filepath"
	"strings"
)

// ignoredDirs are directories that should not be watched
var ignoredDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"bin":          true,
	"build":        true,
	"dist":         true,
	"tmp":          true,
	".idea":        true,
	".vscode":      true,
	"vendor":       true,
}

// ignoredExtensions are file extensions that should be ignored
var ignoredExtensions = map[string]bool{
	".swp": true,
	".tmp": true,
	".log": true,
	".swx": true,
	".swo": true,
}

// ShouldIgnore returns true if the path should be ignored
func ShouldIgnore(path string) bool {
	// Check if any parent directory is in the ignored list
	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		if ignoredDirs[part] {
			return true
		}
	}

	// Check file extension
	ext := filepath.Ext(path)
	if ignoredExtensions[ext] {
		return true
	}

	// Ignore hidden files (starting with .)
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") && base != "." {
		return true
	}

	return false
}

// IsGoFile returns true if the file is a Go source file
func IsGoFile(path string) bool {
	return filepath.Ext(path) == ".go"
}
