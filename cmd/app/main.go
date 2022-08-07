package main

import (
	"analytic-service/config"
	"analytic-service/internal/logger"
	"analytic-service/internal/server"
	"fmt"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// @title Analytic Service API
// @version 1.0
// @description mts teta team21 analytic-service
func main() {
	appConfig, err := config.Init()
	if err != nil {
		log.Fatalf("Ошибка при инициализации конфига: %v", err)
	}
	log, err := logger.New(appConfig)
	if err != nil {
		zap.NewExample().Fatal(fmt.Sprintf("Ошибка при инициализации логгера: %v", err))
	}
	log.Info("LOGGER&CONFIG SUCCESS")

	app := server.New(appConfig, log, grpc.WithBlock())
	defer app.Dispose() //TODO: graceful shutdown
	errChan := app.Run()
	if err := <-errChan; err != nil {
		log.Fatal(fmt.Sprintf("Ошибка при запуске приложения: %v", err.Error()))
	}

	log.Info("Запущено...")

}
