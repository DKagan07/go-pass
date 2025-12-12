package vault

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/testutils"
	"go-pass/utils"
)

func TestDeleteItemInVault(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(t, err)

	cfgFile, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		key,
	)
	assert.NoError(t, err)
	defer cfgFile.Close()

	vaultFile, err := utils.CreateVault(testutils.TEST_VAULT_NAME, key)
	assert.NoError(t, err)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := &model.Config{
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		VaultName:      testutils.TEST_VAULT_NAME,
		LastVisited:    now,
	}

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now, key)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now, key)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now, key)
	assert.NoError(t, err)

	r := strings.NewReader("y\n")
	err = DeleteItemInVault(cfg, vaultEntry2, r, key)
	assert.NoError(t, err)
}

func TestDeleteItemInVaultNotExist(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(t, err)

	cfgFile, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		key,
	)
	assert.NoError(t, err)
	defer cfgFile.Close()

	vaultFile, err := utils.CreateVault(testutils.TEST_VAULT_NAME, key)
	assert.NoError(t, err)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := &model.Config{
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		VaultName:      testutils.TEST_VAULT_NAME,
		LastVisited:    now,
	}

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now, key)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now, key)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now, key)
	assert.NoError(t, err)

	r := strings.NewReader("y\n")
	err = DeleteItemInVault(cfg, "nonExistant", r, key)
	assert.Error(t, err)
}
