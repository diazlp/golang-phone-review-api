package models

import "time"

type (
	Review struct {
		ReviewID  	int      	`gorm:"primaryKey" json:"review_id"`
		PhoneID   	int      	`json:"phone_id"`
		UserID    	int      	`json:"user_id"`
		Rating    	int      	`json:"rating"`
		ReviewText 	string    `json:"review_text"`
		CreatedAt 	time.Time `json:"created_at"`
		Comments 		[]Comment `gorm:"foreignKey:review_id;references:review_id"`
		Likes 			[]Like 		`gorm:"foreignKey:review_id;references:review_id"`
	}
)