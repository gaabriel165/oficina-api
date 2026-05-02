package dto

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type CreateServiceRequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	LaborPrice       float64 `json:"labor_price" binding:"required,gt=0"`
	EstimatedMinutes int     `json:"estimated_minutes" binding:"required,gt=0"`
}

type UpdateServiceRequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	LaborPrice       float64 `json:"labor_price" binding:"required,gt=0"`
	EstimatedMinutes int     `json:"estimated_minutes" binding:"required,gt=0"`
}

type ServiceResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	LaborPrice       float64   `json:"labor_price"`
	EstimatedMinutes int       `json:"estimated_minutes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func ToServiceResponse(s *entity.Service) ServiceResponse {
	return ServiceResponse{
		ID:               s.ID(),
		Name:             s.Name(),
		Description:      s.Description(),
		LaborPrice:       s.LaborPrice(),
		EstimatedMinutes: s.EstimatedMinutes(),
		CreatedAt:        s.CreatedAt(),
		UpdatedAt:        s.UpdatedAt(),
	}
}
