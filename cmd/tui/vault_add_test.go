package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/crypt"
)

func TestValidateAddInput(t *testing.T) {
	tests := []struct {
		name        string
		vaultName   string
		username    string
		password    string
		expectError bool
		errorField  string
	}{
		{
			name:        "valid input",
			vaultName:   "GitHub",
			username:    "user",
			password:    "pass",
			expectError: false,
		},
		{
			name:        "empty name",
			vaultName:   "",
			username:    "user",
			password:    "pass",
			expectError: true,
			errorField:  "Name",
		},
		{
			name:        "empty username",
			vaultName:   "GitHub",
			username:    "",
			password:    "pass",
			expectError: true,
			errorField:  "Username",
		},
		{
			name:        "empty password",
			vaultName:   "GitHub",
			username:    "user",
			password:    "",
			expectError: true,
			errorField:  "Password",
		},
	}

	for _, tt := range tests {
		assert := assert.New(t)
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddInput(tt.vaultName, tt.username, tt.password)

			if tt.expectError {
				assert.Error(err)
				validationErr, ok := err.(*ValidationError)
				assert.True(ok, "expected ValidationError")
				assert.Equal(tt.errorField, validationErr.Field)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestAddToVault(t *testing.T) {
	app, cleanup := NewTestApp(t)
	defer cleanup()

	assert := assert.New(t)

	app.AddToVault("TestName", "", "test_username", "test_password")

	assert.Equal(1, len(app.Vault))
	assert.Equal("TestName", app.Vault[0].Name)
	assert.Equal("test_username", app.Vault[0].Username)
	assert.Equal("", app.Vault[0].Notes)

	encryptedPass := app.Vault[0].Password
	testPass := crypt.DecryptPassword(encryptedPass, app.Keyring, false)
	assert.Equal("test_password", testPass)
}
