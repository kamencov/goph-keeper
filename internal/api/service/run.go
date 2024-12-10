package service

import (
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	handlerAuth "goph-keeper/internal/grpc/auth"
	handlerBinaryData "goph-keeper/internal/grpc/binary_data"
	handlerCards "goph-keeper/internal/grpc/cards"
	handlerCredentials "goph-keeper/internal/grpc/credentials"
	handlerRegister "goph-keeper/internal/grpc/register"
	handlerTextData "goph-keeper/internal/grpc/text_data"
	pd "goph-keeper/internal/proto/v1"
	serviceAuth "goph-keeper/internal/services/auth"
	binaryData "goph-keeper/internal/services/binary_data"
	"goph-keeper/internal/services/cards"
	"goph-keeper/internal/services/credentials"
	textData "goph-keeper/internal/services/text_data"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(log *slog.Logger) error {
	const op = "run.service.app"

	log.With(slog.String("op", op))

	// Парсим флаги
	flags := NewFlags(log)
	flags.Parse()

	// Подключаемся к базе данных
	storage, err := initDB(log, flags)
	if err != nil {
		log.Error("failed to initialize connection to database", "error", err)
		return err
	}

	defer func(storage repo) {
		err := storage.Close()
		if err != nil {

		}
	}(storage)

	// Создаем сервисы
	newServiceAuth := serviceAuth.NewServiceAuth(
		flags.TokenSalt,
		flags.PasswordSalt,
		24*time.Hour,
		log,
		storage,
	)
	newServiceCredentials := credentials.NewService(log, storage)
	newServiceTextData := textData.NewService(log, storage)
	newServiceBinaryData := binaryData.NewService(log, storage)
	newServiceCards := cards.NewServiceCards(log, storage)

	// Создаем grpc
	registerUser := handlerRegister.NewHandlers(log, newServiceAuth)
	authUser := handlerAuth.NewHandlers(log, newServiceAuth)
	postCredentials := handlerCredentials.NewHandlers(log, newServiceCredentials)
	postTextData := handlerTextData.NewHandlers(log, newServiceTextData)
	postBinaryData := handlerBinaryData.NewHandlers(log, newServiceBinaryData)
	postCards := handlerCards.NewHandlers(log, newServiceCards)

	// Создаем GRPC-сервер
	grpcServer := grpc.NewServer()

	// Регистрируем goph-keeper в GRPC-сервере
	pd.RegisterRegisterServer(grpcServer, registerUser)
	pd.RegisterAuthServer(grpcServer, authUser)
	pd.RegisterPostCredentialsServer(grpcServer, postCredentials)
	pd.RegisterPostTextDataServer(grpcServer, postTextData)
	pd.RegisterPostBinaryDataServer(grpcServer, postBinaryData)
	pd.RegisterPostCardsServer(grpcServer, postCards)

	go func() {
		listener, err := net.Listen("tcp", flags.AddrGRPC)
		if err != nil {
			log.Error("failed to listen", "error", err)
			return
		}
		log.Info("application run")
		if err := grpcServer.Serve(listener); err != nil {
			slog.Error("failed to serve", "error", err)
			return
		}

		return
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	grpcServer.GracefulStop()

	log.Info("application stop")

	return nil
}
