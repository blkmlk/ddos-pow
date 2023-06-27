package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/blkmlk/ddos-pow/env"
	"github.com/blkmlk/ddos-pow/services/api"
	"github.com/blkmlk/ddos-pow/services/api/controllers"
	"github.com/blkmlk/ddos-pow/services/pow"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	host, err := env.Get(env.RestHost)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(time.Second)

		challenge, err := getChallenge(ctx, host)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		started := time.Now()
		solution, salt, err := findSolution(challenge)
		elapsed := time.Since(started)
		if err != nil {
			log.Fatal(err)
		}

		enSolution := base64.StdEncoding.EncodeToString(solution)
		enSalt := base64.StdEncoding.EncodeToString(salt)

		if err = sendSolution(ctx, host, challenge.ID, enSolution, enSalt); err != nil {
			log.Fatal(err)
		}

		log.Printf("solution found in %v", elapsed)
	}
}

func getChallenge(ctx context.Context, host string) (*controllers.GetChallengeResponse, error) {
	hostUrl := fmt.Sprintf("http://%s%s", host, api.PathGetChallenge)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, hostUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected code %d", resp.StatusCode)
	}

	var body controllers.GetChallengeResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

func sendSolution(ctx context.Context, host string, id, solution, salt string) error {
	request := controllers.PostChallengeRequest{
		ID:       id,
		Solution: solution,
		Salt:     salt,
	}

	body, err := json.Marshal(&request)
	if err != nil {
		return err
	}

	hostUrl := fmt.Sprintf("http://%s%s", host, api.PathPostChallenge)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hostUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func findSolution(challenge *controllers.GetChallengeResponse) ([]byte, []byte, error) {
	input := pow.GenerateSolutionInput{
		Puzzle: challenge.Puzzle,
		N:      challenge.N,
		R:      challenge.R,
		P:      challenge.P,
		KeyLen: challenge.KeyLen,
	}
	for {
		salt, err := pow.GenerateSalt()
		if err != nil {
			return nil, nil, err
		}

		input.Salt = salt

		solution, err := pow.GenerateSolution(input)
		if err != nil {
			return nil, nil, err
		}

		if pow.VerifySolution(solution, challenge.MinZeroes) {
			return solution, salt, nil
		}
	}
}
