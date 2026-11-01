package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyMessage = errors.New("message is empty")
	ErrPanic        = errors.New("notifier panicked")
)

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

type Notifier interface {
	Send(message string) error
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

type EmailNotifier struct {
	BaseNotifier
	Address string
}

func (n EmailNotifier) Send(message string) error {
	if err := n.check(message); err != nil {
		return err
	}
	fmt.Printf("email to %s: %s\n", n.Address, message)
	return nil
}

type PushNotifier struct {
	BaseNotifier
	DeviceID string
}

func (n PushNotifier) Send(message string) error {
	if err := n.check(message); err != nil {
		return err
	}
	fmt.Printf("push to %s: %s\n", n.DeviceID, message)
	return nil
}

type BrokenNotifier struct{}

func (BrokenNotifier) Send(string) error {
	panic("драйвер канала поврежден")
}

func Dispatch(notifiers []Notifier, message string) error {
	var errs []error
	for _, n := range notifiers {
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

type MultiError struct {
	errs []string
}

func (m *MultiError) Add(msg string) {
	m.errs = append(m.errs, msg)
}

func (m *MultiError) Error() string {
	return strings.Join(m.errs, "; ")
}

type Customer struct {
	Age  int
	Name string
}

func (c Customer) problems() *MultiError {
	var m *MultiError
	if c.Age < 0 {
		m = &MultiError{}
		m.Add("age is negative")
	}
	if c.Name == "" {
		if m == nil {
			m = &MultiError{}
		}
		m.Add("name is empty")
	}
	return m
}

func (c Customer) ValidateTypedNil() error {
	return c.problems()
}

func (c Customer) Validate() error {
	if m := c.problems(); m != nil {
		return m
	}
	return nil
}

func main() {
	defer func() {
		if value := recover(); value != nil {
			fmt.Println("паника перехвачена в main:", value)
		}
	}()

	notifiers := []Notifier{
		EmailNotifier{BaseNotifier: BaseNotifier{Channel: "email"}, Address: "anna@example.com"},
		BrokenNotifier{},
		PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}, DeviceID: "device-1"},
	}

	err := Dispatch(notifiers, "")
	fmt.Println("ошибка:", err)
	fmt.Println("пустое сообщение:", errors.Is(err, ErrEmptyMessage))
	fmt.Println("отказ push:", errors.Is(err, &SendError{Channel: "push"}))
	fmt.Println("паника канала:", errors.Is(err, ErrPanic))
	var sendErr *SendError
	if errors.As(err, &sendErr) {
		fmt.Println("первый SendError, канал:", sendErr.Channel)
	}

	showValidation("ValidateTypedNil", Customer{Age: 33, Name: "John"}.ValidateTypedNil)
	showValidation("Validate", Customer{Age: 33, Name: "John"}.Validate)

	_ = BrokenNotifier{}.Send("без safeSend")
}

func showValidation(name string, validate func() error) {
	err := validate()
	fmt.Printf("%s: err != nil = %t, тип %T\n", name, err != nil, err)
}
