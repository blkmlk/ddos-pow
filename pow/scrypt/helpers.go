package scrypt

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strconv"
)

func challengeToBytes(c *Challenge) []byte {
	result := make([]byte, ChallengeMaxLength)
	p := result

	copy(p, c.Signature)
	p = p[len(c.Signature):]

	copy(p, c.Puzzle)
	p = p[len(c.Puzzle):]

	binary.LittleEndian.PutUint64(p, uint64(c.ExpAt))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(c.N))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(c.R))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(c.P))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(c.KeyLen))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(c.MinZeroes))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(c.Salt))

	return result
}

func challengeFromBytes(data []byte) (*Challenge, error) {
	if len(data) != ChallengeMaxLength {
		return nil, fmt.Errorf("invalid data length")
	}

	var ch Challenge
	ch.Signature = make([]byte, SignatureLength)
	ch.Puzzle = make([]byte, PuzzleLength)

	copy(ch.Signature, data[:SignatureLength])
	data = data[SignatureLength:]

	copy(ch.Puzzle, data[:PuzzleLength])
	data = data[PuzzleLength:]

	ch.ExpAt = int64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	ch.N = int64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	ch.R = int64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	ch.P = int64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	ch.KeyLen = int64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	ch.MinZeroes = int64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	ch.Salt = int64(binary.LittleEndian.Uint64(data[:8]))

	return &ch, nil
}

func generatePuzzle() []byte {
	puzzle := make([]byte, PuzzleLength)
	_, _ = rand.Read(puzzle)
	return puzzle
}

func verifySolution(solution []byte, minZeroes int) bool {
	sumUint64 := binary.BigEndian.Uint64(solution)
	sumBits := strconv.FormatUint(sumUint64, 2)
	zeroes := 64 - len(sumBits)
	return uint(zeroes) >= uint(minZeroes)
}
