package controllers

import (
	"net/http"
	"strconv"
	"time"

	"golang-phone-review-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type(
	/*general comment response*/
	CommentResponse struct {
		CommentID   int      	`json:"comment_id" example:"1"`
		ReviewID  	int      	`json:"review_id" example:"1"`
		UserID    	int      	`json:"user_id" example:"1"`
		CommentText string    `json:"comment_text" example:"sample comment text"`
		CreatedAt 	time.Time `json:"created_at" example:"2030-01-01 00:00:00"`
		TotalLikes	int 			`json:"total_likes" example:"10"`
	}

	/*find-all comment*/
	AllCommentResponse struct {
		CommentID   int      	`json:"comment_id" example:"1"`
		ReviewID  	int      	`json:"review_id" example:"1" extensions:"x-order=0"`
		Review      ReviewResponse `json:"review,omitempty" extensions:"x-order:2"`
		UserID    	int      	`json:"user_id" example:"1" extensions:"x-order=3"`
		CommentText string    `json:"comment_text" example:"sample comment text"`
		CreatedAt 	time.Time `json:"created_at" example:"2030-01-01 00:00:00" extensions:"x-order=6"`
	}

	/*update comment*/
	EditCommentInput struct {
    CommentText  string `json:"comment_text" binding:"required" example:"sample comment text"`	
	}

	EditCommentResponse struct {
		Message string `json:"message" example:"comment updated successfully" extensions:"x-order=0"`
		Comments []models.Comment `json:"comments" extensions:"x-order=1"`
	}

	/*delete review*/
	DeleteCommentResponse struct {
		Message string `json:"message" example:"comment deleted successfully" extensions:"x-order=0"`
	}

	/*get comment likes*/
	GetCommentLikeResponse struct {
		Count		int 	`json:"count" example:"1" extensions:"x-order=0"`
		Rows		[]models.Like 	`json:"rows,omitempty" extensions:"x-order=1"`
	}
)

// @Summary List all comments
// @Description Get a list of comments
// @Tags Comments
// @Produce json
// @Success 200 {object} []AllCommentResponse
// @Router /comments [get]
func GetAllComments(c *gin.Context) {
	var comments []models.Comment
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	var result []AllCommentResponse

	var review models.Review

	for _, cmt := range comments {
		if err := db.Preload("Likes").Find(&review, cmt.ReviewID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}

		result = append(result, AllCommentResponse{
			CommentID: cmt.CommentID,
			ReviewID: cmt.ReviewID,
			Review: ReviewResponse{
				ReviewID: review.ReviewID,
				PhoneID: review.PhoneID,
				UserID: review.UserID,
				Rating: review.Rating,
				ReviewText: review.ReviewText,
				TotalLikes: len(review.Likes),
			},
			UserID: cmt.UserID,
			CommentText: cmt.CommentText,
			CreatedAt: cmt.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Update a comment
// @Description Update a comment
// @Tags Comments
// @Produce json
// @Param comment_id path string true "CommentID"
// @Param Body body EditCommentInput true "the body to edit review comment"
// @Security Bearer
// @Success 200 {object} EditCommentResponse
// @Router /comments/{comment_id} [put]
func EditComment(c *gin.Context) {
	var input EditCommentInput
	db := c.MustGet("db").(*gorm.DB)
	commentID := c.Param("comment_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingComment models.Comment
	if err := db.First(&existingComment, commentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment not found"})
		return
	}

	userRole, _ := c.Get("user_role")
	if userRole != "Writer" {
			c.JSON(http.StatusForbidden, gin.H{"forbidden": `User does not have "Writer" role to update comment`})
			return
	}

	existingComment.CommentText = input.CommentText

	if err := db.Save(&existingComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment updated successfully", "review": existingComment})
}

// @Summary Delete a comment
// @Description Delete a comment
// @Tags Comments
// @Produce json
// @Param comment_id path string true "CommentID"
// @Security Bearer
// @Success 200 {object} DeleteCommentResponse
// @Router /comments/{comment_id} [delete]
func DeleteComment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	commentID := c.Param("comment_id")

	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment not found"})
		return
	}

	userRole, _ := c.Get("user_role")
	if userRole != "Writer" {
		c.JSON(http.StatusForbidden, gin.H{"forbidden": `User does not have "Admin" or "Writer" role to delete comment`})
		return
	}

	db.Delete(&comment)

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}

// @Summary Create a comment like by comment ID
// @Description Create a comment like by comment ID
// @Tags Comments
// @Produce json
// @Param comment_id path string true "CommentID"
// @Security Bearer
// @Success 201 {object} CreateLikeResponse
// @Router /comments/{comment_id}/likes [post]
func CreateCommentLike(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	commentID := c.Param("comment_id")

	var existingComment models.Comment
	if err := db.First(&existingComment, commentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
	}

	userRole, _ := c.Get("user_role")
	if userRole != "Guest" {
			c.JSON(http.StatusForbidden, gin.H{"forbidden": `User does not have "Writer" role to comment a review`})
			return
	}

	parsedCommentID, _ := strconv.Atoi(commentID)
	like := models.Like{
		CommentID: &parsedCommentID,
		UserID: int(userID.(float64)),
		CreatedAt: time.Now(),
	}

	if err := db.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like a comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "like created successfully"})
}

// @Summary Get all comment likes by comment ID
// @Description Get all comment likes by comment ID
// @Tags Comments
// @Produce json
// @Param comment_id path string true "CommentID"
// @Success 200 {object} GetCommentLikeResponse
// @Router /comments/{comment_id}/likes [get]
func GetCommentLike(c *gin.Context) {
	var likes []models.Like
	db := c.MustGet("db").(*gorm.DB)
	commentID := c.Param("comment_id")

	var existingComment models.Comment
	if err := db.First(&existingComment, commentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment not found"})
		return
	}

	if err := db.Where("comment_id= ?", commentID).Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment likes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": len(likes), "rows": likes})
}