package main

import (
	"errors"
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/helpers"
	"github.com/blkmlk/ddos-pow/internal/stream"
	"github.com/blkmlk/ddos-pow/pow"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	host, err := env.Get(env.Host)
	if err != nil {
		log.Fatal(err)
	}

	for {
		getQuote(host)
		time.Sleep(time.Second)
	}
}

func getQuote(host string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal("Connection error", err)
	}

	defer conn.Close()

	inst := stream.New(conn)

	data, err := inst.ReadUntil(pow.ChallengeMaxLength, time.Second*5)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}
		log.Fatal(err)
		return
	}

	challenge, err := helpers.ChallengeFromBytes(data)
	if err != nil {
		log.Fatal(err)
		return
	}

	startedAt := time.Now()
	for {
		solution, err := challenge.GenerateSolution()
		if err != nil {
			log.Fatal(err)
			return
		}
		if pow.VerifySolution(solution, int(challenge.MinZeroes)) {
			break
		}
		challenge.Salt++
	}
	elapsed := time.Since(startedAt)

	log.Printf("found solution in %v", elapsed)

	if err = inst.Write(helpers.ChallengeToBytes(challenge)); err != nil {
		log.Fatal(err)
		return
	}

	rawQuote, err := inst.ReadUntil(0, time.Second*5)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}
		log.Fatal(err)
	}

	log.Printf("found %s", string(rawQuote))
}
