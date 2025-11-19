package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Middleware provides authentication middleware for HTTP handlers
type Middleware struct {
	store          *TokenStore
	challengeStore *ChallengeStore
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(store *TokenStore) *Middleware {
	return &Middleware{
		store:          store,
		challengeStore: NewChallengeStore(),
	}
}

// RequireAuth wraps a handler to require authentication
// Supports both Bearer token and Challenge-Response authentication
func (m *Middleware) RequireAuth(requiredPermission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		var user string
		var permissions []string
		var err error

		// Check if it's challenge-response format: "Challenge <response>;<nonce>;<token_id>"
		if strings.HasPrefix(authHeader, "Challenge ") {
			challengeData := strings.TrimPrefix(authHeader, "Challenge ")
			parts := strings.Split(challengeData, ";")

			if len(parts) != 3 {
				http.Error(w, "Invalid challenge format. Expected: Challenge <response>;<nonce>;<token_id>", http.StatusUnauthorized)
				return
			}

			response, nonce, tokenID := parts[0], parts[1], parts[2]

			// Get token by ID
			token := m.store.GetTokenByID(tokenID)
			if token == nil {
				http.Error(w, "Invalid token ID", http.StatusUnauthorized)
				return
			}

			// Compute expected HMAC: HMAC-SHA256(token_hash, nonce)
			h := hmac.New(sha256.New, []byte(token.TokenHash))
			h.Write([]byte(nonce))
			expectedResponse := hex.EncodeToString(h.Sum(nil))

			// Validate nonce expiry and prevent replay
			_, err := m.challengeStore.ValidateResponse(nonce, response, token.TokenHash)
			if err != nil {
				http.Error(w, fmt.Sprintf("Challenge validation failed: %v", err), http.StatusUnauthorized)
				return
			}

			// Compare responses using constant-time comparison
			if !hmac.Equal([]byte(response), []byte(expectedResponse)) {
				http.Error(w, "Invalid challenge response", http.StatusUnauthorized)
				return
			}

			user = token.User
			permissions = token.Permissions

		} else {
			// Fall back to Bearer token (backward compatibility)
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format. Use: Bearer <token> or Challenge <data>", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate token
			user, permissions, err = m.store.Validate(token)
			if err != nil {
				http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
				return
			}
		}

		// Check permission
		if requiredPermission != "" && !HasPermission(permissions, requiredPermission) {
			http.Error(w, fmt.Sprintf("Permission denied. Required: %s", requiredPermission), http.StatusForbidden)
			return
		}

		// Set user in request context (optional, for logging)
		r.Header.Set("X-Authenticated-User", user)

		// Call the next handler
		next(w, r)
	}
}

// HandleChallenge returns a new authentication challenge
func (m *Middleware) HandleChallenge(w http.ResponseWriter, r *http.Request) {
	challenge, err := m.challengeStore.GenerateChallenge()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate challenge: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(challenge); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode challenge: %v", err), http.StatusInternalServerError)
		return
	}
}

// OptionalAuth wraps a handler to optionally accept authentication
func (m *Middleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				user, _, err := m.store.Validate(parts[1])
				if err == nil {
					r.Header.Set("X-Authenticated-User", user)
				}
			}
		}
		next(w, r)
	}
}
