package main

import (
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/client"
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

	c := client.New(host)

	log.Info("starting the client...")

	for {
		quote, err := c.GetQuote()
		if err != nil {
			log.Error(err)
		}

		log.Infof("successfully got quote: %s", quote)

		time.Sleep(time.Second)
	}
}
