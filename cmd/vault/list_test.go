package vault

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/testutils"
	"go-pass/utils"
)

func TestPrintList(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	c, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		key,
	)
	c.Close()
	assert.NoError(err)

	v, err := utils.CreateVault(testutils.TEST_VAULT_NAME, key)
	v.Close()
	assert.NoError(err)

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		VaultName:      testutils.TEST_VAULT_NAME,
		LastVisited:    now,
	}

	pass1, err := crypt.EncryptPassword([]byte(vaultEntry1), key)
	assert.NoError(err)
	pass2, err := crypt.EncryptPassword([]byte(vaultEntry2), key)
	assert.NoError(err)
	pass3, err := crypt.EncryptPassword([]byte(vaultEntry3), key)
	assert.NoError(err)

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(pass1),
	}, cfg, now, key)
	assert.NoError(err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(pass2),
	}, cfg, now, key)
	assert.NoError(err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(pass3),
	}, cfg, now, key)
	assert.NoError(err)

	err = PrintList("", cfg, key)
	assert.NoError(err)
}

// This simulates having the flag present
func TestPrintList_OneEntry(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	c, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		key,
	)
	c.Close()
	assert.NoError(err)

	v, err := utils.CreateVault(testutils.TEST_VAULT_NAME, key)
	v.Close()
	assert.NoError(err)

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		VaultName:      testutils.TEST_VAULT_NAME,
		LastVisited:    now,
	}

	err1 := AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now, key)
	err2 := AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now, key)
	err3 := AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now, key)
	assert.NoError(err1)
	assert.NoError(err2)
	assert.NoError(err3)

	err = PrintList(vaultEntry2, cfg, key)
	assert.NoError(err)
}
