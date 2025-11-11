package vault

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/testutils"
	"go-pass/utils"
)

func TestSearchVault(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	tests := []struct {
		name  string
		args  string
		error bool
	}{
		{name: "no args", args: "", error: true},
		{name: "one arg", args: "test2", error: false},
	}
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))

	cfgFile, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		key,
	)
	assert.NoError(err)
	defer cfgFile.Close()

	vaultFile, err := utils.CreateVault(testutils.TEST_VAULT_NAME, key)
	assert.NoError(err)
	defer vaultFile.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		VaultName:      testutils.TEST_VAULT_NAME,
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
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

	assert.NoError(err1)
	assert.NoError(err2)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SearchVault(tt.args, cfg, key)
			if tt.error {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
