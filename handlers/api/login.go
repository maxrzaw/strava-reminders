package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/maxrzaw/strava-reminders/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var GENERIC_USER_PASS_ERROR = map[string]string{"message": "Username or password is incorrect"}
var STRAVA_REMINDERS_JWT_KEY = "STRAVA_REMINDERS_JWT"
var JWT_SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))

func Signup(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if username == "" || email == "" || password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Username, email and password are required",
		})
	}

	if err := checkPasswordRequirements(password); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: models.EncryptedString(hashed_password),
	}

	result := models.DB.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Username or email already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	token, err := generateJWT(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	cookie := createCookie(token)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{})
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Username and password are required",
		})
	}

	var user models.User
	result := models.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusUnauthorized, GENERIC_USER_PASS_ERROR)
		}
		return echo.ErrInternalServerError
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, GENERIC_USER_PASS_ERROR)
	}

	token, err := generateJWT(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	cookie := createCookie(token)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{})
}

func Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     STRAVA_REMINDERS_JWT_KEY,
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{})
}

func Validate(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.RegisteredClaims)

	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("User %s is valid", claims.Subject),
	})
}

func checkPasswordRequirements(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	has_upper := false
	has_special := false
	for _, c := range password {
		if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			has_special = true
		}
		if unicode.IsUpper(c) {
			has_upper = true
		}
		if has_upper && has_special {
			break
		}
	}
	if !has_upper {
		return errors.New("Password must contain at least one uppercase letter")
	}
	if !has_special {
		return errors.New("Password must contain at least one special character")
	}
	return nil
}

func generateJWT(user models.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(user.ID),
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
	cookie.SameSite = http.SameSiteStrictMode
	return cookie
}
