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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	f := utils.OpenVault()
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
		updatedPassword, _ = utils.GetPasswordFromUser(os.Stdin)
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

	utils.WriteToVault(f, encryptedCipherText)
}
