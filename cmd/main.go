package main

import (
	"context"
	"log"

	"github.com/SibHelly/TgUrlShorter/internal/actions"
	"github.com/SibHelly/TgUrlShorter/internal/bot"
	"github.com/SibHelly/TgUrlShorter/internal/callbacks"
	"github.com/SibHelly/TgUrlShorter/internal/cfg"
	"github.com/SibHelly/TgUrlShorter/internal/messages"
	"github.com/SibHelly/TgUrlShorter/internal/middleware"
	"github.com/SibHelly/TgUrlShorter/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config := cfg.LoadConfig()
	botApi, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Printf("[ERROR failed to create botAPI")
	}
	urlService := service.NewURLService("http://url-shortener:8080")

	myBot := bot.NewBot(botApi)
	// Регситрация обработчки для команды /start
	myBot.RegisterAction("start", middleware.NotSessionOnlyAction(actions.StartUrlShorter()))
	// Регистрация обработчика для команды /myurls
	myBot.RegisterAction("myurls", middleware.NotSessionOnlyAction(actions.GetMyURLsHandler(urlService)))
	myBot.RegisterAction("shorten", middleware.NotSessionOnlyAction(actions.CreateShortURLHandler()))
	// Регистрация обработчки нажатия кнопок
	myBot.RegisterCallback("urls_", middleware.NotSessionOnlyCallback(callbacks.GetMyURLsHandlerCallback(urlService)))
	myBot.RegisterCallback("delete_", middleware.NotSessionOnlyCallback(callbacks.DeleteURLCallback(urlService)))
	myBot.RegisterCallback("info_", middleware.NotSessionOnlyCallback(callbacks.GetAllInfoUrlCallback(urlService)))
	myBot.RegisterCallback("skip_", middleware.SessionOnlyCallback(callbacks.SkipStepCreateCallback()))
	myBot.RegisterCallback("cancel_", middleware.SessionOnlyCallback(callbacks.CancelCreateCallback()))
	myBot.RegisterCallback("create_", middleware.NotSessionOnlyCallback(callbacks.CreateAliasCallback()))
	myBot.RegisterCallback("help_", middleware.NotSessionOnlyCallback(callbacks.GetHelpInfoCallback()))
	// Регистрация обработчиков ввода
	myBot.RegisterMessageFunc("url", messages.HandleURLStep())
	myBot.RegisterMessageFunc("alias", messages.HandleAliasStep())
	myBot.RegisterMessageFunc("max_visits", messages.HandleMaxVisitsStep())
	myBot.RegisterMessageFunc("expires_at", messages.HandleExpiresAtStep())
	myBot.RegisterMessageFunc("title", messages.HandleTitleStep())
	myBot.RegisterMessageFunc("description", messages.HandleDescriptionStep(urlService))
	// Запуск бота
	if err := myBot.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
