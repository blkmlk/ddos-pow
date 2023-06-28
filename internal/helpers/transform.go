package helpers

import (
	"encoding/binary"
	"fmt"
	"github.com/blkmlk/ddos-pow/pow"
)

func ChallengeToBytes(ch *pow.Challenge) []byte {
	result := make([]byte, pow.ChallengeMaxLength)
	p := result

	copy(p, ch.Signature)
	p = p[len(ch.Signature):]

	copy(p, ch.Puzzle)
	p = p[len(ch.Puzzle):]

	binary.LittleEndian.PutUint64(p, uint64(ch.ExpiresAt))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(ch.N))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(ch.R))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(ch.P))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(ch.KeyLen))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(ch.MinZeroes))
	p = p[8:]
	binary.LittleEndian.PutUint64(p, uint64(ch.Salt))

	return result
}

func ChallengeFromBytes(data []byte) (*pow.Challenge, error) {
	if len(data) != pow.ChallengeMaxLength {
		return nil, fmt.Errorf("invalid data length")
	}

	var ch pow.Challenge
	ch.Signature = make([]byte, pow.SignatureLength)
	ch.Puzzle = make([]byte, pow.PuzzleLength)

	copy(ch.Signature, data[:pow.SignatureLength])
	data = data[pow.SignatureLength:]

	copy(ch.Puzzle, data[:pow.PuzzleLength])
	data = data[pow.PuzzleLength:]

	ch.ExpiresAt = int64(binary.LittleEndian.Uint64(data[:8]))
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
