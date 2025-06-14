package callbacks

import (
	"context"
	"log"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	"github.com/SibHelly/TgUrlShorter/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateAliasCallback() bot.CallbackFunc {
	return func(ctx context.Context, bot *bot.Bot, callback *tgbotapi.CallbackQuery) error {
		answerCallback := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		}
		if _, err := bot.Api.Request(answerCallback); err != nil {
			log.Printf("[ERROR] failed to answer callback: %v", err)
		}

		bot.UserSession[callback.From.ID] = &models.ShortenRequest{
			Step:        "url",
			SkipClicked: false,
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🔗 *Create Short URL*\n\n"+
			"*Step 1/6:* Please send me the URL you want to shorten.\n\n"+
			"Example: `https://example.com/very/long/url`")
		msg.ParseMode = tgbotapi.ModeMarkdown
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel_shorten"),
			),
		)
		msg.ReplyMarkup = keyboard
		if _, err := bot.Api.Send(msg); err != nil {
			log.Printf("Failed to send message about start creating alias in callback: %v", err)
		}
		return nil
	}
}
