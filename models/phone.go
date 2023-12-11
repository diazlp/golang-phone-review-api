package models

import "time"

type (
	Phone struct {
		ID          uint      `gorm:"primaryKey" json:"phone_id"`
		Brand       string    `json:"brand"`
		Model       string    `json:"model"`
		ReleaseDate time.Time `json:"release_date"`
		Price       uint      `json:"price"`
		ImageURL    string    `json:"image_url"`
	}
)