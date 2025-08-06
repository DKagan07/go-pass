/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/utils"
)

// Note: that this can be formatted with a time.Time.Format(DATE_FORMAT_STRING)
// This can also be parsed with time.Parse(DATE_FORMAT_STRING, date_to_parse)
const (
	DATE_FORMAT_STRING = "2006-01-02_15-04-05"
	BACKUP_FILE_NAME   = "backup__%v.json"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup your vault",
	Long: fmt.Sprintf(`%s

'backup' backups your vault to a directory. These backups are encrypted in the
same way as your vault, and can be restored with the 'restore' command. 

Ex.
	$ gopass backup
	Backup YYYY-MM-DD_HH-MM-SS created successfully
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		if err := BackupCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'backup' command: ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// TODO: Think about any sort of flags that could be added

	// ex.
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// BackupCmdHandler is the handler function that encapsulates the logic of
// creating a backup of the vault. This is stored in a different directory
func BackupCmdHandler(cmd *cobra.Command, args []string) error {
	// perhaps check args? Currently there are none
	if len(args) != 0 {
		return errors.New(
			"too many or not enough arugments for 'backup'. see 'help' for correct usage",
		)
	}

	cfg, err := utils.CheckConfig("")
	if err != nil {
		return err
	}

	now := time.Now()
	if !utils.IsAccessBeforeLogin(cfg, now.UnixMilli()) {
		return fmt.Errorf("cannot access, need to login")
	}

	return BackupVault("", cfg.VaultName, now)
}

// BackupVault contains the logic of creating the backup directory, if it
// doesn't exist, create a new backup file following the format of:
// `backup__YYYY-MM-DD_HH-MM-SS.json`. It then copies the contents of the vault
// to the backup file.
// TODO: Look into making better error handling
func BackupVault(configName, vaultName string, now time.Time) error {
	if err := os.MkdirAll(utils.BACKUP_DIR, 0700); err != nil {
		return err
	}

	fn := fmt.Sprintf(BACKUP_FILE_NAME, now.Format(DATE_FORMAT_STRING))

	backupFilePath := path.Join(utils.BACKUP_DIR, fn)
	_, err := os.Create(backupFilePath)
	if err != nil {
		return err
	}

	currentVault, err := utils.OpenVault(vaultName)
	if err != nil {
		return err
	}
	defer currentVault.Close()

	entries := crypt.DecryptVault(currentVault)
	b, err := crypt.EncryptVault(entries)
	if err != nil {
		return err
	}

	if err = os.WriteFile(backupFilePath, b, 0600); err != nil {
		return err
	}

	fmt.Printf("Backup %s created successfully", fn)
	return nil
}
