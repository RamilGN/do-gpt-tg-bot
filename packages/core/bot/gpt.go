package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var errGpt = errors.New("gpt error")

const (
	gptOpenaiAPIEnv  = "OPENAI_API_TOKEN"
	gptCompletionURL = "https://api.openai.com/v1/chat/completions"
	gptModel         = "gpt-3.5-turbo"
	gptRoleSystem    = "system"
	gptRoleUser      = "user"
)

type Gpt struct {
	Client *http.Client
	Body   *GptBody
	Token  string
}

func NewGpt(client *http.Client) (*Gpt, error) {
	token := os.Getenv(gptOpenaiAPIEnv)
	if token == "" {
		return nil, fmt.Errorf("%w. empty OPENAI_API_TOKEN", errGpt)
	}

	gpt := new(Gpt)
	gpt.Token = token
	gpt.Body = &GptBody{
		Model:    gptModel,
		Messages: []GptMessage{},
	}
	gpt.Body.AddMessage(GptMessage{Role: gptRoleSystem, Content: "You are helpful assistant."})

	if client == nil {
		gpt.Client = &http.Client{
			Timeout:       time.Minute,
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
		}
	} else {
		gpt.Client = client
	}

	return gpt, nil
}

func (g *Gpt) Prompt(ctx context.Context, prompt string) (string, error) {
	g.Body.AddMessage(GptMessage{Role: gptRoleUser, Content: prompt})

	body, err := json.Marshal(g.Body)
	if err != nil {
		return "", fmt.Errorf("%w. %w", errGpt, err)
	}

	req, err := NewJSONRequest(ctx, http.MethodPost, gptCompletionURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("%w. %w", errGpt, err)
	}

	req.Header.Set(authorization, fmt.Sprintf("%s %s", bearer, g.Token))

	resp, err := g.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w. %w", errGpt, err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w. %w", errGpt, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w. %s", errGpt, string(data))
	}

	gpr := new(GptResponse)

	err = json.Unmarshal(data, gpr)
	if err != nil {
		return "", fmt.Errorf("%w. %w", errGpt, err)
	}

	if len(gpr.Choices) == 0 {
		return "", fmt.Errorf("%w. no choices", errGpt)
	}

	return gpr.Choices[0].Message.Content, nil
}

type GptBody struct {
	Model    string       `json:"model"`
	Messages []GptMessage `json:"messages"`
}

func (g *GptBody) AddMessage(message GptMessage) {
	g.Messages = append(g.Messages, message)
}

type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
