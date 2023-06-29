package pow_test

import (
	"encoding/binary"
	"github.com/blkmlk/ddos-pow/pow"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"time"
)

func TestPOW(t *testing.T) {
	p := pow.New(pow.Config{
		N:         64,
		R:         2,
		P:         1,
		KeyLen:    16,
		MinZeroes: 12,
	})

	challenge := p.NewSignedChallenge()

	// find solution
	startedAt := time.Now()
	for {
		solution, err := challenge.GenerateSolution()
		require.NoError(t, err)

		if pow.VerifySolution(solution, int(challenge.MinZeroes)) {
			break
		}
		challenge.Salt++
	}
	t.Logf("salt %d found in %v", challenge.Salt, time.Since(startedAt))

	valid, err := p.VerifyChallenge(&challenge)
	require.NoError(t, err)
	require.True(t, valid)

	// check invalid solution
	challenge.Salt--
	solution, err := challenge.GenerateSolution()
	require.NoError(t, err)
	require.False(t, pow.VerifySolution(solution, int(challenge.MinZeroes)))

	// check invalid signature
	challenge.KeyLen = 32
	valid, err = p.VerifyChallenge(&challenge)
	require.NoError(t, err)
	require.False(t, valid)
}

func TestVerifySolution(t *testing.T) {
	sol64 := make([]byte, 8)
	binary.BigEndian.PutUint64(sol64, math.MaxUint64)
	require.False(t, pow.VerifySolution(sol64, 1))

	sol63 := make([]byte, 8)
	binary.BigEndian.PutUint64(sol63, math.MaxInt64)
	require.True(t, pow.VerifySolution(sol63, 1))
}
