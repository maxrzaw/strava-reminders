package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// create a new struct for configuring the middleware
type MissingCookieMiddlewareConfig struct {
	Skipper     func(c echo.Context) bool
	TokenLookup string
	RedirectURL string
}
type ExpiredJWTMiddlewareConfig struct {
	Skipper     func(c echo.Context) bool
	RedirectURL string
	JWTLookup   string
}

func MissingCookieRedirectWithConfig(config MissingCookieMiddlewareConfig) echo.MiddlewareFunc {
	if config.TokenLookup == "" {
		config.TokenLookup = "user"
	}
	if config.RedirectURL == "" {
		config.RedirectURL = "/login"
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}
			_, err := c.Cookie(config.TokenLookup)
			if err != nil {
				fmt.Printf("redirecting from %s due to missing cookie", c.Request().URL.RequestURI())
				return c.Redirect(http.StatusFound, config.RedirectURL)
			}

			return next(c)
		}
	}
}

func ExpiredJWTMiddlewareRedirectWithConfig(config ExpiredJWTMiddlewareConfig) echo.MiddlewareFunc {
	if config.JWTLookup == "" {
		config.JWTLookup = "user"
	}
	if config.RedirectURL == "" {
		config.RedirectURL = "/login"
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}
			// we can use echo to get the claims. This middleware must come after the JWT middleware
			tok := c.Get(config.JWTLookup)
			if tok == nil {
				return c.Redirect(http.StatusFound, config.RedirectURL)
			}
			jwtUser, ok := c.Get(config.JWTLookup).(*jwt.Token)
			if !ok {
				return c.Redirect(http.StatusFound, config.RedirectURL)
			}
			claims, ok := jwtUser.Claims.(*jwt.RegisteredClaims)
			if !ok {
				return c.Redirect(http.StatusFound, config.RedirectURL)
			}
			if claims.ExpiresAt.Unix() < time.Now().Unix() {
				return c.Redirect(http.StatusFound, config.RedirectURL)
			}
			return next(c)
		}
	}
}
