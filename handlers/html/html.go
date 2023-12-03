package html

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maxrzaw/strava-reminders/middleware"
)

func Index(c echo.Context) error {
	athleteContext := c.(*middleware.AthleteContext)
	athlete := athleteContext.Athlete
	return c.Render(http.StatusOK, "index.html", map[string]string{
		"name": athlete.Name,
	})
}

func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]string{})
}
