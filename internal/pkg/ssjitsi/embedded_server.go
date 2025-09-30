package ssjitsi

import (
	"embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed web/*
var embeddedFiles embed.FS

// NewEmbeddedServer создает HTTP сервер с встроенным UI
func NewEmbeddedServer(server *HttpServer) *gin.Engine {
	router := gin.Default()

	// Настройка CORS
	router.Use(corsMiddleware())

	// API маршруты
	api := router.Group("/api/v1")
	{
		api.GET("/bots", server.ListBots)
		api.GET("/:id/screenshot", server.Screenshot)
	}

	// Обработка всех запросов
	router.NoRoute(func(c *gin.Context) {
		// Если запрос к API, возвращаем 404
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		// Если запрос к статическим файлам, обслуживаем их
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			// Файлы находятся в web/assets/...
			filePath := "web" + c.Request.URL.Path
			fileData, err := embeddedFiles.ReadFile(filePath)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "file not found: " + filePath})
				return
			}

			// Определяем Content-Type по расширению файла
			contentType := "application/octet-stream"
			if strings.HasSuffix(filePath, ".css") {
				contentType = "text/css"
			} else if strings.HasSuffix(filePath, ".js") {
				contentType = "application/javascript"
			} else if strings.HasSuffix(filePath, ".html") {
				contentType = "text/html"
			} else if strings.HasSuffix(filePath, ".svg") {
				contentType = "image/svg+xml"
			}

			c.Data(http.StatusOK, contentType, fileData)
			return
		}

		// Для всех остальных запросов отдаем index.html
		indexHTML, err := embeddedFiles.ReadFile("web/index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load UI"})
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	return router
}

// corsMiddleware настраивает CORS для всех запросов
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
