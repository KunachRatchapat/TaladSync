package handler

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanText(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expectText  string
		expectOK    bool
	}{
		{
			description: "normal text",
			input:       "ร้านป้า",
			expectText:  "ร้านป้า",
			expectOK:    true,
		},
		{
			description: "trims surrounding spaces",
			input:       "  ร้านป้า  ",
			expectText:  "ร้านป้า",
			expectOK:    true,
		},
		{
			description: "empty",
			input:       "",
			expectText:  "",
			expectOK:    false,
		},
		{
			description: "spaces only",
			input:       "   ",
			expectText:  "",
			expectOK:    false,
		},
		{
			description: "exactly 100 Thai characters",
			input:       strings.Repeat("ก", 100),
			expectText:  strings.Repeat("ก", 100),
			expectOK:    true,
		},
		{
			description: "101 Thai characters",
			input:       strings.Repeat("ก", 101),
			expectText:  "",
			expectOK:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			text, ok := cleanText(test.input, maxTextLen)
			require.Equal(t, test.expectOK, ok)
			require.Equal(t, test.expectText, text)
		})
	}
}
