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

// LoginCmdHandler is the handler function that encapsulates the LoginUser
// logic.
func LoginCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		fmt.Println("No arguments needed for 'login'. See 'help' for more guidance")
		return fmt.Errorf("no arguments needed for 'login'. see 'help' for more guidance")
	}

	err := LoginUser("", os.Stdin)
	if err != nil {
		return fmt.Errorf("login user: %v", err)
	}
	return nil
}

// LoginUser is the function that logs the user in. It will check if the config
// file exists, and if it does, it will compare the password with the master
// password. If the password is correct, it will set the last visited time and
// return nil. If the password is incorrect, it will return an error.
func LoginUser(cfgName string, input io.Reader) error {
	cfgFile, ok, err := utils.OpenConfig(cfgName)
	if ok && err == nil {
		fmt.Println("A file is not found. Need to 'init'.")
		return errors.New("a file is not found. need to 'init'")
	}
	cfg := crypt.DecryptConfig(cfgFile)

	pass, err := utils.GetPasswordFromUser(true, input)
	if err != nil {
		fmt.Println("Error getting info from user")
		return fmt.Errorf("getting password from user: %v", err)
	}

	if err = bcrypt.CompareHashAndPassword(cfg.MasterPassword, pass); err != nil {
		fmt.Println("Login failed")
		return errors.New("login failed")
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
