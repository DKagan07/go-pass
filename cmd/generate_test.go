package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	tests := []struct{
		name string
		length int
		expected int
	}{
		{
			name: "default length",
			length: 24,
			expected: 24,
		},
		{
			name: "different length",
			length: 30,
			expected: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := GeneratePassword(tt.length)
			assert.Len(t, b, tt.expected)
		})
	}
}
