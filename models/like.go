package models

import "time"

type (
	Like struct {
		LikeID    int      	`gorm:"primaryKey" json:"like_id"`
		ReviewID  int      	`json:"review_id"`
		UserID    int      	`json:"user_id"`
		CommentID int      	`json:"comment_id"`
		CreatedAt time.Time `json:"created_at"`
	}
)