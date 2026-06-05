package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

type Claims struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

func NewManager(secret string, ttlMinutes string) Manager {
	minutes, err := strconv.Atoi(ttlMinutes)
	if err != nil || minutes <= 0 {
		minutes = 1440
	}

	return Manager{
		secret: []byte(secret),
		ttl:    time.Duration(minutes) * time.Minute,
	}
}

func (m Manager) Generate(userID uint, email string, role string) (string, error) {
	now := time.Now()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Iat:    now.Unix(),
		Exp:    now.Add(m.ttl).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal token header: %w", err)
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal token claims: %w", err)
	}

	encodedHeader := encodeSegment(headerJSON)
	encodedClaims := encodeSegment(claimsJSON)
	signingInput := encodedHeader + "." + encodedClaims

	return signingInput + "." + m.sign(signingInput), nil
}

func (m Manager) Verify(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(m.sign(signingInput))) {
		return Claims{}, ErrInvalidToken
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}

	if claims.UserID == 0 || claims.Exp == 0 {
		return Claims{}, ErrInvalidToken
	}

	if time.Now().Unix() > claims.Exp {
		return Claims{}, ErrExpiredToken
	}

	return claims, nil
}

func (m Manager) sign(signingInput string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(signingInput))
	return encodeSegment(mac.Sum(nil))
}

func encodeSegment(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
