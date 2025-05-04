package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

var (
	vaultEntry1 = "test1"
	vaultEntry2 = "test2"
	vaultEntry3 = "test3"
)

func TestDeleteItemInVault(t *testing.T) {
	defer cleanup()

	cfgFile := utils.CreateConfig(TEST_VAULT_NAME, TEST_MASTER_PASSWORD, TEST_CONFIG_NAME)
	defer cfgFile.Close()

	vaultFile := utils.CreateVault(TEST_VAULT_NAME)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: TEST_MASTER_PASSWORD,
		VaultName:      TEST_VAULT_NAME,
		LastVisited:    now,
	}

	AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now)
	AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now)
	AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now)

	err := DeleteItemInVault(cfg, vaultEntry2)
	assert.NoError(t, err)
}

func TestDeleteItemInVaultNotExist(t *testing.T) {
	defer cleanup()

	cfgFile := utils.CreateConfig(TEST_VAULT_NAME, TEST_MASTER_PASSWORD, TEST_CONFIG_NAME)
	defer cfgFile.Close()

	vaultFile := utils.CreateVault(TEST_VAULT_NAME)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: TEST_MASTER_PASSWORD,
		VaultName:      TEST_VAULT_NAME,
		LastVisited:    now,
	}

	AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now)
	AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now)
	AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now)

	err := DeleteItemInVault(cfg, "nonExistant")
	assert.Error(t, err)
}
