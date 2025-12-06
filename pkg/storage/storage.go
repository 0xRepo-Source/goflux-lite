// Package storage provides file storage abstractions for goflux-lite.
// It defines a Storage interface and implements a local filesystem backend
// with path traversal protection.
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xRepo-Source/goflux-lite/pkg/errors"
)

// Storage is an interface for storing and retrieving files.
// Implementations must provide thread-safe operations and protect against
// path traversal attacks.
type Storage interface {
	Put(path string, data []byte) error
	Get(path string) ([]byte, error)
	Exists(path string) bool
	List(path string) ([]string, error)
	Delete(path string) error
	Mkdir(path string) error
}

// Local is a local filesystem storage implementation.
// It stores files under a root directory and validates all paths to prevent
// directory traversal attacks.
type Local struct {
	// Root is the base directory for all storage operations
	Root string
}

// NewLocal creates a new local filesystem storage backend rooted at the specified directory.
// The root directory is created if it doesn't exist.
// Returns an error if the directory cannot be created.
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
		return "", errors.NewStorageError(errors.StorageErrorPathTraversal, path, "path traversal attempt detected")
	}

	return fullPath, nil
}

// Put stores data at the specified path within the storage root.
// Parent directories are created automatically. Returns StorageError if the path
// is invalid or attempts directory traversal.
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

// Get retrieves data from the specified path within the storage root.
// Returns StorageError if the path is invalid or attempts directory traversal.
func (l *Local) Get(path string) ([]byte, error) {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	return os.ReadFile(fullPath)
}

// Exists checks if a file or directory exists at the specified path.
// Returns false if the path is invalid or attempts directory traversal.
func (l *Local) Exists(path string) bool {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return false
	}
	_, err = os.Stat(fullPath)
	return err == nil
}

// List returns the names of all entries in the specified directory.
// Returns StorageError if the path is invalid or the directory cannot be read.
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

// Delete removes a file or directory at the specified path.
// Directories are removed recursively. Returns StorageErrorNotFound if the path doesn't exist.
func (l *Local) Delete(path string) error {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if file/directory exists
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return errors.NewStorageError(errors.StorageErrorNotFound, path, "path does not exist")
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

// Mkdir creates a directory at the specified path, including any necessary parent directories.
// Returns StorageError if the path is invalid or attempts directory traversal.
func (l *Local) Mkdir(path string) error {
	fullPath, err := l.sanitizePath(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Create directory with parent directories
	return os.MkdirAll(fullPath, 0755)
}
