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
	fmt.Printf("%s: ", field)
	username, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}

	return cleanString(username), nil
}

func GetPasswordFromUser(r io.Reader) ([]byte, error) {
	fmt.Print("Password: ")
	fd, ok := (r).(*os.File)
	if !ok {
		return nil, errors.New("cannot read from source")
	}
	return term.ReadPassword(int(fd.Fd()))
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\n")
	return s
}
