package bot

import (
	"github.com/BlackRRR/notion-monitoring/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GlobalBot struct {
	Update tgbotapi.UpdatesChannel
}

func BotInit(channel tgbotapi.UpdatesChannel) *GlobalBot {
	bot := &GlobalBot{Update: channel}
	return bot
}

func ActionWithUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go CheckUpdate(bot, &update)
	}
}

func CheckUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if update.Message == nil && update.CallbackQuery == nil {
		return
	}

	if update.Message != nil {
		if update.Message.Chat.ID != 0 {
			model.AdminSettings.AdminID[update.Message.Chat.ID] = &model.AdminUser{Name: update.Message.From.FirstName}
			model.UploadAdminSettings()
		} else {
			model.AdminSettings.AdminID[update.Message.From.ID] = &model.AdminUser{Name: update.Message.From.FirstName}
			model.UploadAdminSettings()

			message := CreateNewTGMessage("Вы добавлены в notion alerts", update.Message.From.ID)
			SendMsgBot(bot, message)
		}
	}
}
