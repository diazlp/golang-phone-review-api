package models

import "time"

type (
	Review struct {
		ReviewID  	int      	`gorm:"primaryKey" json:"review_id" example:"1" extensions:"x-order=0"`
		PhoneID   	int      	`json:"phone_id" example:"1" extensions:"x-order=1"`
		UserID    	int      	`json:"user_id" example:"1" extensions:"x-order=2"`
		Rating    	int      	`json:"rating" example:"8" extensions:"x-order=3"`
		ReviewText 	string    `json:"review_text" example:"sample review text" extensions:"x-order=4"`
		CreatedAt 	time.Time `json:"created_at" example:"2030-01-01 00:00:00" extensions:"x-order=5"`
		Comments 		[]Comment `gorm:"foreignKey:review_id;references:review_id" extensions:"x-order=5"`
		Likes 			[]Like 		`gorm:"foreignKey:review_id;references:review_id" extensions:"x-order=6"`
	}
)