package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/blkmlk/ddos-pow/services/pow"
	"github.com/blkmlk/ddos-pow/services/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RestController struct {
	storage storage.Storage
}

const (
	ScryptN      = 16
	ScryptR      = 8
	ScryptP      = 1
	ScryptKeyLen = 32
	MinZeroes    = 12
)

func NewUploadController(storage storage.Storage) (*RestController, error) {
	return &RestController{
		storage: storage,
	}, nil
}

func (c *RestController) GetChallenge(ctx *gin.Context) {
	puzzle, err := pow.GeneratePuzzle()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	challenge := storage.NewChallenge(puzzle, ScryptN, ScryptR, ScryptP, ScryptKeyLen)
	if err := c.storage.CreateChallenge(ctx, &challenge); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, &GetChallengeResponse{
		ID:        challenge.ID,
		Puzzle:    challenge.Puzzle,
		N:         challenge.N,
		R:         challenge.R,
		P:         challenge.P,
		KeyLen:    challenge.KeyLen,
		MinZeroes: MinZeroes,
	})
}

func (c *RestController) PostChallenge(ctx *gin.Context) {
	var req PostChallengeRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	challenge, err := c.storage.GetChallenge(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	solution, err := base64.StdEncoding.DecodeString(req.Solution)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	salt, err := base64.StdEncoding.DecodeString(req.Salt)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if !pow.VerifySolution(solution, MinZeroes) {
		fmt.Println("1")
		ctx.Status(http.StatusForbidden)
		return
	}

	localSol, err := pow.GenerateSolution(pow.GenerateSolutionInput{
		Puzzle: challenge.Puzzle,
		Salt:   salt,
		N:      challenge.N,
		R:      challenge.R,
		P:      challenge.P,
		KeyLen: challenge.KeyLen,
	})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if !bytes.Equal(localSol, solution) {
		fmt.Println("2")
		ctx.Status(http.StatusForbidden)
		return
	}

	ctx.Status(http.StatusOK)
}
