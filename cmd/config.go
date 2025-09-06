/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"go-pass/cmd/config"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Parent command for the config sub commands",
	Long: fmt.Sprintf(`%s

'config' is acting as a parent command for config-related commands.
Please use the 'help' flag for additional information
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(config.ChangeMasterpassCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
