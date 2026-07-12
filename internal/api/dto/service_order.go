package dto

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
)

type CreateServiceOrderPartRequest struct {
	PartID   string `json:"part_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type CreateServiceOrderRequest struct {
	CustomerID string                          `json:"customer_id" binding:"required"`
	VehicleID  string                          `json:"vehicle_id" binding:"required"`
	Notes      string                          `json:"notes"`
	Services   []string                        `json:"services"`
	Parts      []CreateServiceOrderPartRequest `json:"parts"`
}

const (
	BudgetDecisionApproved = "approved"
	BudgetDecisionRejected = "rejected"
)

type BudgetApprovalWebhookRequest struct {
	Decision string `json:"decision" binding:"required,oneof=approved rejected"`
}

type AddServiceToOrderRequest struct {
	ServiceID string `json:"service_id" binding:"required"`
}

type AddPartToOrderRequest struct {
	PartID   string `json:"part_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type ServiceOrderItemResponse struct {
	ID          string  `json:"id"`
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	LaborPrice  float64 `json:"labor_price"`
}

type ServiceOrderPartResponse struct {
	ID        string  `json:"id"`
	PartID    string  `json:"part_id"`
	PartName  string  `json:"part_name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Total     float64 `json:"total"`
}

type ServiceOrderResponse struct {
	ID          string                     `json:"id"`
	CustomerID  string                     `json:"customer_id"`
	VehicleID   string                     `json:"vehicle_id"`
	Status      valueobject.OrderStatus    `json:"status"`
	Notes       string                     `json:"notes"`
	TotalAmount float64                    `json:"total_amount"`
	Items       []ServiceOrderItemResponse `json:"items"`
	Parts       []ServiceOrderPartResponse `json:"parts"`
	StartedAt   *time.Time                 `json:"started_at"`
	FinishedAt  *time.Time                 `json:"finished_at"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}

type ServiceMetricRowResponse struct {
	ServiceName     string  `json:"service_name"`
	AvgMinutes      float64 `json:"avg_minutes"`
	CompletedOrders int     `json:"completed_orders"`
}

type ServiceOrderMetricsResponse struct {
	OverallAvgMinutes float64                      `json:"overall_avg_minutes"`
	ByService         []ServiceMetricRowResponse   `json:"by_service"`
}

type ServiceOrderStatusResponse struct {
	ID          string                  `json:"id"`
	Status      valueobject.OrderStatus `json:"status"`
	TotalAmount float64                 `json:"total_amount"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

func ToServiceOrderStatusResponse(o *entity.ServiceOrder) ServiceOrderStatusResponse {
	return ServiceOrderStatusResponse{
		ID:          o.ID(),
		Status:      o.Status(),
		TotalAmount: o.TotalAmount(),
		UpdatedAt:   o.UpdatedAt(),
	}
}

func ToServiceOrderResponse(o *entity.ServiceOrder) ServiceOrderResponse {
	items := make([]ServiceOrderItemResponse, 0, len(o.Items()))
	for _, item := range o.Items() {
		items = append(items, ServiceOrderItemResponse{
			ID:          item.ID(),
			ServiceID:   item.ServiceID(),
			ServiceName: item.ServiceName(),
			LaborPrice:  item.LaborPrice(),
		})
	}

	parts := make([]ServiceOrderPartResponse, 0, len(o.Parts()))
	for _, part := range o.Parts() {
		parts = append(parts, ServiceOrderPartResponse{
			ID:        part.ID(),
			PartID:    part.PartID(),
			PartName:  part.PartName(),
			Quantity:  part.Quantity(),
			UnitPrice: part.UnitPrice(),
			Total:     part.TotalPrice(),
		})
	}

	return ServiceOrderResponse{
		ID:          o.ID(),
		CustomerID:  o.CustomerID(),
		VehicleID:   o.VehicleID(),
		Status:      o.Status(),
		Notes:       o.Notes(),
		TotalAmount: o.TotalAmount(),
		Items:       items,
		Parts:       parts,
		StartedAt:   o.StartedAt(),
		FinishedAt:  o.FinishedAt(),
		CreatedAt:   o.CreatedAt(),
		UpdatedAt:   o.UpdatedAt(),
	}
}
