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
var Hc_uuid uuid.UUID

type HealthCheck struct {
	Id        int `gorm:"primary_key"`
	UUID      uuid.UUID
	UpdatedAt time.Time `gorm:"autoUpdateTime:true"`
}

type Athlete struct {
	Id           int `gorm:"primary_key"`
	Username     string
	AccessToken  EncryptedString
	RefreshToken EncryptedString
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

	DB.Migrator().AutoMigrate(&HealthCheck{})
	Hc_uuid = uuid.New()
	hcm := HealthCheck{UUID: Hc_uuid}
	DB.Create(&hcm)
	DB.Migrator().AutoMigrate(&Athlete{})
}
