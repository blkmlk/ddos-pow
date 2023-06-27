package api

import (
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/services/api/controllers"
	"github.com/gin-gonic/gin"
)

type API interface {
	Start() error
	Stop() error
}

const (
	PathGetChallenge  = "/api/v1/challenge"
	PathPostChallenge = "/api/v1/challenge"
)

type api struct {
	restHost       string
	protocolHost   string
	restController *controllers.RestController
	restServer     *gin.Engine
}

func New(
	restController *controllers.RestController,
) (API, error) {
	restHost, err := env.Get(env.RestHost)
	if err != nil {
		return nil, err
	}

	a := api{
		restHost:       restHost,
		restController: restController,
		restServer:     gin.Default(),
	}

	a.initRest()

	return &a, nil
}

func (a *api) initRest() {
	a.restServer.GET(PathGetChallenge, a.restController.GetChallenge)
	a.restServer.POST(PathPostChallenge, a.restController.PostChallenge)
}

func (a *api) Start() error {
	return a.restServer.Run(a.restHost)
}

func (a *api) Stop() error {
	return nil
}
