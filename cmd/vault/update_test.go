package vault

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/testutils"
	"go-pass/utils"
)

func TestUpdateEntry(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	cF, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		key,
	)
	assert.NoError(err)
	cF.Close()

	vF, err := utils.CreateVault(testutils.TEST_VAULT_NAME, key)
	assert.NoError(err)
	vF.Close()

	cfg := model.Config{
		VaultName:      testutils.TEST_VAULT_NAME,
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		LastVisited:    time.Now().UnixMilli(),
	}

	err1 := AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, time.Now().UnixMilli(), key)
	assert.NoError(err1)

	i := Inputs{
		Source:   true,
		Username: false,
		Password: false,
		Notes:    false,
	}

	is := InputSources{
		Source:   strings.NewReader("newSource\n"),
		Username: strings.NewReader("newUsername\n"),
		Notes:    strings.NewReader("newNotes\n"),
	}
	err = UpdateEntry(i, cfg, vaultEntry1, is, key)
	assert.NoError(err)

	err = PrintList("newSource", cfg, key)
	assert.NoError(err)

	err = PrintList(vaultEntry1, cfg, key)
	assert.Error(err)
}

func TestUpdateVaultEntry(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	now := time.Now().UnixMilli()
	ve := model.VaultEntry{
		Name:      vaultEntry1,
		Username:  vaultEntry1,
		Password:  []byte(vaultEntry1),
		Notes:     "",
		UpdatedAt: now,
	}

	i := Inputs{
		Source:   true,
		Username: false,
		Password: false,
		Notes:    false,
	}

	is := InputSources{
		Source: strings.NewReader("newSource\n"),
	}

	newNow := time.Now().UnixMilli()
	expected := model.VaultEntry{
		Name:      "newSource",
		Username:  vaultEntry1,
		Password:  []byte(vaultEntry1),
		Notes:     "",
		UpdatedAt: newNow,
	}

	updatedVe, err := UpdateVaultEntry(ve, i, is, key)
	assert.NoError(err)
	assert.Equal(expected, updatedVe)
}
