package models

import "errors"

func (activity *Activity) Create() error {
	if err := DB.Create(&activity).Error; err != nil {
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

func (activity *Activity) GetInfo() (*[]User, error) {
	var users []User
	DB.Model(&activity).Association("Users").Find(&users)

	return &users, nil
}

func (activity *Activity) AddParticipants(user *User) error {

	err := DB.Model(&activity).Association("Users").Append(user)
	if err != nil {
		return err
	}
	return nil
}
