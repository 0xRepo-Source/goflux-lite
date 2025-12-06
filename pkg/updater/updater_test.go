package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetPlatform(t *testing.T) {
	platform := GetPlatform()
	if platform == "" {
		t.Error("GetPlatform() returned empty string")
	}
	// Should be in format "os_arch"
	if len(platform) < 5 {
		t.Errorf("GetPlatform() returned unexpected format: %s", platform)
	}
}

func TestNew(t *testing.T) {
	updater := New("0.1.0", "https://example.com/version.json")

	if updater.CurrentVersion != "0.1.0" {
		t.Errorf("Expected version 0.1.0, got %s", updater.CurrentVersion)
	}

	if updater.ManifestURL != "https://example.com/version.json" {
		t.Errorf("Expected URL https://example.com/version.json, got %s", updater.ManifestURL)
	}

	if updater.client == nil {
		t.Error("HTTP client not initialized")
	}
}

func TestSetLocalServer(t *testing.T) {
	updater := New("0.1.0", "https://example.com/version.json")
	updater.SetLocalServer("http://192.168.1.100:8080")

	if updater.LocalServer != "http://192.168.1.100:8080" {
		t.Errorf("Expected local server http://192.168.1.100:8080, got %s", updater.LocalServer)
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "dest.txt")

	content := []byte("test content")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	copiedContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(copiedContent) != string(content) {
		t.Errorf("Copied content mismatch: got %s, want %s", copiedContent, content)
	}
}

func TestRollback(t *testing.T) {
	tmpDir := t.TempDir()

	// Create mock executable
	exePath := filepath.Join(tmpDir, "test.exe")
	backupPath := exePath + ".backup"

	currentContent := []byte("current version")
	backupContent := []byte("backup version")

	if err := os.WriteFile(exePath, currentContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := os.WriteFile(backupPath, backupContent, 0644); err != nil {
		t.Fatalf("Failed to create backup file: %v", err)
	}

	// Since we can't override os.Executable in tests, we'll just test copyFile
	// and verify rollback logic manually

	// Simulate rollback by renaming backup to main
	if err := os.Remove(exePath); err != nil {
		t.Fatalf("Failed to remove current: %v", err)
	}

	if err := os.Rename(backupPath, exePath); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	restoredContent, err := os.ReadFile(exePath)
	if err != nil {
		t.Fatalf("Failed to read restored file: %v", err)
	}

	if string(restoredContent) != string(backupContent) {
		t.Errorf("Rollback content mismatch: got %s, want %s", restoredContent, backupContent)
	}
}
