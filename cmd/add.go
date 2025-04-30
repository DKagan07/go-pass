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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password to the vault",
	Long: fmt.Sprintf(`%s

'add' adds a new password to the vault. The passwords are encrypted and
stored securely. 'add' takes a source, and then you are prompted to add a
username and password, and some notes. This notes section is for extra
information needed for any login. If multiple pieces of information are needed,
the info should be separated by semicolons, as pressing <Enter> will submit the
information.

NOTE: Entries are case sensitive in order to retreive. When you use the list
cmd, that is NOT case sensitive.
Ex.
	$ gopass add github
	Username: me@example.com
	Password: ********
	Notes: <any extra notes, can be empty>
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		addCmdFunc(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// addCmd.Flags().StringP("something", "d", "", "Add something")

	// NOTE: perhaps add a flag here to show the password and not hide it?
}

func addCmdFunc(cmd *cobra.Command, args []string) {
	// the value in GetString has to equal the flag that is created above
	if len(args) != 1 {
		log.Fatal("add::not enough arguments to call 'add'. Please see help")
	}

	username, err := utils.GetInputFromUser(os.Stdin, "Username")
	if err != nil {
		log.Fatalf("add::not a valid input: %v", err)
	}

	passwordBytes, err := utils.GetPasswordFromUser(false, os.Stdin)
	if err != nil {
		log.Fatalf("add::failed reading pword: %v", err)
	}
	fmt.Println()

	hashedPw := crypt.EncryptPassword(passwordBytes)
	if err != nil {
		log.Fatalf("add::error handling username or password: %v\n", err)
	}

	notes, err := utils.GetInputFromUser(os.Stdin, "Notes")
	if err != nil {
		log.Fatalf("add::not valid input for notes: %v", err)
	}

	now := time.Now()
	ve := model.VaultEntry{
		Name:      args[0],
		Username:  username,
		Password:  hashedPw,
		Notes:     notes,
		UpdatedAt: now.UnixMilli(),
	}

	cfgFile := utils.OpenConfig()
	defer cfgFile.Close()

	cfg := crypt.DecryptConfig(cfgFile)
	f := utils.OpenVault(cfg.VaultName)
	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		log.Fatalf("add::stat: %v", err)
	}

	var entries []model.VaultEntry
	if fStat.Size() != 2 {
		entries = crypt.DecryptVault(f)
	}

	entries = append(entries, ve)

	encryptedCipherText, err := crypt.EncryptVault(entries)
	if err != nil {
		log.Fatalf("add::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
}
