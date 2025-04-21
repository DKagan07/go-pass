package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"go-pass/model"
	"go-pass/utils"
)

const NONCE_SIZE = 12

// EncryptPassword encrypts the password with AES-256 GSM.
func EncryptPassword(pw []byte) []byte {
	key := utils.GetAESKey()
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
func EncryptVault(vault []model.VaultEntry) ([]byte, error) {
	b, err := json.Marshal(vault)
	if err != nil {
		log.Fatalf("EncryptVault::Marshal json: %v", err)
		return nil, err
	}
	key := utils.GetAESKey()

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("EncryptVault::creating cipher block: %v", err)
	}

	aesgcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Fatal("EncryptVault::creating aes gcm")
	}

	nonce := generateNonce()

	cipherText := aesgcm.Seal(nil, nonce, b, nil)

	// Appending the nonce to the data bytes
	cipherText = append(nonce, cipherText...)

	dst := make([]byte, hex.EncodedLen(len(cipherText)))
	hex.Encode(dst, cipherText)
	fmt.Println("hex encode in encrypt: ", dst)

	return dst, nil
}

// generateNonce generates a Number Once, used for AES-256 encryption.
func generateNonce() []byte {
	nonce := make([]byte, NONCE_SIZE)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatalf("EncryptVault::creating nonce: %v", err)
	}
	return nonce
}
