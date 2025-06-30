package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	config := newConfig()
	cache := newCache()
	telegramClient := NewClient(config.ApiRoot, config.BotToken)

	server, err := NewServer(config, cache, telegramClient)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handleFileRequest)

	c := cors.Default()
	if config.CorsAllowedOrigin != "" {
		c = cors.New(cors.Options{
			AllowedOrigins: []string{config.CorsAllowedOrigin},
		})
	}
	handler := c.Handler(mux)

	log.Printf("Server starting on %s", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
