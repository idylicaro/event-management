package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

// StateManager manages states for CSRF protection
type StateManager struct {
	states map[string]StateData
	mutex  sync.RWMutex
}

type StateData struct {
	CreatedAt time.Time
	Provider  string
	UserAgent string
	IP        string
}

var stateManager = &StateManager{
	states: make(map[string]StateData),
}

// GenerateState generates a unique state for CSRF protection
func GenerateState(provider, userAgent, ip string) (string, error) {
	// Generate 32 random bytes
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}

	state := base64.URLEncoding.EncodeToString(randomBytes)

	stateManager.mutex.Lock()
	defer stateManager.mutex.Unlock()

	// Store the state with metadata
	stateManager.states[state] = StateData{
		CreatedAt: time.Now(),
		Provider:  provider,
		UserAgent: userAgent,
		IP:        ip,
	}

	return state, nil
}

// ValidateState validates and removes a state
func ValidateState(state, provider, userAgent, ip string) error {
	stateManager.mutex.Lock()
	defer stateManager.mutex.Unlock()

	stateData, exists := stateManager.states[state]
	if !exists {
		return fmt.Errorf("invalid or expired state")
	}

	// Remove the state (one-time use)
	delete(stateManager.states, state)

	// Validate expiration (5 minutes)
	if time.Since(stateData.CreatedAt) > 5*time.Minute {
		return fmt.Errorf("state expired")
	}

	// Validate provider
	if stateData.Provider != provider {
		return fmt.Errorf("provider mismatch")
	}

	// Additional validations can be implemented here
	// For example: validate UserAgent and IP for increased security

	return nil
}

// CleanupExpiredStates removes expired states (call periodically)
func CleanupExpiredStates() {
	stateManager.mutex.Lock()
	defer stateManager.mutex.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	for state, data := range stateManager.states {
		if data.CreatedAt.Before(cutoff) {
			delete(stateManager.states, state)
		}
	}
}

// Nonce for ID Tokens
func GenerateNonce() (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
