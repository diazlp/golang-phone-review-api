package models

import "time"

type (
	Like struct {
		LikeID    int      	`gorm:"primaryKey" json:"like_id" example:"1"`
		ReviewID  *int      `json:"review_id" example:"1"`
		UserID    int      	`json:"user_id" example:"1"`
		CommentID *int      `json:"comment_id" example:"1"`
		CreatedAt time.Time `json:"created_at" example:"2030-01-01 00:00:00"`
	}
)