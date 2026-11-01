package service_test

import (
	"errors"
	"testing"

	"go-course/hw_4/internal/model"
	"go-course/hw_4/internal/repository"
	"go-course/hw_4/internal/service"
)

type recordingNotifier struct {
	messages []string
	err      error
}

func (r *recordingNotifier) Send(message string) error {
	r.messages = append(r.messages, message)
	return r.err
}

type panickingNotifier struct{}

func (panickingNotifier) Send(string) error {
	panic("драйвер канала поврежден")
}

func TestSendDeliversAndSavesNotification(t *testing.T) {
	repo := repository.NewMemoryRepository()
	first, second := &recordingNotifier{}, &recordingNotifier{}
	svc := service.NewService(repo, []service.Notifier{first, second})

	sent, err := svc.Send("Заказ готов")
	if err != nil {
		t.Fatalf("Send: %v", err)
	}

	if len(first.messages) != 1 || len(second.messages) != 1 {
		t.Errorf("messages = %q and %q, want one each", first.messages, second.messages)
	}
	saved, err := svc.Notification(sent.ID)
	if err != nil {
		t.Fatalf("Notification(%d): %v", sent.ID, err)
	}
	if saved != sent || saved.Status != model.StatusSent {
		t.Errorf("saved = %+v, want %+v with status sent", saved, sent)
	}
}

func TestSendAssignsSequentialIDs(t *testing.T) {
	svc := service.NewService(repository.NewMemoryRepository(), nil)

	first, _ := svc.Send("1")
	second, _ := svc.Send("2")

	if first.ID != 1 || second.ID != 2 {
		t.Errorf("ids = %d, %d, want 1, 2", first.ID, second.ID)
	}
}

func TestSendSavesFailedNotificationAndKeepsSending(t *testing.T) {
	errDown := errors.New("channel down")
	after := &recordingNotifier{}
	svc := service.NewService(repository.NewMemoryRepository(),
		[]service.Notifier{&recordingNotifier{err: errDown}, panickingNotifier{}, after})

	sent, err := svc.Send("Заказ готов")

	if !errors.Is(err, errDown) || !errors.Is(err, service.ErrPanic) {
		t.Fatalf("Send error = %v, want channel down and panic", err)
	}
	if len(after.messages) != 1 {
		t.Errorf("notifier after failures got %q, want one message", after.messages)
	}
	saved, _ := svc.Notification(sent.ID)
	if saved.Status != model.StatusFailed {
		t.Errorf("saved status = %q, want %q", saved.Status, model.StatusFailed)
	}
}

func TestNotificationUnknownIDReturnsNotFound(t *testing.T) {
	svc := service.NewService(repository.NewMemoryRepository(), nil)

	if _, err := svc.Notification(42); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Notification(42) error = %v, want %v", err, repository.ErrNotFound)
	}
}
