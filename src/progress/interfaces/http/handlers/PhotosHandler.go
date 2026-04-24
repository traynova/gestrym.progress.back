package handlers

import (
	"gestrym-progress/src/common/middleware"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/application/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhotosHandler struct {
	uploadUseCase *usecases.UploadProgressPhotoUseCase
	getUseCase    *usecases.GetUserPhotosUseCase
}

func NewPhotosHandler(uploadUC *usecases.UploadProgressPhotoUseCase, getUC *usecases.GetUserPhotosUseCase) *PhotosHandler {
	return &PhotosHandler{
		uploadUseCase: uploadUC,
		getUseCase:    getUC,
	}
}

func (h *PhotosHandler) Upload(c *gin.Context) {
	var req dtos.UploadPhotoRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	userID := c.GetUint("user_id")
	roleID := c.GetUint("role_id")
	
	// Si es un cliente, forzamos que el UserID sea el suyo
	if roleID == middleware.RoleCliente {
		req.UserID = userID
	}

	if err := h.uploadUseCase.Execute(c.Request.Context(), req, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Photo uploaded successfully"})
}

func (h *PhotosHandler) GetByUserID(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own photos"})
		return
	}

	response, err := h.getUseCase.Execute(c.Request.Context(), uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
