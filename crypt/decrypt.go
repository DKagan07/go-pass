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
)

// DecryptPassword takes a []byte that we initially stored in the file, decrypt
// it, and return the string form of that password.
func DecryptPassword(passBytes []byte) string {
	b, err := Decrypt(passBytes)
	if err != nil {
		log.Fatalf("DecryptPassword::decrypt: %v", err)
	}

	pass := make([]byte, hex.DecodedLen(len(b)))
	if _, err := hex.Decode(pass, b); err != nil {
		log.Fatalf("DecryptPassword::hex decode password")
	}

	return string(pass)
}

// DecryptVault takes a *os.File (the vault file) and returns a
// []model.VaultEntry. The purpose of this is is to read the contents of the
// file.
func DecryptVault(f *os.File) []model.VaultEntry {
	// Get contents
	contents, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("DecryptVault::reading contents: %v", err)
	}

	b, err := Decrypt(contents)
	if err != nil {
		log.Fatalf("DecryptVault::decrypt: %v", err)
	}

	var entries []model.VaultEntry

	if err = json.Unmarshal(b, &entries); err != nil {
		log.Fatalf("DecryptVault::unmarshal: %v", err)
	}

	return entries
}

// DecryptConfig take the *os.File of the config file. It decrypts it, and
// returns the model.Config.
func DecryptConfig(f *os.File) model.Config {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Fatalf("DecryptConfig::seek: %v", err)
	}

	// Get contents
	contents, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("DecryptConfig::reading contents: %v", err)
	}

	b, err := Decrypt(contents)
	if err != nil {
		log.Fatalf("DecryptVault::decrypt: %v", err)
	}

	var cfg model.Config
	if err = json.Unmarshal(b, &cfg); err != nil {
		log.Fatalf("DecryptConfig::unmarshal: %v", err)
	}

	return cfg
}

// Decrypt hold the decryption logic. It will AES-256 decrypt the contents and
// return the decrypted bytes. It is up to the caller function to then Marshal
// that into the correct struct.
func Decrypt(contents []byte) ([]byte, error) {
	key := GetAESKey()

	// Hex decode first
	hexBuf := make([]byte, hex.DecodedLen(len(contents)))
	_, err := hex.Decode(hexBuf, contents)
	if err != nil {
		log.Fatalf("DecryptVault::decoding hex: %v", err)
	}

	// Nonce is not decrypted in with the AES, so we can grab it after
	// hex-decoding
	nonce := hexBuf[:NONCE_SIZE]

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("DecryptVault::getting cipher block: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("DecryptVault::creating aes gcm: %v", err)
	}

	b, err := aesgcm.Open(nil, nonce, hexBuf[NONCE_SIZE:], nil)
	if err != nil {
		log.Fatalf("DecryptVault::opening gcm: %v", err)
	}
	return b, err
}
