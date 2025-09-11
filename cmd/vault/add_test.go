package vault

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

func TestAddCheckConfig(t *testing.T) {
	utils.TestCleanup()
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
					utils.TEST_VAULT_NAME,
					utils.TEST_MASTER_PASSWORD,
					utils.TEST_CONFIG_NAME,
				)
				assert.NoError(err)
				defer f.Close()

				v, err := utils.CreateVault(utils.TEST_VAULT_NAME)
				assert.NoError(err)
				defer v.Close()
			}

			cfg, err := utils.CheckConfig(utils.TEST_CONFIG_NAME)

			time := cfg.LastVisited
			if tt.configPresent {
				assert.Equal(model.Config{
					MasterPassword: utils.TEST_MASTER_PASSWORD,
					VaultName:      utils.TEST_VAULT_NAME,
					LastVisited:    time,
					Timeout:        utils.THIRTY_MINUTES,
				}, cfg)
				assert.NoError(err)
			} else {
				assert.Error(err)
				assert.Equal(model.Config{}, cfg)
			}
			utils.TestCleanup()
		})

		time.Sleep(time.Millisecond * 100)
	}
}

func TestAddAddToVault(t *testing.T) {
	defer utils.TestCleanup()

	cfgF, err := utils.CreateConfig(
		utils.TEST_VAULT_NAME,
		utils.TEST_MASTER_PASSWORD,
		utils.TEST_CONFIG_NAME,
	)
	assert.NoError(t, err)
	defer cfgF.Close()

	vaultF, err := utils.CreateVault(utils.TEST_VAULT_NAME)
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
		VaultName:      utils.TEST_VAULT_NAME,
		LastVisited:    now - time.Hour.Milliseconds(),
	}
	err = AddToVault(source, ui, cfg, now)
	assert.NoError(t, err)

	fStat, _ := vaultF.Stat()
	assert.Greater(t, fStat.Size(), int64(2))
}
