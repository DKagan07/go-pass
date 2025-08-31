package crypt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
)

func TestEncryptVault(t *testing.T) {
	now := time.Now()
	testVault := []model.VaultEntry{
		{
			Name:      "test",
			Username:  "test_username",
			Password:  []byte("test_password"),
			Notes:     "test_notes",
			UpdatedAt: now.UnixMilli(),
		},
	}

	_, err := EncryptVault(testVault)
	assert.NoError(t, err)
}

func TestEncryptPassword(t *testing.T) {
	testPw := []byte("test_password")

	encPw := EncryptPassword(testPw)
	decPw := DecryptPassword(encPw)

	assert.Equal(t, string(testPw), decPw)
}

func TestEncryptConfig(t *testing.T) {
	now := time.Now()
	testCfg := model.Config{
		MasterPassword: []byte("mastahpassword"),
		VaultName:      "test_vault",
		LastVisited:    now.UnixMilli(),
	}

	_, err := EncryptConfig(testCfg)
	assert.NoError(t, err)
}
