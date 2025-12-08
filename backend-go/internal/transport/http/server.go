package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler - обработчик HTTP запросов.
type Handler struct {
	// Здесь будут зависимости, например, сервис для работы с пользователями.
}

// NewHandler создает новый экземпляр Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// InitRoutes инициализирует роутинг для HTTP-сервера.
func (h *Handler) InitRoutes() *gin.Engine {
	// Создаем новый роутер Gin с настройками по умолчанию (логгер, рекавери).
	router := gin.Default()

	// Настраиваем CORS, чтобы бот мог делать запросы к API.
	// На данном этапе разрешаем все источники, в проде нужно будет указать конкретный.
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Тестовый роут для проверки работоспособности сервера.
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	
	// Группа роутов для API v1.
	apiV1 := router.Group("/api/v1")
	{
		// Здесь будут роуты для пользователей, подписок и т.д.
		apiV1.GET("/status", h.getStatus)
	}

	return router
}

// getStatus - пример обработчика для роута /api/v1/status.
func (h *Handler) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"version": "1.0.0",
	})
}
