package service

import (
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	handlerAuth "goph-keeper/internal/grpc/auth"
	handlerBinaryData "goph-keeper/internal/grpc/binary_data"
	handlerCards "goph-keeper/internal/grpc/cards"
	"goph-keeper/internal/grpc/health"
	handlerRegister "goph-keeper/internal/grpc/register"
	"goph-keeper/internal/grpc/sync"
	handlerTextData "goph-keeper/internal/grpc/text_data"
	"goph-keeper/internal/middleware/auth"
	pd "goph-keeper/internal/proto/v1"
	serviceAuth "goph-keeper/internal/services/server/auth"
	binaryData "goph-keeper/internal/services/server/binary_data"
	"goph-keeper/internal/services/server/cards"
	"goph-keeper/internal/services/server/credentials"
	textData "goph-keeper/internal/services/server/text_data"
	"goph-keeper/internal/storage/postgresql"
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

	// Инициализация подключения к базе данных
	db, err := postgresql.NewPostgresql(log)
	if err != nil {
		log.Error("failed to initialize connection to database", "error", err)
		return err
	}

	defer func(db *postgresql.Postgresql) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Создаем сервисы
	newServiceAuth := serviceAuth.NewServiceAuth(
		flags.TokenSalt,
		flags.PasswordSalt,
		24*time.Hour,
		log,
		db,
	)
	newServiceCredentials := credentials.NewService(log, db)
	newServiceTextData := textData.NewService(log, db)
	newServiceBinaryData := binaryData.NewService(log, db)
	newServiceCards := cards.NewServiceCards(log, db)

	// инициализируем проверку авторизацию.
	authChecker := auth.NewMiddleware(log, newServiceAuth)

	// Создаем grpc
	registerUser := handlerRegister.NewHandlers(log, newServiceAuth)
	authUser := handlerAuth.NewHandlers(log, newServiceAuth)
	//postCredentials := handlerCredentials.NewHandlers(log, newServiceCredentials)
	postTextData := handlerTextData.NewHandlers(log, newServiceTextData)
	postBinaryData := handlerBinaryData.NewHandlers(log, newServiceBinaryData)
	postCards := handlerCards.NewHandlers(log, newServiceCards)
	healthStatus := health.NewHandler(log)
	newSync := sync.NewHandler(log, newServiceCredentials, newServiceTextData, newServiceBinaryData, newServiceCards)

	// Создаем GRPC-сервер
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authChecker.UnaryInterceptor))

	// Регистрируем goph-keeper в GRPC-сервере
	pd.RegisterRegisterServer(grpcServer, registerUser)
	pd.RegisterAuthServer(grpcServer, authUser)
	//pd.RegisterPostCredentialsServer(grpcServer, postCredentials)
	pd.RegisterPostTextDataServer(grpcServer, postTextData)
	pd.RegisterPostBinaryDataServer(grpcServer, postBinaryData)
	pd.RegisterPostCardsServer(grpcServer, postCards)
	pd.RegisterHealthServer(grpcServer, healthStatus)
	pd.RegisterSyncFromClientServer(grpcServer, newSync)

	// канал для ошибки
	errChan := make(chan error, 1)
	go func() {
		listener, err := net.Listen("tcp", flags.AddrGRPC)
		log.Info("application start", "addr:", flags.AddrGRPC)
		if err != nil {
			log.Error("failed to listen", "error", err)
			errChan <- err
			return
		}
		log.Info("application run")
		if err := grpcServer.Serve(listener); err != nil {
			slog.Error("failed to serve", "error", err)
			errChan <- err
			return
		}

		return
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop:
		log.Info("received stop signal, shutting down")
	case err := <-errChan:
		log.Error("server encountered an error", "error", err)
		return err
	}

	grpcServer.GracefulStop()

	log.Info("application stop")

	return nil
}
