package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/maxrzaw/strava-reminders/models"
	"gorm.io/gorm"
)

var STRAVA_REMINDERS_JWT_KEY = "STRAVA_REMINDERS_JWT"
var JWT_SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))

func Login(c echo.Context, user goth.User) error {
	userId, err := strconv.ParseUint(user.UserID, 10, 64)
	if err != nil {
		panic(err)
	}
	var athlete models.Athlete
	result := models.DB.Where("strava_user_id = ?", userId).First(&athlete)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			a := &models.Athlete{
				Name:         user.Name,
				NickName:     user.NickName,
				StravaUserID: userId,
				AccessToken:  models.EncryptedString(user.AccessToken),
				RefreshToken: models.EncryptedString(user.RefreshToken),
				ExpiresAt:    user.ExpiresAt,
			}

			models.DB.Create(&a)
			athlete = *a
		}
	}

	token, err := generateJWT(athlete)
	if err != nil {
		panic(err.Error())
	}

	cookie := createCookie(token)
	c.SetCookie(cookie)

	return c.Redirect(301, "/")
}

func Validate(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.RegisteredClaims)

	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("User %s is valid", claims.Subject),
	})
}

func Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     STRAVA_REMINDERS_JWT_KEY,
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}
	c.SetCookie(cookie)

	return c.Redirect(301, "/login")
}

func generateJWT(athlete models.Athlete) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(athlete.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func createCookie(token string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = STRAVA_REMINDERS_JWT_KEY
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30)
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"
	return cookie
}
