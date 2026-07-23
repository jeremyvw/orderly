package password

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	cost = bcrypt.MinCost
	os.Exit(m.Run())
}

func TestCompare_CorrectPassword_ReturnsNoError(t *testing.T) {
	hashed, err := Hash("s3cret")
	require.NoError(t, err)

	require.NoError(t, Compare("s3cret", hashed))
}

func TestCompare_IncorrectPassword_ReturnsError(t *testing.T) {
	hashed, err := Hash("s3cret")
	require.NoError(t, err)

	require.Error(t, Compare("wrong", hashed))
}

func TestHash_DoesNotReturnPlaintext(t *testing.T) {
	hashed, err := Hash("s3cret")
	require.NoError(t, err)

	require.NotEqual(t, "s3cret", hashed)
}

func TestHash_SamePasswordProducesDifferentHashes(t *testing.T) {
	hashed1, err := Hash("s3cret")
	require.NoError(t, err)

	hashed2, err := Hash("s3cret")
	require.NoError(t, err)

	require.NotEqual(t, hashed1, hashed2)
}

func TestHash_PaswordLongerThan72Bytes_ReturnsError(t *testing.T) {
	long := strings.Repeat("a", 100)
	_, err := Hash(long)
	t.Log(err)
	require.Error(t, err)
}
