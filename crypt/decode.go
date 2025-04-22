package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"

	"go-pass/model"
	"go-pass/utils"
)

func DecodePassword() {}

func DecryptVault(f *os.File) []model.VaultEntry {
	key := utils.GetAESKey()

	// Get contents
	contents, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("DecryptVault::reading contents: %v", err)
	}

	// Hex decode first
	hexBuf := make([]byte, hex.DecodedLen(len(contents)))
	_, err = hex.Decode(hexBuf, contents)
	if err != nil {
		log.Fatalf("DecryptVault::decoding hex: %v", err)
	}

	// Nonce is not decrypted in with the AES, so we can grab it after
	// hex-decoding
	nonce := hexBuf[:utils.NONCE_SIZE]

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("DecryptVault::getting cipher block: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("DecryptVault::creating aes gcm: %v", err)
	}

	b, err := aesgcm.Open(nil, nonce, hexBuf[utils.NONCE_SIZE:], nil)
	if err != nil {
		log.Fatalf("DecryptVault::opening gcm: %v", err)
	}

	var entries []model.VaultEntry

	if err = json.Unmarshal(b, &entries); err != nil {
		log.Fatalf("DecryptVault::unmarshal: %v", err)
	}

	return entries
}
