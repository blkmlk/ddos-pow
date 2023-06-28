package quotes

import (
	"math/rand"
	"time"
)

var quotes = []string{
	"You create your own opportunities",
	"Never break your promises",
	"You are never as stuck as you think you are",
	"Happiness is a choice",
	"Habits develop into character",
	"Be happy with who you are",
	"Don’t seek happiness–create it",
	"If you want to be happy, stop complaining",
	"Asking for help is a sign of strength",
	"Replace every negative thought with a positive one",
	"Accept what is, let go of what was, have faith in what will be",
	"A mind that is stretched by a new experience can never go back to what it was",
	"If you are not willing to learn, no one can help you",
	"Be confident enough to encourage confidence in others",
	"Allow others to figure things out for themselves",
	"Confidence is essential for a successful life",
	"Admit your mistakes and don’t repeat them",
	"Be kind to yourself and forgive yourself",
	"Failures are lessons in progress",
	"Make amends with those who have wronged you",
	"Live your life on your terms",
	"When you don’t know, don’t speak as if you do",
	"Treat others the way you want to be treated",
	"Think before you speak",
	"Cultivate an attitude of gratitude",
	"Life isn’t as serious as our minds make it out to be",
	"Take risks and be bold",
	"Remember that “no” is a complete sentence",
	"Don’t feed yourself only on leftovers",
	"Build on your strengths",
	"Never doubt your instincts",
	"FEAR doesn’t have to stand for Forget Everything and Run",
	"Your attitude will influence your experience",
	"View your life with gentle hindsight",
	"This too shall pass",
}

func getRandomQuote() string {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	return quotes[r.Intn(len(quotes))]
}
