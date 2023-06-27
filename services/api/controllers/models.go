package controllers

type GetChallengeResponse struct {
	ID        string `json:"id"`
	Puzzle    string `json:"puzzle"`
	N         int    `json:"n"`
	R         int    `json:"r"`
	P         int    `json:"p"`
	KeyLen    int    `json:"key_len"`
	MinZeroes int    `json:"min_zeroes"`
}

type PostChallengeRequest struct {
	ID       string `json:"id"`
	Solution string `json:"solution"`
	Salt     string `json:"salt"`
}

type PostChallengeResponse struct {
	Quote string `json:"quote"`
}
