package models

import "time"

type (
	Comment struct {
		ID          uint      `gorm:"primaryKey" json:"comment_id"`
		ReviewID    uint      `json:"review_id"`
		UserID      uint      `json:"user_id"`
		CommentText string    `json:"comment_text"`
		CreatedAt   time.Time `json:"created_at"`
	}
)