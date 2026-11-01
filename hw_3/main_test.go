package main

import (
	"errors"
	"fmt"
	"testing"
)

type recordingNotifier struct {
	messages []string
}

func (r *recordingNotifier) Send(message string) error {
	r.messages = append(r.messages, message)
	return nil
}

func TestSendEmptyMessageReturnsSendError(t *testing.T) {
	err := PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}, DeviceID: "device-1"}.Send("")

	if !errors.Is(err, ErrEmptyMessage) {
		t.Fatalf("errors.Is(%v, ErrEmptyMessage) = false", err)
	}
	var sendErr *SendError
	if !errors.As(err, &sendErr) || sendErr.Channel != "push" {
		t.Errorf("errors.As(%v) = %+v, want SendError for push", err, sendErr)
	}
}

func TestSendErrorIsMatchesChannel(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &SendError{Channel: "push", Err: ErrEmptyMessage})

	if !errors.Is(err, &SendError{Channel: "push"}) {
		t.Error("errors.Is with same channel = false, want true")
	}
	if errors.Is(err, &SendError{Channel: "email"}) {
		t.Error("errors.Is with other channel = true, want false")
	}
}

func TestDispatchSendsToEveryNotifier(t *testing.T) {
	first, second := &recordingNotifier{}, &recordingNotifier{}

	if err := Dispatch([]Notifier{first, second}, "Заказ готов"); err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if len(first.messages) != 1 || len(second.messages) != 1 {
		t.Errorf("messages = %q and %q, want one each", first.messages, second.messages)
	}
}

func TestDispatchIsolatesPanicOfOneChannel(t *testing.T) {
	after := &recordingNotifier{}

	err := Dispatch([]Notifier{BrokenNotifier{}, after}, "Заказ готов")

	if !errors.Is(err, ErrPanic) {
		t.Fatalf("errors.Is(%v, ErrPanic) = false", err)
	}
	if len(after.messages) != 1 {
		t.Errorf("notifier after panic got %q, want one message", after.messages)
	}
}

func TestValidateTypedNilIsNotNil(t *testing.T) {
	customer := Customer{Age: 33, Name: "John"}

	if err := customer.ValidateTypedNil(); err == nil {
		t.Error("ValidateTypedNil() == nil, want typed nil *MultiError")
	}
	if err := customer.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}
