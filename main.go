package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maxrzaw/strava-reminders/handlers"
	"github.com/maxrzaw/strava-reminders/models"
	"github.com/maxrzaw/strava-reminders/template"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	models.InitDb()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	// Tailwind files
	e.Static("/dist", "dist")
	e.Static("/font-awesome", "assets/font-awesome/")

	template.NewTemplateRenderer(e,
		template.TemplateRecipe{
			Name:  "index.html",
			Base:  "base.html",
			Paths: []string{"public/index.html", "public/base.html"},
		},
		template.TemplateRecipe{
			Name:  "signup.html",
			Base:  "base.html",
			Paths: []string{"public/signup.html", "public/base.html"},
		},
		template.TemplateRecipe{
			Name:  "login.html",
			Base:  "base.html",
			Paths: []string{"public/login.html", "public/base.html"},
		},
	)

	handlers.AddRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
