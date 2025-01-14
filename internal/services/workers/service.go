package workers

import (
	"context"
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
