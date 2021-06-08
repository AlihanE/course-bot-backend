package main

import (
	"backend/auth"
	"backend/ourbot"
	"backend/ourbot/receiver"
	"backend/ourbot/sender"
	"backend/utils"
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"sync"
)

func main() {
	
	token := ""
	if !utils.IsTestEnvironment() {
		token = ""
	}

	bot, err := botApi.NewBotAPI(token)
	if err != nil {
		logger.Error("bot init failed", err)
	}

	bot.Debug = true
	var wg sync.WaitGroup

	botSender := sender.New(bot, logger.Named("bot-sender"))
	botReceiver := receiver.New(bot, logger.Named("bot-receiver"), reportService, clientService, botSender)
	err = botReceiver.Start(&wg)
	if err != nil {
		logger.Error("failed to start bot receiver", err)
		os.Exit(1)
	}

	auth.New(logger.Named("auth-service"), e)

	ourbot.New(clientService, botSender, logger.Named("bot-service"), e, clientConv, assetHistoryService)

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.Logger.Fatal(e.Start(":8000"))
}
