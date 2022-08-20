package main

import (
	"context"
	bot2 "github.com/BlackRRR/notion-monitoring/internal/app/bot"
	client2 "github.com/BlackRRR/notion-monitoring/internal/app/client"
	"github.com/BlackRRR/notion-monitoring/internal/app/repository"
	"github.com/BlackRRR/notion-monitoring/internal/app/service"
	"github.com/BlackRRR/notion-monitoring/internal/cfg"
	"github.com/BlackRRR/notion-monitoring/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	model.DownloadAdminSettings()

	bot, update := startBot()
	startHandlers(bot, update)

	token, err := os.ReadFile("./config/token.txt")
	if err != nil {
		log.Println(err)
	}

	config, err := cfg.NewConfig()
	if err != nil {
		log.Println(err)
	}

	dbConn, err := pgxpool.ConnectConfig(ctx, config.PGConfig)
	if err != nil {
		log.Fatalf("failed to init postgres %s", err.Error())
	}

	rep, err := repository.NewRepository(ctx, dbConn)
	if err != nil {
		log.Println(err)
	}

	newService := service.NewService(rep, bot)

	client := client2.NewClient(newService, string(token))

	err = client.StartPages(ctx)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)

	go func() {
		for range time.Tick(time.Second * 5) {
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

func startBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	file, err := os.ReadFile("./config/TG_token.txt")
	if err != nil {
		return nil, nil
	}

	bot, err := tgbotapi.NewBotAPI(string(file))
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
