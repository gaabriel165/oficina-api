package entity_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func makeService(t *testing.T) *entity.Service {
	t.Helper()
	s, err := entity.NewService("Oil Change", "Full oil change", 150.00, 60)
	assert.NoError(t, err)
	return s
}

func makePart(t *testing.T) *entity.Part {
	t.Helper()
	p, err := entity.NewPart("Oil Filter", "Standard filter", 29.90, 10)
	assert.NoError(t, err)
	return p
}

func TestNewServiceOrder_ShouldCreateWithStatusReceived(t *testing.T) {
	order, err := entity.NewServiceOrder("customer-id", "vehicle-id", "noise in engine")
	assert.NoError(t, err)
	assert.NotEmpty(t, order.ID())
	assert.Equal(t, valueobject.OrderStatusReceived, order.Status())
	assert.Empty(t, order.Items())
	assert.Empty(t, order.Parts())
}

func TestNewServiceOrder_ShouldReturnErrorWhenCustomerIDEmpty(t *testing.T) {
	_, err := entity.NewServiceOrder("", "vehicle-id", "")
	assert.ErrorIs(t, err, entity.ErrServiceOrderCustomerRequired)
}

func TestNewServiceOrder_ShouldReturnErrorWhenVehicleIDEmpty(t *testing.T) {
	_, err := entity.NewServiceOrder("customer-id", "", "")
	assert.ErrorIs(t, err, entity.ErrServiceOrderVehicleRequired)
}

func TestServiceOrder_ShouldFollowHappyPath(t *testing.T) {
	order, _ := entity.NewServiceOrder("customer-id", "vehicle-id", "noise in engine")

	assert.NoError(t, order.StartDiagnosis())
	assert.Equal(t, valueobject.OrderStatusInDiagnosis, order.Status())

	order.AddItem(makeService(t))
	assert.NoError(t, order.AddPart(makePart(t), 2))

	assert.NoError(t, order.SendBudget())
	assert.Equal(t, valueobject.OrderStatusWaitingApproval, order.Status())
	assert.Equal(t, 209.80, order.TotalAmount())

	assert.NoError(t, order.ApproveBudget())
	assert.Equal(t, valueobject.OrderStatusInExecution, order.Status())
	assert.NotNil(t, order.StartedAt())

	assert.NoError(t, order.FinishExecution())
	assert.Equal(t, valueobject.OrderStatusFinished, order.Status())
	assert.NotNil(t, order.FinishedAt())

	assert.NoError(t, order.Deliver())
	assert.Equal(t, valueobject.OrderStatusDelivered, order.Status())
}

func TestServiceOrder_ShouldFollowRejectionPath(t *testing.T) {
	order, _ := entity.NewServiceOrder("customer-id", "vehicle-id", "")
	order.StartDiagnosis()
	order.AddItem(makeService(t))
	order.SendBudget()

	assert.NoError(t, order.RejectBudget())
	assert.Equal(t, valueobject.OrderStatusCancelled, order.Status())
}

func TestServiceOrder_SendBudget_ShouldReturnErrorWhenNoItems(t *testing.T) {
	order, _ := entity.NewServiceOrder("customer-id", "vehicle-id", "")
	order.StartDiagnosis()
	err := order.SendBudget()
	assert.ErrorIs(t, err, entity.ErrServiceOrderNoItems)
}

func TestServiceOrder_ShouldReturnErrorOnInvalidTransition(t *testing.T) {
	order, _ := entity.NewServiceOrder("customer-id", "vehicle-id", "")
	err := order.ApproveBudget()
	assert.ErrorIs(t, err, entity.ErrServiceOrderInvalidTransition)
}

func TestServiceOrder_AddPart_ShouldReturnErrorWhenQuantityZero(t *testing.T) {
	order, _ := entity.NewServiceOrder("customer-id", "vehicle-id", "")
	order.StartDiagnosis()
	err := order.AddPart(makePart(t), 0)
	assert.ErrorIs(t, err, entity.ErrServiceOrderItemInvalidQty)
}

func TestServiceOrder_CalculateTotalAmount_ShouldSumItemsAndParts(t *testing.T) {
	order, _ := entity.NewServiceOrder("customer-id", "vehicle-id", "")
	order.StartDiagnosis()
	order.AddItem(makeService(t))
	order.AddPart(makePart(t), 2)
	order.SendBudget()

	expected := 150.00 + (29.90 * 2)
	assert.Equal(t, expected, order.TotalAmount())
}
