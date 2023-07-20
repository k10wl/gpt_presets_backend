package password

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	type TestTable struct {
		name  string
		input string
		valid bool
	}

	var testCases = []TestTable{
		{
			name:  "return nil for valid password",
			input: "Qwerty123",
			valid: true,
		},
		{
			name:  "return error if length is less then 8",
			input: "1234567",
			valid: false,
		},
		{
			name:  "return error if has no characters",
			input: "12345678",
			valid: false,
		},
		{
			name:  "return error if has no symbols",
			input: "qwerty-qwerty",
			valid: false,
		},
		{
			name:  "return error if has no uppercase",
			input: "123qwerty",
			valid: false,
		},
	}

	for _, tt := range testCases {
		func(tt TestTable) {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				err := ValidatePassword(tt.input)

				if tt.valid && err != nil {
					t.Errorf("Expected valid password, got error: %v", err)
				}

				if !tt.valid && err == nil {
					t.Errorf("Expected error in invalid password")
				}
			})
		}(tt)
	}
}

func TestEncryption(t *testing.T) {
	pass := "password"

	hash, err := HashPassword(pass)

	if err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
	}

	if hash == pass {
		t.Errorf("Expected hash to be different from password, got same")
	}

	if err := ComparePassword(hash, pass); err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}
}
