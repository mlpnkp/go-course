package model

const (
	StatusSent   = "sent"
	StatusFailed = "failed"
)

type Notification struct {
	ID      int
	Message string
	Status  string
}
