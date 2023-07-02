package main

import (
	"context"
	"log"
	"regexp"
)

func Main(ctx context.Context, t TgUpdate) {
	gpt, err := NewGpt(nil)
	if err != nil {
		log.Println(err)
		return
	}

	tg, err := NewTg(nil)
	if err != nil {
		log.Println(err)
		return
	}

	if err = RunMain(ctx, t, gpt, tg); err != nil {
		log.Println(err)
		return
	}
}

func RunMain(ctx context.Context, t TgUpdate, gpt *Gpt, tg *Tg) error {
	err := ValidateTgUpdate(&t)
	if err != nil {
		return err
	}

	promptText := regexp.MustCompile("^/ask.* ").ReplaceAllString(t.Message.Text, "")

	text, err := gpt.Prompt(ctx, promptText)
	if err != nil {
		return err
	}

	err = tg.SendMessage(ctx, TgSendMessage{ChatID: t.Message.Chat.ID, Text: text, ReplyToMessageID: t.Message.MessageID})
	if err != nil {
		return err
	}

	return nil
}
