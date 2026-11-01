package notifier_test

import (
	"errors"
	"testing"

	"go-course/hw_4/internal/notifier"
)

func TestSendEmptyMessageReturnsSendError(t *testing.T) {
	for _, n := range []interface{ Send(string) error }{
		notifier.NewEmail("anna@example.com"),
		notifier.NewPush("device-1"),
	} {
		err := n.Send("")

		if !errors.Is(err, notifier.ErrEmptyMessage) {
			t.Errorf("%T: errors.Is(%v, ErrEmptyMessage) = false", n, err)
		}
		var sendErr *notifier.SendError
		if !errors.As(err, &sendErr) {
			t.Errorf("%T: errors.As(%v, *SendError) = false", n, err)
		}
	}
}

func TestSendErrorIsMatchesChannel(t *testing.T) {
	err := notifier.NewPush("device-1").Send("")

	if !errors.Is(err, &notifier.SendError{Channel: notifier.ChannelPush}) {
		t.Error("errors.Is with push channel = false, want true")
	}
	if errors.Is(err, &notifier.SendError{Channel: notifier.ChannelEmail}) {
		t.Error("errors.Is with email channel = true, want false")
	}
}
