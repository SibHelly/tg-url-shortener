package callbacks

import (
	"context"
	"log"
	"strings"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetHelpInfoCallback() bot.CallbackFunc {
	return func(ctx context.Context, bot *bot.Bot, callback *tgbotapi.CallbackQuery) error {
		answerCallback := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		}
		if _, err := bot.Api.Request(answerCallback); err != nil {
			log.Printf("[ERROR] failed to answer callback: %v", err)
		}

		var builder strings.Builder

		builder.WriteString("🔗 *What I can do:*\n")
		builder.WriteString("• Create beautiful custom aliases\n")
		builder.WriteString("• Track click counts and analytics\n")
		builder.WriteString("• Show statistics for your links\n\n")
		builder.WriteString("📝 *How to use:*\n")
		builder.WriteString("• Send me any URL and I'll shorten it instantly\n")
		builder.WriteString("• Use /shorten command to create links with custom settings\n")
		builder.WriteString("• Browse your links through the menu below\n\n")
		builder.WriteString("⚡ *Quick start:* Just send me a link!")

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔗 Create Link", "create_url"),
				tgbotapi.NewInlineKeyboardButtonData("📋 My Links", "urls_myurls"),
			),
			// tgbotapi.NewInlineKeyboardRow(
			// 	tgbotapi.NewInlineKeyboardButtonData("❓ Help", "help_urlshortener"),
			// ),
		)

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, builder.String())
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = keyboard

		if _, err := bot.Api.Send(msg); err != nil {
			log.Printf("Failed to send help info message for help callback: %v", err)
			return err
		}
		return nil
	}
}
