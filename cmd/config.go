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
Please use the 'help' flag for additional information.
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(config.ChangeMasterpassCmd)
	configCmd.AddCommand(config.UpdateTimeoutCmd)
	configCmd.AddCommand(config.ViewCmd)

	initConfigFlags()
}

func initConfigFlags() {
	config.UpdateTimeoutCmd.Flags().IntP("hours", "q", 0, "the hours you want to add to timeout")
	config.UpdateTimeoutCmd.Flags().
		IntP("minutes", "m", 30, "the minutes you want to add to timeout")
}
