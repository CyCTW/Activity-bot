package models

import (
	"errors"
	"log"
)

func (activity *Activity) Create() error {
	if err := DB.Omit("Users").Create(&activity).Error; err != nil {
		return errors.New("Create fail")
	}

	return nil
}

func (activity *Activity) GetByID(activityID string) error {
	err := DB.First(&activity, activityID).Error
	if err != nil {
		return err
	}
	return nil
}

func (activity *Activity) GetByName(activityName string) error {
	err := DB.Where("name = ?", activityName).First(&activity).Error
	if err != nil {
		return err
	}
	return nil
}

type Participation struct {
	ActivityID uint
	UserID     uint
}

func (activity *Activity) GetInfo(activityID string) (*[]APIUser, error) {
	var users []APIUser

	err := DB.Model(&activity).Association("Users").Find(&users).Error
	if err != nil {
		log.Print(err())
		return nil, errors.New("Get user info fail")
	}

	return &users, nil
}

func (activity *Activity) AddParticipants(user *User) error {

	err := DB.Model(&activity).Association("Users").Append(user)
	if err != nil {
		return err
	}
	return nil
}
