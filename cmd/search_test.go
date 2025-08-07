package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

func TestSearchVault(t *testing.T) {
	tests := []struct {
		name  string
		args  string
		error bool
	}{
		{name: "no args", args: "", error: true},
		{name: "one arg", args: "test2", error: false},
	}
	defer cleanup()

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
		VaultName:      utils.TEST_VAULT_NAME,
		MasterPassword: utils.TEST_MASTER_PASSWORD,
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

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SearchVault(tt.args, cfg)
			assert := assert.New(t)
			if tt.error {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
