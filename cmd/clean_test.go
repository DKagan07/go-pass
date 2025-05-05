package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/utils"
)

func TestRemoveConfig(t *testing.T) {
	cfgFile, _ := utils.CreateConfig(TEST_VAULT_NAME, TEST_MASTER_PASSWORD, TEST_CONFIG_NAME)
	defer cfgFile.Close()

	err := RemoveConfig(TEST_CONFIG_NAME)

	assert := assert.New(t)
	assert.NoError(err)

	info, err := os.Stat(path.Join(utils.CONFIG_PATH, TEST_CONFIG_NAME))
	assert.Nil(info)
	assert.Error(err)
}

func TestRemoveVault(t *testing.T) {
	vaultFile, _ := utils.CreateVault(TEST_VAULT_NAME)
	defer vaultFile.Close()

	err := RemoveVault(TEST_VAULT_NAME)

	assert := assert.New(t)
	assert.NoError(err)

	info, err := os.Stat(path.Join(utils.VAULT_PATH, TEST_VAULT_NAME))
	assert.Nil(info)
	assert.Error(err)
}
