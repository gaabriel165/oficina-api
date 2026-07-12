package serviceorder_test

import (
	"errors"
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/serviceorder"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

func TestOrderStatusNotifier_ShouldNotifyCustomerEmail(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	notifications := mocks.NewMockNotificationService(t)

	customer, _ := entity.NewCustomer("John", "529.982.247-25", "11999999999", "john@workshop.com")
	order := makeOrder()

	customerRepo.On("FindByID", order.CustomerID()).Return(customer, nil)
	notifications.On("NotifyOrderStatusChanged", customer, order).Return(nil)

	notifier := serviceorder.NewOrderStatusNotifier(customerRepo, notifications)
	notifier.Notify(order)

	notifications.AssertExpectations(t)
}

func TestOrderStatusNotifier_ShouldSkipWhenCustomerNotFound(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	notifications := mocks.NewMockNotificationService(t)

	order := makeOrder()
	customerRepo.On("FindByID", order.CustomerID()).Return(nil, repository.ErrCustomerNotFound)

	notifier := serviceorder.NewOrderStatusNotifier(customerRepo, notifications)
	notifier.Notify(order)

	notifications.AssertNotCalled(t, "NotifyOrderStatusChanged", mock.Anything, mock.Anything)
}

func TestOrderStatusNotifier_ShouldSwallowNotificationError(t *testing.T) {
	customerRepo := mocks.NewMockCustomerRepository(t)
	notifications := mocks.NewMockNotificationService(t)

	customer, _ := entity.NewCustomer("John", "529.982.247-25", "11999999999", "john@workshop.com")
	order := makeOrder()

	customerRepo.On("FindByID", order.CustomerID()).Return(customer, nil)
	notifications.On("NotifyOrderStatusChanged", mock.Anything, mock.Anything).Return(errors.New("smtp unavailable"))

	notifier := serviceorder.NewOrderStatusNotifier(customerRepo, notifications)
	notifier.Notify(order)
}
