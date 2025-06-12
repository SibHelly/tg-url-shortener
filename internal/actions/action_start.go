package actions

import (
	"context"
	"log"
	"strings"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetStartUrlShorterHandler создает обработчик для команды /start
func StartUrlShorter() bot.ActionFunc {
	return func(ctx context.Context, bot *bot.Bot, update tgbotapi.Update) error {
		var builder strings.Builder

		// Fixed Markdown syntax - each * must have a matching closing *
		builder.WriteString("👋 *Welcome to URL Shortener Bot!*\n\n")
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
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("❓ Help", "help"),
			),
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = keyboard

		if _, err := bot.Api.Send(msg); err != nil {
			log.Printf("Failed to send URL message for start command: %v", err)
			return err
		}
		return nil
	}
}
