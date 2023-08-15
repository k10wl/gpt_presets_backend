package openai

import (
	"github.com/pkoukk/tiktoken-go"
)

type tokenizerClient interface {
	Encode(text string, allowedSpecial []string, disallowedSpecial []string) []int
}

type Tokenizer struct {
	client tokenizerClient
}

func newTokenizer(tokenizer tokenizerClient) *Tokenizer {
	return &Tokenizer{
		client: tokenizer,
	}
}

func tiktokenSetup(model string) tokenizerClient {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		panic("Could not initialize tiktoken")
	}

	return tkm
}
