package service

import (
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	handlerAuth "goph-keeper/internal/grpc/auth"
	handlerRegister "goph-keeper/internal/grpc/register"
	pd "goph-keeper/internal/proto/v1"
	serviceAuth "goph-keeper/internal/services/auth"
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
		log.Error("failed to initialize connection to database", err)
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

	// Создаем grpc
	registerUser := handlerRegister.NewHandlers(log, newServiceAuth)
	authUser := handlerAuth.NewHandlers(log, newServiceAuth)

	// Создаем GRPC-сервер
	grpcServer := grpc.NewServer()

	// Регистрируем goph-keeper в GRPC-сервере
	pd.RegisterRegisterServer(grpcServer, registerUser)
	pd.RegisterAuthServer(grpcServer, authUser)

	go func() {
		listener, err := net.Listen("tcp", flags.AddrGRPC)
		if err != nil {
			log.Error("failed to listen", err)
			return
		}
		log.Info("application run")
		if err := grpcServer.Serve(listener); err != nil {
			slog.Error("failed to serve", err)
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
