package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type AdminKey struct {
	ID        string    `json:"id"`
	AdminID   string    `json:"admin_id"`
	KeyHash   string    `json:"key_hash"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func GenerateSecureToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	token := make([]byte, length)
	for i := range token {
		token[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(token)
}

func ValidateKey(keyFile, inputKey string) bool {
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return false
	}

	var key AdminKey
	if err := json.Unmarshal(data, &key); err != nil {
		return false
	}

	if time.Now().After(key.ExpiresAt) {
		return false
	}

	hash := sha256.Sum256([]byte(inputKey))
	return hex.EncodeToString(hash[:]) == key.KeyHash
}

func GenerateAdminToken(adminID string) string {
	raw := fmt.Sprintf("%s:%d:%s", adminID, time.Now().Unix(), GenerateSecureToken(16))
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func ValidateToken(token, expected string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	return token == expected
}
