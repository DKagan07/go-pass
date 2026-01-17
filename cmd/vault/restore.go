/*
Copyright © 2025 DKagan07
*/
package vault

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// restoreCmd represents the restore command
var RestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore restores a backup to your primary vault",
	Long: `'restore' restores a selected backup to become your vault. This is useful for if
anything were to happen to your primary vault, or if you wanted to restore a
previous state, you can.

Use the arrow keys to navigate the backup to be restore, and then press
'enter' to select it. To cancel, press 'Esc'.

Ex.
	$ gopass vault restore
	Select a backup to restore
		> backup__YYYY-MM-DD_HH-MM-SS.json
		backup__YYYY-MM-DD_HH-MM-SS.json

		Vault restored successfully
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RestoreCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'restore' command: ", err)
			return
		}

		fmt.Println("restore called")
	},
}

// RestoreCmdHandler is the handler for the 'restore' command
func RestoreCmdHandler(cmd *cobra.Command, args []string) error {
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

	return RestoreVault(cfg.VaultName, false, keyring)
}

// RestoreVault encapsulates the logic for the 'restore' command. It ensure that
// the vault is not present, creates one, decrypts the backup, encrypts the
// contents and writes it to the new vault.
func RestoreVault(vaultName string, test bool, key *model.MasterAESKeyManager) error {
	// need to make sure the vault is not present
	_, err := utils.OpenVault(vaultName)
	if err == nil {
		return fmt.Errorf("vault already exists")
	}

	entries, err := os.ReadDir(utils.BACKUP_PATH)
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

	restorePath := path.Join(utils.BACKUP_PATH, selection)
	// Opening selected file
	restoreFp, err := os.Open(restorePath)
	if err != nil {
		return err
	}
	defer restoreFp.Close()

	// Get plaintext from file
	backupEntries, err := crypt.DecryptVault(restoreFp, key, false)
	if err != nil {
		return err
	}
	// Encrypt with new nonce
	backupBytes, err := crypt.EncryptVault(backupEntries, key)
	if err != nil {
		return err
	}

	v, err := utils.CreateVault(vaultName, key)
	if err != nil {
		return err
	}

	if err := utils.WriteToFile(v.Name(), model.FileVault, backupBytes); err != nil {
		return err
	}

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
