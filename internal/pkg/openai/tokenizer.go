package openai

type tokenizerClient interface {
	Encode(text string) []int
}

type Tokenizer struct {
	client tokenizerClient
}

func NewTokenizer(tokenizer tokenizerClient) *Tokenizer {
	return &Tokenizer{
		client: tokenizer,
	}
}
