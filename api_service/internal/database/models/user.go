package models

import (
	"errors"

	"eventryx.api_service/internal/database"
	"eventryx.api_service/internal/utils"
)

type UserRole string

const (
	IsAdmin UserRole = "admin"
	IsUser  UserRole = "user"
)

type User struct {
	Id       *int     `gorm:"primary_key"`
	Name     *string  `gorm:"unique,not null"`
	Password *string  `gorm:"unique,not null"`
	Role     UserRole `sql:"type:user_role"`
}

func (user *User) Create() error {
	if user.Password == nil {
		return errors.New("password is empty")
	}
	encryptedPassword, err := utils.EncryptPassword(*user.Password)
	if err != nil {
		return err
	}
	user.Password = &encryptedPassword

	return database.Connection.Create(user).Error
}

func (user *User) Get() bool {
	return database.Connection.Where(user).First(user).RowsAffected > 0
}

func (user *User) Exists() bool {
	var count int64
	database.Connection.Model(User{}).Where(user).Count(&count)
	return count > 0
}
