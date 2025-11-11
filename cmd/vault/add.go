/*
Copyright © 2025 DKagan07
*/
package vault

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password to the vault",
	Long: `'add' adds a new password to the vault. The passwords are encrypted and
stored securely. 'add' takes a source, and then you are prompted to add a
username and password, and some notes. This notes section is for extra
information needed for any login. If multiple pieces of information are needed,
the info should be separated by semicolons, as pressing <Enter> will submit the
information.

NOTE: Entries are case sensitive in order to retreive. When you use the list
cmd, that is NOT case sensitive.
Ex.
	$ gopass vault add github
	Username: me@example.com
	Password: ********
	Notes: <any extra notes, can be empty>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := AddCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'add' command: ", err)
			return
		}
	},
}

// AddCmdHandler is the handler that orchestrates the 'add' command.
func AddCmdHandler(cmd *cobra.Command, args []string) error {
	// the value in GetString has to equal the flag that is created above
	if len(args) < 1 {
		return errors.New(
			"must have at least one argument",
		)
	}

	totalStr := strings.Join(args, " ")

	pBytes, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	keyring := model.NewMasterAESKeyManager(string(pBytes))

	cfg, err := utils.CheckConfig("", keyring)
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		return fmt.Errorf("cannot access, need to login")
	}

	userInput, err := GetInput(os.Stdin, os.Stdin, os.Stdin, keyring)
	if err != nil {
		return err
	}

	return AddToVault(totalStr, userInput, cfg, time.Now().UnixMilli(), keyring)
}

// GetInput is a function where we get input from the user, and return it in a
// model.UserInput.
func GetInput(us, pw, no io.Reader, key *model.MasterAESKeyManager) (model.UserInput, error) {
	ui := model.UserInput{}

	username, err := utils.GetInputFromUser(us, "Username")
	if err != nil {
		return ui, err
	}
	password, err := utils.GetPasswordFromUser(false, pw)
	if err != nil {
		return ui, err
	}
	notes, err := utils.GetInputFromUser(no, "Notes")
	if err != nil {
		return ui, err
	}

	pas, _ := crypt.EncryptPassword(password, key)

	ui.Username = username
	ui.Password = []byte(pas)
	ui.Notes = notes
	return ui, nil
}

// AddToVault holds the logic that adds encrypts the input from the user, and
// stores it into the vault.
func AddToVault(
	source string,
	ui model.UserInput,
	cfg model.Config,
	t int64,
	key *model.MasterAESKeyManager,
) error {
	ve := model.VaultEntry{
		Name:      source,
		Username:  ui.Username,
		Password:  ui.Password,
		Notes:     ui.Notes,
		UpdatedAt: t,
	}

	f, err := utils.OpenVault(cfg.VaultName)
	if err != nil {
		return fmt.Errorf("opening vault: %v", err)
	}
	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		log.Fatalf("add::stat: %v", err)
	}

	var entries []model.VaultEntry
	if fStat.Size() != 2 {
		entries = crypt.DecryptVault(f, key, false)
	}

	entries = append(entries, ve)

	encryptedCipherText, err := crypt.EncryptVault(entries, key)
	if err != nil {
		return fmt.Errorf("add::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
	return nil
}
