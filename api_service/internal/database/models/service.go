package models

import "eventryx.api_service/internal/database"

type Service struct {
	Id      *int    `gorm:"primary_key" json:"id"`
	Name    *string `gorm:"unique,not null" json:"name"`
	Owner   *User   `gorm:"foreignKey:OwnerId" json:"-"`
	OwnerId *int    `json:"owner_id"`
}

func (service *Service) Create() error {
	return database.Connection.Create(&service).Error
}

func (service *Service) Get() bool {
	return database.Connection.Where(service).First(service).RowsAffected > 0
}

func (service *Service) Exists() bool {
	var count int64
	database.Connection.Model(Service{}).Where(service).Count(&count)
	return count > 0
}
