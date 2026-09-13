package customer

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
)

type ChangeCustomerStatusInput struct {
	ID     string
	Status string
}

type ChangeCustomerStatusUseCase struct {
	customerRepo repository.CustomerRepository
}

func NewChangeCustomerStatusUseCase(customerRepo repository.CustomerRepository) *ChangeCustomerStatusUseCase {
	return &ChangeCustomerStatusUseCase{customerRepo: customerRepo}
}

func (u *ChangeCustomerStatusUseCase) Execute(input ChangeCustomerStatusInput) (*entity.Customer, error) {
	status, err := valueobject.NewCustomerStatus(input.Status)
	if err != nil {
		return nil, err
	}

	customer, err := u.customerRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if status.IsActive() {
		customer.Activate()
	} else {
		customer.Deactivate()
	}

	if err := u.customerRepo.Update(customer); err != nil {
		return nil, err
	}

	return customer, nil
}
