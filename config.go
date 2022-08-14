package main

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	ApiRoot    string
	BotToken   string
	IsApiLocal bool
	ServerAddr string
}

func loadFromEnv() *Config {
	config := &Config{}

	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok || len(botToken) == 0 {
		log.Fatal("BOT_TOKEN environment variable is missing")
	}
	config.BotToken = botToken

	if apiRoot := os.Getenv("API_ROOT"); len(apiRoot) != 0 {
		config.ApiRoot = apiRoot
	} else {
		config.ApiRoot = "https://api.telegram.org"
	}

	if isApiLocal := os.Getenv("API_LOCAL"); isApiLocal == "1" {
		config.IsApiLocal = true
	} else {
		config.IsApiLocal = false
	}

	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	if len(serverPort) == 0 {
		serverPort = "8080"
	}
	config.ServerAddr = strings.Join([]string{serverHost, serverPort}, ":")

	return config
}
