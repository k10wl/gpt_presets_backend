package openai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey           string
	apiUrl           string
	model            string
	httpClient       *http.Client
	tokenizer        *Tokenizer
	maxRequestTokens int
}

const (
	DEFAULT_MODEL            = "gpt-3.5-turbo"
	CHAT_COMPLETIONS_API_URL = "https://api.openai.com/v1/chat/completions"
	HTTP_TIMEOUT             = 2 * time.Minute
	MAX_REQUEST_TOKENS       = 1024
)

func NewClient(apiKey string) (client *Client) {
	return &Client{
		apiKey:           apiKey,
		apiUrl:           CHAT_COMPLETIONS_API_URL,
		model:            DEFAULT_MODEL,
		httpClient:       &http.Client{Timeout: HTTP_TIMEOUT},
		tokenizer:        newTokenizer(tiktokenSetup(DEFAULT_MODEL)),
		maxRequestTokens: MAX_REQUEST_TOKENS,
	}
}

func (c *Client) SetModel(model string) {
	c.model = model
	c.tokenizer = newTokenizer(tiktokenSetup(model))
}

func (c *Client) post(body *[]byte) ([]byte, error) {
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
