package history

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HistoryHandler struct {
	service *HistoryService
}

func NewHistoryHandler(service *HistoryService) *HistoryHandler {
	return &HistoryHandler{service: service}
}

func (h *HistoryHandler) GetByUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	histories, err := h.service.GetByUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, histories)
}
