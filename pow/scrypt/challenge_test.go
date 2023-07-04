package scrypt

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestChallenge_FindSolution(t *testing.T) {
	ch := Challenge{
		Puzzle:    generatePuzzle(),
		N:         64,
		R:         2,
		P:         1,
		KeyLen:    16,
		MinZeroes: 12,
	}

	err := ch.FindSolution(time.Second)
	require.NoError(t, err)

	ch.Salt--
	solution, err := ch.generateSolution()
	require.NoError(t, err)
	require.False(t, verifySolution(solution, int(ch.MinZeroes)))

	ch.Puzzle = generatePuzzle()
	ch.Salt = 0
	err = ch.FindSolution(time.Millisecond)
	require.Error(t, err)
}
