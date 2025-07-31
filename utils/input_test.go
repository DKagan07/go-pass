package utils

import (
	"fmt"
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
	tests := []struct {
		name   string
		master bool
	}{
		{
			name:   "not a master password",
			master: false,
		},
		{
			name:   "master password",
			master: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader("test\n")

			user, err := GetPasswordFromUser(false, r)

			assert := assert.New(t)
			assert.Error(err)
			assert.Nil(user)
		})
	}
}

func TestConfirmPrompt(t *testing.T) {
	tests := []struct {
		name        string
		conf        ConfirmationPrompt
		prompt      string
		input       string
		want        bool
		shouldError bool
	}{
		{
			name:        "delete yes",
			conf:        DeletePrompt,
			prompt:      "test",
			input:       "y",
			want:        true,
			shouldError: false,
		},
		{
			name:        "delete no",
			conf:        DeletePrompt,
			prompt:      "test",
			input:       "n",
			want:        false,
			shouldError: false,
		},
		{
			name:        "delete bad input",
			conf:        DeletePrompt,
			prompt:      "test",
			input:       "bad",
			want:        false,
			shouldError: true,
		},
		{
			name:        "clean yes",
			conf:        CleanPrompt,
			prompt:      "test",
			input:       "y",
			want:        true,
			shouldError: false,
		},
		{
			name:        "clean no",
			conf:        CleanPrompt,
			prompt:      "test",
			input:       "n",
			want:        false,
			shouldError: false,
		},
		{
			name:        "clean bad input",
			conf:        CleanPrompt,
			prompt:      "test",
			input:       "bad",
			want:        false,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(fmt.Sprintf("%s\n", tt.input))

			got, err := ConfirmPrompt(tt.conf, tt.prompt, r)

			assert := assert.New(t)
			if tt.shouldError {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			assert.Equal(tt.want, got)
			fmt.Println()
		})
	}
}
