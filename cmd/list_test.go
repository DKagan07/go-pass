package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

func TestPrintList(t *testing.T) {
	assert := assert.New(t)
	c, err := utils.CreateConfig(TEST_VAULT_NAME, TEST_MASTER_PASSWORD, TEST_CONFIG_NAME)
	c.Close()
	assert.NoError(err)

	v, err := utils.CreateVault(TEST_VAULT_NAME)
	v.Close()
	assert.NoError(err)

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: TEST_MASTER_PASSWORD,
		VaultName:      TEST_VAULT_NAME,
		LastVisited:    now,
	}

	err1 := AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now)
	err2 := AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now)
	err3 := AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now)
	assert.NoError(err1)
	assert.NoError(err2)
	assert.NoError(err3)

	err = PrintList("", cfg)
	assert.NoError(err)
}

// This simulates having the flag present
func TestPrintList_OneEntry(t *testing.T) {
	assert := assert.New(t)
	c, err := utils.CreateConfig(TEST_VAULT_NAME, TEST_MASTER_PASSWORD, TEST_CONFIG_NAME)
	c.Close()
	assert.NoError(err)

	v, err := utils.CreateVault(TEST_VAULT_NAME)
	v.Close()
	assert.NoError(err)

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: TEST_MASTER_PASSWORD,
		VaultName:      TEST_VAULT_NAME,
		LastVisited:    now,
	}

	err1 := AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now)
	err2 := AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now)
	err3 := AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now)
	assert.NoError(err1)
	assert.NoError(err2)
	assert.NoError(err3)

	err = PrintList(vaultEntry2, cfg)
	assert.NoError(err)
}
