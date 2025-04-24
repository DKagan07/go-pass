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

	cipherTextPass := aesgcm.Seal(nil, nonce, hexEncodedPw, nil)
	cipherTextPass = append(nonce, cipherTextPass...)

	dst := make([]byte, hex.EncodedLen(len(cipherTextPass)))
	hex.Encode(dst, cipherTextPass)

	return dst
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

// TODO: There seems to be some replicated code, maybe should break that out?

// func encrypt(b []byte) []byte {
// 	key := utils.GetAESKey()
//
// 	cipherBlock, err := aes.NewCipher(key)
// 	if err != nil {
// 		log.Fatalf("creating cipher block: %v", err)
// 	}
//
// 	aesgcm, err := cipher.NewGCM(cipherBlock)
// 	if err != nil {
// 		log.Fatal("creating aes gcm")
// 	}
//
// 	nonce := utils.GenerateNonce()
//
// 	cipherTextPass := aesgcm.Seal(nil, nonce, b, nil)
// 	cipherTextPass = append(nonce, cipherTextPass...)
//
// 	dst := make([]byte, hex.EncodedLen(len(cipherTextPass)))
// 	hex.Encode(dst, cipherTextPass)
//
// 	return dst
// }
