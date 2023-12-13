package models

import "time"

type (
	Phone struct {
		PhoneID     int      	`gorm:"primaryKey" json:"phone_id" extensions:"x-order=0"`
		Brand       string    `json:"brand" extensions:"x-order=1"`
		Model       string    `json:"model" extensions:"x-order=2"`
		ReleaseDate time.Time `json:"release_date" extensions:"x-order=3"`
		Price       int      	`json:"price" extensions:"x-order=4"`
		ImageURL    string    `json:"image_url" extensions:"x-order=5"`
		Reviews 		[]Review	`gorm:"foreignKey:phone_id;references:phone_id" extensions:"x-order=6"`
	}
)