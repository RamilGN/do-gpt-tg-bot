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
	"regexp"
	"strings"
	"time"
)

var errTelegram = errors.New("telegram error")

const (
	tgAPIURL         = "https://api.telegram.org"
	tgTelegramBotEnv = "TELEGRAM_BOT_TOKEN"
	tgClientTimeout  = time.Second * 7
)

type Tg struct {
	Client *http.Client
	Token  string
	BotURL string
}

func NewTg(client *http.Client) (*Tg, error) {
	token := os.Getenv(tgTelegramBotEnv)
	if token == "" {
		return nil, fmt.Errorf("%w. empty TELEGRAM_BOT_TOKEN", errTelegram)
	}

	tg := new(Tg)
	tg.BotURL = fmt.Sprintf("%s/bot%s", tgAPIURL, token)

	if client == nil {
		tg.Client = &http.Client{
			Timeout:       tgClientTimeout,
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
		}
	} else {
		tg.Client = client
	}

	return tg, nil
}

const sendMessageMethod = "sendMessage"

type TgSendMessage struct {
	ChatID           int    `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageID int    `json:"reply_to_message_id"`
}

func (t *Tg) SendMessage(ctx context.Context, m TgSendMessage) error {
	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("%w. %w", errTelegram, err)
	}

	req, err := NewJSONRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s", t.BotURL, sendMessageMethod),
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("%w. %w", errTelegram, err)
	}

	resp, err := t.Client.Do(req)
	if err != nil {
		return fmt.Errorf("%w. %w", errTelegram, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%w. %w", errTelegram, err)
		}

		return fmt.Errorf("%w. %s", errTelegram, string(body))
	}

	return nil
}

type TgUpdate struct {
	UpdateID int       `json:"update_id"`
	Message  TgMessage `json:"message"`
}

func ValidateTgUpdate(t *TgUpdate) error {
	if t == nil {
		return fmt.Errorf("%w. nil update", errTelegram)
	}

	match := regexp.MustCompile("^/ask .*").MatchString(strings.TrimSpace(t.Message.Text))
	if !match {
		return fmt.Errorf("%w. empty message", errTelegram)
	}

	return nil
}

type TgMessage struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
		IsPremium    bool   `json:"is_premium"`
	} `json:"from"`
	Chat TgChat `json:"chat"`
	Date int    `json:"date"`
	Text string `json:"text"`
}

type TgChat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}
