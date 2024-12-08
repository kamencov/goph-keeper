package auth

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNewServiceAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serv := NewServiceAuth(nil, nil, 0, nil)

	if serv == nil {
		t.Errorf("ServiceAuth is nil")
	}
}
