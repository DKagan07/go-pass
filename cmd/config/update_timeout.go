/*
Copyright © 2025 DKagan07
*/
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// updateTimeoutCmd represents the updateTimeout command
var UpdateTimeoutCmd = &cobra.Command{
	Use:   "update_timeout",
	Short: "Update the login timeout",
	Long: `'update_timeout' updates the timeout in which the need to login is
measured against. There are 2 flags, '--hours' and '--minutes', which have to be
present. 

Ex.
	$ gopass config update_timeout --hours 1 --minutes 45
)`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := UpdateTimeoutCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with update_timeout:", err)
			return
		}
	},
}

func UpdateTimeoutCmdHandler(cmd *cobra.Command, args []string) error {
	passB, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	keyring := model.NewMasterAESKeyManager(string(passB))

	cfg, err := utils.CheckConfig("", keyring)
	if err != nil {
		return err
	}

	if err = updateConfigTimeout(cfg, cmd); err != nil {
		return err
	}

	encryptedConfig, err := crypt.EncryptConfig(cfg, keyring)
	if err != nil {
		return err
	}

	cfgFile, err := os.OpenFile(utils.CONFIG_FILE, os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	utils.WriteToFile(cfgFile, encryptedConfig)

	fmt.Println("Success in updating the timeout time")
	return nil
}

func updateConfigTimeout(cfg *model.Config, cmd *cobra.Command) error {
	hrs, mins, err := getHoursAndMinutes(cmd)
	if err != nil {
		return err
	}

	totalTime := hrs + mins
	cfg.Timeout = totalTime
	return nil
}

func getHoursAndMinutes(cmd *cobra.Command) (hrsMilli int64, minsMilli int64, err error) {
	hours, err := cmd.Flags().GetInt("hours")
	if err != nil {
		return 0, 0, err
	}

	minutes, err := cmd.Flags().GetInt("minutes")
	if err != nil {
		return 0, 0, err
	}

	hrsMilli = time.Hour.Milliseconds() * int64(hours)
	minsMilli = time.Minute.Milliseconds() * int64(minutes)
	return hrsMilli, minsMilli, err
}
