package customer_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/customer"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCustomer_ShouldCreateWithCPF(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	cnpjSvc := mocks.NewMockCNPJValidationService(t)

	customerRepo.On("FindByDocument", mock.Anything).Return(nil, repository.ErrCustomerNotFound)
	customerRepo.On("Create", mock.Anything).Return(nil)

	uc := customer.NewCreateCustomerUseCase(customerRepo, cnpjSvc)
	result, err := uc.Execute(customer.CreateCustomerInput{
		Name:     "John Doe",
		Document: "529.982.247-25",
		Phone:    "11999999999",
		Email:    "john@email.com",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John Doe", result.Name())
	customerRepo.AssertExpectations(t)
}

func TestCreateCustomer_ShouldReturnErrorWhenDocumentAlreadyRegistered(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	cnpjSvc := mocks.NewMockCNPJValidationService(t)

	existing, _ := makeCustomer()
	customerRepo.On("FindByDocument", mock.Anything).Return(existing, nil)

	uc := customer.NewCreateCustomerUseCase(customerRepo, cnpjSvc)
	_, err := uc.Execute(customer.CreateCustomerInput{
		Name:     "John Doe",
		Document: "529.982.247-25",
		Phone:    "11999999999",
		Email:    "john@email.com",
	})

	assert.ErrorIs(t, err, usecase.ErrDocumentAlreadyRegistered)
}

func TestCreateCustomer_ShouldReturnErrorWhenCNPJNotActive(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	cnpjSvc := mocks.NewMockCNPJValidationService(t)

	customerRepo.On("FindByDocument", mock.Anything).Return(nil, repository.ErrCustomerNotFound)
	cnpjSvc.On("Fetch", mock.Anything).Return(&repository.CNPJInfo{IsActive: false}, nil)

	uc := customer.NewCreateCustomerUseCase(customerRepo, cnpjSvc)
	_, err := uc.Execute(customer.CreateCustomerInput{
		Name:     "Acme Corp",
		Document: "11.222.333/0001-81",
		Phone:    "1133334444",
		Email:    "acme@email.com",
	})

	assert.ErrorIs(t, err, usecase.ErrCNPJNotActive)
}
