package serviceorder_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/serviceorder"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeOrder() *entity.ServiceOrder {
	order, _ := entity.NewServiceOrder("customer-id", "vehicle-id", "noise in engine")
	return order
}

func makeOrderInDiagnosis() *entity.ServiceOrder {
	order := makeOrder()
	order.StartDiagnosis()
	return order
}

func makeOrderWaitingApproval() *entity.ServiceOrder {
	order := makeOrderInDiagnosis()
	svc, _ := entity.NewService("Oil Change", "desc", 150.00, 60)
	order.AddItem(svc)
	order.SendBudget()
	return order
}

func TestCreateServiceOrder_ShouldCreateSuccessfully(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	customerRepo := mocks.NewMockCustomerRepository(t)
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	customer, _ := entity.NewCustomer("John", "529.982.247-25", "11999999999", "j@email.com")
	vehicle, _ := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", 2020)

	customerRepo.On("FindByID", "customer-id").Return(customer, nil)
	vehicleRepo.On("FindByID", "vehicle-id").Return(vehicle, nil)
	orderRepo.On("Create", mock.Anything).Return(nil)

	uc := serviceorder.NewCreateServiceOrderUseCase(orderRepo, customerRepo, vehicleRepo)
	result, err := uc.Execute(serviceorder.CreateServiceOrderInput{
		CustomerID: "customer-id",
		VehicleID:  "vehicle-id",
		Notes:      "noise in engine",
	})

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusReceived, result.Status())
	orderRepo.AssertExpectations(t)
}

func TestCreateServiceOrder_ShouldReturnErrorWhenCustomerNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	customerRepo := mocks.NewMockCustomerRepository(t)
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	customerRepo.On("FindByID", "customer-id").Return(nil, repository.ErrCustomerNotFound)

	uc := serviceorder.NewCreateServiceOrderUseCase(orderRepo, customerRepo, vehicleRepo)
	_, err := uc.Execute(serviceorder.CreateServiceOrderInput{
		CustomerID: "customer-id",
		VehicleID:  "vehicle-id",
	})

	assert.ErrorIs(t, err, repository.ErrCustomerNotFound)
}

func TestCreateServiceOrder_ShouldReturnErrorWhenVehicleNotOwnedByCustomer(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	customerRepo := mocks.NewMockCustomerRepository(t)
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	customer, _ := entity.NewCustomer("John", "529.982.247-25", "11999999999", "j@email.com")
	vehicleOfAnotherCustomer, _ := entity.NewVehicle("other-customer-id", "ABC-1234", "Toyota", "Corolla", 2020)

	customerRepo.On("FindByID", "customer-id").Return(customer, nil)
	vehicleRepo.On("FindByID", "vehicle-id").Return(vehicleOfAnotherCustomer, nil)

	uc := serviceorder.NewCreateServiceOrderUseCase(orderRepo, customerRepo, vehicleRepo)
	_, err := uc.Execute(serviceorder.CreateServiceOrderInput{
		CustomerID: "customer-id",
		VehicleID:  "vehicle-id",
	})

	assert.ErrorIs(t, err, entity.ErrServiceOrderVehicleNotOwnedByCustomer)
}

func TestApproveBudget_ShouldDebitStockAndTransitionStatus(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	partRepo := mocks.NewMockPartRepository(t)

	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 10)
	order := makeOrderWaitingApproval()
	order.AddPart(part, 2)

	orderRepo.On("FindByID", "order-id").Return(order, nil)
	partRepo.On("FindByID", part.ID()).Return(part, nil)
	partRepo.On("Update", mock.Anything).Return(nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	uc := serviceorder.NewApproveBudgetUseCase(orderRepo, partRepo)
	result, err := uc.Execute("order-id")

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusInExecution, result.Status())
	assert.Equal(t, 8, part.StockQuantity())
}

func TestRejectBudget_ShouldCancelOrder(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	order := makeOrderWaitingApproval()
	orderRepo.On("FindByID", "order-id").Return(order, nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	uc := serviceorder.NewRejectBudgetUseCase(orderRepo)
	result, err := uc.Execute("order-id")

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusCancelled, result.Status())
}

func TestFullServiceOrderFlow(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	serviceRepo := mocks.NewMockServiceRepository(t)
	partRepo := mocks.NewMockPartRepository(t)

	svc, _ := entity.NewService("Oil Change", "desc", 150.00, 60)
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 10)
	order := makeOrder()

	orderRepo.On("FindByID", mock.Anything).Return(order, nil)
	orderRepo.On("Update", mock.Anything).Return(nil)
	serviceRepo.On("FindByID", svc.ID()).Return(svc, nil)
	partRepo.On("FindByID", part.ID()).Return(part, nil)
	partRepo.On("Update", mock.Anything).Return(nil)

	startUC := serviceorder.NewStartDiagnosisUseCase(orderRepo)
	addSvcUC := serviceorder.NewAddServiceToOrderUseCase(orderRepo, serviceRepo)
	addPartUC := serviceorder.NewAddPartToOrderUseCase(orderRepo, partRepo)
	sendBudgetUC := serviceorder.NewSendBudgetUseCase(orderRepo)
	approveUC := serviceorder.NewApproveBudgetUseCase(orderRepo, partRepo)
	finishExecUC := serviceorder.NewFinishExecutionUseCase(orderRepo)
	deliverUC := serviceorder.NewDeliverVehicleUseCase(orderRepo)

	startUC.Execute("order-id")
	addSvcUC.Execute(serviceorder.AddServiceToOrderInput{OrderID: "order-id", ServiceID: svc.ID()})
	addPartUC.Execute(serviceorder.AddPartToOrderInput{OrderID: "order-id", PartID: part.ID(), Quantity: 2})
	sendBudgetUC.Execute("order-id")
	approveUC.Execute("order-id")
	finishExecUC.Execute("order-id")
	result, err := deliverUC.Execute("order-id")

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusDelivered, result.Status())
	assert.Equal(t, 8, part.StockQuantity())
}
