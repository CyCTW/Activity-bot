package models

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

func (user *User) Create() error {
	// Check exists
	// var new_user User
	log.Print("Before user")
	log.Print(user)
	err := DB.Where("line_user_id = ?", user.LineUserID).First(&user).Error
	log.Print("After user")
	log.Print(user)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := DB.Create(&user).Error; err != nil {
			return errors.New("Create fail")
		}
		log.Print("User Create success")

	} else {
		// Get user
		log.Print("User exists")
	}

	return nil
}

func (user *User) GetByID(userID string) error {
	err := DB.First(&user, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *User) GetByLineID(userLineID string) error {
	log.Println("LineID")
	log.Println(userLineID)
	err := DB.Where("line_user_id = ?", userLineID).First(&user).Error
	if err != nil {
		return err
	}
	// Check user token exists
	if user.AccessToken == "" {
		log.Print("No AC Token!")
		return errors.New("No Access token")
	}
	return nil
}

func (user *User) StoreAccessToken(userLineID string, username string, access_token string) error {
	// Check if user exists

	err := DB.Where("line_user_id = ?", userLineID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User not exists
		user = &User{LineUserID: userLineID, Name: username, AccessToken: access_token}
		if err := DB.Create(&user).Error; err != nil {
			return errors.New("Create fail")
		}
		log.Print("User Create success")

	} else {
		log.Print("User exists")
		user.AccessToken = access_token
		DB.Save(&user)
	}
	return nil
}
