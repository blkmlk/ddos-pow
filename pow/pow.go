package pow

import "time"

type Challenge interface {
	ExpiresAt() time.Time
	FindSolution() error
	Bytes() []byte
}

type POW interface {
	NewChallenge() Challenge
	ParseChallenge(data []byte) (Challenge, error)
	VerifyChallenge(ch Challenge) (bool, error)
}
