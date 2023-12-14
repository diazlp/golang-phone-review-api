package controllers

import (
	"net/http"
	"time"

	"golang-phone-review-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	/*find-all likes*/
	AllLikeResponse struct {
		LikeID    	int      	`json:"like_id" example:"1" extensions:"x-order=0"`
		ReviewID  	*int      	`json:"review_id,omitempty" example:"1" extensions:"x-order=1"`
		Review      *ReviewResponse `json:"review,omitempty" extensions:"x-order=2"`
		UserID    	int      	`json:"user_id" example:"1" extensions:"x-order=3"`
		CommentID   *int      	`json:"comment_id,omitempty" example:"1" extensions:"x-order=4"`
		Comment			*CommentResponse `json:"comment,omitempty" extensions:"x-order=5"`
		CreatedAt 	time.Time `json:"created_at" example:"2030-01-01 00:00:00" extensions:"x-order=6"`
	}

	/*delete like*/
	DeleteLikeResponse struct {
		Message string `json:"message" example:"like deleted successfully" extensions:"x-order=0"`
	}
)

// @Summary List all likes
// @Description Get a list of likes
// @Tags Likes
// @Produce json
// @Success 200 {object} []AllLikeResponse
// @Router /likes [get]
func GetAllLikes(c *gin.Context) {
	var likes []models.Like
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch likes"})
		return
	}

	var result []AllLikeResponse

	var review models.Review
	var comment models.Comment

	for _, l := range likes {
		likeResponse := AllLikeResponse{
			LikeID: l.LikeID,
			ReviewID: nil,
			UserID: l.UserID,
			CommentID: nil,
			CreatedAt: l.CreatedAt,
		}

		
		if l.ReviewID != nil {
			db.Preload("Likes").Find(&review, l.ReviewID)

			likeResponse.ReviewID = &(*l.ReviewID)
			likeResponse.Review = &ReviewResponse{
				ReviewID: review.ReviewID,
				PhoneID: review.PhoneID,
				UserID: review.UserID,
				Rating: review.Rating,
				ReviewText: review.ReviewText,
				TotalLikes: len(review.Likes),
			}
		} else {
			likeResponse.Review = nil
		}

		if l.CommentID != nil {
			db.Preload("Likes").Find(&comment, l.CommentID)

			likeResponse.CommentID = &(*l.CommentID)
			likeResponse.Comment = &CommentResponse{
				CommentID: comment.CommentID,
				ReviewID: comment.ReviewID,
				UserID: comment.UserID,
				CommentText: comment.CommentText,
				CreatedAt: comment.CreatedAt,
				TotalLikes: len(comment.Likes),
			}
		} else {
			likeResponse.Comment = nil
		}

		result = append(result, likeResponse)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Delete a like
// @Description Delete a like
// @Tags Likes
// @Produce json
// @Param like_id path string true "LikeID"
// @Security Bearer
// @Success 200 {object} DeleteLikeResponse
// @Router /likes/{like_id} [delete]
func DeleteLike(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	likeID := c.Param("like_id")

	var like models.Like
	if err := db.First(&like, likeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Like not found"})
		return
	}

	userRole, _ := c.Get("user_role")
	if userRole != "Guest" {
		c.JSON(http.StatusForbidden, gin.H{"forbidden": `User does not have "Guest" role to delete like`})
		return
	}

	db.Delete(&like)

	c.JSON(http.StatusOK, gin.H{"message": "like deleted successfully"})
}