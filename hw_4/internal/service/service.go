package service

import (
	"errors"
	"fmt"

	"go-course/hw_4/internal/model"
)

var ErrPanic = errors.New("notifier panicked")

type Repository interface {
	Save(n model.Notification) error
	ByID(id int) (model.Notification, error)
}

type Notifier interface {
	Send(message string) error
}

type Service struct {
	repo      Repository
	notifiers []Notifier
	nextID    int
}

func NewService(repo Repository, notifiers []Notifier) *Service {
	return &Service{repo: repo, notifiers: notifiers, nextID: 1}
}

func (s *Service) Send(message string) (model.Notification, error) {
	n := model.Notification{ID: s.nextID, Message: message, Status: model.StatusSent}
	sendErr := s.dispatch(message)
	if sendErr != nil {
		n.Status = model.StatusFailed
	}
	if err := s.repo.Save(n); err != nil {
		return model.Notification{}, fmt.Errorf("save notification %d: %w", n.ID, err)
	}
	s.nextID++
	return n, sendErr
}

func (s *Service) Notification(id int) (model.Notification, error) {
	return s.repo.ByID(id)
}

func (s *Service) dispatch(message string) error {
	var errs []error
	for _, n := range s.notifiers {
		if err := safeSend(n, message); err != nil {
			errs = append(errs, fmt.Errorf("dispatch %q: %w", message, err))
		}
	}
	return errors.Join(errs...)
}

func safeSend(n Notifier, message string) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("%w: %v", ErrPanic, value)
		}
	}()
	return n.Send(message)
}
