package client

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"goph-keeper/internal/api/client/cli"
	auth2 "goph-keeper/internal/api/client/handlers/auth"
	"goph-keeper/internal/api/client/handlers/save"
	"goph-keeper/internal/services/client/auth_client"
	"goph-keeper/internal/services/client/binary_data_client"
	"goph-keeper/internal/services/client/cards_client"
	"goph-keeper/internal/services/client/credentials_client"
	"goph-keeper/internal/services/client/text_data_client"
	"goph-keeper/internal/storage/sqlite"
	"log/slog"
	"os"
	"strings"
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
	newCLI := cli.NewCLI(log, newAuthHandler, newSaveHandler, conn)

	// Запуск интерфейса CLI

	newCLI.RunCLI(ctx)

	//runClient(ctx, log, newAuthHandler, newSaveHandler, conn)

	// Инициализация воркера
	//newWorker := workers.NewWorker(nil)

	//go newWorker.Run(ctx)
	return nil
}

func runClient(ctx context.Context, log *slog.Logger, auth *auth2.Handlers, save *save.Handler, conn *grpc.ClientConn) {
	fmt.Println("Welcome to goph_keeper")
	fmt.Println("Please enter your command")
	fmt.Println("1 - Register")
	fmt.Println("2 - Authorize")
	fmt.Println("0 - Exit")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nEnter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			err := register(ctx, log, auth, conn)
			if err != nil {
				log.Error("failed to register user", "error", err)
				fmt.Println("Failed to register user")
				fmt.Println("Please enter your command")
				fmt.Println("1 - Register")
				fmt.Println("2 - Authorize")
				fmt.Println("0 - Exit")
			}
		case "2":
			err := authUser(ctx, log, auth, conn)
			if err != nil {
				log.Error("failed to register user", "error", err)
				fmt.Println("Failed to register user")
				fmt.Println("Please enter your command")
				fmt.Println("1 - Register")
				fmt.Println("2 - Authorize")
				fmt.Println("0 - Exit")
			}
		case "0":
			return
		default:
			fmt.Println("Please enter correct command")
		}
	}
}

func register(ctx context.Context, log *slog.Logger, auth *auth2.Handlers, conn *grpc.ClientConn) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--- Register ---")

	fmt.Print("Enter login: ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	err := auth.RegisterUser(ctx, conn, login, password)
	if err != nil {
		fmt.Println("Failed to register user")
		return err
	}
	return nil
}

func authUser(ctx context.Context, log *slog.Logger, auth *auth2.Handlers, conn *grpc.ClientConn) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--- Authorize ---")

	fmt.Print("Enter login: ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	_, err := auth.AuthUser(ctx, conn, login, password)
	if err != nil {
		log.Error("failed to auth user", "error", err)
		fmt.Println("Failed to auth user")
		return err
	}
	return nil
}
