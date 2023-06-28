package main

import (
	"fmt"
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/internal/client"
	"log"
	"time"
)

func main() {
	host, err := env.Get(env.Host)
	if err != nil {
		log.Fatal(err)
	}

	c := client.New(host)

	for {
		quote, err := c.GetQuote()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("got quote", quote)

		time.Sleep(time.Second)
	}
}
