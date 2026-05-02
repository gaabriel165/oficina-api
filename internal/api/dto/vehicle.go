package dto

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type CreateVehicleRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	Plate      string `json:"plate" binding:"required"`
	Brand      string `json:"brand" binding:"required"`
	Model      string `json:"model" binding:"required"`
	Year       int    `json:"year" binding:"required"`
}

type UpdateVehicleRequest struct {
	Brand string `json:"brand" binding:"required"`
	Model string `json:"model" binding:"required"`
	Year  int    `json:"year" binding:"required"`
}

type VehicleResponse struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Plate      string    `json:"plate"`
	Brand      string    `json:"brand"`
	Model      string    `json:"model"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToVehicleResponse(v *entity.Vehicle) VehicleResponse {
	return VehicleResponse{
		ID:         v.ID(),
		CustomerID: v.CustomerID(),
		Plate:      v.Plate(),
		Brand:      v.Brand(),
		Model:      v.Model(),
		Year:       v.Year(),
		CreatedAt:  v.CreatedAt(),
		UpdatedAt:  v.UpdatedAt(),
	}
}
