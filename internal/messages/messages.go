package messages

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SibHelly/TgUrlShorter/internal/bot"
	"github.com/SibHelly/TgUrlShorter/internal/models"
	"github.com/SibHelly/TgUrlShorter/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleURLStep() bot.MessageFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgbotapi.Update) error {
		// Basic URL validation
		if !strings.HasPrefix(update.Message.Text, "http://") && !strings.HasPrefix(update.Message.Text, "https://") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Invalid URL format. Please provide a valid URL starting with http:// or https://")
			bot.Api.Send(msg)
			return nil
		}

		bot.UserSession[update.Message.From.ID].URL = update.Message.Text
		bot.UserSession[update.Message.From.ID].Step = "alias"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"✅ URL saved!\n\n"+
				"*Step 2/6:* Now provide a custom alias for your short URL.\n\n"+
				"Example: `my-cool-link` (will become: short.ly/my-cool-link)\n\n"+
				"Requirements:\n"+
				"• Only letters, numbers, hyphens, and underscores\n"+
				"• 1-50 characters long")
		msg.ParseMode = tgbotapi.ModeMarkdown

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel_shorten"),
			),
		)
		msg.ReplyMarkup = keyboard

		_, err := bot.Api.Send(msg)
		return err
	}

}

func HandleAliasStep() bot.MessageFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgbotapi.Update) error {
		if len(update.Message.Text) < 1 || len(update.Message.Text) > 50 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Alias must be between 1-50 characters long. Please try again.")
			bot.Api.Send(msg)
			return nil
		}

		// Check for valid characters (letters, numbers, hyphens, underscores)
		for _, char := range update.Message.Text {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '-' || char == '_') {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					"❌ Alias can only contain letters, numbers, hyphens (-), and underscores (_). Please try again.")
				bot.Api.Send(msg)
				return nil
			}
		}

		bot.UserSession[update.Message.From.ID].Alias = update.Message.Text
		bot.UserSession[update.Message.From.ID].Step = "max_visits"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"✅ Alias saved!\n\n"+
				"*Step 3/6:* Set maximum number of visits (optional)\n\n"+
				"After this many clicks, the link will become inactive.\n"+
				"Example: `100` for 100 visits limit")
		msg.ParseMode = tgbotapi.ModeMarkdown

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⏭ Skip", "skip_max_visits"),
				tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel_shorten"),
			),
		)
		msg.ReplyMarkup = keyboard

		_, err := bot.Api.Send(msg)
		return err
	}
}

func HandleMaxVisitsStep() bot.MessageFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgbotapi.Update) error {
		visits, err := strconv.Atoi(update.Message.Text)

		if err != nil || visits <= 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Please enter a valid positive number or use the Skip button.")
			bot.Api.Send(msg)
			return nil
		}

		bot.UserSession[update.Message.From.ID].MaxVisits = visits
		bot.UserSession[update.Message.From.ID].Step = "expires_at"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"✅ Visit limit saved!\n\n"+
				"*Step 4/6:* Set expiration date (optional)\n\n"+
				"After this date, the link will become inactive.\n"+
				"Format: `YYYY-MM-DD` or `DD.MM.YYYY`\n"+
				"Example: `2024-12-31` or `31.12.2024`")
		msg.ParseMode = tgbotapi.ModeMarkdown

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⏭ Skip", "skip_expires_at"),
				tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel_shorten"),
			),
		)
		msg.ReplyMarkup = keyboard

		_, err = bot.Api.Send(msg)
		return err
	}
}

func HandleExpiresAtStep() bot.MessageFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgbotapi.Update) error {
		var expiresAt time.Time
		var err error

		// Try different date formats
		formats := []string{"2006-01-02", "02.01.2006", "2006/01/02"}

		for _, format := range formats {
			expiresAt, err = time.Parse(format, update.Message.Text)
			if err == nil {
				break
			}
		}

		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Invalid date format. Please use YYYY-MM-DD or DD.MM.YYYY format, or use the Skip button.")
			bot.Api.Send(msg)
			return nil
		}

		// Check if date is in the future
		if expiresAt.Before(time.Now()) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Expiration date must be in the future. Please try again.")
			bot.Api.Send(msg)
			return nil
		}

		bot.UserSession[update.Message.From.ID].ExpiresAt = &expiresAt
		bot.UserSession[update.Message.From.ID].Step = "title"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"✅ Expiration date saved!\n\n"+
				"*Step 5/6:* Add a title (optional)\n\n"+
				"This will be shown when someone hovers over your link.\n"+
				"Example: `My Awesome Website`")
		msg.ParseMode = tgbotapi.ModeMarkdown

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⏭ Skip", "skip_title"),
				tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel_shorten"),
			),
		)
		msg.ReplyMarkup = keyboard

		_, err = bot.Api.Send(msg)
		return err
	}
}

func HandleTitleStep() bot.MessageFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgbotapi.Update) error {
		if len(update.Message.Text) > 200 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Title is too long (max 200 characters). Please try again.")
			bot.Api.Send(msg)
			return nil
		}

		bot.UserSession[update.Message.From.ID].Title = update.Message.Text
		bot.UserSession[update.Message.From.ID].Step = "description"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"✅ Title saved!\n\n"+
				"*Step 6/6:* Add a description (optional)\n\n"+
				"This provides additional context about your link.\n"+
				"Example: `Check out this amazing article about technology`")
		msg.ParseMode = tgbotapi.ModeMarkdown

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⏭ Skip", "skip_description"),
				tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel_shorten"),
			),
		)
		msg.ReplyMarkup = keyboard

		_, err := bot.Api.Send(msg)
		return err
	}
}

func HandleDescriptionStep(urlShorter service.UrlShorter) bot.MessageFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgbotapi.Update) error {
		if len(update.Message.Text) > 500 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Description is too long (max 500 characters). Please try again.")
			bot.Api.Send(msg)
			return nil
		}

		bot.UserSession[update.Message.From.ID].Description = update.Message.Text

		err := urlShorter.Create(models.Url{
			Original_url: bot.UserSession[update.Message.From.ID].URL,
			Alias:        bot.UserSession[update.Message.From.ID].Alias,
			Visit_count:  bot.UserSession[update.Message.From.ID].MaxVisits,
			Expires_at:   bot.UserSession[update.Message.From.ID].ExpiresAt,
			Title:        bot.UserSession[update.Message.From.ID].Title,
			Description:  bot.UserSession[update.Message.From.ID].Description,
		})
		if err != nil {
			log.Printf("Failed to get URLs: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Failed to create shorten, err: %v", err))
			_, err := bot.Api.Send(msg)
			return err
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"✅ Your alias for url added!\n\n"+
				fmt.Sprintf("*Usage*\nOld url:%s\nNew url:%s\n\n", bot.UserSession[update.Message.From.ID].URL, bot.UserSession[update.Message.From.ID].Alias))
		msg.ParseMode = tgbotapi.ModeMarkdown
		delete(bot.UserSession, update.Message.From.ID)

		_, err = bot.Api.Send(msg)
		return err

	}
}
