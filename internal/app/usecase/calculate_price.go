package usecase

import (
	"github.com/gabe-frasz/starting-with-go/internal/app/entity"
	"github.com/gabe-frasz/starting-with-go/internal/app/repository"
)

type OrderInputDTO struct {
	ID    string
	Price float64
	Tax   float64
}

type OrderOutputDTO struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

type CalculatePriceUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewCalculatePriceUseCase(orderRepository repository.OrderRepository) *CalculatePriceUseCase {
	return &CalculatePriceUseCase{orderRepository}
}

func (c *CalculatePriceUseCase) Execute(input *OrderInputDTO) (*OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return nil, err
	}
	err = order.CalculateFinalPrice()
	if err != nil {
		return nil, err
	}

	err = c.OrderRepository.Save(order)
	if err != nil {
		return nil, err
	}

	return &OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
