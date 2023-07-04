package scrypt

import (
	"crypto/sha256"
	"encoding/binary"
	"golang.org/x/crypto/scrypt"
	"math/big"
	"time"
)

type Challenge struct {
	Signature []byte
	Puzzle    []byte
	ExpAt     int64
	N         int64
	R         int64
	P         int64
	KeyLen    int64
	MinZeroes int64
	Salt      int64
}

func (c *Challenge) ExpiresAt() time.Time {
	return time.Unix(0, c.ExpAt)
}

func (c *Challenge) Bytes() []byte {
	return challengeToBytes(c)
}

func (c *Challenge) FindSolution() error {
	for {
		solution, err := c.generateSolution()
		if err != nil {
			return err
		}

		// solution is found
		if verifySolution(solution, int(c.MinZeroes)) {
			break
		}
		c.Salt++
	}
	return nil
}

func (c *Challenge) sign(secret []byte) []byte {
	expiresAt := make([]byte, 8)
	binary.LittleEndian.PutUint64(expiresAt, uint64(c.ExpAt))

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

func (c *Challenge) generateSolution() ([]byte, error) {
	salt := big.NewInt(c.Salt)
	return scrypt.Key(c.Puzzle, salt.Bytes(), int(c.N), int(c.R), int(c.P), int(c.KeyLen))
}
