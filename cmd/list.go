/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/utils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List lists all the sources of your login infos",
	Long: fmt.Sprintf(`%s

'list' lists all the sources of login info that's currently in your vault. A
source is, for example, 'Google', when it comes to what username and password
are stored with it.
Ex.
	$ gopass list
	Google
	Github
	...(etc.)
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		listCmdFunc(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// NOTE: Perhaps a flag here to look for a specific item in the list
	listCmd.Flags().StringP("name", "n", "", "Searches your list for the specific source")
}

func listCmdFunc(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		log.Fatal("list::no arguments needed for 'list'. See the help command!")
	}

	sourceName, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Fatalf("list::getting string from flag: %v", err)
	}

	cfgFile, ok, err := utils.OpenConfig()
	if ok && err == nil {
		fmt.Println("A file is not found. Need to init.")
		return
	}
	defer cfgFile.Close()
	cfg := crypt.DecryptConfig(cfgFile)

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		fmt.Println("Cannot access, need to login")
		return
	}

	f := utils.OpenVault(cfg.VaultName)
	defer f.Close()
	entries := crypt.DecryptVault(f)

	if len(entries) == 0 {
		fmt.Println("Nothing in your vault!")
		return
	}

	if sourceName == "" {
		fmt.Println("Entries:")
		for _, v := range entries {
			fmt.Printf("\t%s\n", v.Name)
		}
	} else {
		for _, v := range entries {
			if strings.EqualFold(v.Name, sourceName) {
				fmt.Printf("Yes, %s exists\n", v.Name)
				return
			}
		}
		fmt.Printf("%s does not exist", sourceName)
	}

	encryptedCipherText, err := crypt.EncryptVault(entries)
	if err != nil {
		log.Fatalf("list::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
}
