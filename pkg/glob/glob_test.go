package glob

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpand(t *testing.T) {
	// Create test directory structure
	tmpDir := t.TempDir()

	// Create test files
	files := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
		filepath.Join(tmpDir, "document.pdf"),
		filepath.Join(tmpDir, "image.jpg"),
	}

	for _, f := range files {
		if err := os.WriteFile(f, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	tests := []struct {
		name     string
		patterns []string
		wantLen  int
	}{
		{
			name:     "wildcard all txt files",
			patterns: []string{filepath.Join(tmpDir, "*.txt")},
			wantLen:  2,
		},
		{
			name:     "wildcard all files",
			patterns: []string{filepath.Join(tmpDir, "*")},
			wantLen:  4,
		},
		{
			name:     "literal path",
			patterns: []string{filepath.Join(tmpDir, "file1.txt")},
			wantLen:  1,
		},
		{
			name:     "multiple patterns",
			patterns: []string{filepath.Join(tmpDir, "*.txt"), filepath.Join(tmpDir, "*.pdf")},
			wantLen:  3,
		},
		{
			name:     "question mark wildcard",
			patterns: []string{filepath.Join(tmpDir, "file?.txt")},
			wantLen:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, err := Expand(tt.patterns)
			if err != nil {
				t.Errorf("Expand() error = %v", err)
				return
			}

			if len(matches) != tt.wantLen {
				t.Errorf("Expand() got %d matches, want %d", len(matches), tt.wantLen)
			}

			// Verify all matches are valid files
			for _, match := range matches {
				if _, err := os.Stat(match.Path); err != nil {
					t.Errorf("Match path %s doesn't exist: %v", match.Path, err)
				}
			}
		})
	}
}

func TestExpandSingle(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "single.txt")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("literal path", func(t *testing.T) {
		path, err := ExpandSingle(testFile)
		if err != nil {
			t.Errorf("ExpandSingle() error = %v", err)
			return
		}

		if path != testFile {
			t.Errorf("ExpandSingle() got %s, want %s", path, testFile)
		}
	})

	t.Run("wildcard pattern", func(t *testing.T) {
		pattern := filepath.Join(tmpDir, "*.txt")
		path, err := ExpandSingle(pattern)
		if err != nil {
			t.Errorf("ExpandSingle() error = %v", err)
			return
		}

		if path != testFile {
			t.Errorf("ExpandSingle() got %s, want %s", path, testFile)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		pattern := filepath.Join(tmpDir, "nonexistent.txt")
		_, err := ExpandSingle(pattern)
		if err != os.ErrNotExist {
			t.Errorf("ExpandSingle() error = %v, want %v", err, os.ErrNotExist)
		}
	})
}

func TestContainsWildcard(t *testing.T) {
	tests := []struct {
		pattern string
		want    bool
	}{
		{"*.txt", true},
		{"file?.txt", true},
		{"file[123].txt", true},
		{"file.txt", false},
		{"/path/to/file.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			if got := containsWildcard(tt.pattern); got != tt.want {
				t.Errorf("containsWildcard(%s) = %v, want %v", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestExpandNoDuplicates(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Same file matched by different patterns should not create duplicates
	patterns := []string{
		testFile,                        // Literal
		filepath.Join(tmpDir, "*.txt"),  // Wildcard
		filepath.Join(tmpDir, "test.*"), // Another wildcard
	}

	matches, err := Expand(patterns)
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("Expand() returned %d matches (expected 1 due to deduplication)", len(matches))
	}
}
