package customer

import (
	"errors"

	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/utils"
)

type CreateCustomerInput struct {
	Name     string
	Document string
	Phone    string
	Email    string
}

type CreateCustomerUseCase struct {
	customerRepo      repository.CustomerRepository
	cnpjValidationSvc repository.CNPJValidationService
}

func NewCreateCustomerUseCase(
	customerRepo repository.CustomerRepository,
	cnpjValidationSvc repository.CNPJValidationService,
) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		customerRepo:      customerRepo,
		cnpjValidationSvc: cnpjValidationSvc,
	}
}

func (u *CreateCustomerUseCase) Execute(input CreateCustomerInput) (*entity.Customer, error) {
	_, err := u.customerRepo.FindByDocument(utils.CleanNumericString(input.Document))
	if err == nil {
		return nil, usecase.ErrDocumentAlreadyRegistered
	}
	if !errors.Is(err, repository.ErrCustomerNotFound) {
		return nil, err
	}

	customer, err := entity.NewCustomer(input.Name, input.Document, input.Phone, input.Email)
	if err != nil {
		return nil, err
	}

	if customer.DocumentType() == entity.DocumentTypeCNPJ {
		info, err := u.cnpjValidationSvc.Fetch(customer.Document())
		if err == nil && !info.IsActive {
			return nil, usecase.ErrCNPJNotActive
		}
	}

	if err := u.customerRepo.Create(customer); err != nil {
		return nil, err
	}

	return customer, nil
}
