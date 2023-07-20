package password

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (r string, err error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return password, err
	}

	return string(byte), err
}

func ComparePassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func ValidatePassword(password string) error {
	var (
		hasMinEight  = regexp.MustCompile(`.{8,}`)
		hasUppercase = regexp.MustCompile(`[A-Z]`)
		hasNumber    = regexp.MustCompile(`[0-9]`)
	)

	if !hasMinEight.MatchString(password) {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if !hasUppercase.MatchString(password) {
		return fmt.Errorf("password must include at least one uppercase letter")
	}
	if !hasNumber.MatchString(password) {
		return fmt.Errorf("password must include at least one digit")
	}

	return nil
}
