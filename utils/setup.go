package utils

import (
	"crypto/rand"
	"log"
	"os"
	"path"
)

const (
	NONCE_SIZE = 12
	KEY_SIZE   = 32
)

var (
	home, _    = os.UserHomeDir()
	VAULT_PATH = path.Join(home, ".local", "gopass")
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

// utils.GenerateNonce generates a Number Once, used for AES-256 encryption.
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

// CreateVault creates a file in a default path. If directories aren't created,
// this function will create them.
func CreateVault(name string) *os.File {
	var fName string
	if name == "" {
		fName = "pass.json"
	} else {
		fName = name
	}

	err := os.Mkdir(VAULT_PATH, 0700)
	if !os.IsExist(err) {
		log.Fatalf("CreateVault::Error creating dir: %v\n", err)
	}

	vaultPath := path.Join(VAULT_PATH, fName)
	f, err := os.OpenFile(vaultPath, os.O_RDWR, 0644)
	if !os.IsExist(err) {
		f, err := os.OpenFile(vaultPath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("CreateVault::creating file: %v", err)
		}
		return f
	}
	if err != nil {
		log.Fatalf("CreateVault::Error reading file %s: %v", vaultPath, err)
	}

	return f
}

// OpenVault opens the vault file in which the passwords are stored. It is up to
// the caller to close the opened file.
func OpenVault() *os.File {
	vaultPath := path.Join(VAULT_PATH, "pass.json")
	f, err := os.OpenFile(vaultPath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("OpenVault::Error reading file %s: %v", vaultPath, err)
	}

	return f
}

// WriteToVault takes a *os.File and the contents wanted in the file, in []byte,
// and writes it to the file. It is up to the caller of this function that the
// file is closed.
func WriteToVault(f *os.File, contents []byte) {
	// Reset the file
	if _, err := f.Seek(0, 0); err != nil {
		log.Fatalf("WriteToVault::seek: %v", err)
	}

	if err := f.Truncate(0); err != nil {
		log.Fatalf("WriteToVault::truncate: %v", err)
	}

	// Write to the file
	if _, err := f.Write(contents); err != nil {
		log.Fatalf("WriteToVault::write: %v", err)
	}
}
