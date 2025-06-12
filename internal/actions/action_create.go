package actions

import (
	"context"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	"github.com/SibHelly/TgUrlShorter/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateShortURLHandler(urlShorter service.UrlShorter) bot.ActionFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

		return nil
	}
}
