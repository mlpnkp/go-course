package repository

import (
	"errors"
	"fmt"

	"go-course/hw_4/internal/model"
)

var ErrNotFound = errors.New("notification not found")

type MemoryRepository struct {
	notifications map[int]model.Notification
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{notifications: make(map[int]model.Notification)}
}

func (r *MemoryRepository) Save(n model.Notification) error {
	r.notifications[n.ID] = n
	return nil
}

func (r *MemoryRepository) ByID(id int) (model.Notification, error) {
	n, ok := r.notifications[id]
	if !ok {
		return model.Notification{}, fmt.Errorf("notification %d: %w", id, ErrNotFound)
	}
	return n, nil
}
