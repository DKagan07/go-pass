package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name           string
		field          string
		message        string
		expectedOutput string
	}{
		{
			name:           "name field error",
			field:          "Name",
			message:        "Name cannot be empty",
			expectedOutput: "Error with field Name: Name cannot be empty",
		},
		{
			name:           "username field error",
			field:          "Username",
			message:        "Username cannot be empty",
			expectedOutput: "Error with field Username: Username cannot be empty",
		},
		{
			name:           "password field error",
			field:          "Password",
			message:        "Password cannot be empty",
			expectedOutput: "Error with field Password: Password cannot be empty",
		},
		{
			name:           "custom field and message",
			field:          "Email",
			message:        "Invalid email format",
			expectedOutput: "Error with field Email: Invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			ve := &ValidationError{
				Field:   tt.field,
				Message: tt.message,
			}
			assert.Equal(tt.expectedOutput, ve.Error())
		})
	}
}
