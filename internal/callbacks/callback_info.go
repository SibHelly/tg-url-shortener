package callbacks

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	"github.com/SibHelly/TgUrlShorter/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetAllInfoUrlCallback(urlShorter service.UrlShorter) bot.CallbackFunc {
	return func(ctx context.Context, bot *bot.Bot, callback *tgbotapi.CallbackQuery) error {
		// Сначала отвечаем на callback query, чтобы убрать выделение кнопки
		answerCallback := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		}
		if _, err := bot.Api.Request(answerCallback); err != nil {
			log.Printf("[ERROR] failed to answer callback: %v", err)
		}

		alias := strings.TrimPrefix(callback.Data, "info_")

		answer, err := urlShorter.Info(alias)
		if err != nil {
			log.Printf("Failed to get info about URL %s: %v", alias, err)

			// Отправляем сообщение об ошибке
			errorMsg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Failed to get info about URL")
			if _, err := bot.Api.Send(errorMsg); err != nil {
				log.Printf("Failed to send error message: %v", err)
			}
			return err
		}

		var builder strings.Builder
		// builder.WriteString(fmt.Sprintf("ID: %d\n", answer.Id))
		builder.WriteString(fmt.Sprintf("🔗 *%s*\n", answer.Alias))
		builder.WriteString(fmt.Sprintf("Visits remained: %d | Created: %s\n", answer.Visit_count, answer.Created_at.Format("2006-01-02")))
		if answer.Expires_at != nil {
			builder.WriteString(fmt.Sprintf("Will be valid until %s\n", answer.Expires_at.Format("2006-01-02")))
		}
		if answer.Title != "" {
			builder.WriteString(fmt.Sprintf("Title: %s\n", answer.Title))
		}
		if answer.Description != "" {
			builder.WriteString(fmt.Sprintf("Descriptions: %s\n", answer.Description))
		}
		if len(answer.Visits) > 0 {
			builder.WriteString("History visits:\n")
			for i, visit := range answer.Visits {
				builder.WriteString(fmt.Sprintf("Visit %d: %s\n", i+1, visit.Created_at.Format("2006-01-02 15:04:05")))
			}
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, builder.String())
		msg.ParseMode = "Markdown"
		if _, err := bot.Api.Send(msg); err != nil {
			log.Printf("Failed to send URL message: %v", err)
		}

		return nil
	}
}
