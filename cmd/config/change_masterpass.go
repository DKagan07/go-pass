/*
Copyright © 2025 DKagan07
*/
package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// changeMasterpassCmd represents the changeMasterpass command
var ChangeMasterpassCmd = &cobra.Command{
	Use:   "change_masterpass",
	Short: "Changes your master password, used for logging in",
	Long: `'change_masterpass' changes the master password, used to login

Ex.
	$gopass config change_masterpass
	Master Password: <master_pass>
	Master Password: <new_master_pass>
	Input Master Password again: <new_master_pass>
	Success! Master Password changed.
`,

	Run: func(cmd *cobra.Command, args []string) {
		if err := ChangeMasterpassCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'change_masterpass' command: ", err)
			return
		}
	},
}

// ChangeMasterpassCmdHandler handles the 'change_masterpass' command
func ChangeMasterpassCmdHandler(cmd *cobra.Command, args []string) error {
	passB, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	keyring := model.NewMasterAESKeyManager(string(passB))

	cfg, err := utils.CheckConfig("", keyring)
	if err != nil {
		return err
	}

	err = ChangeMasterpass(cfg, keyring)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}

// ChangeMasterpass handles the business logic of changing the master password
// This will return an error, otherwise it'll print a success message to the
// terminal indicating success.
func ChangeMasterpass(cfg *model.Config, key *model.MasterAESKeyManager) error {
	fmt.Println(strings.Repeat("*", 24))
	fmt.Println("Input current password:")
	fmt.Println(strings.Repeat("*", 24))
	password, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(cfg.MasterPassword, password); err != nil {
		fmt.Println("passwords don't match")
		return errors.New("passwords don't match")
	}

	fmt.Println("Passwords match!")

	fmt.Println()
	fmt.Println(strings.Repeat("*", 20))
	fmt.Println("New Master Password")
	fmt.Println(strings.Repeat("*", 20))

	newPass, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	confirmedNewPass, err := utils.GetPasswordFromUser(true, os.Stdin, true)
	if err != nil {
		return err
	}

	if !bytes.Equal(newPass, confirmedNewPass) {
		return errors.New("passwords do not match")
	}

	bNewPass, err := bcrypt.GenerateFromPassword(newPass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	cfg.MasterPassword = bNewPass
	cfgB, err := crypt.EncryptConfig(cfg, key)
	if err != nil {
		return err
	}

	cfgFile, err := os.OpenFile(utils.CONFIG_FILE, os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	if err := utils.WriteToFile(cfgFile.Name(), model.FileConfig, cfgB); err != nil {
		return err
	}

	fmt.Println("Success! Master Password changed.")
	return nil
}
