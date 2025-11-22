package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Storage is an interface for storing and retrieving files.
type Storage interface {
	Put(path string, data []byte) error
	Get(path string) ([]byte, error)
	Exists(path string) bool
	List(path string) ([]string, error)
	Delete(path string) error
	Mkdir(path string) error
}

// Local is a simple local filesystem storage implementation.
type Local struct {
	Root string
}

// NewLocal creates a new local filesystem storage backend.
func NewLocal(root string) (*Local, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("failed to create root dir: %w", err)
	}
	return &Local{Root: root}, nil
}

// sanitizePath ensures the path cannot escape the root directory
func (l *Local) sanitizePath(path string) (string, error) {
	// Clean the path to resolve . and .. components
	cleanPath := filepath.Clean(path)

	// Remove leading slash to make it relative
	cleanPath = strings.TrimPrefix(cleanPath, "/")
	cleanPath = strings.TrimPrefix(cleanPath, "\\")

	// Join with root and clean again
	fullPath := filepath.Join(l.Root, cleanPath)

	// Get absolute paths to compare
	absRoot, err := filepath.Abs(l.Root)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute root path: %w", err)
	}

	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute file path: %w", err)
	}

	// Ensure the resolved path is still within the root directory
	if !strings.HasPrefix(absPath, absRoot+string(filepath.Separator)) && absPath != absRoot {
		return "", fmt.Errorf("path traversal attempt detected: %s", path)
	}

	return fullPath, nil
}

func (l *Local) Put(path string, data []byte) error {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return os.WriteFile(fullPath, data, 0644)
}

func (l *Local) Get(path string) ([]byte, error) {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	return os.ReadFile(fullPath)
}

func (l *Local) Exists(path string) bool {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return false
	}
	_, err = os.Stat(fullPath)
	return err == nil
}

func (l *Local) List(path string) ([]string, error) {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names, nil
}

func (l *Local) Delete(path string) error {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if file/directory exists
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	// Remove file or directory (recursively)
	if info.IsDir() {
		return os.RemoveAll(fullPath)
	}
	return os.Remove(fullPath)
}

func (l *Local) Mkdir(path string) error {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Create directory with parent directories
	return os.MkdirAll(fullPath, 0755)
}
