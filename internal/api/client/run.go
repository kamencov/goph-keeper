package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"goph-keeper/internal/api/client/cli"
	auth2 "goph-keeper/internal/api/client/handlers/auth"
	"goph-keeper/internal/api/client/handlers/save"
	"goph-keeper/internal/services/client/auth_client"
	"goph-keeper/internal/services/client/binary_data_client"
	"goph-keeper/internal/services/client/cards_client"
	"goph-keeper/internal/services/client/credentials_client"
	"goph-keeper/internal/services/client/get_all_data"
	"goph-keeper/internal/services/client/text_data_client"
	"goph-keeper/internal/storage/sqlite"
	"log/slog"
	"os"
)

const host = "localhost:8081"

func RunClient() error {
	// Создаем или открываем файл
	file, err := os.Create("logfile.log")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Настраиваем slog на запись в файл
	log := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Подключение к базе
	db, err := sqlite.NewSqlStorage(log)
	if err != nil {
		return err
	}

	// Инициализируем сервисы
	newServiceAuth := auth_client.NewService(log, db)
	newServiceCredentials := credentials_client.NewService(log, db)
	newServiceTextData := text_data_client.NewService(log, db)
	newServiceBinaryData := binary_data_client.NewService(log, db)
	newServiceCard := cards_client.NewService(log, db)
	newServiceGet := get_all_data.NewService(log, db)

	conn, err := grpc.Dial(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect client", "error", err)
		return err
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	//Инициализация воркеров пока заглушен
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newAuthHandler := auth2.NewHandlers(log, newServiceAuth)
	newSaveHandler := save.NewHandlers(log, newServiceCredentials, newServiceTextData, newServiceBinaryData, newServiceCard)

	// Инициализация интерфейса CLI
	newCLI := cli.NewCLI(log, newAuthHandler, newSaveHandler, newServiceGet, conn)

	// Запуск интерфейса CLI

	newCLI.RunCLI(ctx)

	// Инициализация воркера
	//newWorker := workers.NewWorker(nil)

	//go newWorker.Run(ctx)
	return nil
}
