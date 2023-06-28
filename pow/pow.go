package pow

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"golang.org/x/crypto/scrypt"
	"math/big"
	"strconv"
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

type POW struct {
	Config Config
}

func New(config Config) *POW {
	return &POW{Config: config}
}

func (p *POW) NewSignedChallenge() Challenge {
	ch := Challenge{
		Puzzle:    GeneratePuzzle(),
		ExpiresAt: time.Now().Add(p.Config.Timeout).UnixNano(),
		N:         p.Config.N,
		R:         p.Config.R,
		P:         p.Config.P,
		KeyLen:    p.Config.KeyLen,
		MinZeroes: p.Config.MinZeroes,
	}
	ch.Signature = ch.Sign(p.Config.Secret)
	return ch
}

func (p *POW) VerifyChallenge(ch *Challenge) (bool, error) {
	validSignature := ch.Sign(p.Config.Secret)
	if !bytes.Equal(validSignature, ch.Signature) {
		return false, nil
	}

	solution, err := ch.GenerateSolution()
	if err != nil {
		return false, err
	}

	if !VerifySolution(solution, int(ch.MinZeroes)) {
		return false, nil
	}
	return true, nil
}

type Challenge struct {
	Signature []byte
	Puzzle    []byte
	ExpiresAt int64
	N         int64
	R         int64
	P         int64
	KeyLen    int64
	MinZeroes int64
	Salt      int64
}

func (c *Challenge) Sign(secret []byte) []byte {
	expiresAt := make([]byte, 8)
	binary.LittleEndian.PutUint64(expiresAt, uint64(c.ExpiresAt))

	N := make([]byte, 8)
	binary.LittleEndian.PutUint64(N, uint64(c.N))

	R := make([]byte, 8)
	binary.LittleEndian.PutUint64(R, uint64(c.R))

	P := make([]byte, 8)
	binary.LittleEndian.PutUint64(P, uint64(c.P))

	KeyLen := make([]byte, 8)
	binary.LittleEndian.PutUint64(P, uint64(c.KeyLen))

	MinZeroes := make([]byte, 8)
	binary.LittleEndian.PutUint64(P, uint64(c.MinZeroes))

	h := sha256.New()
	h.Write(c.Puzzle)
	h.Write(expiresAt)
	h.Write(N)
	h.Write(R)
	h.Write(P)
	h.Write(KeyLen)
	h.Write(MinZeroes)
	h.Write(secret)

	return h.Sum(nil)
}

func GeneratePuzzle() []byte {
	puzzle := make([]byte, PuzzleLength)
	_, _ = rand.Read(puzzle)
	return puzzle
}

func (c *Challenge) GenerateSolution() ([]byte, error) {
	salt := big.NewInt(c.Salt)
	return scrypt.Key(c.Puzzle, salt.Bytes(), int(c.N), int(c.R), int(c.P), int(c.KeyLen))
}

func VerifySolution(solution []byte, minZeroes int) bool {
	sumUint64 := binary.BigEndian.Uint64(solution)
	sumBits := strconv.FormatUint(sumUint64, 2)
	zeroes := 64 - len(sumBits)
	return uint(zeroes) >= uint(minZeroes)
}
