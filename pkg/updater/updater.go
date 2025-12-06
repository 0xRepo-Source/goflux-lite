// Package updater provides self-update functionality for GoFlux Lite binaries.
// It supports downloading updates from remote URLs or local network servers.
package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

// Manifest describes available updates with version and download information.
type Manifest struct {
	Version     string            `json:"version"`      // Semantic version (e.g., "0.2.0")
	ReleaseDate string            `json:"release_date"` // ISO 8601 date
	Binaries    map[string]Binary `json:"binaries"`     // Platform-specific binaries
	Notes       string            `json:"notes"`        // Release notes
}

// Binary contains download information for a specific platform binary.
type Binary struct {
	URL      string `json:"url"`      // Download URL
	Checksum string `json:"checksum"` // SHA-256 checksum
	Size     int64  `json:"size"`     // File size in bytes
}

// Updater manages binary updates from remote or local sources.
type Updater struct {
	CurrentVersion string
	ManifestURL    string
	LocalServer    string // Optional local network server
	client         *http.Client
}

// ProgressFunc reports download progress.
type ProgressFunc func(downloaded, total int64)

// New creates a new Updater with the current version.
func New(currentVersion, manifestURL string) *Updater {
	return &Updater{
		CurrentVersion: currentVersion,
		ManifestURL:    manifestURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetLocalServer configures an optional local network update server.
func (u *Updater) SetLocalServer(serverURL string) {
	u.LocalServer = serverURL
}

// CheckForUpdate fetches the manifest and compares versions.
// Returns the manifest if an update is available, nil otherwise.
func (u *Updater) CheckForUpdate() (*Manifest, error) {
	// Try local server first if configured
	if u.LocalServer != "" {
		manifest, err := u.fetchManifest(u.LocalServer + "/version.json")
		if err == nil && manifest.Version > u.CurrentVersion {
			return manifest, nil
		}
	}

	// Fall back to remote URL
	manifest, err := u.fetchManifest(u.ManifestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}

	if manifest.Version <= u.CurrentVersion {
		return nil, nil // No update available
	}

	return manifest, nil
}

// fetchManifest downloads and parses a version manifest.
func (u *Updater) fetchManifest(url string) (*Manifest, error) {
	resp, err := u.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest request failed: status %d", resp.StatusCode)
	}

	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// DownloadUpdate downloads the binary for the current platform.
func (u *Updater) DownloadUpdate(manifest *Manifest, progress ProgressFunc) (string, error) {
	platform := runtime.GOOS + "_" + runtime.GOARCH
	binary, ok := manifest.Binaries[platform]
	if !ok {
		return "", fmt.Errorf("no binary available for platform: %s", platform)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "gfl-update-*.exe")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Download binary
	resp, err := u.client.Get(binary.URL)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	// Download with progress tracking and checksum calculation
	hash := sha256.New()
	writer := io.MultiWriter(tmpFile, hash)

	var downloaded int64
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := writer.Write(buf[:n]); writeErr != nil {
				os.Remove(tmpFile.Name())
				return "", fmt.Errorf("write failed: %w", writeErr)
			}
			downloaded += int64(n)
			if progress != nil {
				progress(downloaded, binary.Size)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Remove(tmpFile.Name())
			return "", fmt.Errorf("download error: %w", err)
		}
	}

	// Verify checksum
	calculatedChecksum := hex.EncodeToString(hash.Sum(nil))
	if calculatedChecksum != binary.Checksum {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("checksum mismatch: expected %s, got %s", binary.Checksum, calculatedChecksum)
	}

	return tmpFile.Name(), nil
}

// Install replaces the current binary with the downloaded update.
// The current executable must not be running when this is called.
func (u *Updater) Install(downloadedPath string) error {
	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create backup
	backupPath := exePath + ".backup"
	if err := os.Rename(exePath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Copy new binary to executable location
	if err := copyFile(downloadedPath, exePath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, exePath)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Make executable (Unix-like systems)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(exePath, 0755); err != nil {
			return fmt.Errorf("failed to set permissions: %w", err)
		}
	}

	// Clean up
	os.Remove(downloadedPath)

	return nil
}

// Rollback restores the previous version from backup.
func (u *Updater) Rollback() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	backupPath := exePath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("no backup found")
	}

	if err := os.Rename(backupPath, exePath); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}

// GetPlatform returns the current platform identifier.
func GetPlatform() string {
	return runtime.GOOS + "_" + runtime.GOARCH
}
