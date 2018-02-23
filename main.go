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
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = "Please send a picture as a file (wihout compression)"
			bot.Send(msg)
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if update.Message.Document != nil {
			photo := *update.Message.Document
			url, _ := bot.GetFileDirectURL(photo.FileID)

			file := helpers.GetFile(url)
			defer os.Remove(file.Name())

			msg.ReplyToMessageID = update.Message.MessageID
			client.SetImage(file.Name())
			text, _ := client.Text()
			msg.Text = text
		} else {
			msg.Text = "Please send a picture as a file (wihout compression)"
		}
		bot.Send(msg)
	}
}
