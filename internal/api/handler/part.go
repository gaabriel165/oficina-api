package handler

import (
	"net/http"

	"github.com/gabrielcamargo/oficina-api/internal/api/dto"
	"github.com/gabrielcamargo/oficina-api/internal/api/response"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/part"
	"github.com/gin-gonic/gin"
)

type PartHandler struct {
	createUC      *part.CreatePartUseCase
	updateUC      *part.UpdatePartUseCase
	deleteUC      *part.DeletePartUseCase
	getUC         *part.GetPartUseCase
	listUC        *part.ListPartsUseCase
	updateStockUC *part.UpdateStockUseCase
}

func NewPartHandler(
	createUC *part.CreatePartUseCase,
	updateUC *part.UpdatePartUseCase,
	deleteUC *part.DeletePartUseCase,
	getUC *part.GetPartUseCase,
	listUC *part.ListPartsUseCase,
	updateStockUC *part.UpdateStockUseCase,
) *PartHandler {
	return &PartHandler{
		createUC:      createUC,
		updateUC:      updateUC,
		deleteUC:      deleteUC,
		getUC:         getUC,
		listUC:        listUC,
		updateStockUC: updateStockUC,
	}
}

func (h *PartHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/parts", h.Create)
	router.GET("/parts", h.List)
	router.GET("/parts/:id", h.Get)
	router.PUT("/parts/:id", h.Update)
	router.DELETE("/parts/:id", h.Delete)
	router.PATCH("/parts/:id/stock", h.UpdateStock)
}

// Create godoc
// @Summary      Create a new part
// @Tags         parts
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreatePartRequest  true  "Part data"
// @Success      201   {object}  dto.PartResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      422   {object}  response.ErrorResponse
// @Router       /parts [post]
func (h *PartHandler) Create(c *gin.Context) {
	var req dto.CreatePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.createUC.Execute(part.CreatePartInput{
		Name:          req.Name,
		Description:   req.Description,
		UnitPrice:     req.UnitPrice,
		StockQuantity: req.StockQuantity,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToPartResponse(result))
}

// Update godoc
// @Summary      Update a part
// @Tags         parts
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "Part ID"
// @Param        body  body      dto.UpdatePartRequest  true  "Part data"
// @Success      200   {object}  dto.PartResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /parts/{id} [put]
func (h *PartHandler) Update(c *gin.Context) {
	var req dto.UpdatePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.updateUC.Execute(part.UpdatePartInput{
		ID:          c.Param("id"),
		Name:        req.Name,
		Description: req.Description,
		UnitPrice:   req.UnitPrice,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToPartResponse(result))
}

// Delete godoc
// @Summary      Delete a part
// @Tags         parts
// @Security     BearerAuth
// @Param        id  path  string  true  "Part ID"
// @Success      204
// @Failure      404  {object}  response.ErrorResponse
// @Router       /parts/{id} [delete]
func (h *PartHandler) Delete(c *gin.Context) {
	if err := h.deleteUC.Execute(c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary      Get a part by ID
// @Tags         parts
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Part ID"
// @Success      200  {object}  dto.PartResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /parts/{id} [get]
func (h *PartHandler) Get(c *gin.Context) {
	result, err := h.getUC.Execute(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToPartResponse(result))
}

// List godoc
// @Summary      List all parts
// @Tags         parts
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   dto.PartResponse
// @Router       /parts [get]
func (h *PartHandler) List(c *gin.Context) {
	results, err := h.listUC.Execute()
	if err != nil {
		response.Error(c, err)
		return
	}

	responses := make([]dto.PartResponse, 0, len(results))
	for _, r := range results {
		responses = append(responses, dto.ToPartResponse(r))
	}
	c.JSON(http.StatusOK, responses)
}

// UpdateStock godoc
// @Summary      Add stock quantity to a part
// @Tags         parts
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "Part ID"
// @Param        body  body      dto.UpdateStockRequest  true  "Stock data"
// @Success      200   {object}  dto.PartResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /parts/{id}/stock [patch]
func (h *PartHandler) UpdateStock(c *gin.Context) {
	var req dto.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.updateStockUC.Execute(part.UpdateStockInput{
		ID:       c.Param("id"),
		Quantity: req.Quantity,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToPartResponse(result))
}
