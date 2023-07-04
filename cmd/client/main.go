package main

import (
	"errors"
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/client"
	"github.com/blkmlk/ddos-pow/pow/scrypt"
	"go.uber.org/zap"
	"time"
)

func main() {
	logDev, _ := zap.NewDevelopment()
	log := logDev.Sugar()

	host, err := env.Get(env.Host)
	if err != nil {
		log.Fatal(err)
	}

	powVerifier := scrypt.New(scrypt.Config{})

	c := client.New(host, powVerifier)

	log.Info("starting the client...")

	for {
		time.Sleep(time.Second)

		quote, err := c.GetQuote()
		if err != nil {
			if errors.Is(err, client.ErrTerminated) {
				log.Info("connection has been terminated")
			} else {
				log.With("error", err).Error("failed to connect to the server")
			}
			continue
		}

		log.Infof("successfully got a quote: %s", quote)
	}
}
