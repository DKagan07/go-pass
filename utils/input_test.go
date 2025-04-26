package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInputFromUser(t *testing.T) {
	r := strings.NewReader("test\n")

	user, err := GetInputFromUser(r, "Username")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test", user)
}

// Because of the use of github.com/x/term and trying to avoid complexity for
// a fairly simple function, we are just seeing if there's an error. The 'term'
// package relies specifically on the interface with the terminal, and mocking
// that would make this overly complex for what this function actually does.
func TestGetPasswordFromUser(t *testing.T) {
	r := strings.NewReader("test\n")

	user, err := GetPasswordFromUser(r)

	assert := assert.New(t)
	assert.Error(err)
	assert.Nil(user)
}
