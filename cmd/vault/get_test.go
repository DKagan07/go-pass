package vault

import (
	"testing"
	"time"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/testutils"
	"go-pass/utils"
)

func TestGetItemFromVault_HappyPath(t *testing.T) {
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
		VaultName:      testutils.TEST_VAULT_NAME,
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		LastVisited:    now,
	}

	pass1, err := crypt.EncryptPassword([]byte(vaultEntry1), key)
	assert.NoError(t, err)
	pass2, err := crypt.EncryptPassword([]byte(vaultEntry2), key)
	assert.NoError(t, err)
	pass3, err := crypt.EncryptPassword([]byte(vaultEntry3), key)
	assert.NoError(t, err)

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(pass1),
	}, cfg, now, key)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(pass2),
	}, cfg, now, key)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(pass3),
	}, cfg, now, key)
	assert.NoError(t, err)

	err = GetItemFromVault(cfg, vaultEntry2, false, key)
	assert.NoError(t, err)
}

func TestGetItemFromVault_CopyFlag(t *testing.T) {
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
		VaultName:      testutils.TEST_VAULT_NAME,
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		LastVisited:    now,
	}

	pass1, err := crypt.EncryptPassword([]byte(vaultEntry1), key)
	assert.NoError(t, err)

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(pass1),
	}, cfg, now, key)
	assert.NoError(t, err)

	err = GetItemFromVault(cfg, vaultEntry1, true, key)
	assert.NoError(t, err)

	password, err := clipboard.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, vaultEntry1, password)
}

func TestGetItemFromVault_NotExist(t *testing.T) {
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
		VaultName:      testutils.TEST_VAULT_NAME,
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		LastVisited:    now,
	}

	pass1, err := crypt.EncryptPassword([]byte(vaultEntry1), key)
	assert.NoError(t, err)
	pass2, err := crypt.EncryptPassword([]byte(vaultEntry2), key)
	assert.NoError(t, err)
	pass3, err := crypt.EncryptPassword([]byte(vaultEntry3), key)
	assert.NoError(t, err)

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(pass1),
	}, cfg, now, key)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(pass2),
	}, cfg, now, key)
	assert.NoError(t, err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(pass3),
	}, cfg, now, key)
	assert.NoError(t, err)

	err = GetItemFromVault(cfg, "notExist", false, key)
	assert.Error(t, err)
}
