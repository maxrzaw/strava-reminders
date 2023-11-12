package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/maxrzaw/strava-reminders/models"
)

type HealthCheckResponse struct {
	Alive    bool   `json:"alive" xml:"alive"`
	Database string `json:"database" xml:"database"`
}

func Healthz(c echo.Context) error {
	status := http.StatusOK
	r := &HealthCheckResponse{
		Alive:    true,
		Database: "healthy",
	}

	// Check that we can create a record, read it, and delete it
	hc_uuid := uuid.New()
	health_check := models.Healthz{UUID: hc_uuid}

	result := models.DB.Create(&health_check)
	if result.Error != nil {
		r.Database = "unhealthy"
		status = http.StatusServiceUnavailable
	}

	result = models.DB.Where("UUID = ?", hc_uuid).First(&health_check)
	if result.Error != nil {
		r.Database = "unhealthy"
		status = http.StatusServiceUnavailable
	}

	result = models.DB.Delete(&health_check)
	if result.Error != nil {
		r.Database = "unhealthy"
		status = http.StatusServiceUnavailable
	}

	return c.JSON(status, r)
}
