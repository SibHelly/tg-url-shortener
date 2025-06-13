package callbacks

import (
	"context"
	"log"
	"strings"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	"github.com/SibHelly/TgUrlShorter/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// DeleteURLCallback handles delete button clicks
func DeleteURLCallback(urlShorter service.UrlShorter) bot.CallbackFunc {
	return func(ctx context.Context, bot *bot.Bot, callback *tgbotapi.CallbackQuery) error {
		alias := strings.TrimPrefix(callback.Data, "delete_")

		// Delete URL
		err := urlShorter.Delete(alias)
		if err != nil {
			log.Printf("Failed to delete URL %s: %v", alias, err)

			//Показываем ошибку
			answerCallback := tgbotapi.CallbackConfig{
				CallbackQueryID: callback.ID,
				Text:            "Failed to delete URL",
			}
			if _, err := bot.Api.Request(answerCallback); err != nil {
				log.Printf("[ERROR] failed to answer callback: %v", err)
			}
			return err
		}

		// Отвечаем на callback.
		answerCallback := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
			Text:            "URL deleted successfully",
		}
		if _, err := bot.Api.Request(answerCallback); err != nil {
			log.Printf("[ERROR] failed to answer callback: %v", err)
		}

		// Удаляем сообщение
		deleteMsg := tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		if _, err := bot.Api.Send(deleteMsg); err != nil {
			log.Printf("[ERROR] failed to delete message: %v", err)
		}

		return nil
	}
}
