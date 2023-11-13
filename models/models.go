package models

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Healthz struct {
	Id        int `gorm:"primary_key"`
	UUID      uuid.UUID
	UpdatedAt time.Time `gorm:"autoUpdateTime:true"`
}

func (Healthz) TableName() string {
	return "healthz"
}

type User struct {
	gorm.Model
	Username  string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Password  EncryptedString
	Athlete   Athlete
	CreatedAt time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:true"`
}

type Athlete struct {
	gorm.Model
	UserID       uint `gorm:"unique"`
	StravaId     int  `gorm:"unique"`
	AccessToken  EncryptedString
	RefreshToken EncryptedString
	CreatedAt    time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:true"`
}

func InitDb() {
	SetEncryptionKey([]byte(os.Getenv("ENCRYPTION_KEY")))
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_TZ"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.Migrator().AutoMigrate(&Healthz{})
	DB.Migrator().AutoMigrate(&User{})
	DB.Migrator().AutoMigrate(&Athlete{})
}
