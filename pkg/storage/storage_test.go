package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/0xRepo-Source/goflux-lite/pkg/errors"
)

func TestNewLocal(t *testing.T) {
	tmpDir := t.TempDir()
	storageDir := filepath.Join(tmpDir, "storage")

	local, err := NewLocal(storageDir)
	if err != nil {
		t.Fatalf("NewLocal failed: %v", err)
	}

	if local == nil {
		t.Fatal("expected non-nil local storage")
	}

	// Verify directory was created
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		t.Error("expected storage directory to be created")
	}
}

func TestLocal_Put(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	testData := []byte("test content")
	err := local.Put("test.txt", testData)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Verify file exists
	filePath := filepath.Join(tmpDir, "test.txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != string(testData) {
		t.Errorf("expected %s, got %s", testData, data)
	}
}

func TestLocal_Put_WithSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	testData := []byte("nested content")
	err := local.Put("subdir/nested/test.txt", testData)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Verify file exists in nested directory
	filePath := filepath.Join(tmpDir, "subdir", "nested", "test.txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != string(testData) {
		t.Errorf("expected %s, got %s", testData, data)
	}
}

func TestLocal_Get(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	testData := []byte("get test content")
	local.Put("gettest.txt", testData)

	data, err := local.Get("gettest.txt")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(data) != string(testData) {
		t.Errorf("expected %s, got %s", testData, data)
	}
}

func TestLocal_Get_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	_, err := local.Get("nonexistent.txt")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestLocal_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	// Non-existent file
	if local.Exists("nothere.txt") {
		t.Error("expected false for non-existent file")
	}

	// Create file and test again
	local.Put("exists.txt", []byte("data"))
	if !local.Exists("exists.txt") {
		t.Error("expected true for existing file")
	}
}

func TestLocal_List(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	// Create some files
	local.Put("file1.txt", []byte("data1"))
	local.Put("file2.txt", []byte("data2"))
	local.Put("file3.txt", []byte("data3"))

	names, err := local.List("")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(names) != 3 {
		t.Errorf("expected 3 files, got %d", len(names))
	}

	// Check names are present
	expectedNames := map[string]bool{
		"file1.txt": false,
		"file2.txt": false,
		"file3.txt": false,
	}

	for _, name := range names {
		if _, ok := expectedNames[name]; ok {
			expectedNames[name] = true
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("expected to find %s in list", name)
		}
	}
}

func TestLocal_List_Subdirectory(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	// Create files in subdirectory
	local.Put("subdir/file1.txt", []byte("data1"))
	local.Put("subdir/file2.txt", []byte("data2"))

	names, err := local.List("subdir")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(names) != 2 {
		t.Errorf("expected 2 files in subdir, got %d", len(names))
	}
}

func TestLocal_Delete_File(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	// Create and delete file
	local.Put("deleteme.txt", []byte("data"))

	if !local.Exists("deleteme.txt") {
		t.Fatal("file should exist before deletion")
	}

	err := local.Delete("deleteme.txt")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if local.Exists("deleteme.txt") {
		t.Error("file should not exist after deletion")
	}
}

func TestLocal_Delete_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	// Create directory with files
	local.Put("deldir/file1.txt", []byte("data1"))
	local.Put("deldir/file2.txt", []byte("data2"))

	err := local.Delete("deldir")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if local.Exists("deldir") {
		t.Error("directory should not exist after deletion")
	}
}

func TestLocal_Delete_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	err := local.Delete("nonexistent.txt")
	if err == nil {
		t.Error("expected error when deleting non-existent file")
	}
	if errType, ok := errors.GetStorageErrorType(err); ok {
		if errType != errors.StorageErrorNotFound {
			t.Errorf("expected StorageErrorNotFound, got %v", errType)
		}
	}
}

func TestLocal_Mkdir(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	err := local.Mkdir("newdir")
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	// Verify directory exists
	dirPath := filepath.Join(tmpDir, "newdir")
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("directory should exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("expected directory, got file")
	}
}

func TestLocal_Mkdir_Nested(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	err := local.Mkdir("parent/child/grandchild")
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	// Verify nested directory exists
	dirPath := filepath.Join(tmpDir, "parent", "child", "grandchild")
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("nested directory should exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("expected directory, got file")
	}
}

func TestLocal_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	// Attempt path traversal attacks
	dangerousPaths := []string{
		"../etc/passwd",
		"../../secrets.txt",
		"subdir/../../etc/passwd",
	}

	for _, path := range dangerousPaths {
		err := local.Put(path, []byte("malicious"))
		if err == nil {
			t.Errorf("expected error for path traversal attempt: %s", path)
		}
		if errors.IsStorageError(err) {
			if errType, ok := errors.GetStorageErrorType(err); ok {
				if errType != errors.StorageErrorPathTraversal {
					t.Errorf("expected StorageErrorPathTraversal for %s, got %v", path, errType)
				}
			}
		}

		_, err = local.Get(path)
		if err == nil {
			t.Errorf("expected error for path traversal attempt on Get: %s", path)
		}
	}
}

func TestLocal_AbsolutePath(t *testing.T) {
	tmpDir := t.TempDir()
	local, _ := NewLocal(tmpDir)

	// Test with leading slash (should be handled gracefully)
	err := local.Put("/test.txt", []byte("data"))
	if err != nil {
		t.Fatalf("Put with leading slash should work: %v", err)
	}

	// Should be able to retrieve it
	data, err := local.Get("/test.txt")
	if err != nil {
		t.Fatalf("Get with leading slash should work: %v", err)
	}

	if string(data) != "data" {
		t.Errorf("expected 'data', got %s", data)
	}

	// Also should be retrievable without leading slash
	data, err = local.Get("test.txt")
	if err != nil {
		t.Fatalf("Get without leading slash should work: %v", err)
	}

	if string(data) != "data" {
		t.Errorf("expected 'data', got %s", data)
	}
}
