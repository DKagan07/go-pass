package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"log"

	"go-pass/model"
)

// EncryptPassword encrypts the password with AES-256 GCM.
func EncryptPassword(pw []byte) []byte {
	key := GetAESKey()
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

	nonce := GenerateNonce()

	cipherTextPass := aesgcm.Seal(nil, nonce, hexEncodedPw, nil)
	cipherTextPass = append(nonce, cipherTextPass...)

	dst := make([]byte, hex.EncodedLen(len(cipherTextPass)))
	hex.Encode(dst, cipherTextPass)

	return dst
}

// EncryptVault encrypts the whole model.VaultEntry struct to be stored locally on
// disc.
func EncryptVault(vault []model.VaultEntry) ([]byte, error) {
	b, err := json.Marshal(vault)
	if err != nil {
		log.Fatalf("EncryptVault::Marshal json: %v", err)
		return nil, err
	}

	return Encrypt(b)
}

// EncryptConfig encrypts the config with AES-256 GCM
func EncryptConfig(cfg model.Config) ([]byte, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("EncryptConfig::Marshal json: %v", err)
		return nil, err
	}

	return Encrypt(b)
}

// Encrypt holds the encryption logic, encrypting the input bytes with AES-256
// and returns the hex-encoded encrypted bytes.
func Encrypt(b []byte) ([]byte, error) {
	key := GetAESKey()

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("EncryptConfig::creating cipher block: %v", err)
	}

	aesgcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Fatal("EncryptConfig::creating aes gcm")
	}

	nonce := GenerateNonce()

	cipherText := aesgcm.Seal(nil, nonce, b, nil)

	// Appending the nonce to the data bytes
	cipherText = append(nonce, cipherText...)

	dst := make([]byte, hex.EncodedLen(len(cipherText)))
	hex.Encode(dst, cipherText)

	return dst, nil
}
