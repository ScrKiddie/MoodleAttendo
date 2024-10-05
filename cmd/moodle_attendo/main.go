package main

import (
	"context"
	"log"
	"moodle_attendo/internal/initialize"
	"moodle_attendo/internal/model"
	"net/http"
	"os"
	"time"
)

func main() {
	hostname := os.Getenv("HOSTNAME")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	tgChat := os.Getenv("TGCHAT")
	tgBot := os.Getenv("TGBOT")

	if hostname == "" {
		log.Fatal("environment variable HOSTNAME belum diatur")
	}
	if username == "" {
		log.Fatal("environment variable USERNAME belum diatur")
	}
	if password == "" {
		log.Fatal("environment variable PASSWORD belum diatur")
	}
	if tgChat == "" {
		log.Fatal("environment variable TGCHAT belum diatur")
	}
	if tgBot == "" {
		log.Fatal("environment variable TGBOT belum diatur")
	}

	if len(os.Args) < 2 {
		log.Fatal("id course harus diberikan sebagai argumen")
	}

	client := http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	courseID := os.Args[1]
	initialize.App(ctx, client, courseID, model.AccountModel{Hostname: hostname, Username: username, Password: password, BotToken: tgBot, ChatId: tgChat})
}
