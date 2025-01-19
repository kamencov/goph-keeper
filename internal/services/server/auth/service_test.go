package auth

import (
	"github.com/golang/mock/gomock"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestNewServiceAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug}))

	db := NewMockstorageAuth(ctrl)

	serv := NewServiceAuth([]byte("salt"), []byte("salt"), time.Hour, log, db)

	if serv == nil {
		t.Errorf("ServiceAuth is not nil")
	}
}
