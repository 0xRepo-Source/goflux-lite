package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Challenge represents an authentication challenge
type Challenge struct {
	Nonce     string    `json:"nonce"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ChallengeStore manages active authentication challenges
type ChallengeStore struct {
	challenges map[string]*Challenge // nonce -> challenge
	mu         sync.RWMutex
}

// NewChallengeStore creates a new challenge store
func NewChallengeStore() *ChallengeStore {
	store := &ChallengeStore{
		challenges: make(map[string]*Challenge),
	}

	// Start cleanup goroutine
	go store.cleanupExpired()

	return store
}

// GenerateChallenge creates a new random challenge
func (cs *ChallengeStore) GenerateChallenge() (*Challenge, error) {
	// Generate 32 random bytes
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	nonce := hex.EncodeToString(nonceBytes)
	challenge := &Challenge{
		Nonce:     nonce,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5 minute expiry
	}

	cs.mu.Lock()
	cs.challenges[nonce] = challenge
	cs.mu.Unlock()

	return challenge, nil
}

// ValidateResponse validates an HMAC response against a challenge
func (cs *ChallengeStore) ValidateResponse(nonce, response, token string) (bool, error) {
	cs.mu.RLock()
	challenge, exists := cs.challenges[nonce]
	cs.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("invalid or expired nonce")
	}

	if time.Now().After(challenge.ExpiresAt) {
		cs.mu.Lock()
		delete(cs.challenges, nonce)
		cs.mu.Unlock()
		return false, fmt.Errorf("challenge expired")
	}

	// Compute expected HMAC: HMAC-SHA256(token, nonce)
	h := hmac.New(sha256.New, []byte(token))
	h.Write([]byte(nonce))
	expectedResponse := hex.EncodeToString(h.Sum(nil))

	// Compare using constant-time comparison
	valid := hmac.Equal([]byte(response), []byte(expectedResponse))

	// Delete used challenge (prevent replay)
	cs.mu.Lock()
	delete(cs.challenges, nonce)
	cs.mu.Unlock()

	return valid, nil
}

// cleanupExpired removes expired challenges periodically
func (cs *ChallengeStore) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cs.mu.Lock()
		now := time.Now()
		for nonce, challenge := range cs.challenges {
			if now.After(challenge.ExpiresAt) {
				delete(cs.challenges, nonce)
			}
		}
		cs.mu.Unlock()
	}
}
