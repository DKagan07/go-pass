/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: fmt.Sprintf(`%s

'get' gets a specific source, case SENSITIVE, from the vault and returns the
credentials from the name of the source.

If you want to see if a specific source is in the vault, you can use the:
	'gopass list -n <name-of-source>'
command.

Ex.
$ gopass get Google
Name: Google
	Username: <username>
	Password: <human-readable password>
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		getCmdFunc(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getCmdFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("get::no arguments needed for 'list'. See the help command!")
	}

	f := utils.OpenVault()
	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		log.Fatalf("get::getting stat")
	}

	var entries []model.VaultEntry
	if fStat.Size() != 2 {
		entries = crypt.DecryptVault(f)
	}

	if len(entries) == 0 {
		fmt.Println("Nothing in your vault!")
		return
	}

	for _, e := range entries {
		if e.Name == args[0] {
			// decode password
			crypt.DecryptPassword(e.Password)
			fmt.Println("From vault:")
			fmt.Println("Name: ", e.Name)
			fmt.Println("\tUsername: ", e.Username)
			fmt.Println("\tPassword: ", crypt.DecryptPassword(e.Password))
			return
		}
	}
	fmt.Printf("%s not found in vault.\n", args[0])
}
