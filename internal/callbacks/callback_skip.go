package callbacks

import (
	"context"
	"fmt"
	"log"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// var steps = []string{"url", "alias", "max_visits", "expires_at", "title", "description"}

// DeleteURLCallback handles delete button clicks
func SkipStepCreateCallback() bot.CallbackFunc {
	return func(ctx context.Context, bot *bot.Bot, callback *tgbotapi.CallbackQuery) error {
		answerCallback := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		}
		if _, err := bot.Api.Request(answerCallback); err != nil {
			log.Printf("[ERROR] failed to answer callback: %v", err)
		}

		bot.UserSession[callback.From.ID].SkipClicked = true

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("Your skiped input step: %s!\n\n", bot.UserSession[callback.From.ID].Step))

		if _, err := bot.Api.Send(msg); err != nil {
			log.Printf("Failed to send skip message: %v", err)
		}
		// Теперь вызываем соответствующий обработчик для следующего шага
		if handler, exists := bot.GetMessageHandler(bot.UserSession[callback.From.ID].Step); exists {
			// Создаем фейковый update для вызова обработчика
			fakeUpdate := &tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: callback.From,
					Chat: callback.Message.Chat,
					Text: "", // Пустой текст, так как шаг пропускается
				},
			}

			if err := handler(ctx, bot, fakeUpdate); err != nil {
				log.Printf("[ERROR] failed to execute next step handler for %s: %v", bot.UserSession[callback.From.ID].Step, err)
				// Очищаем сессию в случае ошибки
				delete(bot.UserSession, callback.From.ID)
				errorMsg := tgbotapi.NewMessage(callback.Message.Chat.ID, "❌ Error occurred, session cleared.")
				if _, err := bot.Api.Send(errorMsg); err != nil {
					log.Printf("[ERROR] failed to send error message: %v", err)
				}
			}
		} else {
			log.Printf("[ERROR] no handler found for step: %s", bot.UserSession[callback.From.ID].Step)
		}
		return nil
	}
}
