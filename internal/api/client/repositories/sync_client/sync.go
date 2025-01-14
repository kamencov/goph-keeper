//package sync_client
//
//import (
//	"context"
//	"google.golang.org/grpc"
//	v1_pd "goph-keeper/internal/proto/v1"
//	"goph-keeper/internal/services/workers"
//	"log/slog"
//	"time"
//)
//
//type Handler struct {
//	log *slog.Logger
//}
//
//func NewHandlers(log *slog.Logger) *Handler {
//	return &Handler{
//		log: log,
//	}
//}
//
//func (h *Handler) SyncSaveCredentials(ctx context.Context, conn *grpc.ClientConn, data []*workers.Credentials) error {
//	// создаем клиента для регистрации
//	regClient := v1_pd.NewSyncFromClientClient(conn)
//
//	// Преобразуем `workers.Credentials` в `v1_pd.Credentials`
//	var task []*v1_pd.Credentials
//	for _, cred := range data {
//		task = append(task, &v1_pd.Credentials{
//			Id:          int32(cred.ID),
//			IdUser:      int32(cred.UserID),
//			Resource:    cred.Resource,
//			Login:       cred.Login,
//			Password:    cred.Password,
//			UpdatedAt:   cred.UpdatedAt.Format(time.RFC3339),
//			Action:      cred.Action,
//			AccessToken: cred.AccessToken,
//		})
//	}
//
//	req := &v1_pd.SyncFromClientCredentialsRequest{Task: task}
//
//	// записываем токент в контекст
//	ctx = context.WithValue(ctx, "authorization", task[0].AccessToken)
//
//	_, err := regClient.SyncFromClientCredentials(ctx, req)
//	if err != nil {
//		h.log.Error("failed to sync credentials", "error", err)
//		return err
//	}
//	return nil
//}
//
//func (h *Handler) SyncTextData(ctx context.Context, conn *grpc.ClientConn, data []*workers.TextData) error {
//	// создаем клиента для регистрации
//	regClient := v1_pd.NewSyncFromClientClient(conn)
//
//	// Преобразуем `workers.TextData` в `v1_pd.TextData`
//	var task []*v1_pd.TextData
//	for _, cred := range data {
//		task = append(task, &v1_pd.TextData{
//			Id:        int32(cred.ID),
//			IdUser:    int32(cred.UserID),
//			Text:      cred.Text,
//			UpdatedAt: cred.UpdatedAt.Format(time.RFC3339),
//			Action:    cred.Action,
//		})
//	}
//
//	req := &v1_pd.SyncFromClientTextDataRequest{Task: task}
//
//	_, err := regClient.SyncFromClientTextData(ctx, req)
//	if err != nil {
//		h.log.Error("failed to sync text data", "error", err)
//		return err
//	}
//
//	return nil
//}
//
//func (h *Handler) SyncBinaryData(ctx context.Context, conn *grpc.ClientConn, data []*workers.BinaryData) error {
//	// создаем клиента для регистрации
//	regClient := v1_pd.NewSyncFromClientClient(conn)
//
//	// Преобразуем `workers.BinaryData` в `v1_pd.BinaryData`
//	var task []*v1_pd.BinaryData
//	for _, cred := range data {
//		task = append(task, &v1_pd.BinaryData{
//			Id:        int32(cred.ID),
//			IdUser:    int32(cred.UserID),
//			Binary:    cred.Binary,
//			UpdatedAt: cred.UpdatedAt.Format(time.RFC3339),
//			Action:    cred.Action,
//		})
//	}
//
//	req := &v1_pd.SyncFromClientBinaryDataRequest{Task: task}
//
//	_, err := regClient.SyncFromClientBinaryData(ctx, req)
//	if err != nil {
//		h.log.Error("failed to sync binary data", "error", err)
//		return err
//	}
//
//	return nil
//}
//
//func (h *Handler) SyncSaveCards(ctx context.Context, conn *grpc.ClientConn, data []*workers.Cards) error {
//	// создаем клиента для регистрации
//	regClient := v1_pd.NewSyncFromClientClient(conn)
//
//	// Преобразуем `workers.Cards` в `v1_pd.Cards`
//	var task []*v1_pd.Cards
//	for _, cred := range data {
//		task = append(task, &v1_pd.Cards{
//			Id:        int32(cred.ID),
//			IdUser:    int32(cred.UserID),
//			Cards:     cred.Cards,
//			UpdatedAt: cred.UpdatedAt.Format(time.RFC3339),
//			Action:    cred.Action,
//		})
//	}
//
//	req := &v1_pd.SyncFromClientCardsRequest{Task: task}
//
//	_, err := regClient.SyncFromClientCards(ctx, req)
//	if err != nil {
//		h.log.Error("failed to sync cards", "error", err)
//		return err
//	}
//
//	return nil
//}

package sync_client

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	v1_pd "goph-keeper/internal/proto/v1"
	"goph-keeper/internal/services/workers"
	"log/slog"
	"time"
)

type Handler struct {
	log *slog.Logger
}

// NewHandlers создаёт новый Handler
func NewHandlers(log *slog.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

// ключ для контекста авторизации
type contextKey string

const authorizationKey contextKey = "authorization"

var (
	errorNoCredentials = errors.New("no credentials to sync")
	errorNoTextData    = errors.New("no text data to sync")
	errorNoBinaryData  = errors.New("no binary data to sync")
	errorNoCards       = errors.New("no cards to sync")
)

// SyncSaveCredentials синхронизирует учётные данные
func (h *Handler) SyncCredentials(ctx context.Context, conn *grpc.ClientConn, data []*workers.Credentials) error {
	regClient := v1_pd.NewSyncFromClientClient(conn)

	var task []*v1_pd.Credentials
	for _, cred := range data {
		task = append(task, &v1_pd.Credentials{
			Id:          int32(cred.ID),
			IdUser:      int32(cred.UserID),
			Resource:    cred.Resource,
			Login:       cred.Login,
			Password:    cred.Password,
			UpdatedAt:   cred.UpdatedAt.Format(time.RFC3339),
			Action:      cred.Action,
			AccessToken: cred.AccessToken,
		})
	}

	if len(task) == 0 {
		h.log.Warn("no credentials to sync")
		return errorNoCredentials
	}

	req := &v1_pd.SyncFromClientCredentialsRequest{Task: task}

	// записываем токен в контекст
	md := metadata.Pairs(string(authorizationKey), task[0].AccessToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := regClient.SyncFromClientCredentials(ctx, req)
	if err != nil {
		h.log.Error("failed to sync credentials", "error", err)
		return err
	}

	return nil
}

// SyncTextData синхронизирует текстовые данные
func (h *Handler) SyncTextData(ctx context.Context, conn *grpc.ClientConn, data []*workers.TextData) error {
	regClient := v1_pd.NewSyncFromClientClient(conn)

	var task []*v1_pd.TextData
	for _, td := range data {
		task = append(task, &v1_pd.TextData{
			Id:          int32(td.ID),
			IdUser:      int32(td.UserID),
			Text:        td.Text,
			UpdatedAt:   td.UpdatedAt.Format(time.RFC3339),
			Action:      td.Action,
			AccessToken: td.AccessToken,
		})
	}

	if len(task) == 0 {
		h.log.Warn("no text data to sync")
		return errorNoTextData
	}

	req := &v1_pd.SyncFromClientTextDataRequest{Task: task}

	md := metadata.Pairs(string(authorizationKey), task[0].AccessToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := regClient.SyncFromClientTextData(ctx, req)
	if err != nil {
		h.log.Error("failed to sync text data", "error", err)
		return err
	}

	return nil
}

// SyncBinaryData синхронизирует бинарные данные
func (h *Handler) SyncBinaryData(ctx context.Context, conn *grpc.ClientConn, data []*workers.BinaryData) error {
	regClient := v1_pd.NewSyncFromClientClient(conn)

	var task []*v1_pd.BinaryData
	for _, bd := range data {
		task = append(task, &v1_pd.BinaryData{
			Id:          int32(bd.ID),
			IdUser:      int32(bd.UserID),
			Binary:      bd.Binary,
			UpdatedAt:   bd.UpdatedAt.Format(time.RFC3339),
			Action:      bd.Action,
			AccessToken: bd.AccessToken,
		})
	}

	if len(task) == 0 {
		h.log.Warn("no binary data to sync")
		return errorNoBinaryData
	}

	req := &v1_pd.SyncFromClientBinaryDataRequest{Task: task}

	md := metadata.Pairs(string(authorizationKey), task[0].AccessToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := regClient.SyncFromClientBinaryData(ctx, req)
	if err != nil {
		h.log.Error("failed to sync binary data", "error", err)
		return err
	}

	return nil
}

// SyncCards синхронизирует карты
func (h *Handler) SyncCards(ctx context.Context, conn *grpc.ClientConn, data []*workers.Cards) error {
	regClient := v1_pd.NewSyncFromClientClient(conn)

	var task []*v1_pd.Cards
	for _, card := range data {
		task = append(task, &v1_pd.Cards{
			Id:          int32(card.ID),
			IdUser:      int32(card.UserID),
			Cards:       card.Cards,
			UpdatedAt:   card.UpdatedAt.Format(time.RFC3339),
			Action:      card.Action,
			AccessToken: card.AccessToken,
		})
	}

	if len(task) == 0 {
		h.log.Warn("no cards to sync")
		return errorNoCards
	}

	req := &v1_pd.SyncFromClientCardsRequest{Task: task}

	md := metadata.Pairs(string(authorizationKey), task[0].AccessToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := regClient.SyncFromClientCards(ctx, req)
	if err != nil {
		h.log.Error("failed to sync cards", "error", err)
		return err
	}

	return nil
}
