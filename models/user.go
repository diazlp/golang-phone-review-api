package models

import "time"

type (
	User struct {
		ID        uint      `gorm:"primaryKey" json:"user_id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Password  string    `json:"-"`
		CreatedAt time.Time `json:"created_at"`
	}
)