/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"go-pass/cmd"
	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

const (
	SECRET_PASSWORD_KEY string = "SECRET_PASSWORD_KEY"
)

func main() {
	initProgram()
	cmd.Execute()
}

func initProgram() {
	f := utils.CreateVault()

	fileStat, err := f.Stat()
	if err != nil {
		panic("init::getting stat on file")
	}

	if fileStat.Size() == 0 {
		ve := []model.VaultEntry{}
		b, err := crypt.EncryptVault(ve)
		if err != nil {
			panic("init::encrypt ve")
		}
		utils.WriteToVault(f, b)
	}

	f.Close()

	// Ensure that the key exists on startup
	utils.GetAESKey()
}
