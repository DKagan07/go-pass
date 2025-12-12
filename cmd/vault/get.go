/*
Copyright © 2025 DKagan07
*/
package vault

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// getCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get specific information from your vault by source name",
	Long: `'get' gets a specific source, case SENSITIVE, from the vault and returns the
credentials from the name of the source. If the source name has a <Space> in it,
you have to surround the source name with double quotes.

If you want to see if a specific source is in the vault, you can use the:
	'gopass list -n <name-of-source>'
command.

Ex.
$ gopass vault get Google
Name: Google
	Username: <username>
	Password: <human-readable password>
	Notes: <will show if any notes are present>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := GetCmdHandler(cmd, args); err != nil {
			fmt.Printf("Error with 'get' command: %v\n", err)
			return
		}
	},
}

// GetCmdHandler is the handler function that encapsulates the GetItemsFromVault
// logic and runs some checks beforehand.
func GetCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(
			"at least 1 argument needed for the get command. see 'help' for correct usage",
		)
	}

	name := strings.Join(args, " ")

	passB, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}
	keyring := model.NewMasterAESKeyManager(string(passB))

	cfg, err := utils.CheckConfig("", keyring)
	if err != nil {
		return fmt.Errorf("error checking config: %v", err)
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		return fmt.Errorf("cannot access, need to login")
	}

	err = GetItemFromVault(cfg, name, keyring)
	if err != nil {
		return fmt.Errorf("cannot get %s from vault: %v", name, err)
	}

	return nil
}

// GetItemFromVault retreies the 'name' from the vault. If it doesn't exist, an
// error gets returned.
func GetItemFromVault(cfg *model.Config, name string, keyring *model.MasterAESKeyManager) error {
	f, err := utils.OpenVault(cfg.VaultName)
	if err != nil {
		return fmt.Errorf("opening vault: %v", err)
	}
	defer f.Close()

	entries := crypt.DecryptVault(f, keyring, false)

	if len(entries) == 0 {
		return fmt.Errorf("nothing in vault")
	}

	for _, e := range entries {
		if e.Name == name {
			// The \t's are for aligning the text in the terminal
			fmt.Println("From vault:")
			fmt.Println("Name: ", e.Name)
			fmt.Println("\tUsername: \t", e.Username)
			fmt.Println(
				"\tPassword: \t",
				crypt.DecryptPassword(e.Password, keyring, false),
			)

			if len(e.Notes) > 0 {
				fmt.Println("\tNotes: \t\t", e.Notes)
			}
			return nil
		}
	}

	fmt.Printf("'%s' not found.\n", name)

	encryptedCipherText, err := crypt.EncryptVault(entries, keyring)
	if err != nil {
		return fmt.Errorf("add::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
	return fmt.Errorf("'%s' not found in vault", name)
}
