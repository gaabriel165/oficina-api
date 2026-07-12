package integration_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/serviceorder"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	infrarepo "github.com/gabrielcamargo/oficina-api/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopNotifier struct{}

func (noopNotifier) Notify(*entity.ServiceOrder) {}

func TestServiceOrderFullFlow_Integration(t *testing.T) {
	db := newTestDB(t)

	customerRepo := infrarepo.NewGormCustomerRepository(db)
	vehicleRepo := infrarepo.NewGormVehicleRepository(db)
	partRepo := infrarepo.NewGormPartRepository(db)
	serviceRepo := infrarepo.NewGormServiceRepository(db)
	orderRepo := infrarepo.NewGormServiceOrderRepository(db)

	customer, err := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@workshop.com")
	require.NoError(t, err)
	require.NoError(t, customerRepo.Create(customer))

	vehicle, err := entity.NewVehicle(customer.ID(), "ABC-1234", "Toyota", "Corolla", 2020)
	require.NoError(t, err)
	require.NoError(t, vehicleRepo.Create(vehicle))

	svc, err := entity.NewService("Oil Change", "Full oil change", 150.00, 60)
	require.NoError(t, err)
	require.NoError(t, serviceRepo.Create(svc))

	part, err := entity.NewPart("Oil Filter", "Standard filter", 29.90, 10)
	require.NoError(t, err)
	require.NoError(t, partRepo.Create(part))

	createUC := serviceorder.NewCreateServiceOrderUseCase(orderRepo, customerRepo, vehicleRepo, serviceRepo, partRepo)
	order, err := createUC.Execute(serviceorder.CreateServiceOrderInput{
		CustomerID: customer.ID(),
		VehicleID:  vehicle.ID(),
		Notes:      "noise in engine",
	})
	require.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusReceived, order.Status())

	startUC := serviceorder.NewStartDiagnosisUseCase(orderRepo, noopNotifier{})
	order, err = startUC.Execute(order.ID())
	require.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusInDiagnosis, order.Status())

	addSvcUC := serviceorder.NewAddServiceToOrderUseCase(orderRepo, serviceRepo)
	order, err = addSvcUC.Execute(serviceorder.AddServiceToOrderInput{
		OrderID:   order.ID(),
		ServiceID: svc.ID(),
	})
	require.NoError(t, err)
	assert.Len(t, order.Items(), 1)

	addPartUC := serviceorder.NewAddPartToOrderUseCase(orderRepo, partRepo)
	order, err = addPartUC.Execute(serviceorder.AddPartToOrderInput{
		OrderID:  order.ID(),
		PartID:   part.ID(),
		Quantity: 2,
	})
	require.NoError(t, err)
	assert.Len(t, order.Parts(), 1)

	sendBudgetUC := serviceorder.NewSendBudgetUseCase(orderRepo, noopNotifier{})
	order, err = sendBudgetUC.Execute(order.ID())
	require.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusWaitingApproval, order.Status())
	assert.Equal(t, 209.80, order.TotalAmount())

	approveUC := serviceorder.NewApproveBudgetUseCase(orderRepo, partRepo, noopNotifier{})
	order, err = approveUC.Execute(order.ID())
	require.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusInExecution, order.Status())

	updatedPart, err := partRepo.FindByID(part.ID())
	require.NoError(t, err)
	assert.Equal(t, 8, updatedPart.StockQuantity())

	finishUC := serviceorder.NewFinishExecutionUseCase(orderRepo, noopNotifier{})
	order, err = finishUC.Execute(order.ID())
	require.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusFinished, order.Status())

	deliverUC := serviceorder.NewDeliverVehicleUseCase(orderRepo, noopNotifier{})
	order, err = deliverUC.Execute(order.ID())
	require.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusDelivered, order.Status())
	assert.NotNil(t, order.StartedAt())
	assert.NotNil(t, order.FinishedAt())
}

func TestServiceOrderRejectionFlow_Integration(t *testing.T) {
	db := newTestDB(t)

	customerRepo := infrarepo.NewGormCustomerRepository(db)
	vehicleRepo := infrarepo.NewGormVehicleRepository(db)
	serviceRepo := infrarepo.NewGormServiceRepository(db)
	partRepo := infrarepo.NewGormPartRepository(db)
	orderRepo := infrarepo.NewGormServiceOrderRepository(db)

	customer, err := entity.NewCustomer("Jane Doe", "123.456.789-09", "11888888888", "jane@workshop.com")
	require.NoError(t, err)
	require.NoError(t, customerRepo.Create(customer))

	vehicle, err := entity.NewVehicle(customer.ID(), "DEF-5678", "Honda", "Civic", 2021)
	require.NoError(t, err)
	require.NoError(t, vehicleRepo.Create(vehicle))

	svc, err := entity.NewService("Brake Service", "Full brake check", 200.00, 90)
	require.NoError(t, err)
	require.NoError(t, serviceRepo.Create(svc))

	createUC := serviceorder.NewCreateServiceOrderUseCase(orderRepo, customerRepo, vehicleRepo, serviceRepo, partRepo)
	order, err := createUC.Execute(serviceorder.CreateServiceOrderInput{
		CustomerID: customer.ID(),
		VehicleID:  vehicle.ID(),
		Notes:      "brakes squeaking",
	})
	require.NoError(t, err)

	startUC := serviceorder.NewStartDiagnosisUseCase(orderRepo, noopNotifier{})
	order, err = startUC.Execute(order.ID())
	require.NoError(t, err)

	addSvcUC := serviceorder.NewAddServiceToOrderUseCase(orderRepo, serviceRepo)
	_, err = addSvcUC.Execute(serviceorder.AddServiceToOrderInput{
		OrderID:   order.ID(),
		ServiceID: svc.ID(),
	})
	require.NoError(t, err)

	sendBudgetUC := serviceorder.NewSendBudgetUseCase(orderRepo, noopNotifier{})
	_, err = sendBudgetUC.Execute(order.ID())
	require.NoError(t, err)

	rejectUC := serviceorder.NewRejectBudgetUseCase(orderRepo, noopNotifier{})
	order, err = rejectUC.Execute(order.ID())
	require.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusCancelled, order.Status())
}

func TestCreateServiceOrder_ShouldFailWhenVehicleNotOwnedByCustomer_Integration(t *testing.T) {
	db := newTestDB(t)

	customerRepo := infrarepo.NewGormCustomerRepository(db)
	vehicleRepo := infrarepo.NewGormVehicleRepository(db)
	serviceRepo := infrarepo.NewGormServiceRepository(db)
	partRepo := infrarepo.NewGormPartRepository(db)
	orderRepo := infrarepo.NewGormServiceOrderRepository(db)

	customer1, err := entity.NewCustomer("Alice", "529.982.247-25", "11777777777", "alice@workshop.com")
	require.NoError(t, err)
	require.NoError(t, customerRepo.Create(customer1))

	customer2, err := entity.NewCustomer("Bob", "123.456.789-09", "11666666666", "bob@workshop.com")
	require.NoError(t, err)
	require.NoError(t, customerRepo.Create(customer2))

	vehicle, err := entity.NewVehicle(customer2.ID(), "GHI-9012", "Ford", "Focus", 2019)
	require.NoError(t, err)
	require.NoError(t, vehicleRepo.Create(vehicle))

	createUC := serviceorder.NewCreateServiceOrderUseCase(orderRepo, customerRepo, vehicleRepo, serviceRepo, partRepo)
	_, err = createUC.Execute(serviceorder.CreateServiceOrderInput{
		CustomerID: customer1.ID(),
		VehicleID:  vehicle.ID(),
	})

	assert.ErrorIs(t, err, entity.ErrServiceOrderVehicleNotOwnedByCustomer)
}

func TestCreateServiceOrderWithServicesAndParts_Integration(t *testing.T) {
	db := newTestDB(t)

	customerRepo := infrarepo.NewGormCustomerRepository(db)
	vehicleRepo := infrarepo.NewGormVehicleRepository(db)
	serviceRepo := infrarepo.NewGormServiceRepository(db)
	partRepo := infrarepo.NewGormPartRepository(db)
	orderRepo := infrarepo.NewGormServiceOrderRepository(db)

	customer, err := entity.NewCustomer("Carol", "529.982.247-25", "11555555555", "carol@workshop.com")
	require.NoError(t, err)
	require.NoError(t, customerRepo.Create(customer))

	vehicle, err := entity.NewVehicle(customer.ID(), "JKL-3456", "Fiat", "Argo", 2022)
	require.NoError(t, err)
	require.NoError(t, vehicleRepo.Create(vehicle))

	svc, err := entity.NewService("Alignment", "Wheel alignment", 120.00, 45)
	require.NoError(t, err)
	require.NoError(t, serviceRepo.Create(svc))

	part, err := entity.NewPart("Brake Pad", "Front pad", 89.90, 5)
	require.NoError(t, err)
	require.NoError(t, partRepo.Create(part))

	createUC := serviceorder.NewCreateServiceOrderUseCase(orderRepo, customerRepo, vehicleRepo, serviceRepo, partRepo)
	order, err := createUC.Execute(serviceorder.CreateServiceOrderInput{
		CustomerID: customer.ID(),
		VehicleID:  vehicle.ID(),
		Notes:      "full revision",
		Services:   []string{svc.ID()},
		Parts:      []serviceorder.CreateServiceOrderPartInput{{PartID: part.ID(), Quantity: 2}},
	})
	require.NoError(t, err)

	persisted, err := orderRepo.FindByID(order.ID())
	require.NoError(t, err)
	assert.Equal(t, valueobject.OrderStatusReceived, persisted.Status())
	assert.Len(t, persisted.Items(), 1)
	assert.Len(t, persisted.Parts(), 1)
	assert.Equal(t, 2, persisted.Parts()[0].Quantity())
}
