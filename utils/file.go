package utils

import (
	"log"
	"os"
	"path"
)

// OpenVault opens the vault file in which the passwords are stored. It is up to
// the caller to close the opened file.
func OpenVault() *os.File {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error UserHomeDir: %v\n", err)
	}

	err = os.Mkdir(path.Join(home, ".local", "gopass"), 0700)
	if !os.IsExist(err) {
		log.Fatalf("Error creating dir: %v\n", err)
	}

	vaultPath := path.Join(home, ".local", "gopass", "pass.json")
	f, err := os.OpenFile(vaultPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", vaultPath, err)
	}
	return f
}
