package token

import (
	"os"
	"time"

	"gpt_presets_backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type RegisteredClaims = jwt.RegisteredClaims

type JWTPublicUserClaim struct {
	RegisteredClaims
	models.PublicUser
}

var (
	AUTH_SIGNATURE    = os.Getenv("JWT_AUTH_SIGNATURE")
	REFRESH_SIGNATURE = os.Getenv("JWT_REFRESH_SIGNATURE")
)

func GenerateUserToken(payload models.PublicUser, exp time.Time, signature string) (string, error) {
	claims := JWTPublicUserClaim{
		RegisteredClaims: RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		PublicUser: payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signature))
}

func CreateAccessTokens(payload models.PublicUser) (t models.Tokens, err error) {
	var tokens models.Tokens

	auth, err := GenerateUserToken(payload, time.Now().Add(time.Hour), AUTH_SIGNATURE)
	if err != nil {
		return tokens, err
	}

	refresh, err := GenerateUserToken(payload, time.Now().AddDate(0, 1, 0), REFRESH_SIGNATURE)
	if err != nil {
		return tokens, err
	}

	tokens = models.Tokens{
		UserID:       payload.ID,
		AuthToken:    auth,
		RefreshToken: refresh,
	}

	return tokens, nil
}

func ParseUserToken(tokenString string, signature string) (data *JWTPublicUserClaim, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTPublicUserClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(signature), nil
		},
	)

	claims, ok := token.Claims.(*JWTPublicUserClaim)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
