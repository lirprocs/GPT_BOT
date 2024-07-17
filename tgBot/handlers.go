package handlers

import (
	"GPT_BOT/lama"
	model "GPT_BOT/user"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	invalidFormatMessage = "Пожалуйста, предоставьте данные в правильном формате."
)

type Handlers struct {
	userService *model.User
}

func NewHandlers(userService *model.User) *Handlers {
	return &Handlers{
		userService: userService,
	}
}

func (h *Handlers) HelloHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userName := bot.EscapeMarkdown(update.Message.From.FirstName)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Привет, *" + userName + "*",
		ParseMode: models.ParseModeMarkdown,
	})

	err := h.userService.AddUser(update)
	if err != nil {
		log.Println("Ошибка при добавлении пользователя:", err)
	}
}

func (h *Handlers) ClearHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	lama.CleanMessage()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Контекст успешно удалён",
	})
}

func (h *Handlers) ProfileHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	u, err := h.userService.GetUser(update)
	if err != nil {
		log.Println("Ошибка при получении пользователя:", err)
		return
	}

	text := fmt.Sprintf(
		"ID: %d\nИмя: %s\nБаланс: %d\nКоличество запросов: %d/%d\nAPI ключ: %s\nМодель: %s\n",
		u.Id, u.Name, u.Balance, u.NumLama, lama.MaxRequest, u.APIKey, u.ModelID,
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}

func (h *Handlers) ChangeModelHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   invalidFormatMessage,
		})
		return
	}

	modelName := parts[1]
	err := h.userService.UpdateUser(update, "model", modelName)
	if err != nil {
		log.Println("Ошибка при обновлении имени модели:", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Имя модели обновлено.",
	})
}

func (h *Handlers) ChangeApiHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   invalidFormatMessage,
		})
		return
	}

	apiKey := parts[1]
	err := h.userService.UpdateUser(update, "api", apiKey)
	if err != nil {
		log.Println("Ошибка при обновлении API ключа:", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "API ключ обновлен.",
	})
}

func (h *Handlers) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	u, err := h.userService.GetUser(update)
	if err != nil {
		log.Println("Ошибка при получении пользователя:", err)
		return
	}

	lama.SetAPIKey(u.APIKey)
	text := lama.ProcessMessage(update.Message.Text)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}
