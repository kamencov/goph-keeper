package text_data

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

type service interface {
	SaveTextData(ctx context.Context, data string) error
}

type Handlers struct {
	pd.UnimplementedPostTextDataServer
	log     *slog.Logger
	service service
}

func NewHandlers(log *slog.Logger, service service) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

func (h *Handlers) PostTextData(ctx context.Context, in *pd.PostTextDataRequest) (*pd.Empty, error) {
	if in.Data == "" {
		h.log.Error("data is empty")
		return nil, status.Errorf(codes.InvalidArgument, "data is empty")
	}

	err := h.service.SaveTextData(ctx, in.GetData())

	if err != nil {
		h.log.Error("failed to save login and password", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to save login and password")
	}

	return &pd.Empty{
		Message: "save complete",
	}, nil
}
