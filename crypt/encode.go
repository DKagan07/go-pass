package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"log"

	"go-pass/model"
	"go-pass/utils"
)

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

	nonce := utils.GenerateNonce()

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

	nonce := utils.GenerateNonce()

	cipherText := aesgcm.Seal(nil, nonce, b, nil)

	// Appending the nonce to the data bytes
	cipherText = append(nonce, cipherText...)

	dst := make([]byte, hex.EncodedLen(len(cipherText)))
	hex.Encode(dst, cipherText)

	return dst, nil
}
