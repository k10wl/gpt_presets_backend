package token

import (
	"gpt_presets_backend/internal/models"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGenerateAndParseTOken(t *testing.T) {
	signature := "secret-key"
	os.Setenv("JWT_AUTH_SIGNATURE", signature)
	os.Setenv("JWT_REFRESH_SIGNATURE", signature)

	testUser := models.PublicUser{
		ID:   1,
		Name: "james",
	}

	exp := time.Now().Add(time.Hour)

	token, err := GenerateUserToken(testUser, exp, signature)

	if err != nil {
		t.Errorf("Failed to generate token: %v\n:", err)
	}

	parsedUser, err := ParseUserToken(token, signature)

	if err != nil {
		t.Errorf("Failed to parse token: %v\n:", err)
	}

	if !reflect.DeepEqual(testUser, parsedUser.PublicUser) {
		t.Errorf("Test user and parsed user are not equal\nTest user: %v\nParsed used: %v", testUser, parsedUser.PublicUser)
	}
}

func TestTokenExpiration(t *testing.T) {
	signature := "secret-key"
	testUser := models.PublicUser{
		ID:   1,
		Name: "james",
	}

	exp := time.Now().Add(-time.Hour)

	token, err := GenerateUserToken(testUser, exp, signature)

	if err != nil {
		t.Errorf("Failed to generate token: %v\n:", err)
	}

	_, err = ParseUserToken(token, signature)

	if err == nil {
		t.Errorf("Expected token to expire, but it got validated")
	}
}
