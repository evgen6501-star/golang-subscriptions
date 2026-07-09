package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/evgen6501-star/golang-subscriptions/internal/config"
	"github.com/evgen6501-star/golang-subscriptions/internal/handler"
	"github.com/evgen6501-star/golang-subscriptions/internal/logger"
	"github.com/evgen6501-star/golang-subscriptions/internal/middleware"
	"github.com/evgen6501-star/golang-subscriptions/internal/repository"
	"github.com/evgen6501-star/golang-subscriptions/internal/server"
	"github.com/evgen6501-star/golang-subscriptions/internal/service"
	"go.uber.org/zap"
)

// cmd/main.go

// @title           Subscription Service API
// @version         1.0
// @description      Сервис для управления подписками пользователей
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com
// @servers  http://localhost:8080/api/v1
// @host      localhost:8080
// @BasePath  /api/v1

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("failed load config:", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	log, err := logger.NewLogger(config.NewLogConfigMust())
	if err != nil {
		fmt.Println("failed init logger", err)
		os.Exit(1)

	}
	defer log.Close()
	log.Debug("Starting Service")
	dbConn, err := repository.NewPostgresConnection(cfg)
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	log.Debug("Connected to PostgreSQL")
	subRepo := repository.NewSubscriptionRepoPostgres(dbConn, log)
	subservice := service.NewSubscriptionService(subRepo, log)
	subHTTP := handler.NewSubscriptionHandler(subservice, log)
	subRouters := subHTTP.Routes()
	apiVersRouter := server.NewApiVersionRouter(server.ApiVers1)
	apiVersRouter.RegisterRoutes(subRouters...)
	httpServer := server.NewHttpServer(config.NewServConfigMust(), log, middleware.RequestID(), middleware.Logger(log))
	httpServer.RegisterAPIRouters(apiVersRouter)
	if err := httpServer.Run(ctx); err != nil {
		log.Error("Http server run error", zap.Error(err))
	}
}
