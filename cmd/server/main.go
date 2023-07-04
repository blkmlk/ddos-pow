package main

import (
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/server"
	"github.com/blkmlk/ddos-pow/pow/scrypt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logDev, _ := zap.NewDevelopment()
	log := logDev.Sugar()

	host, err := env.Get(env.Host)
	if err != nil {
		log.Fatal(err)
	}

	scryptPow := scrypt.New(scrypt.Config{
		Secret:    []byte("secret"),
		Timeout:   time.Second,
		N:         64,
		R:         2,
		P:         1,
		KeyLen:    16,
		MinZeroes: 12,
	})

	s := server.New(host, scryptPow)

	go func() {
		log.With("host", host).Info("starting the server...")

		if err := s.Start(); err != nil {
			log.With("error", err).Fatal("failed to start the server")
		}
	}()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	log.Infof("got signal to stop the server")

	if err := s.Stop(); err != nil {
		log.With("error", err).Fatal("failed to stop the server")
	}
}
