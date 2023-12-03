package handlers

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/strava"
	"github.com/maxrzaw/strava-reminders/handlers/api"
	apihandlers "github.com/maxrzaw/strava-reminders/handlers/api"
	html "github.com/maxrzaw/strava-reminders/handlers/html"
	middleware "github.com/maxrzaw/strava-reminders/middleware"
)

func authSkipper(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "auth") || strings.Contains(c.Request().URL.Path, "login") {
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
	// Set up Goth
	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("JWT_SECRET")))
	goth.UseProviders(strava.New(os.Getenv("STRAVA_CLIENT_ID"), os.Getenv("STRAVA_CLIENT_SECRET"), "http://localhost:8080/auth/callback?provider=strava"))

	// Must come before the JWT middleware
	e.Use(middleware.MissingCookieRedirectWithConfig(middleware.MissingCookieMiddlewareConfig{
		Skipper:     authSkipper,
		TokenLookup: apihandlers.STRAVA_REMINDERS_JWT_KEY,
		RedirectURL: "/auth?provider=strava",
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
	// The Athlete Context Middleware should come last
	e.Use(middleware.AthleteContextMiddlewareWithConfig(middleware.AthleteContextMiddlewareConfig{
		JWTLookup: "user",
	}))

	// Add Auth Routes
	e.GET("/auth", func(c echo.Context) error {
		gothic.BeginAuthHandler(c.Response(), c.Request())
		return nil
	})

	e.GET("/auth/callback", func(c echo.Context) error {
		user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
		if err != nil {
			return err
		}

		return api.Login(c, user)
	})

	e.GET("/auth/logout", func(c echo.Context) error {
		gothic.Logout(c.Response(), c.Request())
		return apihandlers.Logout(c)
	})

	// Add the Routes
	e.GET("/", html.Index)
	e.GET("/login", html.Login)

	api := e.Group("/api")
	api.GET("/healthz", apihandlers.Healthz)
	api.GET("/validate", apihandlers.Validate)
}
