package cmd

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

var (
	TEST_VAULT_NAME      = "test-vault.json"
	TEST_CONFIG_NAME     = "test-cfg.json"
	TEST_MASTER_PASSWORD = []byte("mastahpass")
)

func TestAddCheckConfig(t *testing.T) {
	cleanup()
	tests := []struct {
		name          string
		configPresent bool
	}{
		{
			name:          "config not created",
			configPresent: false,
		},
		{
			name:          "config is created",
			configPresent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			if tt.configPresent {
				f, err := utils.CreateConfig(
					TEST_VAULT_NAME,
					TEST_MASTER_PASSWORD,
					TEST_CONFIG_NAME,
				)
				assert.NoError(err)
				defer f.Close()

				v, err := utils.CreateVault(TEST_VAULT_NAME)
				assert.NoError(err)
				defer v.Close()
			}

			cfg, err := CheckConfig(TEST_CONFIG_NAME)

			time := cfg.LastVisited
			if tt.configPresent {
				assert.Equal(model.Config{
					MasterPassword: TEST_MASTER_PASSWORD,
					VaultName:      TEST_VAULT_NAME,
					LastVisited:    time,
				}, cfg)
				assert.NoError(err)
			} else {
				assert.Error(err)
				assert.Equal(model.Config{}, cfg)
			}
			cleanup()
		})

		time.Sleep(time.Millisecond * 100)
	}
}

func TestAddAddToVault(t *testing.T) {
	defer cleanup()

	cfgF, err := utils.CreateConfig(TEST_VAULT_NAME, TEST_MASTER_PASSWORD, TEST_CONFIG_NAME)
	assert.NoError(t, err)
	cfgF.Close()

	vaultF, err := utils.CreateVault(TEST_VAULT_NAME)
	assert.NoError(t, err)
	defer vaultF.Close()

	now := time.Now().UnixMilli()
	source := "test"
	ui := model.UserInput{
		Username: "test",
		Password: []byte("test"),
		Notes:    "",
	}
	cfg := model.Config{
		MasterPassword: []byte("yeahnah"),
		VaultName:      TEST_VAULT_NAME,
		LastVisited:    now - time.Hour.Milliseconds(),
	}
	err = AddToVault(source, ui, cfg, now)
	assert.NoError(t, err)

	fStat, _ := vaultF.Stat()
	assert.Greater(t, fStat.Size(), int64(2))
}

func cleanup() {
	os.Remove(path.Join(utils.VAULT_PATH, TEST_VAULT_NAME))
	os.Remove(path.Join(utils.CONFIG_PATH, TEST_CONFIG_NAME))
}
