package crypt

import (
	"encoding/json"
	"log"

	"go-pass/model"
)

// EncryptPassword encrypts the password with AES-256 GCM.
// This function returns the base64 encoded AES-GCM encrypted password.
func EncryptPassword(pw []byte, keychain *model.MasterAESKeyManager) (string, error) {
	return keychain.Encrypt(pw)
}

// EncryptVault encrypts the whole model.VaultEntry struct to be stored locally
// on disc. This returns the base64 encoded AES-GCM encrypted vault.
// For the most part, after this function is called, utils.WriteToFile() gets
// called
func EncryptVault(vault []model.VaultEntry, keychain *model.MasterAESKeyManager) (string, error) {
	b, err := json.Marshal(vault)
	if err != nil {
		log.Fatalf("EncryptVault::Marshal json: %v", err)
		return "", err
	}

	return keychain.Encrypt(b)
}

// EncryptConfig encrypts the config with AES-256 GCM
func EncryptConfig(cfg model.Config, keychain *model.MasterAESKeyManager) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("EncryptConfig::Marshal json: %v", err)
		return "", err
	}

	return keychain.Encrypt(b)
}
