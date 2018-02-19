package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
			msg.Text = "Please send a photo as a file"
			bot.Send(msg)
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		photo := *update.Message.Document
		url, _ := bot.GetFileDirectURL(photo.FileID)
		log.Println(url)

		tmpfile, err := ioutil.TempFile("", "template")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		n, err := io.Copy(tmpfile, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(n)

		msg.ReplyToMessageID = update.Message.MessageID
		client.SetImage(tmpfile.Name())
		text, _ := client.Text()
		msg.Text = text
		bot.Send(msg)
	}
}
