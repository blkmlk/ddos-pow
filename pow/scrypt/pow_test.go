package scrypt

import (
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"time"
)

func TestPOW(t *testing.T) {
	p := New(Config{
		Secret:    []byte("secret"),
		N:         64,
		R:         2,
		P:         1,
		KeyLen:    16,
		MinZeroes: 12,
	})

	challenge := p.NewChallenge()

	// find solution
	startedAt := time.Now()
	require.NoError(t, challenge.FindSolution(time.Second))
	t.Logf("solution found in %v", time.Since(startedAt))

	valid, err := p.VerifyChallenge(challenge)
	require.NoError(t, err)
	require.True(t, valid)
}

func TestVerifySolution(t *testing.T) {
	sol64 := make([]byte, 8)
	binary.BigEndian.PutUint64(sol64, math.MaxUint64)
	require.False(t, verifySolution(sol64, 1))

	sol63 := make([]byte, 8)
	binary.BigEndian.PutUint64(sol63, math.MaxInt64)
	require.True(t, verifySolution(sol63, 1))
}
