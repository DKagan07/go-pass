/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an entry in your vault with specific flags",
	Long: fmt.Sprintf(`%s

'update' updates a current entry in your vault. The command takes in the name
of your entry. To update the entry, at least 1 flag is required. There are 4
flags, each of them to update part of the entry. Minimum of 1, but can have
multiple if multiple fields needs updating.

If you want to update the source name with a name with a <Space>, be careful
that if you want to 'get' this source name, you need to add double quotes around
the name. Ex: gopass get "blah1 blah2" -> is 1 source name with a space.

See help for all the flags available.

Ex.
	$ gopass update github -u -s
	Source name: <updated name for entry>
	Username: <update username>
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		if err := UpdateCmdHandler(cmd, args); err != nil {
			fmt.Printf("Error with 'update' command: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolP("source", "s", false, "Update the source name")
	updateCmd.Flags().BoolP("username", "u", false, "Update the login username")
	updateCmd.Flags().BoolP("password", "p", false, "Update the password")
	updateCmd.Flags().BoolP("notes", "n", false, "Update the notes")
}

type Inputs struct {
	Source   bool
	Username bool
	Password bool
	Notes    bool
}

type InputSources struct {
	Source   io.Reader
	Username io.Reader
	Password io.Reader
	Notes    io.Reader
}

// UpdateCmdHandler is the handler function that encapsulates the update logic
func UpdateCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong number of arguments, need 1. please see 'help'")
	}

	i, err := UpdateFlags(cmd)
	if err != nil {
		return err
	}

	sn := args[0]

	cfgFile, ok, err := utils.OpenConfig("")
	if ok && err == nil {
		return errors.New("need to init")
	}
	defer cfgFile.Close()
	cfg := crypt.DecryptConfig(cfgFile)

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		return errors.New("need to login")
	}

	err = UpdateEntry(i, cfg, sn, InputSources{os.Stdin, os.Stdin, os.Stdin, os.Stdin})
	if err != nil {
		return fmt.Errorf("error updating entry: %v", err)
	}

	return nil
}

// UpdateFlags consolidates the different flags that may or may not be present.
// It returns a type Input and error
func UpdateFlags(cmd *cobra.Command) (Inputs, error) {
	sourceBool, err := cmd.Flags().GetBool("source")
	if err != nil {
		return Inputs{}, err
	}

	usernameBool, err := cmd.Flags().GetBool("username")
	if err != nil {
		return Inputs{}, err
	}

	passwordBool, err := cmd.Flags().GetBool("password")
	if err != nil {
		return Inputs{}, err
	}

	notesBool, err := cmd.Flags().GetBool("notes")
	if err != nil {
		return Inputs{}, err
	}

	if !sourceBool && !usernameBool && !passwordBool && !notesBool {
		fmt.Println("Need at least one flag. See help for more information")
		return Inputs{}, errors.New("need at last 1 flag")
	}

	return Inputs{
		Source:   sourceBool,
		Username: usernameBool,
		Password: passwordBool,
		Notes:    notesBool,
	}, nil
}

// UpdateEntry contains the logic of actually updating the vault entry and
// storing it in the vault.
func UpdateEntry(inputs Inputs, cfg model.Config, sourceName string, is InputSources) error {
	f, err := utils.OpenVault(cfg.VaultName)
	if err != nil {
		return fmt.Errorf("opening vault: %v", err)
	}
	defer f.Close()
	entries := crypt.DecryptVault(f)

	var ve model.VaultEntry

	var idx int
	for i, e := range entries {
		if sourceName == e.Name {
			ve = e
			idx = i
			break
		}
	}

	// Not found
	if ve.Name == "" {
		return fmt.Errorf("'%s' not found", sourceName)
	}

	ve, err = UpdateVaultEntry(ve, inputs, is)
	if err != nil {
		return err
	}

	entries[idx] = ve

	encryptedCipherText, err := crypt.EncryptVault(entries)
	if err != nil {
		log.Fatalf("update::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
	return nil
}

// UpdateVaultEntry takes in the user input and, depending on the flags, update
// the vault entry accordingly. It returns the updated model.VaultEntry
func UpdateVaultEntry(
	ve model.VaultEntry,
	inputs Inputs,
	updateSources InputSources,
) (model.VaultEntry, error) {
	var updatedSourceName string
	var updatedUsername string
	var updatedPassword []byte
	var updatedNotes string
	var err error

	if inputs.Source {
		updatedSourceName, err = utils.GetInputFromUser(updateSources.Source, "Source Name")
		if err != nil {
			return model.VaultEntry{}, err
		}
		ve.Name = updatedSourceName
	}
	if inputs.Username {
		updatedUsername, err = utils.GetInputFromUser(updateSources.Username, "Username")
		if err != nil {
			return model.VaultEntry{}, err
		}
		ve.Username = updatedUsername
	}
	if inputs.Password {
		updatedPassword, err = utils.GetPasswordFromUser(false, updateSources.Password)
		if err != nil {
			return model.VaultEntry{}, err
		}
		ve.Password = crypt.EncryptPassword(updatedPassword)
	}
	if inputs.Notes {
		updatedNotes, err = utils.GetInputFromUser(updateSources.Notes, "Notes")
		if err != nil {
			return model.VaultEntry{}, err
		}
		ve.Notes = updatedNotes
	}

	now := time.Now().UnixMilli()
	ve.UpdatedAt = now

	return ve, nil
}
