package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
		Model:       c.model,
		Messages:    *message,
		Temperature: 1,
	})
	if err != nil {
		return nil, err
	}

	responseData, err := c.post(&payloadData)
	if err != nil {
		fmt.Printf("%+v", err)
		return nil, err
	}

	var resBody Completion
	if err := json.Unmarshal(responseData, &resBody); err != nil {
		return nil, err
	}

	return &resBody, nil
}

func (c *Client) BuildHistory(prevMsg *[]Message) (msg *[]Message, err error) {
	temp := *prevMsg
	messages := []Message{}

	tokensUsage := 0

	for i := len(*prevMsg) - 1; i >= 0; i-- {
		tokens := c.tokenizer.client.Encode(temp[i].Content, nil, nil)

		sum := tokensUsage + len(tokens)
		if sum > c.maxRequestTokens {
			break
		}

		tokensUsage = sum
		messages = append([]Message{temp[i]}, messages...)
	}

	if len(messages) == 0 {
		str := strconv.FormatInt(int64(c.maxRequestTokens), 10)

		return &messages, errors.New("Request exceeds tokens limit: " + str)
	}

	return &messages, nil
}
