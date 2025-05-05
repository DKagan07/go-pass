/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
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
	"golang.org/x/crypto/bcrypt"

	"go-pass/crypt"
	"go-pass/utils"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the app",
	Long: fmt.Sprintf(`%s

'login' logs the user in for 30 minutes. Running 'init' for the first time also
counts as an initial login. To login, you will need your Master Password that
you set when the 'init' command was ran.
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		LoginCmdHandler(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func LoginCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		fmt.Println("No arguments needed for 'login'. See 'help' for more guidance")
		return fmt.Errorf("No arguments needed for 'login'. See 'help' for more guidance")
	}

	err := LoginUser("", os.Stdin)
	if err != nil {
		return fmt.Errorf("Login user: %v", err)
	}
	return nil
}

func LoginUser(cfgName string, input io.Reader) error {
	cfgFile, ok, err := utils.OpenConfig(cfgName)
	if ok && err == nil {
		fmt.Println("A file is not found. Need to 'init'.")
		return errors.New("A file is not found. Need to 'init'")
	}
	cfg := crypt.DecryptConfig(cfgFile)

	pass, err := utils.GetPasswordFromUser(true, input)
	if err != nil {
		fmt.Println("Error getting info from user")
		return fmt.Errorf("Getting password from user: %v", err)
	}

	if err = bcrypt.CompareHashAndPassword(cfg.MasterPassword, pass); err != nil {
		fmt.Println("Login failed")
		return errors.New("Login failed")
	}

	fmt.Println("Success!")

	now := time.Now().UnixMilli()
	cfg.LastVisited = now

	cipherText, err := crypt.EncryptConfig(cfg)
	if err != nil {
		log.Fatalf("login::failed creating ciphertext: %v", err)
	}

	utils.WriteToFile(cfgFile, cipherText)

	return nil
}
