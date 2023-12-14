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
	/*find-all review*/
	AllReviewResponse struct {
		ReviewID  	int      	`json:"review_id" example:"1" extensions:"x-order=0"`
		PhoneID   	int      	`json:"phone_id" example:"1" extensions:"x-order=1"`
		Phone       PhoneResponse `json:"phone,omitempty" extensions:"x-order:2"`
		UserID    	int      	`json:"user_id" example:"1" extensions:"x-order=3"`
		Rating    	int      	`json:"rating" example:"1" extensions:"x-order=4"`
		ReviewText 	string    `json:"review_text" example:"this is sample text" extensions:"x-order=5"`
		CreatedAt 	time.Time `json:"created_at" example:"" extensions:"x-order=6"`
	}

	PhoneResponse struct {
		PhoneID     int       `json:"phone_id" example:"1"`
		Brand       string    `json:"brand" example:"Samsung"`
		Model       string    `json:"model" example:"Galaxy"`
		ReleaseDate time.Time `json:"release_date" example:"2023-11-11"`
		Price       int       `json:"price" example:"10000"`
		ImageURL    string    `json:"image_url" example:""`
	}

	/*update review*/
	EditReviewInput struct {
		Rating      int   `json:"rating" example:"1"`
    ReviewText  string `json:"review_text" example:"sample review text"`	
	}

	EditReviewResponse struct {
		Message string `json:"message" example:"review updated successfully" extensions:"x-order=0"`
		Reviews []models.Review `json:"reviews" extensions:"x-order=1"`
	}

	/*delete review*/
	DeleteReviewResponse struct {
		Message string `json:"message" example:"review deleted successfully" extensions:"x-order=0"`
	}

	/*comment a review*/
	CreateCommentInput struct {
		CommentText string    `json:"comment_text" binding:"required" example:"this review is rigged!"`
	}

	CreateCommentResponse struct {
		Message string `json:"message" example:"comment created successfully" extensions:"x-order=0"`
		Comment models.Comment `json:"comments" extensions:"x-order=1"`
	}

	/*like a review*/
	CreateLikeResponse struct {
		Message string `json:"message" example:"like created successfully" extensions:"x-order=0"`
	}
)

// @Summary List all reviews
// @Description Get a list of Reviews
// @Tags Reviews
// @Produce json
// @Success 200 {object} []AllReviewResponse
// @Router /reviews [get]
func GetAllReviews(c *gin.Context) {
	var reviews []models.Review
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch phones"})
		return
	}

	var result []AllReviewResponse

	var phone models.Phone

	for _, r := range reviews {
		if err := db.Find(&phone, r.PhoneID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Phone not found"})
			return
		}

		result = append(result, AllReviewResponse{
			ReviewID: r.ReviewID,
			PhoneID: r.PhoneID,
			Phone: PhoneResponse{
				PhoneID: phone.PhoneID,
				Brand: phone.Brand,
				Model: phone.Model,
				ReleaseDate: phone.ReleaseDate,
				Price: phone.Price,
				ImageURL: phone.ImageURL,
			},
			UserID: r.UserID,
			Rating: r.Rating,
			ReviewText: r.ReviewText,
			CreatedAt: r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Update a review
// @Description Update a review
// @Tags Reviews
// @Produce json
// @Param review_id path string true "ReviewID"
// @Param Body body EditReviewInput true "the body to edit phone review"
// @Security Bearer
// @Success 200 {object} EditReviewResponse
// @Router /reviews/{review_id} [put]
func EditReview(c *gin.Context) {
	var input EditReviewInput
	db := c.MustGet("db").(*gorm.DB)
	reviewID := c.Param("review_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingReview models.Review
	if err := db.First(&existingReview, reviewID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review not found"})
		return
	}

	userRole, _ := c.Get("user_role")
	if userRole != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"forbidden": `User does not have "Admin" role to update review`})
			return
	}

	if input.Rating != 0 {
		existingReview.Rating = input.Rating
	}
	if input.ReviewText != "" {
		existingReview.ReviewText = input.ReviewText
	}
	
	if err := db.Save(&existingReview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review updated successfully", "review": existingReview})
}

// @Summary Delete a review
// @Description Delete a review
// @Tags Reviews
// @Produce json
// @Param review_id path string true "ReviewID"
// @Security Bearer
// @Success 200 {object} DeleteReviewResponse
// @Router /reviews/{review_id} [delete]
func DeleteReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	reviewID := c.Param("review_id")

	var review models.Review // Change Review to your actual model name
	if err := db.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review not found"})
		return
	}

	userRole, _ := c.Get("user_role")
	if userRole != "Admin" && userRole != "Writer" {
			c.JSON(http.StatusForbidden, gin.H{"forbidden": `User does not have "Admin" or "Writer" role to delete review`})
			return
	}

	db.Delete(&review)

	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}

// @Summary Create a review comment by review ID
// @Description Create a review comment by review ID
// @Tags Reviews
// @Produce json
// @Param review_id path string true "ReviewID"
// @Param Body body CreateCommentInput true "the body to create phone review"
// @Security Bearer
// @Success 200 {object} CreateCommentResponse
// @Router /reviews/{review_id}/comments [post]
func CreateReviewComment(c *gin.Context) {
	var input CreateCommentInput
	db := c.MustGet("db").(*gorm.DB)
	reviewID := c.Param("review_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingReview models.Review
	if err := db.First(&existingReview, reviewID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
	}

	userRole, _ := c.Get("user_role")
	if userRole != "Writer" {
			c.JSON(http.StatusForbidden, gin.H{"forbidden": `User does not have "Writer" role to comment a review`})
			return
	}

	parsedReviewID, _ := strconv.Atoi(reviewID)
	comment := models.Comment{
		ReviewID: parsedReviewID,
		UserID: int(userID.(float64)),
		CommentText: input.CommentText,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to comment a review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "comment created successfully", "comment": comment})
}

// @Summary Get all review comment by review ID
// @Description Get all review comment by review ID
// @Tags Reviews
// @Produce json
// @Param review_id path string true "ReviewID"
// @Success 200 {object} []models.Comment
// @Router /reviews/{review_id}/comments [get]
func GetReviewComments(c *gin.Context) {
	var comments []models.Comment
	db := c.MustGet("db").(*gorm.DB)
	reviewID := c.Param("review_id")

	var existingReview models.Review
	if err := db.First(&existingReview, reviewID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review not found"})
		return
	}

	if err := db.Where("review_id= ?", reviewID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
	
// @Summary Create a review like by review ID
// @Description Create a review like by review ID
// @Tags Reviews
// @Produce json
// @Param review_id path string true "ReviewID"
// @Security Bearer
// @Success 201 {object} CreateLikeResponse
// @Router /reviews/{review_id}/likes [post]
func CreateReviewLike(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	reviewID := c.Param("review_id")

	var existingReview models.Review
	if err := db.First(&existingReview, reviewID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
	}

	userRole, _ := c.Get("user_role")
	if userRole != "Guest" {
			c.JSON(http.StatusForbidden, gin.H{"forbidden": `User does not have "Guest" role to like a review`})
			return
	}

	parsedReviewID, _ := strconv.Atoi(reviewID)
	like := models.Like{
		ReviewID: &parsedReviewID,
		UserID: int(userID.(float64)),
		CreatedAt: time.Now(),
	}

	if err := db.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like a review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "like created successfully"})
}


// @Summary Get all review likes by review ID
// @Description Get all review likes by review ID
// @Tags Reviews
// @Produce json
// @Param review_id path string true "ReviewID"
// @Success 200 {object} []models.Like
// @Router /reviews/{review_id}/likes [get]
func GetReviewLike(c *gin.Context) {
var likes []models.Like
db := c.MustGet("db").(*gorm.DB)
reviewID := c.Param("review_id")

var existingReview models.Review
if err := db.First(&existingReview, reviewID).Error; err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Review not found"})
	return
}

if err := db.Where("review_id= ?", reviewID).Find(&likes).Error; err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review comments"})
	return
}

c.JSON(http.StatusOK, likes)
}