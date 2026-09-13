package handler

import (
	"log/slog"
	"net/http"

	"github.com/gabrielcamargo/oficina-api/internal/api/dto"
	"github.com/gabrielcamargo/oficina-api/internal/api/response"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/observability"
	"github.com/gin-gonic/gin"
)

const (
	eventServiceOrderCreated       = "service_order.created"
	eventServiceOrderStatusChanged = "service_order.status_changed"
	eventServiceOrderFailed        = "service_order.failed"
)

func respondCreatedOrder(c *gin.Context, order *entity.ServiceOrder, err error) {
	if err != nil {
		failOrder(c, "create", c.Param("id"), err)
		return
	}

	observability.Event(c.Request.Context(), eventServiceOrderCreated,
		slog.String("order_id", order.ID()),
		slog.String("customer_id", order.CustomerID()),
		slog.String("vehicle_id", order.VehicleID()),
		slog.Int("services", len(order.Items())),
		slog.Int("parts", len(order.Parts())),
	)
	c.JSON(http.StatusCreated, dto.ToServiceOrderResponse(order))
}

func respondOrderTransition(c *gin.Context, action string, order *entity.ServiceOrder, err error) {
	if err != nil {
		failOrder(c, action, c.Param("id"), err)
		return
	}

	attrs := []any{
		slog.String("order_id", order.ID()),
		slog.String("action", action),
		slog.String("status", order.Status().String()),
		slog.Float64("total_amount", order.TotalAmount()),
	}
	if transition := order.LastTransition(); transition != nil {
		attrs = append(attrs,
			slog.String("from_status", transition.From.String()),
			slog.String("to_status", transition.To.String()),
			slog.Float64("duration_minutes", transition.Duration.Minutes()),
		)
	}

	observability.Event(c.Request.Context(), eventServiceOrderStatusChanged, attrs...)
	c.JSON(http.StatusOK, dto.ToServiceOrderResponse(order))
}

func failOrder(c *gin.Context, action, orderID string, err error) {
	observability.Failure(c.Request.Context(), eventServiceOrderFailed, err,
		slog.String("action", action),
		slog.String("order_id", orderID),
	)
	response.Error(c, err)
}
