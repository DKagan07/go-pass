/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"go-pass/cmd"
	"go-pass/crypt"
)

const (
	SECRET_PASSWORD_KEY string = "SECRET_PASSWORD_KEY"
)

func main() {
	// We run this here because the key needs to be in place for any sort of
	// encryption. This function panics if th ekey isn't found in the env vars.
	crypt.GetAESKey()

	cmd.Execute()
}
