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

func TestRestoreVault(t *testing.T) {
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

	err = BackupVault(utils.TEST_CONFIG_NAME, utils.TEST_VAULT_NAME, utils.TEST_BACKUP_NAME, now)
	assert.NoError(err)

	backupFn := fmt.Sprintf(utils.TEST_BACKUP_NAME, now.Format(DATE_FORMAT_STRING))
	// cleanup the backup file after the test is done
	defer os.Remove(path.Join(utils.BACKUP_DIR, backupFn))

	// need to remove the test vault
	testVault := path.Join(utils.VAULT_PATH, utils.TEST_VAULT_NAME)
	err = os.Remove(testVault)
	assert.NoError(err)

	err = RestoreVault(utils.TEST_VAULT_NAME, true)
	assert.NoError(err)
}
