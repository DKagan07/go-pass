package vault

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/testutils"
	"go-pass/utils"
)

func TestRestoreVault(t *testing.T) {
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

	vaultFile, err := utils.CreateVault(testutils.TEST_VAULT_NAME, key)
	assert.NoError(err)
	defer vaultFile.Close()

	now := time.Now()
	cfg := model.Config{
		MasterPassword: testutils.TEST_MASTER_PASSWORD,
		VaultName:      testutils.TEST_VAULT_NAME,
		LastVisited:    now.UnixMilli(),
	}

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now.UnixMilli(), key)
	assert.NoError(err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now.UnixMilli(), key)
	assert.NoError(err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now.UnixMilli(), key)
	assert.NoError(err)

	err = BackupVault(
		testutils.TEST_CONFIG_NAME,
		testutils.TEST_VAULT_NAME,
		testutils.TEST_BACKUP_NAME,
		now,
		key,
	)
	assert.NoError(err)

	backupFn := fmt.Sprintf(testutils.TEST_BACKUP_NAME, now.Format(DATE_FORMAT_STRING))
	// cleanup the backup file after the test is done
	defer os.Remove(path.Join(utils.BACKUP_DIR, backupFn))

	// need to remove the test vault
	testVault := path.Join(utils.VAULT_PATH, testutils.TEST_VAULT_NAME)
	err = os.Remove(testVault)
	assert.NoError(err)

	err = RestoreVault(testutils.TEST_VAULT_NAME, true, key)
	assert.NoError(err)
}
