package customer

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type UpdateCustomerInput struct {
	ID    string
	Name  string
	Phone string
	Email string
}

type UpdateCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

func NewUpdateCustomerUseCase(customerRepo repository.CustomerRepository) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{customerRepo: customerRepo}
}

func (u *UpdateCustomerUseCase) Execute(input UpdateCustomerInput) (*entity.Customer, error) {
	customer, err := u.customerRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := customer.Update(input.Name, input.Phone, input.Email); err != nil {
		return nil, err
	}

	if err := u.customerRepo.Update(customer); err != nil {
		return nil, err
	}

	return customer, nil
}
