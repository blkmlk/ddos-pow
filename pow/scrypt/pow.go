package scrypt

import (
	"bytes"
	"fmt"
	"github.com/blkmlk/ddos-pow/pow"
	"time"
)

const (
	SignatureLength = 32
	// ChallengeMaxLength sha256 + PuzzleLength + ExpiresAt + N + R + P + KeyLen + MinZeroes + Salt
	ChallengeMaxLength = SignatureLength + PuzzleLength*8 + 7*8
	PuzzleLength       = 20
)

type Config struct {
	Secret    []byte
	Timeout   time.Duration
	N         int64
	R         int64
	P         int64
	KeyLen    int64
	MinZeroes int64
}

func New(config Config) pow.POW {
	return &scryptPOW{
		config: config,
	}
}

type scryptPOW struct {
	config Config
}

func (p *scryptPOW) NewChallenge() pow.Challenge {
	ch := Challenge{
		Puzzle:    generatePuzzle(),
		ExpAt:     time.Now().Add(p.config.Timeout).UnixNano(),
		N:         p.config.N,
		R:         p.config.R,
		P:         p.config.P,
		KeyLen:    p.config.KeyLen,
		MinZeroes: p.config.MinZeroes,
	}
	ch.Signature = ch.sign(p.config.Secret)
	return &ch
}

func (p *scryptPOW) ParseChallenge(data []byte) (pow.Challenge, error) {
	return challengeFromBytes(data)
}

func (p *scryptPOW) VerifyChallenge(powCh pow.Challenge) (bool, error) {
	ch, ok := powCh.(*Challenge)
	if !ok {
		return false, fmt.Errorf("wrong challenge")
	}

	validSignature := ch.sign(p.config.Secret)
	if !bytes.Equal(validSignature, ch.Signature) {
		return false, nil
	}

	solution, err := ch.generateSolution()
	if err != nil {
		return false, err
	}

	if !verifySolution(solution, int(ch.MinZeroes)) {
		return false, nil
	}
	return true, nil
}
