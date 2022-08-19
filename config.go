package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ApiRoot    string
	BotToken   string
	IsApiLocal bool
	ServerAddr string
}

func newConfig() *Config {
	var config *Config

	config, ok := loadConfigFromArgs()
	if !ok {
		config = loadConfigFromEnv()
	}

	if len(config.BotToken) == 0 {
		log.Fatal("Bot token is missing")
	}

	return config
}

func loadConfigFromArgs() (*Config, bool) {
	var (
		serverHost string
		serverPort int
	)

	config := &Config{}

	flag.StringVar(&config.BotToken, "bot-token", "", "bot token")
	flag.StringVar(&config.ApiRoot, "api-root", "https://api.telegram.org", "Bot API root")
	flag.BoolVar(&config.IsApiLocal, "api-local", false, "allow providing files from the file system")
	flag.StringVar(&serverHost, "server-host", "", "server host")
	flag.IntVar(&serverPort, "server-port", 8080, "server port")

	flag.Parse()

	config.ServerAddr = strings.Join([]string{serverHost, strconv.Itoa(serverPort)}, ":")

	return config, flag.NFlag() > 0
}

func loadConfigFromEnv() *Config {
	config := &Config{}

	config.BotToken = os.Getenv("BOT_TOKEN")

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
