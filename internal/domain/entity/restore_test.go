package entity_test

import (
	"testing"
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestRestoreCustomer_ShouldRestoreAllFields(t *testing.T) {
	now := time.Now()
	c := entity.RestoreCustomer("id-1", "John Doe", string(entity.DocumentTypeCPF), "52998224725", "11999999999", "john@email.com", now, now)

	assert.Equal(t, "id-1", c.ID())
	assert.Equal(t, "John Doe", c.Name())
	assert.Equal(t, "52998224725", c.Document())
	assert.Equal(t, now, c.CreatedAt())
	assert.Equal(t, now, c.UpdatedAt())
}

func TestRestoreVehicle_ShouldRestoreAllFields(t *testing.T) {
	now := time.Now()
	v := entity.RestoreVehicle("id-1", "cust-1", "ABC1234", "Toyota", "Corolla", 2020, now, now)

	assert.Equal(t, "id-1", v.ID())
	assert.Equal(t, "cust-1", v.CustomerID())
	assert.Equal(t, "ABC1234", v.Plate())
	assert.Equal(t, now, v.CreatedAt())
	assert.Equal(t, now, v.UpdatedAt())
}

func TestRestorePart_ShouldRestoreAllFields(t *testing.T) {
	now := time.Now()
	p := entity.RestorePart("id-1", "Filtro", "Filtro de óleo", 29.90, 10, now, now)

	assert.Equal(t, "id-1", p.ID())
	assert.Equal(t, "Filtro de óleo", p.Description())
	assert.Equal(t, now, p.CreatedAt())
	assert.Equal(t, now, p.UpdatedAt())
}

func TestRestoreService_ShouldRestoreAllFields(t *testing.T) {
	now := time.Now()
	s := entity.RestoreService("id-1", "Troca de Óleo", "Inclui filtro", 120.00, 60, now, now)

	assert.Equal(t, "id-1", s.ID())
	assert.Equal(t, "Inclui filtro", s.Description())
	assert.Equal(t, now, s.CreatedAt())
	assert.Equal(t, now, s.UpdatedAt())
}

func TestRestoreServiceOrder_ShouldRestoreAllFields(t *testing.T) {
	now := time.Now()
	o := entity.RestoreServiceOrder(
		"id-1", "cust-1", "veh-1",
		valueobject.OrderStatusReceived,
		"observações",
		0,
		nil, nil,
		nil, nil,
		now, now,
	)

	assert.Equal(t, "id-1", o.ID())
	assert.Equal(t, "cust-1", o.CustomerID())
	assert.Equal(t, "veh-1", o.VehicleID())
	assert.Equal(t, "observações", o.Notes())
	assert.Equal(t, now, o.CreatedAt())
	assert.Equal(t, now, o.UpdatedAt())
}

func TestRestoreServiceOrderItem_ShouldRestoreAllFields(t *testing.T) {
	item := entity.RestoreServiceOrderItem("id-1", "order-1", "svc-1", "Troca de Óleo", 120.00)

	assert.Equal(t, "id-1", item.ID())
	assert.Equal(t, "order-1", item.ServiceOrderID())
	assert.Equal(t, "svc-1", item.ServiceID())
	assert.Equal(t, "Troca de Óleo", item.ServiceName())
	assert.Equal(t, 120.00, item.LaborPrice())
}

func TestRestoreServiceOrderPart_ShouldRestoreAllFields(t *testing.T) {
	part := entity.RestoreServiceOrderPart("id-1", "order-1", "part-1", "Filtro", 2, 29.90)

	assert.Equal(t, "id-1", part.ID())
	assert.Equal(t, "order-1", part.ServiceOrderID())
	assert.Equal(t, "part-1", part.PartID())
	assert.Equal(t, "Filtro", part.PartName())
	assert.Equal(t, 2, part.Quantity())
	assert.Equal(t, 29.90, part.UnitPrice())
	assert.Equal(t, 59.80, part.TotalPrice())
}

func TestNewUser_ShouldCreateSuccessfully(t *testing.T) {
	u, err := entity.NewUser("admin@test.com", "hashed_password")

	assert.NoError(t, err)
	assert.NotEmpty(t, u.ID())
	assert.Equal(t, "admin@test.com", u.Email())
	assert.Equal(t, "hashed_password", u.PasswordHash())
	assert.False(t, u.CreatedAt().IsZero())
	assert.False(t, u.UpdatedAt().IsZero())
}

func TestNewUser_ShouldReturnErrorWhenEmailEmpty(t *testing.T) {
	_, err := entity.NewUser("", "hashed_password")
	assert.ErrorIs(t, err, entity.ErrUserEmailRequired)
}

func TestNewUser_ShouldReturnErrorWhenPasswordEmpty(t *testing.T) {
	_, err := entity.NewUser("admin@test.com", "")
	assert.ErrorIs(t, err, entity.ErrUserPasswordRequired)
}

func TestRestoreUser_ShouldRestoreAllFields(t *testing.T) {
	now := time.Now()
	u := entity.RestoreUser("id-1", "admin@test.com", "hashed_password", now, now)

	assert.Equal(t, "id-1", u.ID())
	assert.Equal(t, "admin@test.com", u.Email())
	assert.Equal(t, "hashed_password", u.PasswordHash())
	assert.Equal(t, now, u.CreatedAt())
	assert.Equal(t, now, u.UpdatedAt())
}
