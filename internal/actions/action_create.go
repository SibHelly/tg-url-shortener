package actions

import (
	"context"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	"github.com/SibHelly/TgUrlShorter/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CreateShorten создает обработчик для команды /shorten .
func CreateShortURLHandler() bot.ActionFunc {
	return func(ctx context.Context, bot *bot.Bot, update tgbotapi.Update) error {
		bot.UserSession[update.Message.From.ID] = &models.ShortenRequest{
			Step: "url",
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"🔗 *Create Short URL*\n\n"+
				"*Step 1/6:* Please send me the URL you want to shorten.\n\n"+
				"Example: `https://example.com/very/long/url`")
		msg.ParseMode = tgbotapi.ModeMarkdown

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel_shorten"),
			),
		)
		msg.ReplyMarkup = keyboard

		_, err := bot.Api.Send(msg)
		return err
	}
}
