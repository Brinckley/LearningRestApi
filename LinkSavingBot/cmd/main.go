package main

import (
	event_consumer "MusicBot/consumer/event-consumer"
	"MusicBot/pkg/clients/telegram"
	telegram2 "MusicBot/pkg/events/telegram"
	"MusicBot/pkg/storage"
	"errors"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = ""
	batchSize   = 100
)

func main() {
	eventProcessor := telegram2.NewTgProcessor(
		telegram.NewClient(tgBotHost, mustToken()),
		storage.NewFileStorage(storagePath),
	)

	log.Println("service started....")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalln("service stopped : ", err)
	}
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"check token for telegram bot access",
	)

	if *token == "" {
		log.Fatal(errors.New("no token found"))
		return ""
	}

	return *token
}
