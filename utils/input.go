package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type ConfirmationPrompt string

const (
	DeletePrompt ConfirmationPrompt = "DELETE"
	CleanPrompt  ConfirmationPrompt = "CLEAN"
)

func (c ConfirmationPrompt) String() string {
	return string(c)
}

// GetInputFromUser reads the input from the user and returns the string form
// of that input
func GetInputFromUser(r io.Reader, field string) (string, error) {
	br := bufio.NewReader(r)
	fmt.Printf("%s: ", field)
	username, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}

	return cleanString(username), nil
}

// GetPasswordFromUser reads the password using the 'term' package from r.
// 'again' is an optional parameter specifically used by change_masterpass to
// change the prompt.
func GetPasswordFromUser(master bool, r io.Reader, again ...bool) ([]byte, error) {
	phrase := "Password: "
	if master {
		phrase = "Master Password: "

		if len(again) > 0 && again[0] {
			phrase = "New Master Password again: "
		}
	}

	fmt.Print(phrase)
	fd, ok := (r).(*os.File)
	if !ok {
		return nil, errors.New("cannot read from source")
	}
	b, err := term.ReadPassword(int(fd.Fd()))
	fmt.Println()
	if len(b) == 0 {
		return nil, errors.New("must enter a password")
	}
	return b, err
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\n")
	return s
}

func ConfirmPrompt(confType ConfirmationPrompt, prompt string, r io.Reader) (bool, error) {
	switch confType.String() {
	case "DELETE":
		br := bufio.NewReader(r)
		response := fmt.Sprintf("Are you sure you want to delete '%s'? (y/n) ", prompt)
		fmt.Print(response)
		confirm, err := br.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("failed to read input: %v", err)
		}

		confirm = cleanString(confirm)
		if strings.EqualFold(confirm, "y") {
			return true, nil
		} else if strings.EqualFold(confirm, "n") {
			return false, nil
		} else {
			return false, fmt.Errorf("invalid input")
		}
	case "CLEAN":
		response := "Are you sure you want to delete everything? This includes your config and vault? (y/n) "

		br := bufio.NewReader(r)
		fmt.Print(response)

		confirm, err := br.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("failed to read input: %v", err)
		}

		confirm = cleanString(confirm)
		if strings.EqualFold(confirm, "y") {
			return true, nil
		} else if strings.EqualFold(confirm, "n") {
			return false, nil
		} else {
			return false, fmt.Errorf("invalid input")
		}
	}

	return false, nil
}
