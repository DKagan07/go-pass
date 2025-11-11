package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/testutils"
	"go-pass/utils"
)

func TestRemoveConfig(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	cfgFile, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		key,
	)
	assert.NoError(err)
	defer cfgFile.Close()

	err = RemoveConfig(testutils.TEST_CONFIG_NAME)

	assert.NoError(err)

	info, err := os.Stat(path.Join(utils.CONFIG_PATH, testutils.TEST_CONFIG_NAME))
	assert.Nil(info)
	assert.Error(err)
}

func TestRemoveVault(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	vaultFile, err := utils.CreateVault(testutils.TEST_VAULT_NAME, key)
	assert.NoError(err)
	defer vaultFile.Close()

	err = RemoveVault(testutils.TEST_VAULT_NAME)

	assert.NoError(err)

	info, err := os.Stat(path.Join(utils.VAULT_PATH, testutils.TEST_VAULT_NAME))
	assert.Nil(info)
	assert.Error(err)
}
