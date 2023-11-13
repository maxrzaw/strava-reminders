package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo) {
	jwt_config := echojwt.Config{
		SigningKey:  JWT_SECRET_KEY,
		TokenLookup: "cookie:" + STRAVA_REMINDERS_JWT_KEY,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
	}

	e.GET("/", Index)

	api := e.Group("/api")
	api.POST("/signup", Signup)
	api.POST("/login", Login)
	api.GET("/healthz", Healthz)

	validate := api.Group("/validate")
	validate.Use(echojwt.WithConfig(jwt_config))
	validate.GET("", Validate)
}
