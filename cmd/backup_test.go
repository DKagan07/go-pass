package cmd

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

func TestBackupVault(t *testing.T) {
	// removes the config and vault
	defer cleanup()
	assert := assert.New(t)

	cfgFile, err := utils.CreateConfig(
		utils.TEST_VAULT_NAME,
		utils.TEST_MASTER_PASSWORD,
		utils.TEST_CONFIG_NAME,
	)
	assert.NoError(err)
	defer cfgFile.Close()

	vaultFile, err := utils.CreateVault(utils.TEST_VAULT_NAME)
	assert.NoError(err)
	defer vaultFile.Close()

	now := time.Now()
	cfg := model.Config{
		MasterPassword: utils.TEST_MASTER_PASSWORD,
		VaultName:      utils.TEST_VAULT_NAME,
		LastVisited:    now.UnixMilli(),
	}

	err = AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, now.UnixMilli())
	assert.NoError(err)
	err = AddToVault(vaultEntry2, model.UserInput{
		Username: vaultEntry2,
		Password: []byte(vaultEntry2),
	}, cfg, now.UnixMilli())
	assert.NoError(err)
	err = AddToVault(vaultEntry3, model.UserInput{
		Username: vaultEntry3,
		Password: []byte(vaultEntry3),
	}, cfg, now.UnixMilli())
	assert.NoError(err)

	vaultStat, err := vaultFile.Stat()
	assert.NoError(err)
	vaultSize := vaultStat.Size()

	err = BackupVault(utils.TEST_CONFIG_NAME, utils.TEST_VAULT_NAME, now)
	assert.NoError(err)

	backupFileName := fmt.Sprintf(BACKUP_FILE_NAME, now.Format(DATE_FORMAT_STRING))
	backupFile, err := os.Open(
		path.Join(utils.BACKUP_DIR, backupFileName),
	)
	assert.NoError(err)
	defer backupFile.Close()

	backupStat, err := backupFile.Stat()
	assert.NoError(err)
	assert.Greater(backupStat.Size(), int64(0))
	assert.Equal(backupStat.Size(), vaultSize)

	// backup cleanup
	err = os.Remove(path.Join(utils.BACKUP_DIR, backupFileName))
	assert.NoError(err)
}
