package usecase

import "github.com/gabe-frasz/starting-with-go/internal/app/repository"

type GetTotalOutputDTO struct {
	Total int
}

type GetTotalUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewGetTotalUseCase(orderRepository repository.OrderRepository) *GetTotalUseCase {
	return &GetTotalUseCase{orderRepository}
}

func (g *GetTotalUseCase) Execute() (*GetTotalOutputDTO, error) {
	total, err := g.OrderRepository.GetTotal()
	if err != nil {
		return nil, err
	}

	return &GetTotalOutputDTO{total}, nil
}
