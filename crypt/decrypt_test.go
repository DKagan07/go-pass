package crypt

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/testutils"
)

func TestDecryptVault(t *testing.T) {
	originalVault := []model.VaultEntry{
		{
			Name:      testEntry1,
			Username:  testEntry1,
			Password:  []byte(testEntry1),
			Notes:     "Test notes 1",
			UpdatedAt: time.Now().UnixMilli(),
		},
		{
			Name:      testEntry2,
			Username:  testEntry2,
			Password:  []byte(testEntry2),
			Notes:     "",
			UpdatedAt: time.Now().UnixMilli(),
		},
	}

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)
	assert.NoError(err)

	t.Run("decrypt vault from file", func(t *testing.T) {
		// Encrypt vault
		encrypted, err := EncryptVault(originalVault, key)
		assert.NoError(err)
		assert.NotNil(encrypted)

		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test_vault_*.json")
		assert.NoError(err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// Write encrypted data to file
		_, err = tmpFile.Write([]byte(encrypted))
		assert.NoError(err)

		// Close and reopen file for reading
		tmpFile.Close()
		file, err := os.Open(tmpFile.Name())
		assert.NoError(err)
		defer file.Close()

		// Decrypt vault from file
		decryptedVault := DecryptVault(file, key, false)

		// Verify the decrypted vault matches the original
		assert.Len(decryptedVault, len(originalVault))
		for i, entry := range originalVault {
			assert.Equal(entry.Name, decryptedVault[i].Name)
			assert.Equal(entry.Username, decryptedVault[i].Username)
			assert.Equal(entry.Password, decryptedVault[i].Password)
			assert.Equal(entry.Notes, decryptedVault[i].Notes)
			assert.Equal(entry.UpdatedAt, decryptedVault[i].UpdatedAt)
		}
	})
}

func TestDecryptConfig(t *testing.T) {
	// Create test config
	originalConfig := model.Config{
		MasterPassword: []byte("mastahpass"),
		VaultName:      "test-vault.json",
		LastVisited:    time.Now().UnixMilli(),
	}

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)
	assert.NoError(err)

	t.Run("decrypt config from file", func(t *testing.T) {
		// Encrypt config
		encrypted, err := EncryptConfig(originalConfig, key)
		assert.NoError(err)
		assert.NotNil(encrypted)

		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test_config_*.json")
		assert.NoError(err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// Write encrypted data to file
		_, err = tmpFile.Write([]byte(encrypted))
		assert.NoError(err)

		// Close and reopen file for reading
		tmpFile.Close()
		file, err := os.Open(tmpFile.Name())
		assert.NoError(err)
		defer file.Close()

		// Decrypt config from file
		decryptedConfig := DecryptConfig(file, key, false)

		// Verify the decrypted config matches the original
		assert.Equal(originalConfig.MasterPassword, decryptedConfig.MasterPassword)
		assert.Equal(originalConfig.VaultName, decryptedConfig.VaultName)
		assert.Equal(originalConfig.LastVisited, decryptedConfig.LastVisited)
	})
}

func TestDecryptVault_Empty(t *testing.T) {
	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)
	assert.NoError(err)

	t.Run("decrypt empty vault", func(t *testing.T) {
		// Create empty vault
		emptyVault := []model.VaultEntry{}

		// Encrypt empty vault
		encrypted, err := EncryptVault(emptyVault, key)
		assert.NoError(err)
		assert.NotNil(encrypted)

		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test_empty_vault_*.json")
		assert.NoError(err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// Write encrypted data to file
		_, err = tmpFile.Write([]byte(encrypted))
		assert.NoError(err)

		// Close and reopen file for reading
		tmpFile.Close()
		file, err := os.Open(tmpFile.Name())
		assert.NoError(err)
		defer file.Close()

		// Decrypt vault from file
		decryptedVault := DecryptVault(file, key, false)

		// Verify the decrypted vault is empty
		assert.Len(decryptedVault, 0)
	})
}

func TestDecryptVault_VaultEntryLargeData(t *testing.T) {
	// Create vault with large data
	largeVault := []model.VaultEntry{
		{
			Name:      "Large Data Site",
			Username:  testEntry1,
			Password:  bytes.Repeat([]byte("password"), 100),                  // Large password
			Notes:     string(bytes.Repeat([]byte("Very long notes "), 1000)), // Very long notes
			UpdatedAt: time.Now().UnixMilli(),
		},
	}

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)
	assert.NoError(err)

	t.Run("decrypt vault with large data", func(t *testing.T) {
		// Encrypt vault
		encrypted, err := EncryptVault(largeVault, key)
		assert.NoError(err)
		assert.NotNil(encrypted)

		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test_large_vault_*.json")
		assert.NoError(err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// Write encrypted data to file
		_, err = tmpFile.Write([]byte(encrypted))
		assert.NoError(err)

		// Close and reopen file for reading
		tmpFile.Close()
		file, err := os.Open(tmpFile.Name())
		assert.NoError(err)
		defer file.Close()

		// Decrypt vault from file
		decryptedVault := DecryptVault(file, key, false)

		// Verify the decrypted vault matches the original
		assert.Len(decryptedVault, len(largeVault))
		assert.Equal(largeVault[0].Name, decryptedVault[0].Name)
		assert.Equal(largeVault[0].Username, decryptedVault[0].Username)
		assert.Equal(largeVault[0].Password, decryptedVault[0].Password)
		assert.Equal(largeVault[0].Notes, decryptedVault[0].Notes)
		assert.Equal(largeVault[0].UpdatedAt, decryptedVault[0].UpdatedAt)
	})
}

func TestDecryptPasswordEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "very long password",
			password: string(bytes.Repeat([]byte("a"), 10000)),
		},
		{
			name: "password with all byte values",
			password: func() string {
				b := make([]byte, 256)
				for i := range b {
					b[i] = byte(i)
				}
				return string(b)
			}(),
		},
		{
			name:     "password with null bytes",
			password: "hello\x00world\x00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
			assert := assert.New(t)
			assert.NoError(err)

			// Encrypt password
			encrypted, err := EncryptPassword([]byte(tt.password), key)
			assert.NoError(err)
			assert.NotNil(encrypted)

			// Decrypt password
			decrypted := DecryptPassword([]byte(encrypted), key, false)
			assert.Equal(tt.password, decrypted)
		})
	}
}
