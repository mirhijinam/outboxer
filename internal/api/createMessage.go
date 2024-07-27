package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mirhijinam/outboxer/internal/model"
	"go.uber.org/zap"
)

type createdMsgReq struct {
	Content string `json:"content"`
}

func (h *Handler) CreateMessage(c *gin.Context) {
	var req createdMsgReq

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("message was not created:", zap.Error(err))
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	createdMsg := model.Message{
		Content: req.Content,
	}

	if err := h.messageService.Create(c, createdMsg); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong, try again later.",
		})
		h.logger.Error("message was not created:", zap.Error(err))
		return
	}

	h.logger.Info("message was created successfully")
	c.Status(http.StatusOK)
}
