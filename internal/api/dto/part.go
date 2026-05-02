package dto

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type CreatePartRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	UnitPrice     float64 `json:"unit_price" binding:"required,gt=0"`
	StockQuantity int     `json:"stock_quantity" binding:"min=0"`
}

type UpdatePartRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
}

type UpdateStockRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

type PartResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	UnitPrice     float64   `json:"unit_price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ToPartResponse(p *entity.Part) PartResponse {
	return PartResponse{
		ID:            p.ID(),
		Name:          p.Name(),
		Description:   p.Description(),
		UnitPrice:     p.UnitPrice(),
		StockQuantity: p.StockQuantity(),
		CreatedAt:     p.CreatedAt(),
		UpdatedAt:     p.UpdatedAt(),
	}
}
