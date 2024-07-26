package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirhijinam/outboxer/internal/model"
)

type msgCreatedReq struct {
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"created_at" binding:"required,date"`
}

func (h *Handler) CreateMessage(c *gin.Context) {
	var req msgCreatedReq

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	createdMsg := model.Message{
		Content:   req.Content,
		CreatedAt: req.CreatedAt,
	}

	if err := h.messageService.Create(c, createdMsg); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong, try again later.",
		})

		h.logger.Error(err.Error())

		return
	}

	c.Status(http.StatusOK)
}
