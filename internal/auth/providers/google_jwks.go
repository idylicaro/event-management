package providers

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type GoogleJWKS struct {
	Keys []GoogleJWK `json:"keys"`
}

type GoogleJWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type GoogleCerts struct {
	certs     map[string]*rsa.PublicKey
	expiresAt time.Time
}

var googleCertsCache *GoogleCerts

// GetGooglePublicKey retrieves and caches Google's public keys
func GetGooglePublicKey(kid string) (*rsa.PublicKey, error) {
	// Checks if the cache is still valid
	if googleCertsCache != nil && time.Now().Before(googleCertsCache.expiresAt) {
		if key, exists := googleCertsCache.certs[kid]; exists {
			return key, nil
		}
	}

	// Fetches new keys
	if err := refreshGoogleCerts(); err != nil {
		return nil, err
	}

	// Returns the requested key
	if key, exists := googleCertsCache.certs[kid]; exists {
		return key, nil
	}

	return nil, fmt.Errorf("key with ID %s not found", kid)
}

func refreshGoogleCerts() error {
	const certsURL = "https://www.googleapis.com/oauth2/v3/certs"

	resp, err := http.Get(certsURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Google certs: %w", err)
	}
	defer resp.Body.Close()

	var jwks GoogleJWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	certs := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" && key.Use == "sig" {
			pubKey, err := parseRSAPublicKey(key.N, key.E)
			if err != nil {
				continue // Skip invalid keys
			}
			certs[key.Kid] = pubKey
		}
	}

	// Cache for 1 hour (Google rotates keys infrequently)
	googleCertsCache = &GoogleCerts{
		certs:     certs,
		expiresAt: time.Now().Add(time.Hour),
	}

	return nil
}

func parseRSAPublicKey(n, e string) (*rsa.PublicKey, error) {
	// Implementation to convert n and e values into an RSA key
	// This is a simplification - you should use a library like go-jose
	return nil, fmt.Errorf("RSA key parsing not implemented - use go-jose library")
}

// ValidateGoogleIDToken validates a Google ID token
func ValidateGoogleIDToken(idToken, clientID string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		// Checks the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Gets the Kid from the header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		// Fetches the corresponding public key
		return GetGooglePublicKey(kid)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token or claims")
	}

	// Required OIDC validations
	if !claims.VerifyIssuer("https://accounts.google.com", true) &&
		!claims.VerifyIssuer("accounts.google.com", true) {
		return nil, fmt.Errorf("invalid issuer")
	}

	if !claims.VerifyAudience(clientID, true) {
		return nil, fmt.Errorf("invalid audience")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
