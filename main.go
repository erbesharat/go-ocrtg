package main

import (
	"log"
	"os"

	"github.com/erbesharat/go-ocrtg/helpers"
	"github.com/joho/godotenv"
	"github.com/otiai10/gosseract"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := gosseract.NewClient()
	defer client.Close()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := helpers.SetUpdate(0, 60)

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			msg := helpers.CreateMessage(update, "Please send a picture as a file (wihout compression)")
			bot.Send(msg)
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Document != nil {
			photo := *update.Message.Document
			url, _ := bot.GetFileDirectURL(photo.FileID)

			file := helpers.GetFile(url)
			defer os.Remove(file.Name())

			client.SetImage(file.Name())
			text, _ := client.Text()

			msg := helpers.CreateMessage(update, text)
			bot.Send(msg)
		} else {
			msg := helpers.CreateMessage(update, "Please send a picture as a file (wihout compression)")
			bot.Send(msg)
		}
	}
}
