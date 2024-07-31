package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getEventStatsResp struct {
	NumOfAll  int `json:"num_of_all"`
	NumOfNew  int `json:"num_of_new"`
	NumOfDone int `json:"num_of_done"`
}

func (h *Handler) GetEventStats(c *gin.Context) {
	// Создаем контекст для запроса
	ctx := context.Background()

	// Получаем статистику через сервисный слой
	stats, err := h.messageService.GetStats(ctx)
	if err != nil {
		// Если произошла ошибка, возвращаем 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Формируем структуру для ответа
	response := getEventStatsResp{
		NumOfAll:  stats["all"],
		NumOfNew:  stats["new"],
		NumOfDone: stats["done"],
	}

	// Возвращаем статистику в ответе
	c.JSON(http.StatusOK, response)
}
