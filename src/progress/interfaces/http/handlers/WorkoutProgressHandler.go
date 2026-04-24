package handlers

import (
	"gestrym-progress/src/common/middleware"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/application/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkoutProgressHandler struct {
	markUseCase *usecases.MarkWorkoutProgressUseCase
	getUseCase  *usecases.GetWorkoutProgressUseCase
}

func NewWorkoutProgressHandler(markUC *usecases.MarkWorkoutProgressUseCase, getUC *usecases.GetWorkoutProgressUseCase) *WorkoutProgressHandler {
	return &WorkoutProgressHandler{
		markUseCase: markUC,
		getUseCase:  getUC,
	}
}

func (h *WorkoutProgressHandler) Mark(c *gin.Context) {
	var req dtos.MarkWorkoutProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	req.UserID = userID

	if err := h.markUseCase.Execute(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Workout progress marked successfully"})
}

func (h *WorkoutProgressHandler) GetByUserID(c *gin.Context) {
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

	if roleID == middleware.RoleCliente && currentUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own workout progress"})
		return
	}

	response, err := h.getUseCase.Execute(c.Request.Context(), uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
