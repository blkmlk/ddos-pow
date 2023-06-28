package main

import (
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/helpers"
	"github.com/blkmlk/ddos-pow/internal/stream"
	"github.com/blkmlk/ddos-pow/pow"
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

	data, err := inst.Read(pow.ChallengeMaxLength)
	if err != nil {
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

	if err = inst.Write(helpers.ChallengeToBytes(challenge)); err != nil {
		log.Fatal(err)
		return
	}

	var quote string
	for {
		rawQuote, err := inst.Read(0)
		if err != nil {
			log.Fatal(err)
		}
		if len(rawQuote) > 0 {
			quote = string(rawQuote)
			break
		}
	}

	log.Printf("found (%s) solution in %v", quote, elapsed)
}
