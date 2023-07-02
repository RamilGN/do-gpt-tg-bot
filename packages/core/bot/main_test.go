package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRunMain(t *testing.T) {
	t.Setenv(gptOpenaiAPIEnv, "open_ai_test_token")
	t.Setenv(tgTelegramBotEnv, "tg_telegram_bot_token")

	ctx := context.Background()

	tg, err := NewTg(&http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Request:    r,
			}, nil
		}),
	})
	if err != nil {
		t.Fatal(err)
	}

	errorTests := []struct {
		title string
		tu    TgUpdate
		err   error
	}{
		{
			title: "empty message",
			tu:    TgUpdate{},
			err:   errTelegram,
		},
		{
			title: "random command",
			tu:    TgUpdate{Message: TgMessage{Text: "/rand"}},
			err:   errTelegram,
		},
		{
			title: "normal command. gpt no choices",
			tu:    TgUpdate{Message: TgMessage{Text: "/ask Hello"}},
			err:   errGpt,
		},
	}

	t.Run("errors", func(t *testing.T) {
		gpt, err := NewGpt(&http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{}`)),
					Request:    r,
				}, nil
			}),
		})
		if err != nil {
			t.Fatal(err)
		}

		for _, v := range errorTests {
			t.Run(v.title, func(t *testing.T) {
				err := RunMain(ctx, v.tu, gpt, tg)
				if !errors.Is(err, v.err) {
					t.Fatal("not telegram message")
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		gpt, err := NewGpt(&http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				body := `{"choices": [{ "message": { "role": "assistant", "content": "Hello"}}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    r,
				}, nil
			}),
		})
		if err != nil {
			t.Fatal(err)
		}

		err = RunMain(ctx, TgUpdate{Message: TgMessage{Text: "/ask Hello"}}, gpt, tg)
		if err != nil {
			t.Fatal(err)
		}
	})
}

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}
