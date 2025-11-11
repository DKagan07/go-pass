package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/testutils"
	"go-pass/utils"
)

func TestDoesConfigExist(t *testing.T) {
	tests := []struct {
		name             string
		doesConfigExists bool
	}{
		{
			name:             "config does not exist",
			doesConfigExists: false,
		},
		{
			name:             "config exists",
			doesConfigExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
			defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
			assert := assert.New(t)

			key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
			assert.NoError(err)

			if tt.doesConfigExists {
				f, err := utils.CreateConfig(
					testutils.TEST_VAULT_NAME,
					testutils.TEST_MASTER_PASSWORD,
					testutils.TEST_CONFIG_NAME,
					key,
				)
				assert.NoError(err)
				f.Close()
			}
			assert.Equal(tt.doesConfigExists, DoesConfigExist(testutils.TEST_CONFIG_NAME))
		})
	}
}

func TestCreateFiles(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	err = CreateFiles(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_CONFIG_NAME,
		testutils.TEST_MASTER_PASSWORD,
		key,
	)
	assert.NoError(err)
}

func TestEnsureVaultName(t *testing.T) {
	tests := []struct {
		name      string
		vaultName string
		good      bool
	}{
		{
			name:      "good vault name",
			vaultName: "good.json",
			good:      true,
		},
		{
			name:      "bad vault name",
			vaultName: "bad",
			good:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := EnsureVaultName(tt.vaultName)
			if tt.good {
				assert.Equal(t, s, tt.vaultName)
			} else {
				assert.Equal(t, fmt.Sprintf("%s.json", tt.vaultName), s)
			}
		})
	}
}
