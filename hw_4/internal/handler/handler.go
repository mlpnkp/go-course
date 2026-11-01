package handler

import "go-course/hw_4/internal/model"

type NotificationService interface {
	Send(message string) (model.Notification, error)
}

type SendRequest struct {
	Message string
}

type Handler struct {
	notifications NotificationService
}

func New(notifications NotificationService) *Handler {
	return &Handler{notifications: notifications}
}

func (h *Handler) Send(request SendRequest) (model.Notification, error) {
	return h.notifications.Send(request.Message)
}
