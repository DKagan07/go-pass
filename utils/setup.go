package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"go-pass/model"
)

// GetAESKey is a helper function that gets the key from the environment
// variable called 'SECRET_PASSWORD_KEY'. 'SECRET_PASSWORD_KEY' should be a
// strong, 32-byte password that's unique. You can generate a strong password
// from something like: 'https://passwords-generator.org/32-character'
func GetAESKey() []byte {
	return []byte(os.Getenv("SECRET_PASSWORD_KEY"))
}

// OpenVault opens the vault file in which the passwords are stored. It is up to
// the caller to close the opened file.
func OpenVault() *os.File {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("OpenVault::Error UserHomeDir: %v\n", err)
	}

	err = os.Mkdir(path.Join(home, ".local", "gopass"), 0700)
	if !os.IsExist(err) {
		log.Fatalf("OpenVault::Error creating dir: %v\n", err)
	}

	vaultPath := path.Join(home, ".local", "gopass", "pass.json")
	f, err := os.OpenFile(vaultPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("OpenVault::Error reading file %s: %v", vaultPath, err)
	}

	fileStat, err := f.Stat()
	if err != nil {
		log.Fatal("OpenVault::getting stat on file")
	}

	if fileStat.Size() == 0 {
		if _, err = f.Write([]byte("[]")); err != nil {
			log.Fatal("OpenVault::appending empty array to new file")
		}
	}

	return f
}

func GetCurrentVaultEntries(f *os.File) []model.VaultEntry {
	currContents := []model.VaultEntry{}

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&currContents); err != nil {
		log.Fatalf("GetCurrentVaultEntries::decoding: %v", err)
	}
	fmt.Println("currContents: ", currContents)

	return currContents
}

func WriteToVault(f *os.File, contents []byte) {
	// Reset the file
	if _, err := f.Seek(0, 0); err != nil {
		log.Fatalf("WriteToVault::seek: %v", err)
	}

	if err := f.Truncate(0); err != nil {
		log.Fatalf("WriteToVault::truncate: %v", err)
	}

	fmt.Println("contents in WriteToVault: ", contents)
	// Write to the file
	if _, err := f.Write(contents); err != nil {
		log.Fatalf("WriteToVault::write: %v", err)
	}
}
