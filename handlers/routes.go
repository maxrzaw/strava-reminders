package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	apihandlers "github.com/maxrzaw/strava-reminders/handlers/api"
	html "github.com/maxrzaw/strava-reminders/handlers/html"
)

func AddRoutes(e *echo.Echo) {
	jwt_config := echojwt.Config{
		SigningKey:  apihandlers.JWT_SECRET_KEY,
		TokenLookup: "cookie:" + apihandlers.STRAVA_REMINDERS_JWT_KEY,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
	}

	e.GET("/", html.Index)
	e.GET("/login", html.LoginPage)
	e.POST("/login-form", html.LoginForm)
	e.GET("/signup", html.SignupPage)
	e.POST("/signup-form", html.SignupForm)

	api := e.Group("/api")
	api.POST("/signup", apihandlers.Signup)
	api.POST("/login", apihandlers.Login)
	api.GET("/healthz", apihandlers.Healthz)

	validate := api.Group("/validate")
	validate.Use(echojwt.WithConfig(jwt_config))
	validate.GET("", apihandlers.Validate)
}
