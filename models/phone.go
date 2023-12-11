package models

import "time"

type (
	Phone struct {
		PhoneID     int      	`gorm:"primaryKey" json:"phone_id"`
		Brand       string    `json:"brand"`
		Model       string    `json:"model"`
		ReleaseDate time.Time `json:"release_date"`
		Price       int      	`json:"price"`
		ImageURL    string    `json:"image_url"`
		Reviews 		[]Review	`gorm:"foreignKey:phone_id;references:phone_id"`
	}
)