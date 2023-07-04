package pow

import (
	"time"
)

type Challenge interface {
	Equals(ch Challenge) bool
	ExpiresAt() time.Time
	FindSolution(timeout time.Duration) error
	Bytes() []byte
}

type POW interface {
	NewChallenge() Challenge
	ParseChallenge(data []byte) (Challenge, error)
	VerifyChallenge(ch Challenge) (bool, error)
}
