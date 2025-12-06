package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/0xRepo-Source/goflux-lite/pkg/errors"
)

func TestNewTokenStore(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	store, err := NewTokenStore(tokenFile)
	if err != nil {
		t.Fatalf("NewTokenStore failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected non-nil store")
	}

	if store.filename != tokenFile {
		t.Errorf("expected filename %s, got %s", tokenFile, store.filename)
	}
}

func TestTokenStore_LoadEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	// Create empty file
	if err := os.WriteFile(tokenFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}

	store, err := NewTokenStore(tokenFile)
	if err != nil {
		t.Fatalf("NewTokenStore failed: %v", err)
	}

	if len(store.tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(store.tokens))
	}
}

func TestTokenStore_LoadWithTokens(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	// Create token file with test data
	testToken := Token{
		ID:          "test-id",
		TokenHash:   "test-hash",
		User:        "testuser",
		Permissions: []string{"read", "write"},
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Revoked:     false,
	}

	storeFile := TokenStoreFile{
		Tokens: []Token{testToken},
	}

	data, err := json.Marshal(storeFile)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	if err := os.WriteFile(tokenFile, data, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	store, err := NewTokenStore(tokenFile)
	if err != nil {
		t.Fatalf("NewTokenStore failed: %v", err)
	}

	if len(store.tokens) != 1 {
		t.Errorf("expected 1 token, got %d", len(store.tokens))
	}

	token, exists := store.tokens[testToken.TokenHash]
	if !exists {
		t.Fatal("expected token to exist in store")
	}

	if token.User != testToken.User {
		t.Errorf("expected user %s, got %s", testToken.User, token.User)
	}
}

func TestTokenStore_GetTokenByID(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	testToken := Token{
		ID:          "test-id",
		TokenHash:   "test-hash",
		User:        "testuser",
		Permissions: []string{"read"},
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Revoked:     false,
	}

	storeFile := TokenStoreFile{
		Tokens: []Token{testToken},
	}

	data, _ := json.Marshal(storeFile)
	os.WriteFile(tokenFile, data, 0644)

	store, _ := NewTokenStore(tokenFile)

	// Test valid token
	token := store.GetTokenByID("test-id")
	if token == nil {
		t.Fatal("expected to find token by ID")
	}

	if token.User != "testuser" {
		t.Errorf("expected user testuser, got %s", token.User)
	}

	// Test non-existent token
	token = store.GetTokenByID("non-existent")
	if token != nil {
		t.Error("expected nil for non-existent token")
	}
}

func TestTokenStore_GetTokenByID_Expired(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	expiredToken := Token{
		ID:          "expired-id",
		TokenHash:   "expired-hash",
		User:        "testuser",
		Permissions: []string{"read"},
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		ExpiresAt:   time.Now().Add(-24 * time.Hour), // Expired
		Revoked:     false,
	}

	storeFile := TokenStoreFile{
		Tokens: []Token{expiredToken},
	}

	data, _ := json.Marshal(storeFile)
	os.WriteFile(tokenFile, data, 0644)

	store, _ := NewTokenStore(tokenFile)

	token := store.GetTokenByID("expired-id")
	if token != nil {
		t.Error("expected nil for expired token")
	}
}

func TestTokenStore_GetTokenByID_Revoked(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	revokedToken := Token{
		ID:          "revoked-id",
		TokenHash:   "revoked-hash",
		User:        "testuser",
		Permissions: []string{"read"},
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Revoked:     true, // Revoked
	}

	storeFile := TokenStoreFile{
		Tokens: []Token{revokedToken},
	}

	data, _ := json.Marshal(storeFile)
	os.WriteFile(tokenFile, data, 0644)

	store, _ := NewTokenStore(tokenFile)

	token := store.GetTokenByID("revoked-id")
	if token != nil {
		t.Error("expected nil for revoked token")
	}
}

func TestTokenStore_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	// Create a token with known hash
	rawToken := "test-token-secret"
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	testToken := Token{
		ID:          "test-id",
		TokenHash:   tokenHash,
		User:        "testuser",
		Permissions: []string{"read", "write"},
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Revoked:     false,
	}

	storeFile := TokenStoreFile{
		Tokens: []Token{testToken},
	}

	data, _ := json.Marshal(storeFile)
	os.WriteFile(tokenFile, data, 0644)

	store, _ := NewTokenStore(tokenFile)

	// Test valid token
	user, perms, err := store.Validate(rawToken)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	if user != "testuser" {
		t.Errorf("expected user testuser, got %s", user)
	}

	if len(perms) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(perms))
	}

	// Test invalid token
	_, _, err = store.Validate("wrong-token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
	if !errors.IsAuthError(err) {
		t.Error("expected AuthError for invalid token")
	}
	if errType, ok := errors.GetAuthErrorType(err); ok {
		if errType != errors.AuthErrorInvalidToken {
			t.Errorf("expected AuthErrorInvalidToken, got %v", errType)
		}
	}
}

func TestTokenStore_Validate_Expired(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	rawToken := "expired-token-secret"
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	expiredToken := Token{
		ID:          "expired-id",
		TokenHash:   tokenHash,
		User:        "testuser",
		Permissions: []string{"read"},
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		ExpiresAt:   time.Now().Add(-24 * time.Hour),
		Revoked:     false,
	}

	storeFile := TokenStoreFile{
		Tokens: []Token{expiredToken},
	}

	data, _ := json.Marshal(storeFile)
	os.WriteFile(tokenFile, data, 0644)

	store, _ := NewTokenStore(tokenFile)

	_, _, err := store.Validate(rawToken)
	if err == nil {
		t.Error("expected error for expired token")
	}
	if errType, ok := errors.GetAuthErrorType(err); ok {
		if errType != errors.AuthErrorExpiredToken {
			t.Errorf("expected AuthErrorExpiredToken, got %v", errType)
		}
	}
}

func TestTokenStore_Validate_Revoked(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "tokens.json")

	rawToken := "revoked-token-secret"
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	revokedToken := Token{
		ID:          "revoked-id",
		TokenHash:   tokenHash,
		User:        "testuser",
		Permissions: []string{"read"},
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Revoked:     true,
	}

	storeFile := TokenStoreFile{
		Tokens: []Token{revokedToken},
	}

	data, _ := json.Marshal(storeFile)
	os.WriteFile(tokenFile, data, 0644)

	store, _ := NewTokenStore(tokenFile)

	_, _, err := store.Validate(rawToken)
	if err == nil {
		t.Error("expected error for revoked token")
	}
	if errType, ok := errors.GetAuthErrorType(err); ok {
		if errType != errors.AuthErrorRevokedToken {
			t.Errorf("expected AuthErrorRevokedToken, got %v", errType)
		}
	}
}

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name        string
		permissions []string
		required    string
		expected    bool
	}{
		{
			name:        "exact match",
			permissions: []string{"read", "write"},
			required:    "read",
			expected:    true,
		},
		{
			name:        "wildcard permission",
			permissions: []string{"*"},
			required:    "anything",
			expected:    true,
		},
		{
			name:        "no match",
			permissions: []string{"read"},
			required:    "write",
			expected:    false,
		},
		{
			name:        "empty permissions",
			permissions: []string{},
			required:    "read",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPermission(tt.permissions, tt.required)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
