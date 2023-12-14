package models

import "time"

type (
	Comment struct {
		CommentID   int      	`gorm:"primaryKey" json:"comment_id" example:"1"`
		ReviewID    int      	`json:"review_id" example:"1"`
		UserID      int      	`json:"user_id" example:"1"`
		CommentText string    `json:"comment_text" example:"sample comment text"`
		CreatedAt   time.Time `json:"created_at" example:"2030-01-01 00:00:00"`
		Likes 			[]Like 		`gorm:"foreignKey:comment_id;references:comment_id"`
	}
)