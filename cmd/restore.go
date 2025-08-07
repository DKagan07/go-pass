/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/utils"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore restores a backup to your pirmary vault",
	// TODO: Add a long description here.
	Long: fmt.Sprintf(`%s

restore restores a selected backup to become your vault. This is useful for if 
anything were to happen to your primary vault, or if you wanted to restore a
previous state, you can.

Use the arrow keys to navigate the backup to be restore, and then press
'enter' to select it. To cancel, press 'Esc'.

Ex.
	$ gopass restore
	Select a backup to restore
		> backup__YYYY-MM-DD_HH-MM-SS.json
		backup__YYYY-MM-DD_HH-MM-SS.json

		Vault restored successfully
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		if err := RestoreCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'restore' command: ", err)
			return
		}

		fmt.Println("restore called")
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}

// RestoreCmdHandler is the handler for the 'restore' command
func RestoreCmdHandler(cmd *cobra.Command, args []string) error {
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

	return RestoreVault(cfg.VaultName, false)
}

// RestoreVault encapsulates the logic for the 'restore' command. It ensure that
// the vault is not present, creates one, decrypts the backup, encrypts the
// contents and writes it to the new vault.
func RestoreVault(vaultName string, test bool) error {
	// need to make sure the vault is not present
	_, err := utils.OpenVault(vaultName)
	if err == nil {
		return fmt.Errorf("vault already exists")
	}

	entries, err := os.ReadDir(utils.BACKUP_DIR)
	if err != nil {
		return err
	}

	backupFileNames := []string{}
	for _, entry := range entries {
		backupFileNames = append(backupFileNames, entry.Name())
	}

	var selection string
	// This feels jank, but because I dont' know how to mock the testing of
	// 'huh', I'm just doing this for now.
	if !test {
		selection, err = getSelection(backupFileNames)
		if err != nil {
			return err
		}
	} else {
		for _, entry := range backupFileNames {
			if strings.Contains(entry, "test") {
				selection = entry
			}
		}
	}

	restorePath := path.Join(utils.BACKUP_DIR, selection)
	// Opening selected file
	restoreFp, err := os.Open(restorePath)
	if err != nil {
		return err
	}
	defer restoreFp.Close()

	// Get plaintext from file
	backupEntries := crypt.DecryptVault(restoreFp)
	// Encrypt with new nonce
	backupBytes, err := crypt.EncryptVault(backupEntries)
	if err != nil {
		return err
	}

	v, err := utils.CreateVault(vaultName)
	if err != nil {
		return err
	}

	utils.WriteToFile(v, backupBytes)

	fmt.Printf("Vault restored successfully")
	return nil
}

// getSelection is a helper function that handles the 'huh' functionality for
// selecting the backup file to restore
func getSelection(entryNames []string) (string, error) {
	var selection string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a backup to restore").
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(entryNames...)
				}, &selection).
				Value(&selection).
				Height(10),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}

	return selection, nil
}
