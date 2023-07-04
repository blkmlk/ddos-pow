package scrypt

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChallengeFromToBytes(t *testing.T) {
	puzzle := generatePuzzle()
	origin := Challenge{
		Puzzle:    puzzle,
		ExpAt:     2000,
		N:         64,
		R:         2,
		P:         1,
		KeyLen:    16,
		MinZeroes: 15,
		Salt:      100,
	}
	origin.Signature = origin.sign([]byte("secret"))

	data := challengeToBytes(&origin)

	recovered, err := challengeFromBytes(data)
	require.NoError(t, err)

	require.Equal(t, &origin, recovered)
}
