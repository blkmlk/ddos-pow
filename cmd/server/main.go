package main

import (
	"github.com/blkmlk/ddos-pow/services/api"
	"github.com/blkmlk/ddos-pow/services/api/controllers"
	"github.com/blkmlk/ddos-pow/services/storage"
	"go.uber.org/dig"
	"log"
)

func main() {
	container := dig.New()

	container.Provide(controllers.NewUploadController)
	container.Provide(storage.NewMapStorage)
	container.Provide(api.New)

	var listener api.API
	err := container.Invoke(func(a api.API) {
		listener = a
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := listener.Start(); err != nil {
		panic(err)
	}
}
