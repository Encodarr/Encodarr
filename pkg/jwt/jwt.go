package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidFormat    = fmt.Errorf("invalid token format")
	ErrInvalidSignature = fmt.Errorf("invalid signature")
	ErrTokenExpired     = fmt.Errorf("token has expired")
	ErrInvalidAlgorithm = fmt.Errorf("invalid algorithm")
)

type Claims map[string]interface{}

type Token struct {
	Header    map[string]string
	Claims    Claims
	Signature string
}

func NewToken(claims Claims) *Token {
	if claims == nil {
		claims = Claims{}
	}
	return &Token{
		Header: map[string]string{
			"alg": "HS256",
			"typ": "JWT",
		},
		Claims: claims,
	}
}

func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func base64URLDecode(s string) ([]byte, error) {
	if s == "" {
		return nil, fmt.Errorf("empty input")
	}
	if l := len(s) % 4; l > 0 {
		s += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(s)
}

func (t *Token) Sign(secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf("empty secret")
	}

	headerJSON, err := json.Marshal(t.Header)
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}

	claimsJSON, err := json.Marshal(t.Claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}

	var b strings.Builder
	b.WriteString(base64URLEncode(headerJSON))
	b.WriteByte('.')
	b.WriteString(base64URLEncode(claimsJSON))

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(b.String()))
	signature := base64URLEncode(h.Sum(nil))

	b.WriteByte('.')
	b.WriteString(signature)
	return b.String(), nil
}

func Parse(tokenString string, secret []byte) (*Token, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("empty token")
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("empty secret")
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidFormat
	}

	headerBytes, err := base64URLDecode(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}

	claimsBytes, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode claims: %w", err)
	}

	var header map[string]string
	var claims Claims

	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("invalid header JSON: %w", err)
	}

	if alg := header["alg"]; alg != "HS256" {
		return nil, ErrInvalidAlgorithm
	}

	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("invalid claims JSON: %w", err)
	}

	// Verify expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, ErrTokenExpired
		}
	}

	// Verify signature using constant-time comparison
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(parts[0] + "." + parts[1]))
	expectedSignature := base64URLEncode(h.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(parts[2])) != 1 {
		return nil, ErrInvalidSignature
	}

	return &Token{
		Header:    header,
		Claims:    claims,
		Signature: parts[2],
	}, nil
}
