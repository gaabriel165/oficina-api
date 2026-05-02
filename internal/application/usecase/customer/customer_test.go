package customer_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/customer"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCustomer_ShouldReturnCustomer(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	c, _ := makeCustomer()

	customerRepo.On("FindByID", "customer-id").Return(c, nil)

	uc := customer.NewGetCustomerUseCase(customerRepo)
	result, err := uc.Execute("customer-id")

	assert.NoError(t, err)
	assert.Equal(t, "John Doe", result.Name())
}

func TestGetCustomer_ShouldReturnErrorWhenNotFound(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)

	customerRepo.On("FindByID", "customer-id").Return(nil, repository.ErrCustomerNotFound)

	uc := customer.NewGetCustomerUseCase(customerRepo)
	_, err := uc.Execute("customer-id")

	assert.ErrorIs(t, err, repository.ErrCustomerNotFound)
}

func TestListCustomers_ShouldReturnAll(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	c, _ := makeCustomer()

	customerRepo.On("FindAll").Return([]*entity.Customer{c, c}, nil)

	uc := customer.NewListCustomersUseCase(customerRepo)
	results, err := uc.Execute()

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestDeleteCustomer_ShouldDeleteSuccessfully(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	c, _ := makeCustomer()

	customerRepo.On("FindByID", "customer-id").Return(c, nil)
	customerRepo.On("Delete", "customer-id").Return(nil)

	uc := customer.NewDeleteCustomerUseCase(customerRepo)
	err := uc.Execute("customer-id")

	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
}

func TestDeleteCustomer_ShouldReturnErrorWhenNotFound(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)

	customerRepo.On("FindByID", "customer-id").Return(nil, repository.ErrCustomerNotFound)

	uc := customer.NewDeleteCustomerUseCase(customerRepo)
	err := uc.Execute("customer-id")

	assert.ErrorIs(t, err, repository.ErrCustomerNotFound)
}

func TestUpdateCustomer_ShouldUpdateSuccessfully(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	c, _ := makeCustomer()

	customerRepo.On("FindByID", "customer-id").Return(c, nil)
	customerRepo.On("Update", mock.Anything).Return(nil)

	uc := customer.NewUpdateCustomerUseCase(customerRepo)
	result, err := uc.Execute(customer.UpdateCustomerInput{
		ID:    "customer-id",
		Name:  "Jane Doe",
		Phone: "11988888888",
		Email: "jane@email.com",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", result.Name())
	assert.Equal(t, "jane@email.com", result.Email())
}

func TestUpdateCustomer_ShouldReturnErrorWhenNotFound(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)

	customerRepo.On("FindByID", "customer-id").Return(nil, repository.ErrCustomerNotFound)

	uc := customer.NewUpdateCustomerUseCase(customerRepo)
	_, err := uc.Execute(customer.UpdateCustomerInput{
		ID:    "customer-id",
		Name:  "Jane Doe",
		Phone: "11988888888",
		Email: "jane@email.com",
	})

	assert.ErrorIs(t, err, repository.ErrCustomerNotFound)
}

func TestUpdateCustomer_ShouldReturnErrorWhenNameEmpty(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	c, _ := makeCustomer()

	customerRepo.On("FindByID", "customer-id").Return(c, nil)

	uc := customer.NewUpdateCustomerUseCase(customerRepo)
	_, err := uc.Execute(customer.UpdateCustomerInput{
		ID:    "customer-id",
		Name:  "",
		Phone: "11988888888",
		Email: "jane@email.com",
	})

	assert.ErrorIs(t, err, entity.ErrCustomerNameRequired)
}
