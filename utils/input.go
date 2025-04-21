package utils

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

func GetInputFromUser(field string) (string, error) {
	br := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", field)
	return br.ReadString('\n')
}

func GetPasswordFromUser() ([]byte, error) {
	fmt.Print("Password: ")
	return term.ReadPassword(int(os.Stdin.Fd()))
}
