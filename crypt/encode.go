package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"

	"go-pass/model"
)

// HashPassword hashes the password via a hashing algo.
func HashPassword(pw string) []byte {
	key := generateAESKey()
	src := []byte(pw)
	hexEncodedPw := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(hexEncodedPw, src)

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("creating cipher block: %v", err)
	}

	aesgcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Fatal("creating aes gcm")
	}

	nonce := generateNonce()

	return aesgcm.Seal(nil, nonce, hexEncodedPw, nil)
}

// Encrypt encrypts the whole model.VaultEntry struct to be stored locally on
// disc.
func Encrypt(vault model.VaultEntry) ([]byte, error) {
	b, err := json.Marshal(vault)
	if err != nil {
		log.Fatalf("Marshal json: %v", err)
		return nil, err
	}
	key := generateAESKey()

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("creating cipher block: %v", err)
	}

	aesgcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Fatal("creating aes gcm")
	}

	nonce := generateNonce()

	cipherText := aesgcm.Seal(nil, nonce, b, nil)

	return cipherText, nil
}

// generateAESKey is a helper function that generates the key from an
// environment variable, which should be unique to the user.
// PASSWORD_KEY should be a strong, 32-byte password that's unique. You can
// generate a strong password from something like:
// 'https://passwords-generator.org/32-character'
func generateAESKey() []byte {
	secretPassword := os.Getenv("PASSWORD_KEY")
	b := []byte(secretPassword)

	if _, err := rand.Read(b); err != nil {
		log.Fatalf("creating key: %v", err)
	}
	return b
}

func generateNonce() []byte {
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatalf("creating nonce: %v", err)
	}
	return nonce
}
