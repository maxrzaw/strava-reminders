package handlers

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	apihandlers "github.com/maxrzaw/strava-reminders/handlers/api"
	html "github.com/maxrzaw/strava-reminders/handlers/html"
	middleware "github.com/maxrzaw/strava-reminders/middleware"
)

func authSkipper(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "login") || strings.Contains(c.Request().URL.Path, "signup") {
		return true
	}
	if c.Request().URL.Path == "/api/healtz" {
		return true
	}
	if strings.HasPrefix(c.Request().URL.Path, "/dist") || strings.HasPrefix(c.Request().URL.Path, "/font-awesome") {
		return true
	}
	return false
}

func AddRoutes(e *echo.Echo) {
	// Must come before the JWT middleware
	e.Use(middleware.MissingCookieRedirectWithConfig(middleware.MissingCookieMiddlewareConfig{
		Skipper:     authSkipper,
		TokenLookup: apihandlers.STRAVA_REMINDERS_JWT_KEY,
		RedirectURL: "/login",
	}))

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  apihandlers.JWT_SECRET_KEY,
		TokenLookup: "cookie:" + apihandlers.STRAVA_REMINDERS_JWT_KEY,
		Skipper:     authSkipper,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
	}))

	// Must come after the JWT middleware because it uses the "user" key
	e.Use(middleware.ExpiredJWTMiddlewareRedirectWithConfig(middleware.ExpiredJWTMiddlewareConfig{
		Skipper:     authSkipper,
		JWTLookup:   "user",
		RedirectURL: "/login",
	}))
	e.GET("/", html.Index)
	e.GET("/login", html.LoginPage)
	e.POST("/login-form", html.LoginForm)
	e.GET("/signup", html.SignupPage)
	e.POST("/signup-form", html.SignupForm)

	api := e.Group("/api")
	api.POST("/signup", apihandlers.Signup)
	api.POST("/login", apihandlers.Login)
	api.GET("/healthz", apihandlers.Healthz)
	api.GET("/validate", apihandlers.Validate)
}
