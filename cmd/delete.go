/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"slices"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/utils"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: fmt.Sprintf(`%s

'delete' deletes a specific source name from your vault. This HAS to be case
sensitive.
Ex.
	$ gopass delete google
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		deleteCmdFunc(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteCmdFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("delete::incorrect number of args. See help for more information!")
	}

	name := args[0]

	f := utils.OpenVault("")
	defer f.Close()

	entries := crypt.DecryptVault(f)

	if len(entries) == 0 {
		fmt.Println("Nothing in your vault!")
		return
	}

	found := false
	for i, v := range entries {
		if name == v.Name {
			entries = slices.Delete(entries, i, i+1)
			fmt.Printf("Deleted %s from your vault\n", name)
			found = true
		}
	}

	if !found {
		fmt.Printf("%s not found\n", name)
	}

	b, err := crypt.EncryptVault(entries)
	if err != nil {
		log.Fatalf("delete::failed to encrypt: %v", err)
	}

	utils.WriteToFile(f, b)
}
