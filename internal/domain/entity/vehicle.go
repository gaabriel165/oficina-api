package entity

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/google/uuid"
)

type Vehicle struct {
	id         string
	customerID string
	plate      string
	brand      string
	model      string
	year       int
	createdAt  time.Time
	updatedAt  time.Time
}

func NewVehicle(customerID, plate, brand, model string, year int) (*Vehicle, error) {
	if customerID == "" {
		return nil, ErrVehicleCustomerRequired
	}
	if brand == "" {
		return nil, ErrVehicleBrandRequired
	}
	if model == "" {
		return nil, ErrVehicleModelRequired
	}
	if year < 1886 || year > time.Now().Year()+1 {
		return nil, ErrVehicleInvalidYear
	}

	plateVO, err := valueobject.NewPlate(plate)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Vehicle{
		id:         uuid.NewString(),
		customerID: customerID,
		plate:      plateVO.Value(),
		brand:      brand,
		model:      model,
		year:       year,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func RestoreVehicle(id, customerID, plate, brand, model string, year int, createdAt, updatedAt time.Time) *Vehicle {
	return &Vehicle{
		id:         id,
		customerID: customerID,
		plate:      plate,
		brand:      brand,
		model:      model,
		year:       year,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

func (v *Vehicle) ID() string           { return v.id }
func (v *Vehicle) CustomerID() string   { return v.customerID }
func (v *Vehicle) Plate() string        { return v.plate }
func (v *Vehicle) Brand() string        { return v.brand }
func (v *Vehicle) Model() string        { return v.model }
func (v *Vehicle) Year() int            { return v.year }
func (v *Vehicle) CreatedAt() time.Time { return v.createdAt }
func (v *Vehicle) UpdatedAt() time.Time { return v.updatedAt }

func (v *Vehicle) Update(brand, model string, year int) error {
	if brand == "" {
		return ErrVehicleBrandRequired
	}
	if model == "" {
		return ErrVehicleModelRequired
	}
	if year < 1886 || year > time.Now().Year()+1 {
		return ErrVehicleInvalidYear
	}
	v.brand = brand
	v.model = model
	v.year = year
	v.updatedAt = time.Now()
	return nil
}
