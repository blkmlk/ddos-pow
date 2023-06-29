package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/blkmlk/ddos-pow/internal/helpers"
	"github.com/blkmlk/ddos-pow/internal/quotes"
	"github.com/blkmlk/ddos-pow/internal/stream"
	"github.com/blkmlk/ddos-pow/pow"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	host      string
	powClient *pow.POW
	listener  net.Listener
	wg        sync.WaitGroup
	log       *zap.SugaredLogger
}

func New(host string, config pow.Config) *Server {
	logDev, _ := zap.NewDevelopment()

	return &Server{
		host:      host,
		powClient: pow.New(config),
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

	challenge := s.powClient.NewSignedChallenge()
	data := helpers.ChallengeToBytes(&challenge)

	// sending a new generated challenge
	if err := strm.Write(data, time.Second); err != nil {
		if clientNetworkErr(err) {
			return nil
		}
		return fmt.Errorf("failed to send challenge to solve: %v", err)
	}

	// puzzle timeout + network delay
	timeout := s.powClient.Config.Timeout + time.Millisecond*500
	received, err := strm.Read(pow.ChallengeMaxLength, timeout)
	if err != nil {
		if clientNetworkErr(err) {
			return nil
		}
		return fmt.Errorf("failed to read challenge: %v", err)
	}

	solvedChallenge, err := helpers.ChallengeFromBytes(received)
	if err != nil {
		return fmt.Errorf("failed to get challenge from bytes: %v", err)
	}

	// checking if the challenge is not what we just sent
	if !bytes.Equal(challenge.Signature, solvedChallenge.Signature) {
		return fmt.Errorf("received a wrong challenge")
	}

	// checking the solution for expiration
	challengeExp := time.Unix(0, solvedChallenge.ExpiresAt)
	if time.Now().After(challengeExp) {
		return fmt.Errorf("solution has been expired")
	}

	valid, err := s.powClient.VerifyChallenge(solvedChallenge)
	if err != nil {
		return fmt.Errorf("failed to verify challenge %v", err)
	}

	if !valid {
		return fmt.Errorf("solution is not valid")
	}

	// sending the quote
	quote := quotes.GetRandomQuote()
	if err = strm.Write([]byte(quote), time.Second); err != nil {
		if clientNetworkErr(err) {
			return nil
		}
		return fmt.Errorf("failed to send quote: %v", err)
	}

	return nil
}

func clientNetworkErr(err error) bool {
	// client closed the connection themselves - close the connection
	if errors.Is(err, io.EOF) {
		return true
	}
	// client is too slow - close the connection
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return true
	}
	return false
}
