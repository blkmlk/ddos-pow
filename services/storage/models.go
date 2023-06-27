package storage

import (
	"github.com/google/uuid"
)

type Challenge struct {
	ID     string
	Puzzle string
	N      int
	R      int
	P      int
	KeyLen int
}

func NewChallenge(puzzle string, n, r, p, keyLen int) Challenge {
	return Challenge{
		ID:     uuid.NewString(),
		Puzzle: puzzle,
		N:      n,
		R:      r,
		P:      p,
		KeyLen: keyLen,
	}
}
