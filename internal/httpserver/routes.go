package httpserver

import (
	"github.com/labstack/echo/v4"

	"url-shortener/internal/httpserver/handlers"
)

func RegisterRoutes(e *echo.Echo, urlHandler *handlers.URLHandler) {
	e.GET("/health", urlHandler.Health)
	e.POST("/api/v1/shorten", urlHandler.CreateShortURL)
	e.GET("/:code", urlHandler.Redirect)
}
