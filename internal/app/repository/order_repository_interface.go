package repository

import "github.com/gabe-frasz/starting-with-go/internal/app/entity"

type OrderRepository interface {
	Save(order *entity.Order) error
	GetTotal() (int, error)
}
