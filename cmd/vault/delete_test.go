package vault

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

func TestDeleteItemInVault(t *testing.T) {
	defer utils.TestCleanup()

	cfgFile, err := utils.CreateConfig(
		utils.TEST_VAULT_NAME,
		utils.TEST_MASTER_PASSWORD,
		utils.TEST_CONFIG_NAME,
	)
	assert.NoError(t, err)
	defer cfgFile.Close()

	vaultFile, err := utils.CreateVault(utils.TEST_VAULT_NAME)
	assert.NoError(t, err)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: utils.TEST_MASTER_PASSWORD,
		VaultName:      utils.TEST_VAULT_NAME,
		LastVisited:    now,
	}

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now)
	assert.NoError(t, err)

	r := strings.NewReader("y\n")
	err = DeleteItemInVault(cfg, vaultEntry2, r)
	assert.NoError(t, err)
}

func TestDeleteItemInVaultNotExist(t *testing.T) {
	defer utils.TestCleanup()

	cfgFile, err := utils.CreateConfig(
		utils.TEST_VAULT_NAME,
		utils.TEST_MASTER_PASSWORD,
		utils.TEST_CONFIG_NAME,
	)
	assert.NoError(t, err)
	defer cfgFile.Close()

	vaultFile, err := utils.CreateVault(utils.TEST_VAULT_NAME)
	assert.NoError(t, err)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: utils.TEST_MASTER_PASSWORD,
		VaultName:      utils.TEST_VAULT_NAME,
		LastVisited:    now,
	}

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now)
	assert.NoError(t, err)

	r := strings.NewReader("y\n")
	err = DeleteItemInVault(cfg, "nonExistant", r)
	assert.Error(t, err)
}
