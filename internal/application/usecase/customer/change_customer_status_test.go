package customer_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/customer"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangeCustomerStatus_ShouldDeactivateCustomer(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	c, _ := makeCustomer()

	customerRepo.On("FindByID", "customer-id").Return(c, nil)
	customerRepo.On("Update", mock.Anything).Return(nil)

	uc := customer.NewChangeCustomerStatusUseCase(customerRepo)
	result, err := uc.Execute(customer.ChangeCustomerStatusInput{ID: "customer-id", Status: "inactive"})

	assert.NoError(t, err)
	assert.False(t, result.IsActive())
}

func TestChangeCustomerStatus_ShouldActivateCustomer(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	c, _ := makeCustomer()
	c.Deactivate()

	customerRepo.On("FindByID", "customer-id").Return(c, nil)
	customerRepo.On("Update", mock.Anything).Return(nil)

	uc := customer.NewChangeCustomerStatusUseCase(customerRepo)
	result, err := uc.Execute(customer.ChangeCustomerStatusInput{ID: "customer-id", Status: "active"})

	assert.NoError(t, err)
	assert.True(t, result.IsActive())
}

func TestChangeCustomerStatus_ShouldReturnErrorWhenStatusInvalid(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)

	uc := customer.NewChangeCustomerStatusUseCase(customerRepo)
	_, err := uc.Execute(customer.ChangeCustomerStatusInput{ID: "customer-id", Status: "blocked"})

	assert.ErrorIs(t, err, valueobject.ErrInvalidCustomerStatus)
}

func TestChangeCustomerStatus_ShouldReturnErrorWhenCustomerNotFound(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)

	customerRepo.On("FindByID", "customer-id").Return(nil, repository.ErrCustomerNotFound)

	uc := customer.NewChangeCustomerStatusUseCase(customerRepo)
	_, err := uc.Execute(customer.ChangeCustomerStatusInput{ID: "customer-id", Status: "inactive"})

	assert.ErrorIs(t, err, repository.ErrCustomerNotFound)
}
