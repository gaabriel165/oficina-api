package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type ListServiceOrdersUseCase struct {
	orderRepo repository.ServiceOrderRepository
}

func NewListServiceOrdersUseCase(orderRepo repository.ServiceOrderRepository) *ListServiceOrdersUseCase {
	return &ListServiceOrdersUseCase{orderRepo: orderRepo}
}

func (u *ListServiceOrdersUseCase) Execute() ([]*entity.ServiceOrder, error) {
	return u.orderRepo.FindAll()
}
