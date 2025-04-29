package crypt

import (
	"crypto/rand"
	"log"
	"os"
)

const (
	NONCE_SIZE = 12
	KEY_SIZE   = 32
)

// GetAESKey is a helper function that gets the key from the environment
// variable called 'SECRET_PASSWORD_KEY'. 'SECRET_PASSWORD_KEY' should be a
// strong, 32-byte password that's unique. You can generate a strong password
// from something like: 'https://passwords-generator.org/32-character'
func GetAESKey() []byte {
	key := []byte(os.Getenv("SECRET_PASSWORD_KEY"))
	if len(key) != KEY_SIZE {
		log.Fatal("GetAESKey::Key not appropriate length")
	}
	return key
}

// GenerateNonce generates a Number Once, used for AES-256 encryption.
func GenerateNonce() []byte {
	nonce := make([]byte, NONCE_SIZE)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatalf("EncryptVault::creating nonce: %v", err)
	}

	if len(nonce) != NONCE_SIZE {
		log.Fatal("GenerateNonce::nonce not correct length")
	}
	return nonce
}
