package main

import (
	"GPT_BOT/config"
	"GPT_BOT/database"
	handlers "GPT_BOT/tgBot"
	model "GPT_BOT/user"
	"context"
	"github.com/go-telegram/bot"
	"log"
	"os"
	"os/signal"
)

var user *model.User

func main() {
	conf := config.New()

	db, err := database.InitDB(conf.StoragePath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	user = model.NewUser(db)
	handler := handlers.NewHandlers(user)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Handler),
	}

	b, err := bot.New(conf.BotToken, opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handler.HelloHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/deletecontext", bot.MatchTypeExact, handler.ClearHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/profile", bot.MatchTypeExact, handler.ProfileHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/model", bot.MatchTypePrefix, handler.ChangeModelHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/api", bot.MatchTypePrefix, handler.ChangeApiHandler)
	b.Start(ctx)
}
