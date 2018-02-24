package helpers

import (
	"gopkg.in/telegram-bot-api.v4"
)

func CreateMessage(update tgbotapi.Update, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.Text = text
	return msg
}

func SetUpdate(offset int, timeout int) tgbotapi.UpdateConfig {
	u := tgbotapi.NewUpdate(offset)
	u.Timeout = timeout
	return u
}
