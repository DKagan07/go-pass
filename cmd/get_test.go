package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

func TestGetItemFromVault_HappyPath(t *testing.T) {
	defer cleanup()
	cfgFile := utils.CreateConfig(TEST_VAULT_NAME, TEST_MASTER_PASSWORD, TEST_CONFIG_NAME)
	defer cfgFile.Close()
	vaultFile := utils.CreateVault(TEST_VAULT_NAME)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		VaultName:      TEST_VAULT_NAME,
		MasterPassword: TEST_MASTER_PASSWORD,
		LastVisited:    now,
	}

	AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: crypt.EncryptPassword([]byte(vaultEntry1)),
	}, cfg, now)
	AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: crypt.EncryptPassword([]byte(vaultEntry2)),
	}, cfg, now)
	AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: crypt.EncryptPassword([]byte(vaultEntry3)),
	}, cfg, now)

	err := GetItemFromVault(cfg, vaultEntry2)
	assert.NoError(t, err)
}

func TestGetItemFromVault_NotExist(t *testing.T) {
	defer cleanup()
	cfgFile := utils.CreateConfig(TEST_VAULT_NAME, TEST_MASTER_PASSWORD, TEST_CONFIG_NAME)
	defer cfgFile.Close()
	vaultFile := utils.CreateVault(TEST_VAULT_NAME)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		VaultName:      TEST_VAULT_NAME,
		MasterPassword: TEST_MASTER_PASSWORD,
		LastVisited:    now,
	}

	AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: crypt.EncryptPassword([]byte(vaultEntry1)),
	}, cfg, now)
	AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: crypt.EncryptPassword([]byte(vaultEntry2)),
	}, cfg, now)
	AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: crypt.EncryptPassword([]byte(vaultEntry3)),
	}, cfg, now)

	err := GetItemFromVault(cfg, "notExist")
	assert.Error(t, err)
}
