package serviceorder_test

import (
	"testing"
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/serviceorder"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetServiceOrder_ShouldReturnOrder(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrder(), nil)

	uc := serviceorder.NewGetServiceOrderUseCase(orderRepo)
	result, err := uc.Execute("order-id")

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusReceived, result.Status())
}

func TestGetServiceOrder_ShouldReturnErrorWhenNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewGetServiceOrderUseCase(orderRepo)
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func makeOrderInExecution() *entity.ServiceOrder {
	order := makeOrderWaitingApproval()
	order.ApproveBudget()
	return order
}

func TestListServiceOrders_ShouldSortByStatusPriority(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	received := makeOrder()
	inDiagnosis := makeOrderInDiagnosis()
	waitingApproval := makeOrderWaitingApproval()
	inExecution := makeOrderInExecution()

	orderRepo.On("FindByStatuses", mock.Anything).Return(
		[]*entity.ServiceOrder{received, waitingApproval, inExecution, inDiagnosis}, nil,
	)

	uc := serviceorder.NewListServiceOrdersUseCase(orderRepo)
	results, err := uc.Execute()

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusInExecution, results[0].Status())
	assert.Equal(t, valueobject.OrderStatusWaitingApproval, results[1].Status())
	assert.Equal(t, valueobject.OrderStatusInDiagnosis, results[2].Status())
	assert.Equal(t, valueobject.OrderStatusReceived, results[3].Status())
}

func TestListServiceOrders_ShouldSortOldestFirstWithinSameStatus(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	older := entity.RestoreServiceOrder(
		"older", "customer-id", "vehicle-id", valueobject.OrderStatusReceived, "", 0,
		nil, nil, nil, nil, time.Now().Add(-2*time.Hour), time.Now(),
	)
	newer := entity.RestoreServiceOrder(
		"newer", "customer-id", "vehicle-id", valueobject.OrderStatusReceived, "", 0,
		nil, nil, nil, nil, time.Now().Add(-1*time.Hour), time.Now(),
	)

	orderRepo.On("FindByStatuses", mock.Anything).Return(
		[]*entity.ServiceOrder{newer, older}, nil,
	)

	uc := serviceorder.NewListServiceOrdersUseCase(orderRepo)
	results, err := uc.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "older", results[0].ID())
	assert.Equal(t, "newer", results[1].ID())
}

func TestStartDiagnosis_ShouldTransitionStatus(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrder(), nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	uc := serviceorder.NewStartDiagnosisUseCase(orderRepo, noopNotifier{})
	result, err := uc.Execute("order-id")

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusInDiagnosis, result.Status())
}

func TestStartDiagnosis_ShouldReturnErrorWhenNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewStartDiagnosisUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func TestStartDiagnosis_ShouldReturnErrorOnInvalidTransition(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)

	uc := serviceorder.NewStartDiagnosisUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, entity.ErrServiceOrderInvalidTransition)
}

func TestAddServiceToOrder_ShouldAddSuccessfully(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	serviceRepo := mocks.NewMockServiceRepository(t)
	svc, _ := entity.NewService("Oil Change", "desc", 150.00, 60)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)
	serviceRepo.On("FindByID", svc.ID()).Return(svc, nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	uc := serviceorder.NewAddServiceToOrderUseCase(orderRepo, serviceRepo)
	result, err := uc.Execute(serviceorder.AddServiceToOrderInput{
		OrderID:   "order-id",
		ServiceID: svc.ID(),
	})

	assert.NoError(t, err)
	assert.Len(t, result.Items(), 1)
}

func TestAddServiceToOrder_ShouldReturnErrorWhenOrderNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	serviceRepo := mocks.NewMockServiceRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewAddServiceToOrderUseCase(orderRepo, serviceRepo)
	_, err := uc.Execute(serviceorder.AddServiceToOrderInput{
		OrderID:   "order-id",
		ServiceID: "service-id",
	})

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func TestAddServiceToOrder_ShouldReturnErrorWhenServiceNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	serviceRepo := mocks.NewMockServiceRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)
	serviceRepo.On("FindByID", "service-id").Return(nil, repository.ErrServiceNotFound)

	uc := serviceorder.NewAddServiceToOrderUseCase(orderRepo, serviceRepo)
	_, err := uc.Execute(serviceorder.AddServiceToOrderInput{
		OrderID:   "order-id",
		ServiceID: "service-id",
	})

	assert.ErrorIs(t, err, repository.ErrServiceNotFound)
}

func TestAddPartToOrder_ShouldAddSuccessfully(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	partRepo := mocks.NewMockPartRepository(t)
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 10)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)
	partRepo.On("FindByID", part.ID()).Return(part, nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	uc := serviceorder.NewAddPartToOrderUseCase(orderRepo, partRepo)
	result, err := uc.Execute(serviceorder.AddPartToOrderInput{
		OrderID:  "order-id",
		PartID:   part.ID(),
		Quantity: 2,
	})

	assert.NoError(t, err)
	assert.Len(t, result.Parts(), 1)
}

func TestAddPartToOrder_ShouldReturnErrorWhenOrderNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	partRepo := mocks.NewMockPartRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewAddPartToOrderUseCase(orderRepo, partRepo)
	_, err := uc.Execute(serviceorder.AddPartToOrderInput{
		OrderID:  "order-id",
		PartID:   "part-id",
		Quantity: 1,
	})

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func TestAddPartToOrder_ShouldReturnErrorWhenPartNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	partRepo := mocks.NewMockPartRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)
	partRepo.On("FindByID", "part-id").Return(nil, repository.ErrPartNotFound)

	uc := serviceorder.NewAddPartToOrderUseCase(orderRepo, partRepo)
	_, err := uc.Execute(serviceorder.AddPartToOrderInput{
		OrderID:  "order-id",
		PartID:   "part-id",
		Quantity: 1,
	})

	assert.ErrorIs(t, err, repository.ErrPartNotFound)
}

func TestAddPartToOrder_ShouldReturnErrorWhenInsufficientStock(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	partRepo := mocks.NewMockPartRepository(t)
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 2)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)
	partRepo.On("FindByID", part.ID()).Return(part, nil)

	uc := serviceorder.NewAddPartToOrderUseCase(orderRepo, partRepo)
	_, err := uc.Execute(serviceorder.AddPartToOrderInput{
		OrderID:  "order-id",
		PartID:   part.ID(),
		Quantity: 10,
	})

	assert.ErrorIs(t, err, entity.ErrPartInsufficientStock)
}

func TestAddPartToOrder_ShouldReturnErrorWhenQuantityZero(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	partRepo := mocks.NewMockPartRepository(t)
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 10)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)
	partRepo.On("FindByID", part.ID()).Return(part, nil)

	uc := serviceorder.NewAddPartToOrderUseCase(orderRepo, partRepo)
	_, err := uc.Execute(serviceorder.AddPartToOrderInput{
		OrderID:  "order-id",
		PartID:   part.ID(),
		Quantity: 0,
	})

	assert.ErrorIs(t, err, entity.ErrServiceOrderItemInvalidQty)
}

func TestSendBudget_ShouldTransitionStatus(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	svc, _ := entity.NewService("Oil Change", "desc", 150.00, 60)
	order := makeOrderInDiagnosis()
	order.AddItem(svc)
	orderRepo.ExpectedCalls = nil
	orderRepo.On("FindByID", "order-id").Return(order, nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	uc := serviceorder.NewSendBudgetUseCase(orderRepo, noopNotifier{})
	result, err := uc.Execute("order-id")

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusWaitingApproval, result.Status())
}

func TestSendBudget_ShouldReturnErrorWhenOrderNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewSendBudgetUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func TestSendBudget_ShouldReturnErrorWhenNoItems(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrderInDiagnosis(), nil)

	uc := serviceorder.NewSendBudgetUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, entity.ErrServiceOrderNoItems)
}

func TestFinishExecution_ShouldTransitionStatus(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	order := makeOrderWaitingApproval()
	order.ApproveBudget()

	orderRepo.On("FindByID", "order-id").Return(order, nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	uc := serviceorder.NewFinishExecutionUseCase(orderRepo, noopNotifier{})
	result, err := uc.Execute("order-id")

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusFinished, result.Status())
}

func TestFinishExecution_ShouldReturnErrorWhenOrderNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewFinishExecutionUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func TestFinishExecution_ShouldReturnErrorOnInvalidTransition(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrder(), nil)

	uc := serviceorder.NewFinishExecutionUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, entity.ErrServiceOrderInvalidTransition)
}

func TestDeliverVehicle_ShouldTransitionStatus(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	order := makeOrderWaitingApproval()
	order.ApproveBudget()
	order.FinishExecution()

	orderRepo.On("FindByID", "order-id").Return(order, nil)
	orderRepo.On("Update", mock.Anything).Return(nil)

	uc := serviceorder.NewDeliverVehicleUseCase(orderRepo, noopNotifier{})
	result, err := uc.Execute("order-id")

	assert.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusDelivered, result.Status())
}

func TestDeliverVehicle_ShouldReturnErrorWhenOrderNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewDeliverVehicleUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func TestDeliverVehicle_ShouldReturnErrorOnInvalidTransition(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrder(), nil)

	uc := serviceorder.NewDeliverVehicleUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, entity.ErrServiceOrderInvalidTransition)
}

func TestApproveBudget_ShouldReturnErrorWhenOrderNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	partRepo := mocks.NewMockPartRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewApproveBudgetUseCase(orderRepo, partRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func TestApproveBudget_ShouldReturnErrorOnInvalidTransition(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)
	partRepo := mocks.NewMockPartRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrder(), nil)

	uc := serviceorder.NewApproveBudgetUseCase(orderRepo, partRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, entity.ErrServiceOrderInvalidTransition)
}

func TestRejectBudget_ShouldReturnErrorWhenOrderNotFound(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(nil, repository.ErrServiceOrderNotFound)

	uc := serviceorder.NewRejectBudgetUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, repository.ErrServiceOrderNotFound)
}

func TestRejectBudget_ShouldReturnErrorOnInvalidTransition(t *testing.T) {
	orderRepo := mocks.NewMockServiceOrderRepository(t)

	orderRepo.On("FindByID", "order-id").Return(makeOrder(), nil)

	uc := serviceorder.NewRejectBudgetUseCase(orderRepo, noopNotifier{})
	_, err := uc.Execute("order-id")

	assert.ErrorIs(t, err, entity.ErrServiceOrderInvalidTransition)
}
