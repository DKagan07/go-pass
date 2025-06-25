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
	cfgFile, err := utils.CreateConfig(utils.TEST_VAULT_NAME, utils.TEST_MASTER_PASSWORD, utils.TEST_CONFIG_NAME)
	assert.NoError(t, err)
	defer cfgFile.Close()

	vaultFile, err := utils.CreateVault(utils.TEST_VAULT_NAME)
	assert.NoError(t, err)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		VaultName:      utils.TEST_VAULT_NAME,
		MasterPassword: utils.TEST_MASTER_PASSWORD,
		LastVisited:    now,
	}

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: crypt.EncryptPassword([]byte(vaultEntry1)),
	}, cfg, now)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: crypt.EncryptPassword([]byte(vaultEntry2)),
	}, cfg, now)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: crypt.EncryptPassword([]byte(vaultEntry3)),
	}, cfg, now)
	assert.NoError(t, err)

	err = GetItemFromVault(cfg, vaultEntry2)
	assert.NoError(t, err)
}

func TestGetItemFromVault_NotExist(t *testing.T) {
	defer cleanup()
	cfgFile, err := utils.CreateConfig(utils.TEST_VAULT_NAME, utils.TEST_MASTER_PASSWORD, utils.TEST_CONFIG_NAME)
	assert.NoError(t, err)
	defer cfgFile.Close()
	vaultFile, err := utils.CreateVault(utils.TEST_VAULT_NAME)
	assert.NoError(t, err)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		VaultName:      utils.TEST_VAULT_NAME,
		MasterPassword: utils.TEST_MASTER_PASSWORD,
		LastVisited:    now,
	}

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: crypt.EncryptPassword([]byte(vaultEntry1)),
	}, cfg, now)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: crypt.EncryptPassword([]byte(vaultEntry2)),
	}, cfg, now)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: crypt.EncryptPassword([]byte(vaultEntry3)),
	}, cfg, now)
	assert.NoError(t, err)

	err = GetItemFromVault(cfg, "notExist")
	assert.Error(t, err)
}
