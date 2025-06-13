package middleware

import (
	"context"
	"log"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NotSessionOnlyAction(next bot.ActionFunc) bot.ActionFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgbotapi.Update) error {
		if _, ok := bot.UserSession[update.Message.From.ID]; !ok {
			return next(ctx, bot, update)
		}

		if _, err := bot.Api.Send(tgbotapi.NewMessage(
			update.FromChat().ID,
			"Finish creating the alias before executing other commands.",
		)); err != nil {
			return err
		}
		return nil
	}
}

func NotSessionOnlyCallback(next bot.CallbackFunc) bot.CallbackFunc {
	return func(ctx context.Context, bot *bot.Bot, callback *tgbotapi.CallbackQuery) error {
		if _, ok := bot.UserSession[callback.From.ID]; !ok {
			return next(ctx, bot, callback)
		}

		answerCallback := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		}
		if _, err := bot.Api.Request(answerCallback); err != nil {
			log.Printf("[ERROR] failed to answer callback: %v", err)
		}

		if _, err := bot.Api.Send(tgbotapi.NewMessage(
			callback.Message.Chat.ID,
			"Finish creating the alias before executing other commands.",
		)); err != nil {
			return err
		}
		return nil
	}
}

func SessionOnlyCallback(next bot.CallbackFunc) bot.CallbackFunc {
	return func(ctx context.Context, bot *bot.Bot, callback *tgbotapi.CallbackQuery) error {

		if _, ok := bot.UserSession[callback.From.ID]; ok {
			return next(ctx, bot, callback)
		}
		answerCallback := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		}
		if _, err := bot.Api.Request(answerCallback); err != nil {
			log.Printf("[ERROR] failed to answer callback: %v", err)
		}
		if _, err := bot.Api.Send(tgbotapi.NewMessage(
			callback.Message.Chat.ID,
			"Start creating alias for use this buttons",
		)); err != nil {
			return err
		}
		return nil
	}
}
