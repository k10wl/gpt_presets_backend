package openai

import (
	"encoding/json"
	"errors"
)

type Usage struct {
	PromptTokens     uint `json:"prompt_tokens"`
	CompletionTokens uint `json:"completion_tokens"`
	TotalTokens      uint `json:"total_tokens"`
}

type Choices struct {
	Index        uint    `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Completion struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
	Usage   `          json:"usage"`
}

type Message struct {
	Role    string `json:"role,omitempty"    binding:"required"`
	Content string `json:"content,omitempty" binding:"required"`
}

type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
}

func (c *Client) TextCompletion(message *[]Message) (*Completion, error) {
	payloadData, err := json.Marshal(Request{
		Model:       DefaultModel,
		Messages:    *message,
		Temperature: 1,
	})
	if err != nil {
		return nil, err
	}

	responseData, err := c.makePostRequest(&payloadData)
	if err != nil {
		return nil, err
	}

	var resBody Completion
	if err := json.Unmarshal(responseData, &resBody); err != nil {
		return nil, err
	}

	return &resBody, nil
}

func (c *Client) HasTokensOverflow(message *Message) bool {
	tokens := c.CountMessageTokens(message)

	return len(tokens) > MaxRequestTokens
}

func (c *Client) CountMessageTokens(message *Message) []int {
	return c.tokenizer.client.Encode(message.Content)
}

func (c *Client) BuildHistory(prevMsg *[]Message) (msg *[]Message, err error) {
	temp := *prevMsg
	messages := []Message{}

	tokensUsage := 0

	for i := len(*prevMsg) - 1; i >= 0; i-- {
		tokens := c.CountMessageTokens(&temp[i])

		sum := tokensUsage + len(tokens)
		if sum > MaxRequestTokens {
			break
		}

		tokensUsage = sum
		messages = append([]Message{temp[i]}, messages...)
	}

	if len(messages) == 0 {
		return &messages, errors.New(TokensOverflowError)
	}

	return &messages, nil
}
