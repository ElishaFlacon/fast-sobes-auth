package jwtmanager

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Claims struct {
	UserID          int64  `json:"sub"`
	Email           string `json:"email"`
	PermissionLevel int32  `json:"perm"`
	ExpiresAt       int64  `json:"exp"`
	IssuedAt        int64  `json:"iat"`
}

type Manager struct {
	secret []byte
}

func New(secret string) *Manager {
	if secret == "" {
		secret = "dev-secret"
	}

	return &Manager{secret: []byte(secret)}
}

func (m *Manager) Sign(claims Claims) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	unsigned := header + "." + payload
	signature := m.sign(unsigned)

	return unsigned + "." + signature, nil
}

func (m *Manager) Verify(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	unsigned := strings.Join(parts[:2], ".")
	expectedSignature := m.sign(unsigned)

	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, errors.New("invalid token signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("unmarshal claims: %w", err)
	}

	return &claims, nil
}

func (m *Manager) sign(unsigned string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
