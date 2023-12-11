package models

import "time"

type (
	Like struct {
		ID        uint      `gorm:"primaryKey" json:"like_id"`
		ReviewID  uint      `json:"review_id"`
		UserID    uint      `json:"user_id"`
		CommentID uint      `json:"comment_id"`
		CreatedAt time.Time `json:"created_at"`
	}
)