package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This is going to run into the same problem as the test for
// utils.GetPasswordFromUser() --> the golang.org/x/term package can only read
// from the terminal, not from a buffer like the std library can with fmt.Scanf
func TestLoginUser(t *testing.T) {
	bytesReader := bytes.NewReader(TEST_MASTER_PASSWORD)
	err := LoginUser(TEST_CONFIG_NAME, bytesReader)
	assert.Error(t, err)
}
