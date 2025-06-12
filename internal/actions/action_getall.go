package actions

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	"github.com/SibHelly/TgUrlShorter/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetMyURLsHandler создает обработчик для команды /myurls .
func GetMyURLsHandler(urlShorter service.UrlShorter) bot.ActionFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		// Получаем список URL
		urls, err := urlShorter.GetAll()
		if err != nil {
			log.Printf("Failed to get URLs: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to get your URLs. Please try again later.")
			_, err := bot.Send(msg)
			return err
		}

		// Если список пуст
		if len(urls) == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You don't have any shortened URLs yet.")
			_, err := bot.Send(msg)
			return err
		}

		for _, url := range urls {
			var builder strings.Builder
			builder.WriteString(fmt.Sprintf("ID: %d\n", url.Id))
			builder.WriteString(fmt.Sprintf("🔗 *%s*\n", url.Alias))
			builder.WriteString(fmt.Sprintf("Original: %s\n", url.Original_url))
			builder.WriteString(fmt.Sprintf("Visits: %d | Created: %s\n", url.Visit_count, url.Created_at.Format("2006-01-02")))

			if url.Title != "" {
				builder.WriteString(fmt.Sprintf("Title: %s\n", url.Title))
			}

			deleteButton := tgbotapi.NewInlineKeyboardButtonData("🗑️ Delete", fmt.Sprintf("delete_%s", url.Alias))
			infoButton := tgbotapi.NewInlineKeyboardButtonData("Info", fmt.Sprintf("info_%s", url.Alias))
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(deleteButton, infoButton),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
			msg.ParseMode = tgbotapi.ModeMarkdown
			msg.ReplyMarkup = keyboard

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send URL message: %v", err)
				return err
			}
		}

		return nil
	}
}
