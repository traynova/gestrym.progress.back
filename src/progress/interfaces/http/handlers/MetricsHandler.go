package handlers

import (
	"gestrym-progress/src/common/middleware"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/application/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	createUseCase *usecases.CreateBodyMetricsUseCase
	getUseCase    *usecases.GetUserMetricsUseCase
	chartUseCase  *usecases.GetWeightChartUseCase
}

func NewMetricsHandler(createUC *usecases.CreateBodyMetricsUseCase, getUC *usecases.GetUserMetricsUseCase, chartUC *usecases.GetWeightChartUseCase) *MetricsHandler {
	return &MetricsHandler{
		createUseCase: createUC,
		getUseCase:    getUC,
		chartUseCase:  chartUC,
	}
}

func (h *MetricsHandler) Create(c *gin.Context) {
	var req dtos.CreateMetricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	roleID := c.GetUint("role_id")
	
	// Si es un cliente, forzamos que el UserID sea el suyo
	if roleID == middleware.RoleCliente {
		req.UserID = userID
	}

	if err := h.createUseCase.Execute(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Metrics created successfully"})
}

func (h *MetricsHandler) GetByUserID(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	currentUserID := c.GetUint("user_id")
	roleID := c.GetUint("role_id")

	// Los clientes solo pueden ver sus propios datos
	if roleID == middleware.RoleCliente && currentUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own metrics"})
		return
	}

	response, err := h.getUseCase.Execute(c.Request.Context(), uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MetricsHandler) GetWeightChart(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	currentUserID := c.GetUint("user_id")
	roleID := c.GetUint("role_id")

	if roleID == middleware.RoleCliente && currentUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own weight chart"})
		return
	}

	response, err := h.chartUseCase.Execute(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
