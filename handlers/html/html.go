package html

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.Redirect(http.StatusFound, "/login")
	}
	claims := user.Claims.(*jwt.RegisteredClaims)
	return c.Render(http.StatusOK, "index.html", map[string]string{
		"subject": claims.Subject,
	})
}
func LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}
func SignupPage(c echo.Context) error {
	return c.Render(http.StatusOK, "signup.html", nil)
}
