/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
	Short: "A brief description of your command",
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
		updateCmdFunc(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().BoolP("source", "s", false, "Update the source name")
	updateCmd.Flags().BoolP("username", "u", false, "Update the login username")
	updateCmd.Flags().BoolP("password", "p", false, "Update the password")
	updateCmd.Flags().BoolP("notes", "n", false, "Update the notes")
}

func updateCmdFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("update::not enough arguments to call 'update'. Please see help")
	}

	sn := args[0]

	sourceBool, _ := cmd.Flags().GetBool("source")
	usernameBool, _ := cmd.Flags().GetBool("username")
	passwordBool, _ := cmd.Flags().GetBool("password")
	notesBool, _ := cmd.Flags().GetBool("notes")

	if !sourceBool && !usernameBool && !passwordBool && !notesBool {
		fmt.Println("Need at least one flag. See help for more information")
		return
	}

	f := utils.OpenVault("")
	defer f.Close()

	var ve model.VaultEntry

	entries := crypt.DecryptVault(f)
	var idx int
	for i, e := range entries {
		if sn == e.Name {
			ve = e
			idx = i
			break
		}
	}

	// Not found
	if ve.Name == "" {
		fmt.Println("not found")
		return
	}

	var updatedSourceName string
	var updatedUsername string
	var updatedPassword []byte
	var updatedNotes string

	if sourceBool {
		updatedSourceName, _ = utils.GetInputFromUser(os.Stdin, "Source Name")
		ve.Name = updatedSourceName
	}
	if usernameBool {
		updatedUsername, _ = utils.GetInputFromUser(os.Stdin, "Username")
		ve.Username = updatedUsername
	}
	if passwordBool {
		updatedPassword, _ = utils.GetPasswordFromUser(false, os.Stdin)
		ve.Password = crypt.EncryptPassword(updatedPassword)
	}
	if notesBool {
		updatedNotes, _ = utils.GetInputFromUser(os.Stdin, "Notes")
		ve.Notes = updatedNotes
	}

	now := time.Now().UnixMilli()
	ve.UpdatedAt = now

	entries[idx] = ve

	encryptedCipherText, err := crypt.EncryptVault(entries)
	if err != nil {
		log.Fatalf("update::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
}
