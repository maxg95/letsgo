package mocks

import (
	"time"

	"new.mod/internal/models"
)

var mockOrder = models.Order{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type OrderModel struct{}

func (m *OrderModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *OrderModel) Get(id int) (models.Order, error) {
	switch id {
	case 1:
		return mockOrder, nil
	default:
		return models.Order{}, models.ErrNoRecord
	}
}

func (m *OrderModel) Latest() ([]models.Order, error) {
	return []models.Order{mockOrder}, nil
}
