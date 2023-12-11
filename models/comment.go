package models

import "time"

type (
	Comment struct {
		CommentID   int      	`gorm:"primaryKey" json:"comment_id"`
		ReviewID    int      	`json:"review_id"`
		UserID      int      	`json:"user_id"`
		CommentText string    `json:"comment_text"`
		CreatedAt   time.Time `json:"created_at"`
		Likes 			[]Like 		`gorm:"foreignKey:comment_id;references:comment_id"`
	}
)