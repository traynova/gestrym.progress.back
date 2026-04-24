package handlers

import (
	"gestrym-progress/src/common/middleware"
	"gestrym-progress/src/progress/application/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ComparisonHandler struct {
	comparisonUseCase *usecases.GetProgressComparisonUseCase
}

func NewComparisonHandler(comparisonUC *usecases.GetProgressComparisonUseCase) *ComparisonHandler {
	return &ComparisonHandler{
		comparisonUseCase: comparisonUC,
	}
}

func (h *ComparisonHandler) GetComparison(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	currentUserID := c.GetUint("user_id")
	roleID := c.GetUint("role_id")

	if roleID == middleware.RoleCliente && currentUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own progress comparison"})
		return
	}

	response, err := h.comparisonUseCase.Execute(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
