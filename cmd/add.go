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
	"golang.org/x/term"

	"go-pass/crypt"
	"go-pass/model"
)

// TODO: For 'add', I need to make sure to add the password securely as well

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password to the vault",
	Long: fmt.Sprintf(`%s

'add' adds a new password to the vault. The passwords are encrypted and
stored securely. 'add' takes a source, and then you are prompted to add a
username and password.
Ex.
	$ gopass add github
	Username: me@example.com
	Password: ********
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
}

func addCmdFunc(cmd *cobra.Command, args []string) {
	// the value in GetString has to equal the flag that is created above
	fmt.Println("add called")
	if len(args) != 1 {
		log.Fatal("Not enough arguments to call 'add'. Please see help")
	}
	fmt.Print("Username:")
	var username string
	_, err := fmt.Scan(&username)
	if err != nil {
		log.Fatalf("Error getting username: %v", err)
	}
	fmt.Printf("Input password:")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed reading pword: %v", err)
	}
	hashedPw := crypt.HashPassword(string(passwordBytes))
	if err != nil {
		log.Fatalf("Error handling username or password: %v\n", err)
	}

	fmt.Println("hashedPw: ", hashedPw)

	now := time.Now()
	ve := model.VaultEntry{
		Name:      args[0],
		Username:  username,
		Password:  hashedPw,
		CreatedAt: now.UnixMilli(),
	}

	// crypt.Encrypt(ve)

	fmt.Printf("vault entry: %+v\n", ve)
}
