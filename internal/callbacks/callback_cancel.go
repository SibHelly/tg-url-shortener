package callbacks

import (
	"context"
	"log"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// DeleteURLCallback handles delete button clicks
func CancelCreateCallback() bot.CallbackFunc {
	return func(ctx context.Context, bot *bot.Bot, callback *tgbotapi.CallbackQuery) error {
		answerCallback := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		}
		if _, err := bot.Api.Request(answerCallback); err != nil {
			log.Printf("[ERROR] failed to answer callback: %v", err)
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Your canceled creating alias for url!\n\n")
		msg.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Api.Send(msg); err != nil {
			log.Printf("Failed to send message about canceled creating alias: %v", err)
		}
		delete(bot.UserSession, callback.From.ID)
		return nil
	}
}
