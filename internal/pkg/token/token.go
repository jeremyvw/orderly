package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 24 * time.Hour

var secret = []byte("secret-key")

// Init sets the signing secret. Must be called at startup before any
// token is issued; the package default is for tests only.
func Init(s []byte) error {
	if len(s) < 32 {
		return errors.New("jwt secret must be at least 32 bytes")
	}
	secret = s
	return nil
}

func Generate(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": time.Now().Add(tokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}

func Parse(token string) (int64, error) {
	jwtToken, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) { return secret, nil },
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return 0, fmt.Errorf("parse token: %w", err)
	}

	sub, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("read subject: %w", err)
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid subject %q: %w", sub, err)
	}

	return userID, nil
}
