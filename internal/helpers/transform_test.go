package helpers

import (
	"github.com/blkmlk/ddos-pow/pow"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChallengeFromToBytes(t *testing.T) {
	puzzle := pow.GeneratePuzzle()
	origin := pow.Challenge{
		Puzzle:    puzzle,
		ExpiresAt: 2000,
		N:         64,
		R:         2,
		P:         1,
		KeyLen:    16,
		MinZeroes: 15,
		Salt:      100,
	}
	origin.Signature = origin.Sign([]byte("secret"))

	data := ChallengeToBytes(&origin)

	recovered, err := ChallengeFromBytes(data)
	require.NoError(t, err)

	require.Equal(t, &origin, recovered)
}
