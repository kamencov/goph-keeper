package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"goph-keeper/internal/api/client/cli"
	"goph-keeper/internal/api/client/handlers"
	auth2 "goph-keeper/internal/api/client/repositories/auth"
	"goph-keeper/internal/api/client/repositories/health"
	"goph-keeper/internal/api/client/repositories/sync_client"
	"goph-keeper/internal/services/client/auth_client"
	"goph-keeper/internal/services/client/binary_data_client"
	"goph-keeper/internal/services/client/cards_client"
	"goph-keeper/internal/services/client/credentials_client"
	"goph-keeper/internal/services/client/get_and_deleted_data"
	"goph-keeper/internal/services/client/text_data_client"
	workers2 "goph-keeper/internal/services/workers"
	"goph-keeper/internal/storage/sqlite"
	"goph-keeper/internal/workers"
	"log/slog"
	"os"
)

const defaultHost = "localhost:8081"

// RunClient - запускает приложение client
func RunClient() error {
	// Создаем или открываем файл
	file, err := os.Create("logfile.log")
	if err != nil {
		return err
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

	defer func(db *sqlite.Storage) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Инициализируем сервисы
	newServiceAuth := auth_client.NewService(log, db)
	newServiceCredentials := credentials_client.NewService(log, db)
	newServiceTextData := text_data_client.NewService(log, db)
	newServiceBinaryData := binary_data_client.NewService(log, db)
	newServiceCard := cards_client.NewService(log, db)
	newServiceData := get_and_deleted_data.NewService(log, db)

	repoSync := sync_client.NewHandlers(log)
	newServiceWorker := workers2.NewService(log, db, repoSync)

	conn, err := grpc.Dial(
		defaultHost,
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

	//Инициализация воркера
	newWorker := workers.NewWorker(log, newServiceWorker, conn)

	go newWorker.Run(ctx)

	newAuthHandler := auth2.NewHandlers(log, newServiceAuth)
	newSaveHandler := handlers.NewHandlers(log, newServiceCredentials, newServiceTextData, newServiceBinaryData, newServiceCard)
	newHealth := health.NewHandlers(log)

	// Инициализация интерфейса CLI
	newCLI := cli.NewCLI(log, newAuthHandler, newSaveHandler, newServiceData, newServiceData, newHealth, conn)

	// Запуск интерфейса CLI
	newCLI.RunCLI(ctx)

	return nil
}
