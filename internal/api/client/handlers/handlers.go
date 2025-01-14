package handlers

import (
	"context"
	"fmt"
	"log/slog"
)

var (
	ErrNotEmpty = fmt.Errorf("empty data entered")
)

type serviceCredentials interface {
	SaveLoginAndPassword(ctx context.Context, token, resource, login, password string) error
}
type serviceTextData interface {
	SaveTextData(ctx context.Context, token, data string) error
}
type serviceBinaryData interface {
	SaveBinaryData(ctx context.Context, token, data string) error
}
type serviceCards interface {
	SaveCards(ctx context.Context, token, data string) error
}

type Handler struct {
	log                *slog.Logger
	serviceCredentials serviceCredentials
	serviceTextData    serviceTextData
	serviceBinaryData  serviceBinaryData
	serviceCards       serviceCards
}

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
