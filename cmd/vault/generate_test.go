package vault

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	defaultSpecialChars := "!@#$%^&*"
	tests := []struct {
		name     string
		length   int
		expected int
		special  string
	}{
		{
			name:     "default length",
			length:   24,
			expected: 24,
			special:  "",
		},
		{
			name:     "different length",
			length:   30,
			expected: 30,
			special:  "",
		},
		{
			name: "special characters",
			// Make the length arbitrarily large to ensure that the characters
			// are in the password
			length:   40,
			expected: 40,
			special:  "()<>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specialChars := defaultSpecialChars
			if tt.special != "" {
				specialChars = tt.special
			}

			b := GeneratePassword(tt.length, specialChars)

			doesContain := false
			for _, letter := range strings.Split(specialChars, "") {
				if strings.Contains(string(b), letter) {
					doesContain = true
					break
				}
			}

			time.Sleep(time.Millisecond * 30)
			assert.Len(t, b, tt.expected)
			assert.True(t, doesContain)
		})
	}
}
