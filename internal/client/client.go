package client

import (
	"errors"
	"fmt"
	"github.com/blkmlk/ddos-pow/internal/helpers"
	"github.com/blkmlk/ddos-pow/internal/stream"
	"github.com/blkmlk/ddos-pow/pow"
	"io"
	"net"
	"time"
)

var (
	ErrTerminated = errors.New("terminated")
)

const (
	AwaitTimeout = time.Second * 5
)

type Client struct {
	host string
}

func New(host string) *Client {
	return &Client{
		host: host,
	}
}

func (c *Client) GetQuote() (string, error) {
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the server: %v", err)
	}

	defer conn.Close()

	strm := stream.New(conn)

	// waiting for a new generated challenge
	data, err := strm.ReadUntil(pow.ChallengeMaxLength, AwaitTimeout)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "", ErrTerminated
		}
		return "", fmt.Errorf("failed to read challenge from the stream: %v", err)
	}

	challenge, err := helpers.ChallengeFromBytes(data)
	if err != nil {
		return "", fmt.Errorf("failed to get challenge from bytes: %v", err)
	}

	// looking for the solution
	for {
		solution, err := challenge.GenerateSolution()
		if err != nil {
			return "", err
		}

		// solution is found
		if pow.VerifySolution(solution, int(challenge.MinZeroes)) {
			break
		}
		challenge.Salt++
	}

	// sending the solution for verification
	if err = strm.Write(helpers.ChallengeToBytes(challenge)); err != nil {
		return "", fmt.Errorf("failed to send the solution: %v", err)
	}

	// waiting for a quote
	rawQuote, err := strm.ReadUntil(0, time.Second*5)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "", ErrTerminated
		}
		return "", fmt.Errorf("failed to read a quote from the stream: %v", err)
	}

	return string(rawQuote), nil
}
