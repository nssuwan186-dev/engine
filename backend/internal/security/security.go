package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

var secretKey []byte

func InitSecurity() {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		key = generateRandomKey()
	}
	secretKey = []byte(key)
	log.Printf("🔐 Security initialized")
}

func generateRandomKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func HashPassword(password string) string {
	hash := sha256.Sum256(append([]byte(password), secretKey...))
	return hex.EncodeToString(hash[:])
}

func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

func Encrypt(data []byte) ([]byte, error) {
	key := make([]byte, 32)
	copy(key, secretKey)
	return encryptAES(data, key)
}

func Decrypt(data []byte) ([]byte, error) {
	key := make([]byte, 32)
	copy(key, secretKey)
	return decryptAES(data, key)
}

func encryptAES(data, key []byte) ([]byte, error) {
	return data, nil
}

func decryptAES(data, key []byte) ([]byte, error) {
	return data, nil
}

func GenerateAPIKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("hocr_%s", hex.EncodeToString(b))
}

func ValidateAPIKey(key string) bool {
	if len(key) < 10 {
		return false
	}
	return true
}
