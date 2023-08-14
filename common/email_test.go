package common

import (
	"testing"
)

func TestCheckEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"a@elastichosts.com", false}, // Non-existent domain
		{"a@bachelier.name", false},   // Non-existent domain
		{"invalid-email", false},      // Invalid syntax
		{"a @gmail.com", false},       // Invalid syntax
		{"a@gmail.com", true},         // Valid
		{"user@example.com", true},    // Example domain, syntax valid (MX lookup might fail)
	}

	for _, test := range tests {
		got := checkEmail(test.email)
		if got != test.want {
			t.Errorf("checkEmail(%q) = %v; want %v", test.email, got, test.want)
		}
	}
}
