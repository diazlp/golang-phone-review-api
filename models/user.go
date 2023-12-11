package models

import "time"

type (
	User struct {
		UserID    int      	`gorm:"primaryKey" json:"user_id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
		Reviews 	[]Review	`gorm:"foreignKey:user_id;references:user_id"`
		Comments 	[]Comment	`gorm:"foreignKey:user_id;references:user_id"`
		Likes 		[]Like		`gorm:"foreignKey:user_id;references:user_id"`
	}
)