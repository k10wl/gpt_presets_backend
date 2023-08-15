package openai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	apiUrl     string
	httpClient *http.Client
	tokenizer  *Tokenizer
}

const (
	DefaultModel         = "gpt-4"
	ChatCompletionAPIUrl = "https://api.openai.com/v1/chat/completions"
	HTTPTimeout          = 2 * time.Minute
	MaxRequestTokens     = 2048
)

func NewClient(apiKey string) *Client {
	tokenizerClient, err := GetTokenizerForModel(DefaultModel)
	if err != nil {
		panic("Could not initialize tokenizer, cannot proceed")
	}

	return &Client{
		apiKey:     apiKey,
		apiUrl:     ChatCompletionAPIUrl,
		httpClient: &http.Client{Timeout: HTTPTimeout},
		tokenizer:  tokenizerClient,
	}
}

func (c *Client) makePostRequest(body *[]byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.apiUrl, bytes.NewReader(*body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP Error: %d", res.StatusCode)
	}

	return io.ReadAll(res.Body)
}
