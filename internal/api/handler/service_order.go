package handler

import (
	"net/http"

	"github.com/gabrielcamargo/oficina-api/internal/api/dto"
	"github.com/gabrielcamargo/oficina-api/internal/api/middleware"
	"github.com/gabrielcamargo/oficina-api/internal/api/response"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/serviceorder"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type ServiceOrderHandler struct {
	createUC     *serviceorder.CreateServiceOrderUseCase
	startDiagUC  *serviceorder.StartDiagnosisUseCase
	addServiceUC *serviceorder.AddServiceToOrderUseCase
	addPartUC    *serviceorder.AddPartToOrderUseCase
	sendBudgetUC *serviceorder.SendBudgetUseCase
	approveUC    *serviceorder.ApproveBudgetUseCase
	rejectUC     *serviceorder.RejectBudgetUseCase
	finishExecUC *serviceorder.FinishExecutionUseCase
	deliverUC    *serviceorder.DeliverVehicleUseCase
	getUC        *serviceorder.GetServiceOrderUseCase
	listUC       *serviceorder.ListServiceOrdersUseCase
	metricsUC    *serviceorder.GetExecutionMetricsUseCase
}

func NewServiceOrderHandler(
	createUC *serviceorder.CreateServiceOrderUseCase,
	startDiagUC *serviceorder.StartDiagnosisUseCase,
	addServiceUC *serviceorder.AddServiceToOrderUseCase,
	addPartUC *serviceorder.AddPartToOrderUseCase,
	sendBudgetUC *serviceorder.SendBudgetUseCase,
	approveUC *serviceorder.ApproveBudgetUseCase,
	rejectUC *serviceorder.RejectBudgetUseCase,
	finishExecUC *serviceorder.FinishExecutionUseCase,
	deliverUC *serviceorder.DeliverVehicleUseCase,
	getUC *serviceorder.GetServiceOrderUseCase,
	listUC *serviceorder.ListServiceOrdersUseCase,
	metricsUC *serviceorder.GetExecutionMetricsUseCase,
) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		createUC:     createUC,
		startDiagUC:  startDiagUC,
		addServiceUC: addServiceUC,
		addPartUC:    addPartUC,
		sendBudgetUC: sendBudgetUC,
		approveUC:    approveUC,
		rejectUC:     rejectUC,
		finishExecUC: finishExecUC,
		deliverUC:    deliverUC,
		getUC:        getUC,
		listUC:       listUC,
		metricsUC:    metricsUC,
	}
}

func (h *ServiceOrderHandler) RegisterOperatorRoutes(router *gin.RouterGroup) {
	router.GET("/service-orders/metrics", h.GetMetrics)
	router.POST("/service-orders", h.Create)
	router.POST("/service-orders/:id/services", h.AddService)
	router.POST("/service-orders/:id/parts", h.AddPart)
	router.PATCH("/service-orders/:id/start-diagnosis", h.StartDiagnosis)
	router.PATCH("/service-orders/:id/send-budget", h.SendBudget)
	router.PATCH("/service-orders/:id/finish-execution", h.FinishExecution)
	router.PATCH("/service-orders/:id/deliver", h.DeliverVehicle)
}

func (h *ServiceOrderHandler) RegisterSharedRoutes(router *gin.RouterGroup) {
	router.GET("/service-orders", h.List)
	router.GET("/service-orders/:id", h.Get)
	router.PATCH("/service-orders/:id/approve-budget", h.ApproveBudget)
	router.PATCH("/service-orders/:id/reject-budget", h.RejectBudget)
}

func (h *ServiceOrderHandler) ensureCustomerOwnsOrder(c *gin.Context) bool {
	subject, role := middleware.Principal(c)
	if role != middleware.RoleCustomer {
		return true
	}
	if _, err := h.getUC.ExecuteForCustomer(c.Param("id"), subject); err != nil {
		response.Error(c, err)
		return false
	}
	return true
}

func (h *ServiceOrderHandler) RegisterPublicRoutes(router *gin.RouterGroup) {
	router.GET("/service-orders/:id/status", h.GetStatus)
}

func (h *ServiceOrderHandler) RegisterWebhookRoutes(router *gin.RouterGroup) {
	router.POST("/service-orders/:id/budget-approval", h.BudgetApprovalWebhook)
}

// GetMetrics godoc
// @Summary      Get average execution time metrics for service orders
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.ServiceOrderMetricsResponse
// @Router       /service-orders/metrics [get]
func (h *ServiceOrderHandler) GetMetrics(c *gin.Context) {
	result, err := h.metricsUC.Execute()
	if err != nil {
		response.Error(c, err)
		return
	}

	byService := make([]dto.ServiceMetricRowResponse, 0, len(result.ByService))
	for _, row := range result.ByService {
		byService = append(byService, dto.ServiceMetricRowResponse{
			ServiceName:     row.ServiceName,
			AvgMinutes:      row.AvgMinutes,
			CompletedOrders: row.CompletedOrders,
		})
	}

	c.JSON(http.StatusOK, dto.ServiceOrderMetricsResponse{
		OverallAvgMinutes: result.OverallAvgMinutes,
		ByService:         byService,
	})
}

// GetStatus godoc
// @Summary      Get the current status of a service order (public)
// @Tags         service-orders
// @Produce      json
// @Param        id  path      string  true  "Service Order ID"
// @Success      200  {object}  dto.ServiceOrderStatusResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /service-orders/{id}/status [get]
func (h *ServiceOrderHandler) GetStatus(c *gin.Context) {
	result, err := h.getUC.Execute(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToServiceOrderStatusResponse(result))
}

// Create godoc
// @Summary      Create a new service order
// @Tags         service-orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateServiceOrderRequest  true  "Service order data"
// @Success      201   {object}  dto.ServiceOrderResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /service-orders [post]
func (h *ServiceOrderHandler) Create(c *gin.Context) {
	var req dto.CreateServiceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parts := make([]serviceorder.CreateServiceOrderPartInput, 0, len(req.Parts))
	for _, part := range req.Parts {
		parts = append(parts, serviceorder.CreateServiceOrderPartInput{
			PartID:   part.PartID,
			Quantity: part.Quantity,
		})
	}

	result, err := h.createUC.Execute(serviceorder.CreateServiceOrderInput{
		CustomerID: req.CustomerID,
		VehicleID:  req.VehicleID,
		Notes:      req.Notes,
		Services:   req.Services,
		Parts:      parts,
	})
	respondCreatedOrder(c, result, err)
}

// Get godoc
// @Summary      Get a service order by ID
// @Description  Operators can read any order; customers authenticated by CPF can only read their own.
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Service Order ID"
// @Success      200  {object}  dto.ServiceOrderResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /service-orders/{id} [get]
func (h *ServiceOrderHandler) Get(c *gin.Context) {
	subject, role := middleware.Principal(c)
	var result *entity.ServiceOrder
	var err error
	if role == middleware.RoleCustomer {
		result, err = h.getUC.ExecuteForCustomer(c.Param("id"), subject)
	} else {
		result, err = h.getUC.Execute(c.Param("id"))
	}
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToServiceOrderResponse(result))
}

// List godoc
// @Summary      List active service orders
// @Description  Operators see every active order sorted by urgency; customers authenticated by CPF see only their own.
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   dto.ServiceOrderResponse
// @Router       /service-orders [get]
func (h *ServiceOrderHandler) List(c *gin.Context) {
	subject, role := middleware.Principal(c)
	var results []*entity.ServiceOrder
	var err error
	if role == middleware.RoleCustomer {
		results, err = h.listUC.ExecuteForCustomer(subject)
	} else {
		results, err = h.listUC.Execute()
	}
	if err != nil {
		response.Error(c, err)
		return
	}

	responses := make([]dto.ServiceOrderResponse, 0, len(results))
	for _, r := range results {
		responses = append(responses, dto.ToServiceOrderResponse(r))
	}
	c.JSON(http.StatusOK, responses)
}

// StartDiagnosis godoc
// @Summary      Start diagnosis for a service order
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Service Order ID"
// @Success      200  {object}  dto.ServiceOrderResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Router       /service-orders/{id}/start-diagnosis [patch]
func (h *ServiceOrderHandler) StartDiagnosis(c *gin.Context) {
	result, err := h.startDiagUC.Execute(c.Param("id"))
	respondOrderTransition(c, "start_diagnosis", result, err)
}

// AddService godoc
// @Summary      Add a service to a service order
// @Tags         service-orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                       true  "Service Order ID"
// @Param        body  body      dto.AddServiceToOrderRequest true  "Service data"
// @Success      200   {object}  dto.ServiceOrderResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /service-orders/{id}/services [post]
func (h *ServiceOrderHandler) AddService(c *gin.Context) {
	var req dto.AddServiceToOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.addServiceUC.Execute(serviceorder.AddServiceToOrderInput{
		OrderID:   c.Param("id"),
		ServiceID: req.ServiceID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToServiceOrderResponse(result))
}

// AddPart godoc
// @Summary      Add a part to a service order
// @Tags         service-orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "Service Order ID"
// @Param        body  body      dto.AddPartToOrderRequest true  "Part data"
// @Success      200   {object}  dto.ServiceOrderResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Failure      422   {object}  response.ErrorResponse
// @Router       /service-orders/{id}/parts [post]
func (h *ServiceOrderHandler) AddPart(c *gin.Context) {
	var req dto.AddPartToOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.addPartUC.Execute(serviceorder.AddPartToOrderInput{
		OrderID:  c.Param("id"),
		PartID:   req.PartID,
		Quantity: req.Quantity,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToServiceOrderResponse(result))
}

// SendBudget godoc
// @Summary      Send budget to customer for approval
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Service Order ID"
// @Success      200  {object}  dto.ServiceOrderResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Router       /service-orders/{id}/send-budget [patch]
func (h *ServiceOrderHandler) SendBudget(c *gin.Context) {
	result, err := h.sendBudgetUC.Execute(c.Param("id"))
	respondOrderTransition(c, "send_budget", result, err)
}

// ApproveBudget godoc
// @Summary      Approve the budget and start execution
// @Description  Customers authenticated by CPF can only approve their own order.
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Service Order ID"
// @Success      200  {object}  dto.ServiceOrderResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Router       /service-orders/{id}/approve-budget [patch]
func (h *ServiceOrderHandler) ApproveBudget(c *gin.Context) {
	if !h.ensureCustomerOwnsOrder(c) {
		return
	}
	result, err := h.approveUC.Execute(c.Param("id"))
	respondOrderTransition(c, "approve_budget", result, err)
}

// RejectBudget godoc
// @Summary      Reject the budget and cancel the order
// @Description  Customers authenticated by CPF can only reject their own order.
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Service Order ID"
// @Success      200  {object}  dto.ServiceOrderResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Router       /service-orders/{id}/reject-budget [patch]
func (h *ServiceOrderHandler) RejectBudget(c *gin.Context) {
	if !h.ensureCustomerOwnsOrder(c) {
		return
	}
	result, err := h.rejectUC.Execute(c.Param("id"))
	respondOrderTransition(c, "reject_budget", result, err)
}

// BudgetApprovalWebhook godoc
// @Summary      Receive an external budget approval or rejection notification
// @Tags         service-orders
// @Accept       json
// @Produce      json
// @Param        X-Webhook-Secret  header  string                            true  "Webhook shared secret"
// @Param        id                path    string                            true  "Service Order ID"
// @Param        body              body    dto.BudgetApprovalWebhookRequest  true  "Budget decision"
// @Success      200  {object}  dto.ServiceOrderResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Router       /service-orders/{id}/budget-approval [post]
func (h *ServiceOrderHandler) BudgetApprovalWebhook(c *gin.Context) {
	var req dto.BudgetApprovalWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Decision == dto.BudgetDecisionRejected {
		result, err := h.rejectUC.Execute(c.Param("id"))
		respondOrderTransition(c, "webhook_reject_budget", result, err)
		return
	}

	result, err := h.approveUC.Execute(c.Param("id"))
	respondOrderTransition(c, "approve_budget", result, err)
}

// FinishExecution godoc
// @Summary      Mark execution as finished
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Service Order ID"
// @Success      200  {object}  dto.ServiceOrderResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Router       /service-orders/{id}/finish-execution [patch]
func (h *ServiceOrderHandler) FinishExecution(c *gin.Context) {
	result, err := h.finishExecUC.Execute(c.Param("id"))
	respondOrderTransition(c, "finish_execution", result, err)
}

// DeliverVehicle godoc
// @Summary      Mark vehicle as delivered to customer
// @Tags         service-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Service Order ID"
// @Success      200  {object}  dto.ServiceOrderResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Router       /service-orders/{id}/deliver [patch]
func (h *ServiceOrderHandler) DeliverVehicle(c *gin.Context) {
	result, err := h.deliverUC.Execute(c.Param("id"))
	respondOrderTransition(c, "deliver", result, err)
}
