/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"crypto/rand"
	"fmt"
	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate generates a secure password for you",
	Long: fmt.Sprintf(`%s

'generate' is a helper command that helps generate a strong password for you!
It will print it out to the terminal, and is then copy-pastable.

***
Note: There is always a chance that this generator doesn't return out a password
that satisfies a password input because these are created with a
cryptographically secure RNG. Please modify and change that if needed.
(If using the -a flag, run 'gopass update <source_name> to update the password')
***
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		GenerateCmdHandler(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntP("length", "l", 24, "Decides length of new password")
	generateCmd.Flags().StringP("add", "a", "", "Add a newly generated password to your vault")
}

func GenerateCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		fmt.Println("No arguments needed for 'list'. See 'help' for more guidance")
		return fmt.Errorf("No arguments needed for 'list'. See 'help' for more guidance")
	}

	cfg, err := CheckConfig("")
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		fmt.Println("Cannot access, need to login")
		return fmt.Errorf("Cannot access, need to login")
	}

	length, err := cmd.Flags().GetInt("length")
	if err != nil {
		fmt.Println("getting length flag")
		return err
	}

	strongPasswordBytes := GeneratePassword(length)
	fmt.Println("Generated Password: ", string(strongPasswordBytes))

	source, err := cmd.Flags().GetString("add")
	if err != nil {
		fmt.Println("getting add flag")
		return err
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
		return  err
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
func GeneratePassword(l int) []byte {
	byteSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+=_-[]{}|;<>:~`0123456789"
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

