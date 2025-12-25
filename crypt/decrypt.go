package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go-pass/model"
)

// DecryptPassword takes a []byte that we initially stored in the file, decrypt
// it, and return the string form of that password.
func DecryptPassword(
	passBytes []byte,
	keychain *model.MasterAESKeyManager,
	isOld bool,
) (string, error) {
	var ciphertext []byte
	var err error
	if !isOld {
		ciphertext, err = keychain.Decrypt(string(passBytes))
		if err != nil {
			return "", err
		}
		return string(ciphertext), nil
	} else {
		ciphertext, err = Decrypt(passBytes)
		if err != nil {
			return "", err
		}

		pass := make([]byte, hex.DecodedLen(len(ciphertext)))
		if _, err := hex.Decode(pass, ciphertext); err != nil {
			return "", fmt.Errorf("hex decode password")
		}
		return string(pass), nil
	}
}

// DecryptVault takes a *os.File (the vault file) and returns a
// []model.VaultEntry. The purpose of this is is to read the contents of the
// file.
func DecryptVault(
	f *os.File,
	keychain *model.MasterAESKeyManager,
	isOld bool,
) ([]model.VaultEntry, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seeking for vault: %w", err)
	}

	// Get contents
	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading vault contents: %v", err)
	}

	var ciphertext []byte
	if !isOld {
		ciphertext, err = keychain.Decrypt(string(contents))
		if err != nil {
			return nil, err
		}
	} else {
		ciphertext, err = Decrypt(contents)
		if err != nil {
			return nil, fmt.Errorf("decryping vault: %v", err)
		}
	}

	var entries []model.VaultEntry

	if err = json.Unmarshal(ciphertext, &entries); err != nil {
		return nil, fmt.Errorf("unmarshaling: %v", err)
	}

	return entries, nil
}

// DecryptConfig take the *os.File of the config file. It decrypts it, and
// returns the model.Config.
func DecryptConfig(
	f *os.File,
	keychain *model.MasterAESKeyManager,
	isOld bool,
) (*model.Config, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seeking for config: %w", err)
	}

	// Get contents
	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var b []byte
	if isOld {
		b, err = Decrypt(contents)
		if err != nil {
			return nil, fmt.Errorf("decrypting file: %w", err)
		}
	} else {
		b, err = keychain.Decrypt(string(contents))
		if err != nil {
			return nil, fmt.Errorf("decrypting contents: %w", err)
		}
	}

	var cfg model.Config
	if err = json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling decrypted config: %w", err)
	}

	return &cfg, nil
}

// Decrypt hold the decryption logic. It will AES-256 decrypt the contents and
// return the decrypted bytes. It is up to the caller function to then Marshal
// that into the correct struct.
func Decrypt(contents []byte) ([]byte, error) {
	key := model.GetSalt()

	// Hex decode first
	hexBuf := make([]byte, hex.DecodedLen(len(contents)))
	_, err := hex.Decode(hexBuf, contents)
	if err != nil {
		return nil, fmt.Errorf("decoding hex: %v", err)
	}

	// Nonce is not decrypted in with the AES, so we can grab it after
	// hex-decoding
	nonce := hexBuf[:model.NONCE_SIZE]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("getting cipher block: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating aes gcm: %v", err)
	}

	b, err := aesgcm.Open(nil, nonce, hexBuf[model.NONCE_SIZE:], nil)
	if err != nil {
		return nil, fmt.Errorf("opening gcm: %v", err)
	}
	return b, nil
}
