package configuration

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type BotConfiguration struct {
	DiscordBotToken string
}

var Global BotConfiguration = load()

// Load the current configuration
func load() BotConfiguration {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return BotConfiguration {
		DiscordBotToken: os.Getenv("DISCORD_BOT_TOKEN"),
	}
}

