/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get specific information from your vault by source name",
	Long: fmt.Sprintf(`%s

'get' gets a specific source, case SENSITIVE, from the vault and returns the
credentials from the name of the source. If the source name has a <Space> in it,
you have to surround the source name with double quotes.

If you want to see if a specific source is in the vault, you can use the:
	'gopass list -n <name-of-source>'
command.

Ex.
$ gopass get Google
Name: Google
	Username: <username>
	Password: <human-readable password>
	Notes: <will show if any notes are present>
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		GetCmdHandler(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

// GetCmdHandler is the handler function that encapsulates the GetItemsFromVault
// logic and runs some checks beforehand.
func GetCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(
			"Only 1 argument needed for the get command. See 'help' for correct usage",
		)
	}

	name := args[0]

	cfg, err := CheckConfig("")
	if err != nil {
		fmt.Println("Error checking config: ", err.Error())
		return err
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		fmt.Println("Cannot access, need to login")
		return fmt.Errorf("Cannot access, need to login")
	}

	err = GetItemFromVault(cfg, name)
	if err != nil {
		return fmt.Errorf("Cannot get %s from vault: %v", name, err)
	}

	return nil
}

// GetItemFromVault retreies the 'name' from the vault. If it doesn't exist, an
// error gets returned.
func GetItemFromVault(cfg model.Config, name string) error {
	f := utils.OpenVault(cfg.VaultName)
	defer f.Close()

	entries := crypt.DecryptVault(f)

	if len(entries) == 0 {
		fmt.Println("Nothing in your vault!")
		return fmt.Errorf("Nothing in vault.")
	}

	for _, e := range entries {
		if e.Name == name {
			// The \t's are for aligning the text in the terminal
			fmt.Println("From vault:")
			fmt.Println("Name: ", e.Name)
			fmt.Println("\tUsername: \t", e.Username)
			fmt.Println("\tPassword: \t", crypt.DecryptPassword(e.Password))

			if len(e.Notes) > 0 {
				fmt.Println("\tNotes: \t\t", e.Notes)
			}
			return nil
		}
	}

	fmt.Printf("'%s' not found.\n", name)

	encryptedCipherText, err := crypt.EncryptVault(entries)
	if err != nil {
		fmt.Println("Error with 'add' command")
		return fmt.Errorf("add::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
	return fmt.Errorf("'%s' not found in vault.\n", name)
}
