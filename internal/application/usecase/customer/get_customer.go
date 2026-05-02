package customer

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type GetCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

func NewGetCustomerUseCase(customerRepo repository.CustomerRepository) *GetCustomerUseCase {
	return &GetCustomerUseCase{customerRepo: customerRepo}
}

func (u *GetCustomerUseCase) Execute(id string) (*entity.Customer, error) {
	return u.customerRepo.FindByID(id)
}
