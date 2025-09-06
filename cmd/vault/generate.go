/*
Copyright © 2025 DKagan07
*/
package vault

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var DefaultChars = "!@#$%^&*"

// generateCmd represents the generate command
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate generates a secure password for you",
	Long: `'generate' is a helper command that helps generate a strong password for you!
It will print it out to the terminal, and is then copy-pastable. The default for
special characters are '%s'. You can adjust this in any way you like by
using the -c flag.

***
Note: There is always a chance that this generator doesn't return out a password
that satisfies a password input because these are created with a
cryptographically secure RNG. Please modify and change that if needed.
(If using the -a flag, run 'gopass update <source_name> to update the password')
***
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := GenerateCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'generate' command: ", err)
			return
		}
	},
}

// GenerateCmdHandler is the handler function that encapsulates the GeneratePassword
// and if the flag is provided, it will prompt the user for the required information
// and add it to the vault.
func GenerateCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("no arguments needed for 'generate'. see 'help' for more guidance")
	}

	cfg, err := utils.CheckConfig("")
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		return fmt.Errorf("cannot access, need to login")
	}

	length, err := cmd.Flags().GetInt("length")
	if err != nil {
		return fmt.Errorf("getting length flag: %v", err)
	}

	special, err := cmd.Flags().GetString("specialChars")
	if err != nil {
		return fmt.Errorf("getting specialChar flag: %v", err)
	}

	strongPasswordBytes := GeneratePassword(length, special)
	fmt.Println("Generated Password: ", string(strongPasswordBytes))

	source, err := cmd.Flags().GetString("add")
	if err != nil {
		return fmt.Errorf("getting add flag: %v", err)
	}

	if source != "" {
		return AddGeneratedPasswordToVault(source, strongPasswordBytes, cfg, now)
	}

	return nil
}

// AddGeneratedPasswordToVault contains the logic of getting the information
// from the user and storing the information in the vault. The source is
// obtained from the '-a' flag, and the password is generated from the
// 'GeneratePassword' function.
func AddGeneratedPasswordToVault(source string, password []byte, cfg model.Config, t int64) error {
	userInput := model.UserInput{
		Password: crypt.EncryptPassword(password),
	}

	username, err := utils.GetInputFromUser(os.Stdin, "Username")
	if err != nil {
		return err
	}
	notes, err := utils.GetInputFromUser(os.Stdin, "Notes")
	if err != nil {
		return err
	}

	userInput.Username = username
	userInput.Notes = notes
	return AddToVault(source, userInput, cfg, t)
}

// GeneratePassword generates a strong password of default length 24 consisting
// of lower case, uppercase, numbers, and special characters using the
// crypto/rand package for cryptographically secure RNG.
//
// Note: There is always a chance that these passwords will not satisfy password
// inputs, so double check that it does.
func GeneratePassword(l int, special string) []byte {
	baseByteSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	byteSet := baseByteSet + special
	setLength := big.NewInt(int64(len(byteSet)))

	b := make([]byte, l)
	for i := range b {
		idx, err := rand.Int(rand.Reader, setLength)
		if err != nil {
			log.Fatalf("failed to get random number: %v", err)
		}

		b[i] = byteSet[idx.Int64()]
	}
	return b
}
