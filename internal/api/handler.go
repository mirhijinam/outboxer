package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mirhijinam/outboxer/internal/model"
	"go.uber.org/zap"
)

type messageService interface {
	Create(ctx context.Context, msg model.Message) error
}

type Handler struct {
	messageService messageService
	logger         *zap.Logger
}

func New(ms messageService, zl *zap.Logger) (*gin.Engine, error) {
	h := Handler{
		messageService: ms,
		logger:         zl,
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/api/messages", h.CreateMessage)

	return r, nil
}
