package workers

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

type SyncModel struct {
	ID        int
	UserID    int
	TableName string
	TaskID    int
	Action    string
	UpdatedAt time.Time
}

type Credentials struct {
	ID          int
	UserID      int
	Resource    string
	Login       string
	Password    string
	UpdatedAt   time.Time
	DeleteTask  int
	Action      string
	AccessToken string
}

type TextData struct {
	ID          int
	UserID      int
	Text        string
	UpdatedAt   time.Time
	DeleteTask  int
	Action      string
	AccessToken string
}

type BinaryData struct {
	ID          int
	UserID      int
	Binary      string
	UpdatedAt   time.Time
	DeleteTask  int
	Action      string
	AccessToken string
}
type Cards struct {
	ID          int
	UserID      int
	Cards       string
	UpdatedAt   time.Time
	DeleteTask  int
	Action      string
	AccessToken string
}

type Status struct {
	credentials bool
	textData    bool
	binaryData  bool
	cards       bool
}

//go:generate mockgen -source=service.go -destination=service_mock.go -package=workers
type storage interface {
	GetAllSync() ([]*SyncModel, error)
	GetTokenWithUserID(ctx context.Context, userID int) (string, error)
	GetDataCredentials(userID, taskID int) (*Credentials, error)
	GetDataTextData(userID, taskID int) (*TextData, error)
	GetDataBinaryData(userID, taskID int) (*BinaryData, error)
	GetDataCards(userID, taskID int) (*Cards, error)
	ClearSyncCredentials() error
	ClearSyncTextData() error
	ClearSyncBinaryData() error
	ClearSyncCards() error
}

type repoSync interface {
	SyncCredentials(ctx context.Context, conn *grpc.ClientConn, data []*Credentials) error
	SyncTextData(ctx context.Context, conn *grpc.ClientConn, data []*TextData) error
	SyncBinaryData(ctx context.Context, conn *grpc.ClientConn, data []*BinaryData) error
	SyncCards(ctx context.Context, conn *grpc.ClientConn, data []*Cards) error
}

var (
	errLen          = errors.New("len sync = 0")
	errClearSyncErr = errors.New("failed to clear sync")
)

type Service struct {
	log      *slog.Logger
	storage  storage
	repoSync repoSync
}

func NewService(log *slog.Logger, storage storage, repoSync repoSync) *Service {
	return &Service{
		log:      log,
		storage:  storage,
		repoSync: repoSync,
	}
}

func (s *Service) PushData(ctx context.Context, conn *grpc.ClientConn) error {
	var (
		credentials []*Credentials
		textData    []*TextData
		binaryData  []*BinaryData
		cards       []*Cards

		statusBool Status
	)

	// собираем все данные из таблицы sync_client
	syncsData, err := s.storage.GetAllSync()
	if err != nil {
		s.log.Error("failed to get all data from database", "error", err)
		return err
	}

	// проверяем наличие данных
	if len(syncsData) == 0 {
		s.log.Error("workers service", "error", errLen)
		return errLen
	}

	// запускам цикл
	for _, syncData := range syncsData {
		// распределяем по таблицам сбор данных
		switch syncData.TableName {
		case "credentials":
			data, err := s.storage.GetDataCredentials(syncData.UserID, syncData.TaskID)
			if err != nil {
				return err
			}

			token, err := s.storage.GetTokenWithUserID(ctx, syncData.UserID)
			if err != nil {
				return err
			}

			data.AccessToken = token

			// добавляем действие
			data.Action = syncData.Action

			// добавляем в слайс
			credentials = append(credentials, data)

		case "text_data":
			data, err := s.storage.GetDataTextData(syncData.UserID, syncData.TaskID)
			if err != nil {
				return err
			}

			// добавляем действие
			data.Action = syncData.Action

			token, err := s.storage.GetTokenWithUserID(ctx, syncData.UserID)
			if err != nil {
				return err
			}

			data.AccessToken = token
			// добавляем в слайс
			textData = append(textData, data)

		case "binary_data":
			data, err := s.storage.GetDataBinaryData(syncData.UserID, syncData.TaskID)
			if err != nil {
				return err
			}

			// добавляем действие
			data.Action = syncData.Action

			token, err := s.storage.GetTokenWithUserID(ctx, syncData.UserID)
			if err != nil {
				return err
			}

			data.AccessToken = token

			// добавляем в слайс
			binaryData = append(binaryData, data)

		case "cards":
			data, err := s.storage.GetDataCards(syncData.UserID, syncData.TaskID)
			if err != nil {
				return err
			}
			// добавляем действие
			data.Action = syncData.Action

			token, err := s.storage.GetTokenWithUserID(ctx, syncData.UserID)
			if err != nil {
				return err
			}

			data.AccessToken = token

			// добавляем в слайс
			cards = append(cards, data)
		}
	}

	// отправляем данные на сервер
	if err := s.repoSync.SyncCredentials(ctx, conn, credentials); err == nil {
		statusBool.credentials = true
	}
	if err := s.repoSync.SyncTextData(ctx, conn, textData); err == nil {
		statusBool.textData = true
	}
	if err := s.repoSync.SyncBinaryData(ctx, conn, binaryData); err == nil {
		statusBool.binaryData = true
	}
	if err := s.repoSync.SyncCards(ctx, conn, cards); err == nil {
		statusBool.cards = true
	}

	// чистим таблицу sync_client
	if err := s.ClearSync(&statusBool); err != nil {
		return err
	}

	return nil
}

func (s *Service) ClearSync(status *Status) error {
	var errs []error

	if status.credentials {
		if err := s.storage.ClearSyncCredentials(); err != nil {
			s.log.Error("failed to clear sync credentials", "error", err)
			errs = append(errs, fmt.Errorf("failed to clear sync credentials: %w", err))
		}
	}

	if status.textData {
		if err := s.storage.ClearSyncTextData(); err != nil {
			s.log.Error("failed to clear sync text data", "error", err)
			errs = append(errs, fmt.Errorf("failed to clear sync text data: %w", err))
		}
	}

	if status.binaryData {
		if err := s.storage.ClearSyncBinaryData(); err != nil {
			s.log.Error("failed to clear sync binary data", "error", err)
			errs = append(errs, fmt.Errorf("failed to clear sync binary data: %w", err))
		}
	}

	if status.cards {
		if err := s.storage.ClearSyncCards(); err != nil {
			s.log.Error("failed to clear sync cards", "error", err)
			errs = append(errs, fmt.Errorf("failed to clear sync cards: %w", err))
		}
	}

	if len(errs) > 0 {
		// Объединяем ошибки в одну
		return errClearSyncErr
	}

	return nil
}
