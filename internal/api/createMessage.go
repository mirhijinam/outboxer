package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirhijinam/outboxer/internal/model"
)

type createdMsgReq struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) CreateMessage(c *gin.Context) {
	var req createdMsgReq

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("message was not created 1")
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	createdMsg := model.Message{
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := h.messageService.Create(c, createdMsg); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong, try again later.",
		})
		log.Println("message was not created 2")
		h.logger.Error(err.Error())

		return
	}

	log.Println("message was created")
	c.Status(http.StatusOK)
}
