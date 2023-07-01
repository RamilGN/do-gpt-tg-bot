package main

import (
	"context"
	"log"
)

func Main(ctx context.Context, t TgUpdate) {
	err := ValidateTgUpdate(&t)
	if err != nil {
		log.Println(err)
		return
	}

	gpt, err := NewGpt(nil)
	if err != nil {
		log.Println(err)
		return
	}

	text, err := gpt.Prompt(ctx, t.Message.Text)
	if err != nil {
		log.Println(err)
		return
	}

	tg, err := NewTg(nil)
	if err != nil {
		log.Println(err)
		return
	}

	err = tg.SendMessage(ctx, TgSendMessage{ChatID: t.Message.Chat.ID, Text: text, ReplyToMessageID: t.Message.MessageID})
	if err != nil {
		log.Println(err)
		return
	}
}
