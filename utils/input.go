package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func GetInputFromUser(field string) (string, error) {
	br := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", field)
	username, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}

	return cleanString(username), nil
}

func GetPasswordFromUser() ([]byte, error) {
	fmt.Print("Password: ")
	return term.ReadPassword(int(os.Stdin.Fd()))
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\n")
	return s
}
