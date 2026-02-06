package config

import "os"

type Config struct {
	DiscordWebhookURL string
	TwitchClientID    string
	TwitchClientSecret string
}

func Load() Config {
	return Config{
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		TwitchClientID:    os.Getenv("TWITCH_CLIENT_ID"),
		TwitchClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
	}
}
