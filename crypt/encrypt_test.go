package crypt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/testutils"
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
	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(t, err)

	_, err = EncryptVault(testVault, key)
	assert.NoError(t, err)
}

func TestEncryptPassword(t *testing.T) {
	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(t, err)

	testPw := []byte("test_password")

	encPw, err := EncryptPassword(testPw, key)
	assert.NoError(t, err)
	decPw, err := DecryptPassword([]byte(encPw), key, false)
	assert.NoError(t, err)

	assert.Equal(t, string(testPw), decPw)
}

func TestEncryptConfig(t *testing.T) {
	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(t, err)

	now := time.Now()
	testCfg := model.Config{
		MasterPassword: []byte("mastahpassword"),
		VaultName:      "test_vault",
		LastVisited:    now.UnixMilli(),
	}

	_, err = EncryptConfig(&testCfg, key)
	assert.NoError(t, err)
}
