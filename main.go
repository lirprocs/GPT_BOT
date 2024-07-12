package main

import (
	"GPT_BOT/config"
	"GPT_BOT/database"
	"GPT_BOT/lama"
	model "GPT_BOT/user"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

//const (
//	storagePath = "./database/users.db"
//	botToken    = "1060049923:AAFK_hdwWz63hy_LPFCZtpCjqJwm0qm7rrA"
//)

var user *model.User

func main() {
	conf := config.New()

	db, err := database.InitDB(conf.StoragePath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	user = model.NewUser(db)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(conf.BotToken, opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, helloHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/deletecontext", bot.MatchTypeExact, clearHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/profile", bot.MatchTypeExact, profileHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/model", bot.MatchTypePrefix, changeModelHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/api", bot.MatchTypePrefix, changeApiHandler)
	b.Start(ctx)
}

func helloHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func clearHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	lama.CleanMessage()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Контекст успешно удалён",
	})
}

func profileHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func changeModelHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func changeApiHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
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
