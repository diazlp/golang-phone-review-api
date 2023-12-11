package models

import "time"

type (
	Review struct {
		ID        uint      `gorm:"primaryKey" json:"review_id"`
		PhoneID   uint      `json:"phone_id"`
		UserID    uint      `json:"user_id"`
		Rating    uint      `json:"rating"`
		ReviewText string    `json:"review_text"`
		CreatedAt time.Time `json:"created_at"`
	}
)