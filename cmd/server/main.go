package main

import (
	"errors"
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/helpers"
	"github.com/blkmlk/ddos-pow/internal/quotes"
	"github.com/blkmlk/ddos-pow/internal/stream"
	"github.com/blkmlk/ddos-pow/pow"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	powClient *pow.POW
}

func NewServer(config pow.Config) *Server {
	return &Server{
		powClient: pow.New(config),
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	inst := stream.New(conn)

	challenge := s.powClient.NewSignedChallenge()
	data := helpers.ChallengeToBytes(&challenge)

	if err := inst.Write(data); err != nil {
		log.Fatal(err)
		return
	}

	received, err := inst.ReadUntil(pow.ChallengeMaxLength, time.Second)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}
		log.Fatal(err)
		return
	}

	solvedChallenge, err := helpers.ChallengeFromBytes(received)
	if err != nil {
		log.Fatal(err)
		return
	}

	challengeExp := time.Unix(0, solvedChallenge.ExpiresAt)
	now := time.Now()
	if now.After(challengeExp) {
		log.Println("expired", now.Sub(challengeExp))
		return
	}

	verified, err := s.powClient.VerifyChallenge(solvedChallenge)
	if err != nil {
		log.Fatal(err)
		return
	}

	if !verified {
		log.Fatal("no verified")
		return
	}

	quote := quotes.GetRandomQuote()
	if err := inst.Write([]byte(quote)); err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	host, err := env.Get(env.Host)

	server := NewServer(pow.Config{
		Secret:    []byte("secret"),
		Timeout:   time.Millisecond * 500,
		N:         64,
		R:         2,
		P:         1,
		KeyLen:    16,
		MinZeroes: 12,
	})

	log.Printf("listening to %v", host)

	ln, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.handleConnection(conn)
	}
}
