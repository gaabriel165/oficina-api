package customer

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type ListCustomersUseCase struct {
	customerRepo repository.CustomerRepository
}

func NewListCustomersUseCase(customerRepo repository.CustomerRepository) *ListCustomersUseCase {
	return &ListCustomersUseCase{customerRepo: customerRepo}
}

func (u *ListCustomersUseCase) Execute() ([]*entity.Customer, error) {
	return u.customerRepo.FindAll()
}
