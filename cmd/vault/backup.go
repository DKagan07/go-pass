/*
Copyright © 2025 DKagan07
*/
package vault

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// Note: that this can be formatted with a time.Time.Format(DATE_FORMAT_STRING)
// This can also be parsed with time.Parse(DATE_FORMAT_STRING, date_to_parse)
const (
	DATE_FORMAT_STRING = "2006-01-02_15-04-05"
	BACKUP_FILE_NAME   = "backup__%v.json"
)

// backupCmd represents the backup command
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup your vault",
	Long: `'backup' backups your vault to a directory. These backups are encrypted in the
same way as your vault, and can be restored with the 'restore' command. 

Ex.
	$ gopass vault backup
	Backup YYYY-MM-DD_HH-MM-SS created successfully
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := BackupCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'backup' command: ", err)
			return
		}
	},
}

// BackupCmdHandler is the handler function that encapsulates the logic of
// creating a backup of the vault. This is stored in a different directory
func BackupCmdHandler(cmd *cobra.Command, args []string) error {
	now := time.Now()
	// perhaps check args? Currently there are none
	if len(args) != 0 {
		return errors.New(
			"too many or not enough arugments for 'backup'. see 'help' for correct usage",
		)
	}

	passB, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	keyring := model.NewMasterAESKeyManager(string(passB))

	cfg, err := utils.CheckConfig("", keyring)
	if err != nil {
		return err
	}

	successString, err := BackupVault("", cfg.VaultName, "", now, keyring)
	if err != nil {
		return err
	}
	fmt.Println(successString)
	return nil
}

// BackupVault contains the logic of creating the backup directory, if it
// doesn't exist, create a new backup file following the format of:
// `backup__YYYY-MM-DD_HH-MM-SS.json`. It then copies the contents of the vault
// to the backup file.
func BackupVault(
	configName, vaultName, backupName string,
	now time.Time,
	key *model.MasterAESKeyManager,
) (string, error) {
	if err := os.MkdirAll(utils.BACKUP_DIR, 0o700); err != nil {
		return "", err
	}

	var fn string
	if backupName == "" {
		fn = fmt.Sprintf(BACKUP_FILE_NAME, now.Format(DATE_FORMAT_STRING))
	} else {
		fn = fmt.Sprintf(backupName, now.Format(DATE_FORMAT_STRING))
	}

	backupFilePath := path.Join(utils.BACKUP_DIR, fn)
	_, err := os.Create(backupFilePath)
	if err != nil {
		return "", err
	}

	currentVault, err := utils.OpenVault(vaultName)
	if err != nil {
		return "", err
	}
	defer currentVault.Close()

	entries, err := crypt.DecryptVault(currentVault, key, false)
	if err != nil {
		return "", err
	}

	b, err := crypt.EncryptVault(entries, key)
	if err != nil {
		return "", err
	}

	if err = os.WriteFile(backupFilePath, []byte(b), 0o600); err != nil {
		return "", err
	}

	return fmt.Sprintf("Backup '%s' created successfully", fn), nil
}
