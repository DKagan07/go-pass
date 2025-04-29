package utils

import (
	"log"
	"os"
	"path"

	"go-pass/crypt"
	"go-pass/model"
)

var (
	home, _    = os.UserHomeDir()
	VAULT_PATH = path.Join(home, ".local", "gopass")
)

// CreateVault creates a file in a default path. If directories aren't created,
// this function will create them.
func CreateVault(name string) *os.File {
	fName := name
	if name == "" {
		fName = "pass.json"
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

		fileStat, err := f.Stat()
		if err != nil {
			panic("init::getting stat on file")
		}

		if fileStat.Size() == 0 {
			ve := []model.VaultEntry{}
			b, err := crypt.EncryptVault(ve)
			if err != nil {
				panic("init::encrypt ve")
			}
			WriteToVault(f, b)
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
func OpenVault(name string) *os.File {
	fName := name
	if name == "" {
		fName = "pass.json"
	}
	vaultPath := path.Join(VAULT_PATH, fName)
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
