package main

import (
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/helpers"
	"github.com/blkmlk/ddos-pow/internal/protocol"
	"github.com/blkmlk/ddos-pow/pow"
	"github.com/blkmlk/ddos-pow/quotes"
	"log"
	"net"
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
	inst := protocol.New(conn)

	ch := s.powClient.NewSignedChallenge()
	data := helpers.ChallengeToBytes(&ch)

	if err := inst.Write(data); err != nil {
		log.Fatal(err)
		return
	}

	var inData []byte
	for {
		received, err := inst.Read(pow.ChallengeMaxLength)
		if err != nil {
			log.Fatal(err)
			return
		}
		if len(received) > 0 {
			inData = received
			break
		}
	}

	solution, err := helpers.ChallengeFromBytes(inData)
	if err != nil {
		log.Fatal(err)
		return
	}

	verified, err := s.powClient.VerifyChallenge(solution)
	if err != nil {
		log.Fatal(err)
		return
	}

	if !verified {
		log.Fatal("no verified")
		return
	}

	quote := quotes.getRandomQuote()
	if err := inst.Write([]byte(quote)); err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	host, err := env.Get(env.Host)

	server := NewServer(pow.Config{
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
