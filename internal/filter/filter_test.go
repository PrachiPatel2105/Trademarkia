package filter

import (
	"testing"
)

func TestShouldIgnore_IgnoredDirectories(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{".git/config", true},
		{"node_modules/package/index.js", true},
		{"bin/server", true},
		{"build/output.txt", true},
		{"dist/bundle.js", true},
		{"tmp/cache.dat", true},
		{".idea/workspace.xml", true},
		{".vscode/settings.json", true},
		{"vendor/package/file.go", true},
		{"src/main.go", false},
		{"internal/app/app.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ShouldIgnore(tt.path)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestShouldIgnore_IgnoredExtensions(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"file.swp", true},
		{"file.tmp", true},
		{"file.log", true},
		{"file.swx", true},
		{"file.swo", true},
		{"file.go", false},
		{"file.txt", false},
		{"file.json", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ShouldIgnore(tt.path)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestShouldIgnore_HiddenFiles(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{".hidden", true},
		{".DS_Store", true},
		{"dir/.gitignore", true},
		{".", false}, // Current directory is not ignored
		{"normal.go", false},
		{"src/.hidden", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ShouldIgnore(tt.path)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestShouldIgnore_NestedPaths(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"project/node_modules/package/index.js", true},
		{"project/.git/objects/abc", true},
		{"project/src/bin/data", true}, // 'bin' anywhere in path
		{"project/src/main.go", false},
		{"deep/nested/path/file.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ShouldIgnore(tt.path)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"main.go", true},
		{"package/file.go", true},
		{"main.c", false},
		{"file.txt", false},
		{"go.mod", false},
		{".go", true}, // Edge case: file named ".go"
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := IsGoFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsGoFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestShouldIgnore_RealWorldPaths(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
		reason   string
	}{
		{"cmd/server/main.go", false, "normal Go file"},
		{"internal/app/app.go", false, "internal package"},
		{".git/HEAD", true, "git directory"},
		{"node_modules/express/index.js", true, "node_modules"},
		{"bin/server.exe", true, "binary output"},
		{"build/output/server", true, "build directory"},
		{"tmp/cache/data.tmp", true, "tmp directory and extension"},
		{".vscode/launch.json", true, "editor config"},
		{"vendor/github.com/pkg/errors/errors.go", true, "vendor directory"},
		{"main.go~", false, "backup file (not in ignored extensions)"},
		{"src/.env", true, "hidden file"},
		{"dist/bundle.min.js", true, "dist directory"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ShouldIgnore(tt.path)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q) = %v, want %v (%s)", 
					tt.path, result, tt.expected, tt.reason)
			}
		})
	}
}

func TestShouldIgnore_EdgeCases(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"", false},
		{".", false},
		{"/", false},
		{"bin", true}, // Just the directory name
		{".git", true},
		{"normal", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ShouldIgnore(tt.path)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
