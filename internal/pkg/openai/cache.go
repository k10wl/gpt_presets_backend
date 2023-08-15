package openai

import "sync"

var (
	cachedModels = make(map[string]*Tokenizer)
	cacheMutex   = &sync.Mutex{}
)

func GetTokenizerForModel(model string) (*Tokenizer, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	tokenizer, exists := cachedModels[model]

	if exists {
		return tokenizer, nil
	}

	tiktoken, err := NewTiktokenDecorator(model)
	if err != nil {
		return nil, err
	}

	newTokenizer := NewTokenizer(tiktoken)

	cachedModels[model] = newTokenizer

	return newTokenizer, nil
}
