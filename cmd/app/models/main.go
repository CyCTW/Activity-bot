package models

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	Name   string
	Date   time.Time
	Place  string
	UserID uint
	User   User    // Activity Holder
	Users  []*User `gorm:"many2many:participations"` // Attendee
}

type User struct {
	gorm.Model
	LineUserID  string `gorm:"uniqueIndex"`
	AccessToken string
	Name        string
	Activities  []*Activity `gorm:"many2many:participations"`
}

type APIUser struct {
	ID         uint
	LineUserID string
	Name       string
}

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect db")
	}
	if err := db.Migrator().DropTable(&User{}).Error; err != nil {
		log.Fatal("Drop table user fail")
	}
	if err := db.Migrator().DropTable(&Activity{}); err != nil {
		log.Fatal("Drop table activity fail")
	}
	if err := db.Migrator().DropTable(&Participation{}); err != nil {
		log.Fatal("Drop table participation fail")

	}

	if err := db.AutoMigrate(&User{}).Error; err != nil {
		log.Fatal("Migrate table User fail")
	}
	if err := db.AutoMigrate(&Activity{}); err != nil {
		log.Fatal("Migrate table Activity fail")
	}

	DB = db
}
