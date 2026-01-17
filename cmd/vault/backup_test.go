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

var (
	vaultEntry1 = "test1"
	vaultEntry2 = "test2"
	vaultEntry3 = "test3"
)

func TestBackupVault(t *testing.T) {
	// removes the config and vault
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
	cfg := &model.Config{
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

	vf, err := os.OpenFile(path.Join(utils.VAULT_PATH, testutils.TEST_VAULT_NAME), os.O_RDWR, 0o600)
	assert.NoError(err)
	defer vf.Close()

	vaultStat, err := vf.Stat()
	assert.NoError(err)
	vaultSize := vaultStat.Size()

	_, err = BackupVault(
		testutils.TEST_CONFIG_NAME,
		testutils.TEST_VAULT_NAME,
		testutils.TEST_BACKUP_NAME,
		now,
		key,
	)
	assert.NoError(err)

	backupFileName := fmt.Sprintf(testutils.TEST_BACKUP_NAME, now.Format(DATE_FORMAT_STRING))
	backupFile, err := os.Open(
		path.Join(utils.BACKUP_PATH, backupFileName),
	)
	assert.NoError(err)
	defer backupFile.Close()

	backupStat, err := backupFile.Stat()
	assert.NoError(err)
	assert.Greater(backupStat.Size(), int64(0))
	assert.Equal(backupStat.Size(), vaultSize)

	// backup cleanup
	err = os.Remove(path.Join(utils.BACKUP_PATH, backupFileName))
	assert.NoError(err)
}
