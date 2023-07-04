package server

import (
	"errors"
	"fmt"
	"github.com/blkmlk/ddos-pow/internal/quotes"
	"github.com/blkmlk/ddos-pow/internal/stream"
	"github.com/blkmlk/ddos-pow/pow"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

const (
	NetworkDelay = time.Millisecond * 500
)

type Server struct {
	host      string
	powClient pow.POW
	listener  net.Listener
	wg        sync.WaitGroup
	log       *zap.SugaredLogger
}

func New(host string, powClient pow.POW) *Server {
	logDev, _ := zap.NewDevelopment()

	return &Server{
		host:      host,
		powClient: powClient,
		log:       logDev.Sugar(),
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.host)
	if err != nil {
		return err
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// server is stopped
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			err = s.handleConnection(conn)
			if err != nil {
				s.log.With("error", err.Error()).Info("failed to handle connection")
			}
		}()
	}
}

func (s *Server) Stop() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	s.wg.Wait()
	return nil
}

func (s *Server) handleConnection(conn net.Conn) error {
	defer conn.Close()
	strm := stream.New(conn)

	challenge := s.powClient.NewChallenge()
	data := challenge.Bytes()

	// sending a new generated challenge
	if err := strm.Write(data, NetworkDelay); err != nil {
		if clientNetworkErr(err) {
			return nil
		}
		return fmt.Errorf("failed to send challenge to solve: %v", err)
	}

	// puzzle timeout + network delay
	timeToSolve := challenge.ExpiresAt().Sub(time.Now()) + NetworkDelay
	received, err := strm.Read(len(data), timeToSolve)
	if err != nil {
		if clientNetworkErr(err) {
			return nil
		}
		return fmt.Errorf("failed to read challenge: %v", err)
	}

	solvedChallenge, err := s.powClient.ParseChallenge(received)
	if err != nil {
		return fmt.Errorf("failed to get challenge from bytes: %v", err)
	}

	// checking if the challenge is not what we just sent
	if !challenge.Equals(solvedChallenge) {
		return fmt.Errorf("received a wrong challenge")
	}

	// checking if the solution is correct
	valid, err := s.powClient.VerifyChallenge(solvedChallenge)
	if err != nil {
		return fmt.Errorf("failed to verify challenge %v", err)
	}

	if !valid {
		return fmt.Errorf("solution is not valid")
	}

	// sending the quote
	quote := quotes.GetRandomQuote()
	if err = strm.Write([]byte(quote), NetworkDelay); err != nil {
		if clientNetworkErr(err) {
			return nil
		}
		return fmt.Errorf("failed to send quote: %v", err)
	}

	return nil
}

func clientNetworkErr(err error) bool {
	// client closed the connection themselves - close the connection
	if errors.Is(err, stream.ErrClosed) {
		return true
	}
	// client is too slow - close the connection
	if errors.Is(err, stream.ErrExpired) {
		return true
	}
	// client is trying to send more than it's required
	if errors.Is(err, stream.ErrMaxLenExceeded) {
		return true
	}
	return false
}
