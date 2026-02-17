package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ConnectionString string
var Connection *gorm.DB

func InitConnectionString(host, username, password, port, dbName string) {
	ConnectionString = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		host, username, password, dbName, port,
	)
}

func CreateConnection() {
	db, err := gorm.Open(postgres.Open(ConnectionString), &gorm.Config{})
	if err != nil {
		panic("Error while connecting to database:\n" + err.Error())
	}
	Connection = db
}
