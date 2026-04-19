package unit

import (
	"regexp"
	"testing"

	"url-shortener/internal/app"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortCode(t *testing.T) {
	code := app.GenerateShortCode(6)
	assert.Len(t, code, 6)
	// Should contain only alphanumeric characters
	match, _ := regexp.MatchString("^[A-Za-z0-9]+$", code)
	assert.True(t, match, "short code should be alphanumeric")

	code2 := app.GenerateShortCode(8)
	assert.Len(t, code2, 8)

	// Ensure randomness (very low probability of collision)
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		c := app.GenerateShortCode(6)
		assert.False(t, codes[c], "duplicate short code generated")
		codes[c] = true
	}
}

func TestCustomCodeValidation(t *testing.T) {
	// Valid codes
	valid := []string{
		"hello",
		"HELLO",
		"123abc",
		"test_123",
		"test-code",
		"a",
		"a_b-c",
	}
	// Invalid codes
	invalid := []string{
		"",
		"test@code",
		"test code",
		"test.code",
		"test?",
		"test#",
		"test/",
		"test\\",
		"test|",
		"test!",
	}

	for _, code := range valid {
		assert.True(t, isValidCustomCode(code), "expected valid: %s", code)
	}
	for _, code := range invalid {
		assert.False(t, isValidCustomCode(code), "expected invalid: %s", code)
	}
}

// Helper function copied from app.go logic
func isValidCustomCode(code string) bool {
	if code == "" {
		return false
	}
	for _, ch := range code {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-') {
			return false
		}
	}
	return true
}
