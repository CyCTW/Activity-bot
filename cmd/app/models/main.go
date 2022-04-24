package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	Name  string
	Time  string
	Place string
	Users []*User `gorm:"many2many:participations"`
}

type User struct {
	gorm.Model
	user_id    string
	token      string
	name       string
	Activities []*Activity `gorm:"many2many:participations"`
}

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:root@tcp(127.0.0.1:3306)/gocrud?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect db")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&User{})

	DB = db
}
