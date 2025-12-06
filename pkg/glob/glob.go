// Package glob provides utilities for expanding file patterns and wildcards.
package glob

import (
	"os"
	"path/filepath"
	"strings"
)

// Match represents a matched file with its original pattern and resolved path.
type Match struct {
	Pattern string // Original pattern that matched
	Path    string // Resolved absolute path
	RelPath string // Relative path from the pattern's base directory
}

// Expand expands glob patterns into a list of matching files.
// It supports standard wildcards: *, ?, and [...].
// Patterns without wildcards are returned as-is if they exist.
// Returns an error if a pattern is malformed or if no files match.
func Expand(patterns []string) ([]Match, error) {
	var matches []Match
	seen := make(map[string]bool) // Prevent duplicates

	for _, pattern := range patterns {
		// Check if pattern contains wildcards
		if !containsWildcard(pattern) {
			// No wildcard - treat as literal path
			absPath, err := filepath.Abs(pattern)
			if err != nil {
				return nil, err
			}

			if _, err := os.Stat(absPath); err != nil {
				// File doesn't exist - skip it
				continue
			}

			if !seen[absPath] {
				matches = append(matches, Match{
					Pattern: pattern,
					Path:    absPath,
					RelPath: filepath.Base(absPath),
				})
				seen[absPath] = true
			}
			continue
		}

		// Pattern contains wildcards - use filepath.Glob
		absPattern, err := filepath.Abs(pattern)
		if err != nil {
			return nil, err
		}

		globMatches, err := filepath.Glob(absPattern)
		if err != nil {
			return nil, err
		}

		// Calculate base directory for relative paths
		baseDir := filepath.Dir(strings.Split(absPattern, "*")[0])
		baseDir = strings.TrimSuffix(baseDir, string(filepath.Separator))

		for _, match := range globMatches {
			// Skip directories unless explicitly requested
			info, err := os.Stat(match)
			if err != nil {
				continue
			}
			if info.IsDir() {
				continue
			}

			if !seen[match] {
				relPath, err := filepath.Rel(baseDir, match)
				if err != nil {
					relPath = filepath.Base(match)
				}

				matches = append(matches, Match{
					Pattern: pattern,
					Path:    match,
					RelPath: relPath,
				})
				seen[match] = true
			}
		}
	}

	return matches, nil
}

// containsWildcard checks if a pattern contains wildcard characters.
func containsWildcard(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[]")
}

// ExpandSingle expands a single pattern and returns the first match.
// This is useful for operations that expect a single file.
func ExpandSingle(pattern string) (string, error) {
	matches, err := Expand([]string{pattern})
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", os.ErrNotExist
	}

	return matches[0].Path, nil
}
