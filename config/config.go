package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Config struct {
	StoragePath string
	BotToken    string
	GROQAPIKey  string
	ModelID     string
}

func New() *Config {
	return &Config{
		StoragePath: os.Getenv("STORAGE_PATH"),
		BotToken:    os.Getenv("BOT_TOKEN"),
		GROQAPIKey:  os.Getenv("GROQ_API_KEY"),
		ModelID:     os.Getenv("MODEL_ID"),
	}
}
