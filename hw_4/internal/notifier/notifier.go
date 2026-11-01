package notifier

import (
	"errors"
	"fmt"
)

const (
	ChannelEmail = "email"
	ChannelPush  = "push"
)

var ErrEmptyMessage = errors.New("message is empty")

type SendError struct {
	Channel string
	Err     error
}

func (e *SendError) Error() string {
	return fmt.Sprintf("channel %s: %v", e.Channel, e.Err)
}

func (e *SendError) Unwrap() error {
	return e.Err
}

func (e *SendError) Is(target error) bool {
	t, ok := target.(*SendError)
	return ok && t != nil && t.Channel == e.Channel
}

type BaseNotifier struct {
	Channel string
}

func (b BaseNotifier) check(message string) error {
	if message == "" {
		return &SendError{Channel: b.Channel, Err: ErrEmptyMessage}
	}
	return nil
}

type Email struct {
	BaseNotifier
	Address string
}

func NewEmail(address string) Email {
	return Email{BaseNotifier: BaseNotifier{Channel: ChannelEmail}, Address: address}
}

func (n Email) Send(message string) error {
	if err := n.check(message); err != nil {
		return err
	}
	fmt.Printf("email to %s: %s\n", n.Address, message)
	return nil
}

type Push struct {
	BaseNotifier
	DeviceID string
}

func NewPush(deviceID string) Push {
	return Push{BaseNotifier: BaseNotifier{Channel: ChannelPush}, DeviceID: deviceID}
}

func (n Push) Send(message string) error {
	if err := n.check(message); err != nil {
		return err
	}
	fmt.Printf("push to %s: %s\n", n.DeviceID, message)
	return nil
}
