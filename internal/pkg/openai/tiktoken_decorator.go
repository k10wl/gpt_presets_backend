package openai

import (
	"github.com/pkoukk/tiktoken-go"
)

type TiktokenDecorator struct {
	model *tiktoken.Tiktoken
}

func NewTiktokenDecorator(model string) (*TiktokenDecorator, error) {
	tkn, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return nil, err
	}

	return &TiktokenDecorator{
		model: tkn,
	}, nil
}

func (tkn *TiktokenDecorator) Encode(text string) []int {
	return tkn.model.Encode(text, nil, nil)
}
