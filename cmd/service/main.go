package main

import (
	"goph-keeper/internal/api/service"
	"log/slog"
	"os"
)

//@title GophKeeper API
//@version 1.0
//@description API for GophKeeper service

//@host localhost:8081
//@BasePath /api/v1

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	// создаем логгер
	log := setupLogger(envLocal)

	// запускаем приложение
	if err := service.Run(log); err != nil {
		panic(err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo}))

	}

	return log
}
