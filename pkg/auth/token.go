// Package auth provides authentication and authorization functionality for goflux-lite.
// It implements token-based authentication with support for permissions, expiration,
// and revocation.
package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/0xRepo-Source/goflux-lite/pkg/errors"
)

// Token represents an authentication token with associated metadata.
// Tokens are identified by a hash of the secret value and include
// user information, permissions, and validity period.
type Token struct {
	ID          string    `json:"id"`
	TokenHash   string    `json:"token_hash"`
	User        string    `json:"user"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Revoked     bool      `json:"revoked"`
}

// TokenStore manages authentication tokens with thread-safe access.
// It persists tokens to a JSON file and provides methods for validation,
// loading, and retrieval.
type TokenStore struct {
	mu       sync.RWMutex
	tokens   map[string]*Token // key is token hash
	filename string
}

// TokenStoreFile represents the JSON file format for persisting tokens.
// This structure is used for serialization and deserialization of the token store.
type TokenStoreFile struct {
	Tokens []Token `json:"tokens"`
}

// NewTokenStore creates a new token store that persists to the specified file.
// It automatically loads existing tokens from the file if it exists.
// Returns an error if the file cannot be read or parsed.
func NewTokenStore(filename string) (*TokenStore, error) {
	ts := &TokenStore{
		tokens:   make(map[string]*Token),
		filename: filename,
	}

	if err := ts.Load(); err != nil {
		return nil, err
	}

	return ts, nil
}

// Load reads tokens from the configured file and populates the token store.
// If the file doesn't exist, this is not an error and returns nil.
// Returns an error if the file cannot be read or contains invalid JSON.
func (ts *TokenStore) Load() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	data, err := os.ReadFile(ts.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's okay
			return nil
		}
		return fmt.Errorf("error reading token file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	var storeFile TokenStoreFile
	if err := json.Unmarshal(data, &storeFile); err != nil {
		return fmt.Errorf("error parsing token file: %w", err)
	}

	// Build token map
	ts.tokens = make(map[string]*Token)
	for i := range storeFile.Tokens {
		token := &storeFile.Tokens[i]
		ts.tokens[token.TokenHash] = token
	}

	return nil
}

// Reload reloads tokens from the file, replacing the current in-memory store.
// This is useful for picking up external changes to the token file.
func (ts *TokenStore) Reload() error {
	return ts.Load()
}

// GetTokenByID retrieves a token by its ID for challenge-response authentication.
// Returns nil if the token is not found, revoked, or expired.
func (ts *TokenStore) GetTokenByID(tokenID string) *Token {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	for _, token := range ts.tokens {
		if token.ID == tokenID {
			if token.Revoked || time.Now().After(token.ExpiresAt) {
				return nil
			}
			return token
		}
	}
	return nil
}

// Validate checks if a token string is valid and returns the associated user and permissions.
// The token is hashed before lookup. Returns AuthError types for invalid, revoked, or expired tokens.
func (ts *TokenStore) Validate(tokenStr string) (string, []string, error) {
	// Hash the provided token
	hash := sha256.Sum256([]byte(tokenStr))
	tokenHash := hex.EncodeToString(hash[:])

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	token, exists := ts.tokens[tokenHash]
	if !exists {
		return "", nil, errors.NewAuthError(errors.AuthErrorInvalidToken, "invalid token")
	}

	if token.Revoked {
		return "", nil, errors.NewAuthError(errors.AuthErrorRevokedToken, "token has been revoked")
	}

	if time.Now().After(token.ExpiresAt) {
		return "", nil, errors.NewAuthError(errors.AuthErrorExpiredToken, "token has expired")
	}

	return token.User, token.Permissions, nil
}

// HasPermission checks if a user has a specific permission.
// Returns true if the permissions list contains the required permission or the wildcard "*".
func HasPermission(permissions []string, required string) bool {
	for _, perm := range permissions {
		if perm == required || perm == "*" {
			return true
		}
	}
	return false
}
