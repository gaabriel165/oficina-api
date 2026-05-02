package dto

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type CreateCustomerRequest struct {
	Name     string `json:"name" binding:"required"`
	Document string `json:"document" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type CustomerResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	DocumentType string    `json:"document_type"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToCustomerResponse(c *entity.Customer) CustomerResponse {
	return CustomerResponse{
		ID:           c.ID(),
		Name:         c.Name(),
		Document:     c.Document(),
		DocumentType: string(c.DocumentType()),
		Phone:        c.Phone(),
		Email:        c.Email(),
		CreatedAt:    c.CreatedAt(),
		UpdatedAt:    c.UpdatedAt(),
	}
}
