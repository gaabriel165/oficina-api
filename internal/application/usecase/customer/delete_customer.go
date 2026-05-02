package customer

import "github.com/gabrielcamargo/oficina-api/internal/domain/repository"

type DeleteCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

func NewDeleteCustomerUseCase(customerRepo repository.CustomerRepository) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{customerRepo: customerRepo}
}

func (u *DeleteCustomerUseCase) Execute(id string) error {
	if _, err := u.customerRepo.FindByID(id); err != nil {
		return err
	}
	return u.customerRepo.Delete(id)
}
