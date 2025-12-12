package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"go-pass/crypt"
	"go-pass/model"
)

var (
	home, _        = os.UserHomeDir()
	VAULT_PATH     = path.Join(home, ".local", "gopass")
	CONFIG_PATH    = path.Join(home, ".config", "gopass")
	CONFIG_FILE    = path.Join(CONFIG_PATH, "gopass-cfg.json")
	BACKUP_DIR     = path.Join(home, ".local", "gopass-backup")
	THIRTY_MINUTES = time.Minute.Milliseconds() * 30
)

// CreateVault creates a file in a default path. If directories aren't created,
// this function will create them.
func CreateVault(name string, key *model.MasterAESKeyManager) (*os.File, error) {
	fName := name
	if name == "" {
		fName = "pass.json"
	}

	err := os.MkdirAll(VAULT_PATH, 0o700)
	if err != nil {
		return nil, fmt.Errorf("CreateVault::Error creating dir: %v\n", err)
	}

	vaultPath := path.Join(VAULT_PATH, fName)
	f, err := os.OpenFile(vaultPath, os.O_RDWR, 0o644)
	if !os.IsExist(err) {
		f, err := os.OpenFile(vaultPath, os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			return nil, fmt.Errorf("CreateVault::creating file: %v", err)
		}

		fileStat, err := f.Stat()
		if err != nil {
			panic("init::getting stat on file")
		}

		if fileStat.Size() == 0 {
			ve := []model.VaultEntry{}
			b, err := json.Marshal(ve)
			if err != nil {
				return nil, err
			}
			plaintext, err := key.Encrypt(b)
			if err != nil {
				return nil, err
			}
			WriteToFile(f, plaintext)
		}

		return f, nil
	}
	if err != nil {
		return nil, fmt.Errorf("CreateVault::Error reading file %s: %v", vaultPath, err)
	}

	return f, nil
}

// OpenVault opens the vault file in which the passwords are stored. It is up to
// the caller to close the opened file.
func OpenVault(name string) (*os.File, error) {
	fName := name
	if name == "" {
		fName = "pass.json"
	}
	vaultPath := path.Join(VAULT_PATH, fName)
	f, err := os.OpenFile(vaultPath, os.O_RDWR, 0o644)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("OpenVault::vault file does not exist")
	}
	if err != nil {
		return nil, fmt.Errorf("OpenVault::Error reading file %s: %v", vaultPath, err)
	}

	return f, nil
}

// WriteToFile takes a *os.File and the contents wanted in the file, in []byte,
// and writes it to the file. It is up to the caller of this function that the
// file is closed.
func WriteToFile(f *os.File, contents string) {
	// Reset the file
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Fatalf("WriteToFile::seek: %v", err)
	}

	if err := f.Truncate(0); err != nil {
		log.Fatalf("WriteToFile::truncate: %v", err)
	}

	// Write to the file
	if _, err := f.WriteString(contents); err != nil {
		log.Fatalf("WriteToFile::write: %v", err)
	}
}

// Caller should close these open files
func CreateConfig(
	vaultName string,
	mPass []byte,
	configName string,
	keychain *model.MasterAESKeyManager,
) (*os.File, error) {
	err := os.MkdirAll(CONFIG_PATH, 0o700)
	if err != nil {
		return nil, fmt.Errorf("CreateConfig::Err creating dir: %v", err)
	}

	if configName != "" {
		configName = path.Join(CONFIG_PATH, configName)
	} else {
		configName = CONFIG_FILE
	}

	f, err := os.OpenFile(configName, os.O_RDWR, 0o644)
	if !os.IsExist(err) {
		f, err := os.OpenFile(configName, os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			return nil, fmt.Errorf("CreateVault::creating file: %v", err)
		}

		now := time.Now().UnixMilli()
		cfg := model.Config{
			MasterPassword: mPass,
			VaultName:      vaultName,
			LastVisited:    now,
			Timeout:        THIRTY_MINUTES,
		}

		cipherText, err := crypt.EncryptConfig(&cfg, keychain)
		if err != nil {
			return nil, fmt.Errorf("err in creating cfg ciphertext: %v", err)
		}

		WriteToFile(f, cipherText)

		return f, nil
	}
	if err != nil {
		// log.Fatalf("CreateVault::Error reading file %s: %v", CONFIG_FILE, err)
		return nil, fmt.Errorf("CreateVault::Error reading file %s: %v", CONFIG_FILE, err)
	}

	return f, nil
}

// OpenConfig opens the config file. It returns the file, a boolean whether or
// not the file does not exist (true if it doesn't exist), and an error.
//
// It is up to the caller to close the file
// TODO: This probably isn't the most correct way to do this, but this is ok
// for now
func OpenConfig(fn string) (*os.File, bool, error) {
	if fn != "" {
		fn = path.Join(CONFIG_PATH, fn)
	} else {
		fn = CONFIG_FILE
	}

	f, err := os.OpenFile(fn, os.O_RDWR, 0o644)
	if os.IsNotExist(err) {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, err
	}
	return f, false, nil
}

// IsAccessBeforeLogin returns true if the command being run is before the
// timeout time in the config, in other words, 'logged in', false if otherwise
func IsAccessBeforeLogin(cfg *model.Config, now int64) bool {
	return now <= (cfg.LastVisited + cfg.Timeout)
}

// CheckConfig checks to see if the config file exists. If it does, we return
// the model.Config.
func CheckConfig(fn string, key *model.MasterAESKeyManager) (*model.Config, error) {
	if key == nil {
		return nil, nil
	}

	cfgFile, ok, err := OpenConfig(fn)
	if ok && err == nil {
		fmt.Println("A file is not found. Need to init.")
		return nil, fmt.Errorf("file needs to be created")
	}
	defer cfgFile.Close()
	cfg, err := crypt.DecryptConfig(cfgFile, key, false)
	if err != nil {
		return nil, fmt.Errorf("decryping config: %w", err)
	}
	return cfg, nil
}
