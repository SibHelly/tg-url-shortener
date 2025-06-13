package bot

import (
	"context"
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/SibHelly/TgUrlShorter/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Api         *tgbotapi.BotAPI
	actions     map[string]ActionFunc
	callbacks   map[string]CallbackFunc
	handlersMsg map[string]MessageFunc
	UserSession map[int64]*models.ShortenRequest
}

type ActionFunc func(ctx context.Context, bot *Bot, update *tgbotapi.Update) error
type CallbackFunc func(ctx context.Context, bot *Bot, callback *tgbotapi.CallbackQuery) error
type MessageFunc func(ctx context.Context, bot *Bot, update *tgbotapi.Update) error

func NewBot(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		Api:         api,
		UserSession: make(map[int64]*models.ShortenRequest),
	}
}

func (b *Bot) RegisterAction(nameAction string, action ActionFunc) {
	if b.actions == nil {
		b.actions = make(map[string]ActionFunc)
	}
	b.actions[nameAction] = action
}

func (b *Bot) RegisterCallback(callbackPrefix string, callback CallbackFunc) {
	if b.callbacks == nil {
		b.callbacks = make(map[string]CallbackFunc)
	}
	b.callbacks[callbackPrefix] = callback
}

func (b *Bot) RegisterMessageFunc(step string, input MessageFunc) {
	if b.handlersMsg == nil {
		b.handlersMsg = make(map[string]MessageFunc)
	}
	b.handlersMsg[step] = input
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.Api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Minute)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	// Обработка callback'ов
	if update.CallbackQuery != nil {
		b.handleCallback(ctx, update.CallbackQuery)
		return
	}

	// Проверяем, что это сообщение и Message не nil
	if update.Message == nil {
		return
	}

	// Обработка обычных сообщений (не команд)
	if !update.Message.IsCommand() {
		b.handleMessage(ctx, &update)
		return
	}

	// Обработка команд
	cmd := update.Message.Command()
	actionView, ok := b.actions[cmd]
	if !ok {
		// Команда не найдена
		return
	}

	// Выполняем команду (middleware уже встроен в actionView)
	if err := actionView(ctx, b, &update); err != nil {
		log.Printf("[ERROR] failed to execute action: %v", err)

		if _, err := b.Api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error")); err != nil {
			log.Printf("[ERROR] failed to send error message: %v", err)
		}
	}
}

func (b *Bot) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered in callback: %v\n%s", p, string(debug.Stack()))
		}
	}()

	for prefix, handler := range b.callbacks {
		if strings.HasPrefix(callback.Data, prefix) {
			if err := handler(ctx, b, callback); err != nil {
				log.Printf("[ERROR] failed to execute callback: %v", err)

				// Send error message to user
				answerCallback := tgbotapi.NewCallback(callback.ID, "An error occurred")
				if _, err := b.Api.Request(answerCallback); err != nil {
					log.Printf("[ERROR] failed to answer callback: %v", err)
				}
			}
			return
		}
	}

	// No handler found
	log.Printf("[WARNING] no handler found for callback: %s", callback.Data)
	answerCallback := tgbotapi.NewCallback(callback.ID, "Unknown action")
	if _, err := b.Api.Request(answerCallback); err != nil {
		log.Printf("[ERROR] failed to answer callback: %v", err)
	}
}

func (b *Bot) handleMessage(ctx context.Context, update *tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered in msg: %v\n%s", p, string(debug.Stack()))
		}
	}()

	// Проверяем, есть ли у пользователя активная сессия (шаг ввода)
	if currentStep, ok := b.UserSession[update.Message.From.ID]; ok {
		// Если есть активная сессия, ищем обработчик для текущего шага
		if handler, exists := b.handlersMsg[currentStep.Step]; exists {
			// Создаем фейковый update для передачи в обработчик

			if err := handler(ctx, b, update); err != nil {
				log.Printf("[ERROR] failed to execute message handler for step %s: %v", currentStep.Step, err)
				// В случае ошибки очищаем сессию
				delete(b.UserSession, update.Message.From.ID)
				if _, err := b.Api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error, session cleared")); err != nil {
					log.Printf("[ERROR] failed to send error message: %v", err)
				}
			}
			return
		}
	}

	// Если нет активной сессии или команды, просто логируем сообщение
	log.Printf("Received message from %d: %s", update.Message.Chat.ID, update.Message.Text)
}

func (b *Bot) GetMessageHandler(step string) (MessageFunc, bool) {
	handler, exists := b.handlersMsg[step]
	return handler, exists
}

// func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
// 	defer func() {
// 		if p := recover(); p != nil {
// 			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
// 		}
// 	}()

// 	var action ActionFunc

// 	if update.CallbackQuery != nil {
// 		b.handleCallback(ctx, update.CallbackQuery)
// 		return
// 	}

// 	// if _, ok := b.UserSession[update.Message.From.ID]; ok {
// 	// 	b.handleMessage(ctx, &update)
// 	// 	return
// 	// }
// 	if update.CallbackQuery == nil && !update.Message.IsCommand() {
// 		b.handleMessage(ctx, &update)
// 		return
// 	}

// 	cmd := update.Message.Command()

// 	actionView, ok := b.actions[cmd]
// 	if !ok {
// 		return
// 	}

// 	action = actionView

// 	if err := action(ctx, b, &update); err != nil {
// 		log.Printf("[ERROR] failed to execute action: %v", err)

// 		if _, err := b.Api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error")); err != nil {
// 			log.Printf("[ERROR] failed to send error message: %v", err)
// 		}
// 	}

// }
