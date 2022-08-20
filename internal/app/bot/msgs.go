package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func CreateNewTGMessage(text string, chatID int64) tgbotapi.Chattable {
	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Text: text,
	}

	return message
}

func SendMsgBot(bot *tgbotapi.BotAPI, msg tgbotapi.Chattable) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
