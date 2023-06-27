package pow

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"golang.org/x/crypto/scrypt"
	"strconv"
)

const (
	PuzzleLength = 50
)

type GenerateSolutionInput struct {
	Puzzle string
	Salt   []byte
	N      int
	R      int
	P      int
	KeyLen int
}

func GeneratePuzzle() (string, error) {
	puzzle := make([]byte, PuzzleLength)
	if _, err := rand.Read(puzzle); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(puzzle), nil
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, PuzzleLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateSolution(input GenerateSolutionInput) ([]byte, error) {
	return scrypt.Key([]byte(input.Puzzle), input.Salt, input.N, input.R, input.P, input.KeyLen)
}

func VerifySolution(solution []byte, minZeroes int) bool {
	sumUint64 := binary.BigEndian.Uint64(solution)
	sumBits := strconv.FormatUint(sumUint64, 2)
	zeroes := 64 - len(sumBits)
	return uint(zeroes) >= uint(minZeroes)
}
