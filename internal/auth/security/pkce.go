package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type PKCEManager struct {
	verifiers map[string]PKCEData
	mutex     sync.RWMutex
}

type PKCEData struct {
	CodeVerifier string
	CreatedAt    time.Time
	State        string
}

var pkceManager = &PKCEManager{
	verifiers: make(map[string]PKCEData),
}

// GeneratePKCE generates code_verifier and code_challenge for PKCE
func GeneratePKCE(state string) (codeVerifier, codeChallenge string, err error) {
	// Generate code_verifier (43-128 characters)
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	codeVerifier = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)

	// Generate code_challenge using SHA256
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])

	// Temporarily store
	pkceManager.mutex.Lock()
	defer pkceManager.mutex.Unlock()

	pkceManager.verifiers[state] = PKCEData{
		CodeVerifier: codeVerifier,
		CreatedAt:    time.Now(),
		State:        state,
	}

	return codeVerifier, codeChallenge, nil
}

// ValidateAndRetrievePKCE validates and retrieves the code_verifier
func ValidateAndRetrievePKCE(state string) (string, error) {
	pkceManager.mutex.Lock()
	defer pkceManager.mutex.Unlock()

	data, exists := pkceManager.verifiers[state]
	if !exists {
		return "", fmt.Errorf("PKCE verifier not found for state")
	}

	// Remove after use (one-time use)
	delete(pkceManager.verifiers, state)

	// Check expiration (5 minutes)
	if time.Since(data.CreatedAt) > 5*time.Minute {
		return "", fmt.Errorf("PKCE verifier expired")
	}

	return data.CodeVerifier, nil
}

// CleanupExpiredPKCE removes expired verifiers
func CleanupExpiredPKCE() {
	pkceManager.mutex.Lock()
	defer pkceManager.mutex.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	for state, data := range pkceManager.verifiers {
		if data.CreatedAt.Before(cutoff) {
			delete(pkceManager.verifiers, state)
		}
	}
}
