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

func GetInputFromUser(r io.Reader, field string) (string, error) {
	br := bufio.NewReader(r)
	username, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}

	return cleanString(username), nil
}

func GetPasswordFromUser(master bool, r io.Reader) ([]byte, error) {
	phrase := "Password: "
	if master {
		phrase = "Master Password: "
	}

	fmt.Print(phrase)
	fd, ok := (r).(*os.File)
	if !ok {
		return nil, errors.New("cannot read from source")
	}
	b, err := term.ReadPassword(int(fd.Fd()))
	fmt.Println()
	if len(b) == 0 {
		return nil, errors.New("Must enter a password")
	}
	return b, err
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\n")
	return s
}
