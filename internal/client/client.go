package client

import (
	"errors"
	"fmt"
	"github.com/blkmlk/ddos-pow/internal/stream"
	"github.com/blkmlk/ddos-pow/pow"
	"net"
	"time"
)

var (
	ErrTerminated = errors.New("terminated")
)

const (
	AwaitTimeout        = time.Second * 5
	FindSolutionTimeout = time.Second
)

type Client struct {
	host string
	pow  pow.POW
}

func New(host string, pow pow.POW) *Client {
	return &Client{
		host: host,
		pow:  pow,
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
	data, err := strm.Read(0, AwaitTimeout)
	if err != nil {
		if errors.Is(err, stream.ErrClosed) {
			return "", ErrTerminated
		}
		return "", fmt.Errorf("failed to read challenge from the stream: %v", err)
	}

	challenge, err := c.pow.ParseChallenge(data)
	if err != nil {
		return "", fmt.Errorf("failed to get challenge from bytes: %v", err)
	}

	// looking for the solution
	if err = challenge.FindSolution(FindSolutionTimeout); err != nil {
		return "", fmt.Errorf("failed to find a solution")
	}

	// sending the solution for verification
	if err = strm.Write(challenge.Bytes(), time.Second); err != nil {
		return "", fmt.Errorf("failed to send the solution: %v", err)
	}

	// waiting for a quote
	rawQuote, err := strm.Read(0, time.Second*5)
	if err != nil {
		if errors.Is(err, stream.ErrClosed) {
			return "", ErrTerminated
		}
		return "", fmt.Errorf("failed to read a quote from the stream: %v", err)
	}

	return string(rawQuote), nil
}
