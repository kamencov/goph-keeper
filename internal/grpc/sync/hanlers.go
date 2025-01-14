package sync

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1_pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

type serviceCred interface {
	SyncSaveCredentials(ctx context.Context, accessToken, resource, login, password string) error
	SyncDelCredentials(ctx context.Context, accessToken, resource string) error
}

type serviceTextData interface {
	SyncSaveText(ctx context.Context, accessToken, data string) error
	SyncDelText(ctx context.Context, accessToken, resource string) error
}

type serviceBinaryData interface {
	SyncSaveBinary(ctx context.Context, accessToken, data string) error
	SyncDelBinary(ctx context.Context, accessToken, data string) error
}

type serviceCards interface {
	SyncSaveCards(ctx context.Context, accessToken, cards string) error
	SyncDelBinary(ctx context.Context, accessToken, data string) error
}

type Handler struct {
	v1_pd.UnimplementedSyncFromClientServer
	log           *slog.Logger
	serviceCred   serviceCred
	serviceText   serviceTextData
	serviceBinary serviceBinaryData
	serviceCards  serviceCards
}

func NewHandler(log *slog.Logger,
	serviceSync serviceCred,
	serviceText serviceTextData,
	serviceBinary serviceBinaryData,
	serviceCards serviceCards) *Handler {
	return &Handler{
		log:           log,
		serviceCred:   serviceSync,
		serviceText:   serviceText,
		serviceBinary: serviceBinary,
		serviceCards:  serviceCards,
	}
}

func (h *Handler) SyncFromClientCredentials(ctx context.Context, in *v1_pd.SyncFromClientCredentialsRequest) (*v1_pd.Empty, error) {

	for _, credential := range in.Task {
		switch credential.Action {
		case "save":
			if err := h.serviceCred.SyncSaveCredentials(ctx, credential.AccessToken, credential.Resource, credential.Login, credential.Password); err != nil {
				h.log.Error("failed to save login and password", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to save login and password")
			}
		case "deleted":
			if err := h.serviceCred.SyncDelCredentials(ctx, credential.AccessToken, credential.Resource); err != nil {
				h.log.Error("failed to deleted", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to save login and password")
			}
		}

	}

	return &v1_pd.Empty{Message: "completed"}, nil
}

func (h *Handler) SyncFromClientTextData(ctx context.Context, in *v1_pd.SyncFromClientTextDataRequest) (*v1_pd.Empty, error) {
	for _, textData := range in.Task {
		switch textData.Action {
		case "save":
			if err := h.serviceText.SyncSaveText(ctx, textData.AccessToken, textData.Text); err != nil {
				h.log.Error("failed to save text data", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to save login and password")
			}
		case "deleted":
			if err := h.serviceText.SyncDelText(ctx, textData.AccessToken, textData.Text); err != nil {
				h.log.Error("failed to deleted", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to save login and password")
			}
		}
	}

	return &v1_pd.Empty{Message: "completed"}, nil
}

func (h *Handler) SyncFromClientBinaryData(ctx context.Context, in *v1_pd.SyncFromClientBinaryDataRequest) (*v1_pd.Empty, error) {

	for _, binaryData := range in.Task {
		switch binaryData.Action {
		case "save":
			if err := h.serviceBinary.SyncSaveBinary(ctx, binaryData.AccessToken, binaryData.Binary); err != nil {
				h.log.Error("failed to save binary data", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to save login and password")
			}

		case "deleted":
			if err := h.serviceBinary.SyncDelBinary(ctx, binaryData.AccessToken, binaryData.Binary); err != nil {
				h.log.Error("failed to deleted", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to save login and password")
			}
		}
	}

	return &v1_pd.Empty{Message: "completed"}, nil
}

func (h *Handler) SyncFromClientCards(ctx context.Context, in *v1_pd.SyncFromClientCardsRequest) (*v1_pd.Empty, error) {

	for _, card := range in.Task {
		switch card.Action {
		case "save":
			if err := h.serviceCards.SyncSaveCards(ctx, card.AccessToken, card.Cards); err != nil {
				h.log.Error("failed to save cards", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to save login and password")
			}
		case "deleted":
			if err := h.serviceCards.SyncDelBinary(ctx, card.AccessToken, card.Cards); err != nil {
				h.log.Error("failed to deleted", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to save login and password")
			}
		}
	}

	return &v1_pd.Empty{Message: "completed"}, nil
}
