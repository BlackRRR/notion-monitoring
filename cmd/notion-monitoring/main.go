package main

import (
	"context"
	"database/sql"
	bot2 "github.com/BlackRRR/notion-monitoring/internal/app/bot"
	client2 "github.com/BlackRRR/notion-monitoring/internal/app/client"
	"github.com/BlackRRR/notion-monitoring/internal/app/repository"
	"github.com/BlackRRR/notion-monitoring/internal/app/service"
	"github.com/BlackRRR/notion-monitoring/internal/cfg"
	"github.com/BlackRRR/notion-monitoring/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	model.DownloadAdminSettings()

	config, err := cfg.NewConfig()
	if err != nil {
		log.Println(err)
	}

	bot, update := startBot(config.TGConfig)
	startHandlers(bot, update)

	dbConn, err := sql.Open("mysql", config.DBConn)
	if err != nil {
		return
	}

	rep, err := repository.NewRepository(ctx, dbConn)
	if err != nil {
		log.Println(err)
	}

	newService := service.NewService(rep, bot)

	client := client2.NewClient(newService, config.NotionSecretKey)

	err = client.StartPages(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for range time.Tick(time.Second * 30) {
			err := client.StartClient(ctx)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	log.Println("Client Started")

	sig := <-subscribeToSystemSignals()

	log.Printf("shutdown all process on '%s' system signal\n", sig.String())
}

func startBot(token string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)

	channel := bot.GetUpdatesChan(u)

	globalBot := bot2.BotInit(channel)

	log.Println("The bot is running")

	return bot, globalBot.Update
}

func startHandlers(bot *tgbotapi.BotAPI, update tgbotapi.UpdatesChannel) {
	go func() {
		bot2.ActionWithUpdates(bot, update)
	}()

	log.Println("bot handler is running")
}

func subscribeToSystemSignals() chan os.Signal {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	return ch
}
