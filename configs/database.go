package configs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"golang-phone-review-api/models"
)

func SetupDatabase() (*gorm.DB, error) {
	username := "root"
	password := ""
	host := "tcp(127.0.0.1:3306)"
	database := "phone-review"

	dsn := fmt.Sprintf("%v:%v@%v/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Phone{}, &models.Review{}, &models.Comment{}, &models.Like{})
	if err != nil {
			return nil, err
	}

	return db, nil
}