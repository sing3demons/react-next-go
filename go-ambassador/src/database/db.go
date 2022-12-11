package database

import (
	"github.com/sing3demons/ambassador/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	var err error
	db, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/ambassador"), &gorm.Config{})
	if err != nil {
		panic("Could not connect with the database")
	}
}

func AutoMigrate() {
	db.AutoMigrate(models.User{}, models.Product{}, models.Link{}, models.Order{}, models.OrderItem{})
}

func GetDB() *gorm.DB {
	return db
}
