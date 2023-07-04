package scrypt

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/blkmlk/ddos-pow/pow"
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

func (c *Challenge) Equals(inCh pow.Challenge) bool {
	ch, ok := inCh.(*Challenge)
	if !ok {
		return false
	}

	return bytes.Equal(c.Signature, ch.Signature)
}

func (c *Challenge) ExpiresAt() time.Time {
	return time.Unix(0, c.ExpAt)
}

func (c *Challenge) Bytes() []byte {
	return challengeToBytes(c)
}

func (c *Challenge) FindSolution(timeout time.Duration) error {
	startedAt := time.Now()
	for {
		solution, err := c.generateSolution()
		if err != nil {
			return err
		}

		if timeout > 0 && time.Since(startedAt) > timeout {
			return fmt.Errorf("time is out")
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
