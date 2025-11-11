package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/testutils"
)

// This is going to run into the same problem as the test for
// utils.GetPasswordFromUser() --> the golang.org/x/term package can only read
// from the terminal, not from a buffer like the std library can with fmt.Scanf
func TestLoginUser(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	bytesReader := bytes.NewReader(testutils.TEST_MASTER_PASSWORD)
	err = LoginUser(testutils.TEST_CONFIG_NAME, bytesReader, key, testutils.TEST_MASTER_PASSWORD)
	assert.Error(err)
}
