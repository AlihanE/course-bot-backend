package main

import (
	"backend/asset"
	asset_history "backend/asset-history"
	"backend/client"
	clientConversation "backend/client-conversation"
	"backend/db"
	"backend/files"
	"backend/reports"
	"backend/step"
	userCourse "backend/user-course"
	"backend/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

var (
	logger = hclog.New(&hclog.LoggerOptions{
		Name:  "bot-backend",
		Level: hclog.LevelFromString("DEBUG"),
	})
	e             = echo.New()
	reportService *reports.Service
	clientService *client.Service
	clientConv *clientConversation.Service
	assetHistoryService *asset_history.Service
)

func init() {

	dbConnString := "postgres://postgres:example@localhost:5432/postgres?sslmode=disable"

	if !utils.IsTestEnvironment() {
		dbConnString = "postgres://postgres:179b3a5b-a59a-4ca4-856d-b54a926c1378@localhost:5433/postgres?sslmode=disable"
	}

	m, err := migrate.New(
		"file://migrations",
		dbConnString)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	logger.Info("Initializing db...")
	dbConn, err := db.Connect(logger.Named("db-connect"), "pgx", dbConnString)
	if err != nil {
		logger.Error("failed to connect to db", err)
		os.Exit(1)
	}

	logger.Info("Initializing client store...")
	store := client.NewStore(logger.Named("client-store"), dbConn)

	logger.Info("Initializing client asset...")
	assetStore := asset.NewStore(logger.Named("asset-store"), dbConn)

	logger.Info("Initializing client report...")
	reportStore := reports.NewStore(logger.Named("report-store"), dbConn)

	logger.Info("Initializing client step...")
	stepStore := step.NewStore(logger.Named("step-store"), dbConn)

	logger.Info("Initializing client user course...")
	userCourseStore := userCourse.NewStore(logger.Named("user-course-store"), dbConn)

	logger.Info("Initializing client conversation store...")
	clientConvStore := clientConversation.NewStore(logger.Named("client-conversation-store"), dbConn)

	logger.Info("Initializing asset history store...")
	assetHistory := asset_history.NewStore(logger.Named("asset-history-store"), dbConn)

	logger.Info("Initializing client service...")
	clientService = client.NewService(logger.Named("client-service"), store, e)

	logger.Info("Initializing file service...")
	fileService, err := files.NewService(logger.Named("file-service"), e)
	if err != nil {
		logger.Error("failed initialize file service", err)
		os.Exit(1)
	}

	logger.Info("Initializing asset service...")
	asset.NewService(logger.Named("asset-service"), assetStore, fileService, e)

	logger.Info("Initializing report service...")
	reportService = reports.NewService(logger.Named("report-service"), reportStore, e)

	logger.Info("Initializing step service...")
	step.NewService(logger.Named("step-service"), stepStore, e)

	logger.Info("Initializing user-course service...")
	userCourse.NewService(logger.Named("user-course-service"), userCourseStore, e)

	logger.Info("Initializing client conversation service...")
	clientConv = clientConversation.NewService(logger.Named("client-conversation-service"), clientConvStore, e)

	logger.Info("Initializing asset history service...")
	assetHistoryService = asset_history.NewService(logger.Named("asset-history-service"), assetHistory, e)
}
