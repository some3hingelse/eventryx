package models

import (
	"time"

	"eventryx.api_service/config"
	"eventryx.api_service/internal/database"
	"eventryx.api_service/internal/utils"
)

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

type ServiceToken struct {
	Id        *int       `gorm:"primary_key" json:"id"`
	Value     *string    `gorm:"not null,unique" json:"value"`
	ServiceId *int       `json:"service_id"`
	Service   *Service   `gorm:"foreignKey:ServiceId" json:"-"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (serviceToken *ServiceToken) Create() error {
	for {
		tokenValue := utils.GenerateRandomString(config.Config.ServiceTokenValueLength)
		serviceToken.Value = &tokenValue
		if !serviceToken.ExistsWithValue() {
			break
		}
	}

	return database.Connection.Create(&serviceToken).Error
}

func (serviceToken *ServiceToken) Get() bool {
	return database.Connection.Where(serviceToken).Preload("Service").First(serviceToken).RowsAffected > 0
}

func (serviceToken *ServiceToken) ExistsWithValue() bool {
	var count int64
	database.Connection.Model(ServiceToken{}).Where("value = ?", serviceToken.Value).Count(&count)
	return count > 0
}
