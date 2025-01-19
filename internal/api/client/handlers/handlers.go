package handlers

import (
	"context"
	"fmt"
	"log/slog"
)

var (
	ErrNotEmpty = fmt.Errorf("empty data entered")
)

// serviceCredentials - интерфейс на сервисный слой Credentials.
//go:generate mockgen -source=handlers.go -destination=handlers_mock.go -package=handlers
type serviceCredentials interface {
	SaveLoginAndPassword(ctx context.Context, token, resource, login, password string) error
}

// serviceTextData - интерфейс на сервисный слой TextData.
type serviceTextData interface {
	SaveTextData(ctx context.Context, token, data string) error
}

// serviceBinaryData - интерфейс на сервисный слой BinaryData.
type serviceBinaryData interface {
	SaveBinaryData(ctx context.Context, token, data string) error
}

// serviceCards - интерфейс на сервисный слой Cards.
type serviceCards interface {
	SaveCards(ctx context.Context, token, data string) error
}


// Handler - обработчик запросов.
type Handler struct {
	log                *slog.Logger
	serviceCredentials serviceCredentials
	serviceTextData    serviceTextData
	serviceBinaryData  serviceBinaryData
	serviceCards       serviceCards
}

// NewHandlers - конструктор обработчика.
func NewHandlers(log *slog.Logger,
	serviceCredentials serviceCredentials,
	serviceTextData serviceTextData,
	serviceBinaryData serviceBinaryData,
	serviceCards serviceCards) *Handler {
	return &Handler{
		log:                log,
		serviceCredentials: serviceCredentials,
		serviceTextData:    serviceTextData,
		serviceBinaryData:  serviceBinaryData,
		serviceCards:       serviceCards,
	}
}


// PostLoginAndPassword - сохраняет логин и пароль.
func (h *Handler) PostLoginAndPassword(ctx context.Context, token, resource, login, password string) error {
	if resource == "" || login == "" || password == "" {
		h.log.Error("resource or login or password is empty")
		return ErrNotEmpty
	}

	err := h.serviceCredentials.SaveLoginAndPassword(ctx, token, resource, login, password)
	if err != nil {
		h.log.Error("failed to handlers login and password", "error", err)
		return err
	}

	return nil
}


// PostTextData - сохраняет текстовые данные.
func (h *Handler) PostTextData(ctx context.Context, token, data string) error {

	if data == "" {
		fmt.Println("data is empty")
		return ErrNotEmpty
	}

	err := h.serviceTextData.SaveTextData(ctx, token, data)
	if err != nil {
		return err
	}

	return nil
}


// PostBinaryData - сохраняет бинарные данные.
func (h *Handler) PostBinaryData(ctx context.Context, token, data string) error {

	if data == "" {
		fmt.Println("data is empty")
		return ErrNotEmpty
	}

	err := h.serviceBinaryData.SaveBinaryData(ctx, token, data)
	if err != nil {
		return err
	}

	return nil
}


// PostCards - сохраняет карты.
func (h *Handler) PostCards(ctx context.Context, token, data string) error {

	if data == "" {
		fmt.Println("data is empty")
		return ErrNotEmpty
	}

	err := h.serviceCards.SaveCards(ctx, token, data)
	if err != nil {
		return err
	}

	return nil
}
