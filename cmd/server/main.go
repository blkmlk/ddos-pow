package main

import (
	"fmt"
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/server"
	"github.com/blkmlk/ddos-pow/pow"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	host, err := env.Get(env.Host)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(host, pow.Config{
		Secret:    []byte("secret"),
		Timeout:   time.Millisecond * 500,
		N:         64,
		R:         2,
		P:         1,
		KeyLen:    16,
		MinZeroes: 12,
	})

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		<-ch
		fmt.Println("stopping the server")

		if err := s.Stop(); err != nil {
			log.Fatal(err)
		}
	}()
}
