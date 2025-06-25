package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/utils"
)

func TestRemoveConfig(t *testing.T) {
	assert := assert.New(t)
	cfgFile, err := utils.CreateConfig(utils.TEST_VAULT_NAME, utils.TEST_MASTER_PASSWORD, utils.TEST_CONFIG_NAME)
	assert.NoError(err)
	defer cfgFile.Close()

	err = RemoveConfig(utils.TEST_CONFIG_NAME)

	assert.NoError(err)

	info, err := os.Stat(path.Join(utils.CONFIG_PATH, utils.TEST_CONFIG_NAME))
	assert.Nil(info)
	assert.Error(err)
}

func TestRemoveVault(t *testing.T) {
	assert := assert.New(t)
	vaultFile, err := utils.CreateVault(utils.TEST_VAULT_NAME)
	assert.NoError(err)
	defer vaultFile.Close()

	err = RemoveVault(utils.TEST_VAULT_NAME)

	assert.NoError(err)

	info, err := os.Stat(path.Join(utils.VAULT_PATH, utils.TEST_VAULT_NAME))
	assert.Nil(info)
	assert.Error(err)
}
