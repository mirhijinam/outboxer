package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mirhijinam/outboxer/internal/model"
	"go.uber.org/zap"
)

type messageService interface {
	Create(ctx context.Context, msg model.Message) (int, error)
	GetStats(ctx context.Context) (map[string]int, error)
}

type Handler struct {
	messageService messageService
	log            *zap.Logger
}

func New(ms messageService, zl *zap.Logger) (*gin.Engine, error) {
	h := Handler{
		messageService: ms,
		log:            zl,
	}

	r := gin.New()
	r.POST("/api/messages", h.CreateMessage)
	r.GET("/api/event_stats", h.GetEventStats)

	return r, nil
}
