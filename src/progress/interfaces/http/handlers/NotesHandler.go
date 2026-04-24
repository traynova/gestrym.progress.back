package handlers

import (
	"gestrym-progress/src/common/middleware"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/application/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotesHandler struct {
	createUseCase *usecases.CreateCoachNoteUseCase
	getUseCase    *usecases.GetUserNotesUseCase
}

func NewNotesHandler(createUC *usecases.CreateCoachNoteUseCase, getUC *usecases.GetUserNotesUseCase) *NotesHandler {
	return &NotesHandler{
		createUseCase: createUC,
		getUseCase:    getUC,
	}
}

func (h *NotesHandler) Create(c *gin.Context) {
	var req dtos.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trainerID := c.GetUint("user_id")
	req.TrainerID = trainerID

	if err := h.createUseCase.Execute(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Note created successfully"})
}

func (h *NotesHandler) GetByUserID(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own notes"})
		return
	}

	response, err := h.getUseCase.Execute(c.Request.Context(), uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
