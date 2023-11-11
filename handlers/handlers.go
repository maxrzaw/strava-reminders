package handlers

import (
	"github.com/labstack/echo/v4"
)

func AddHandlers(e *echo.Echo) {
	e.Logger.Warn("Adding handlers")

	e.GET("/", Index)

	api := e.Group("/api")
	api.GET("/healthz", Healthz)
}
