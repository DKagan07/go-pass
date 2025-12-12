package vault

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/testutils"
	"go-pass/utils"
)

func TestAddCheckConfig(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
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

			keyManager, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
			assert.NoError(err, "Failed to initialize test keyring")

			if tt.configPresent {
				f, err := utils.CreateConfig(
					testutils.TEST_VAULT_NAME,
					testutils.TEST_MASTER_PASSWORD,
					testutils.TEST_CONFIG_NAME,
					keyManager,
				)
				assert.NoError(err)
				defer f.Close()

				v, err := utils.CreateVault(testutils.TEST_VAULT_NAME, keyManager)
				assert.NoError(err)
				defer v.Close()
			}

			cfg, err := utils.CheckConfig(testutils.TEST_CONFIG_NAME, keyManager)

			if tt.configPresent {
				assert.NotNil(cfg)
				assert.Equal(&model.Config{
					MasterPassword: testutils.TEST_MASTER_PASSWORD,
					VaultName:      testutils.TEST_VAULT_NAME,
					LastVisited:    cfg.LastVisited,
					Timeout:        utils.THIRTY_MINUTES,
				}, cfg)
				assert.NoError(err)
			} else {
				assert.Error(err)
				assert.Nil(cfg)
			}

			testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
		})

		time.Sleep(time.Millisecond * 100)
	}
}

func TestAddAddToVault(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))

	keyManager, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(t, err, "Failed to initialize test keyring")

	cfgF, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		keyManager,
	)
	assert.NoError(t, err)
	defer cfgF.Close()

	vaultF, err := utils.CreateVault(testutils.TEST_VAULT_NAME, keyManager)
	assert.NoError(t, err)
	defer vaultF.Close()

	now := time.Now().UnixMilli()
	source := "test"
	ui := model.UserInput{
		Username: "test",
		Password: []byte("test"),
		Notes:    "",
	}
	cfg := &model.Config{
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		VaultName:      testutils.TEST_VAULT_NAME,
		LastVisited:    now - time.Hour.Milliseconds(),
	}

	err = AddToVault(source, ui, cfg, now, keyManager)
	assert.NoError(t, err)

	fStat, _ := vaultF.Stat()
	assert.Greater(t, fStat.Size(), int64(2), "Vault should contain encrypted data")
}
