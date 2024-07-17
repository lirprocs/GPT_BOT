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

var user *model.User

func HelloHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Hello, *" + bot.EscapeMarkdown(update.Message.From.FirstName) + "*",
		ParseMode: models.ParseModeMarkdown,
	})

	err := user.AddUser(update)
	if err != nil {
		log.Println(err)
	}
}

func ClearHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	lama.CleanMessage()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Контекст успешно удалён",
	})
}

func ProfileHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	u, err := user.GetUser(update)
	if err != nil {
		log.Println(err)
	}
	text := "ID: " + fmt.Sprintf("%d", u.Id) + "\n" +
		"Name: " + fmt.Sprintf("%s", u.Name) + "\n" +
		"Баланс: " + fmt.Sprintf("%d", u.Balance) + "\n" +
		"Количество запросов: " + fmt.Sprintf("%d", u.NumLama) + "/" + fmt.Sprintf("%d", lama.MaxRequest) + "\n" +
		"API key: " + fmt.Sprintf("%s", u.APIKey) + "\n" +
		"Model: " + fmt.Sprintf("%s", u.ModelID) + "\n"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}

func ChangeModelHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пожалуйста, предоставьте имя модели в формате /model <имя модели>",
		})
		return
	}

	infoName := parts[0]
	modelName := parts[1]

	err := user.UpdateUser(update, infoName, modelName)
	if err != nil {
		log.Println(err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Имя модели обновлено.",
	})
}

func ChangeApiHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пожалуйста, предоставьте API ключ в формате /api <ключ>",
		})
		return
	}

	infoName := parts[0]
	apiKey := parts[1]

	err := user.UpdateUser(update, infoName, apiKey)
	if err != nil {
		log.Println(err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "API ключ обновлен.",
	})
}

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	u, err := user.GetUser(update)
	if err != nil {
		log.Println(err)
	}

	lama.SetAPIKey(u.APIKey)

	text := lama.ProcessMessage(update.Message.Text)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}
