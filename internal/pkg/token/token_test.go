package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGenerate_ReturnsNonEmptyToken(t *testing.T) {
	token, err := Generate(42)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestGenerate_ThenParse_ReturnsUserID(t *testing.T) {
	const userID int64 = 42

	token, err := Generate(userID)
	require.NoError(t, err)

	got, err := Parse(token)
	require.NoError(t, err)
	require.Equal(t, userID, got)
}

func TestGenerate_SetsExpiryTo24Hours(t *testing.T) {
	token, err := Generate(42)
	require.NoError(t, err)

	parsed, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return secret, nil
	})
	require.NoError(t, err)

	exp, err := parsed.Claims.GetExpirationTime()
	require.NoError(t, err)
	require.WithinDuration(t, time.Now().Add(tokenTTL), exp.Time, time.Minute)
}

func TestParse_ExpiredToken_ReturnsError(t *testing.T) {
	token := signWith(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "42",
		"exp": time.Now().Add(-time.Hour).Unix(),
	})

	_, err := Parse(token)
	require.Error(t, err)
}

func TestParse_TokenSignedWithDifferentSecret_ReturnsError(t *testing.T) {
	token := signWith(t, []byte("a-completely-different-secret"), jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "42",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	_, err := Parse(token)
	require.Error(t, err)
}

func TestParse_TokenWithoutExpiry_ReturnsError(t *testing.T) {
	token := signWith(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "42",
	})

	_, err := Parse(token)
	require.Error(t, err)
}

func TestParse_TokenWithNoneAlgorithm_ReturnsError(t *testing.T) {
	unsigned := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "42",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	token, err := unsigned.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = Parse(token)
	require.Error(t, err)
}

func TestParse_MalformedToken_ReturnsError(t *testing.T) {
	_, err := Parse("not-a-jwt")
	require.Error(t, err)
}

func TestParse_NonNumericSubject_ReturnsError(t *testing.T) {
	token := signWith(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "banana",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	_, err := Parse(token)
	require.Error(t, err)
}

func TestInit_RejectsShortSecret(t *testing.T) {
	defer restoreSecret(secret)

	require.Error(t, Init([]byte("too-short")))
}

func TestInit_AcceptsSecretOfSufficientLength(t *testing.T) {
	defer restoreSecret(secret)

	long := []byte("0123456789012345678901234567890123456789")
	require.NoError(t, Init(long))
	require.Equal(t, long, secret)
}

func signWith(t *testing.T, key []byte, method jwt.SigningMethod, claims jwt.MapClaims) string {
	t.Helper()

	token, err := jwt.NewWithClaims(method, claims).SignedString(key)
	require.NoError(t, err)
	return token
}

func restoreSecret(original []byte) {
	secret = original
}
