package controllers

import (
	"net/http"
	"time"
	"strconv"

	"golang-phone-review-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	AllPhoneResponse struct {
		PhoneID     int       `json:"phone_id" example:"1" extensions:"x-order=0"`
		Brand       string    `json:"brand" example:"Samsung" extensions:"x-order=1"`
		Model       string    `json:"model" example:"Galaxy" extensions:"x-order=2"`
		ReleaseDate time.Time `json:"release_date" example:"2023-11-11T00:00:00+07:00" extensions:"x-order=3"`
		Price       int       `json:"price" example:"10000" extensions:"x-order=4"`
		ImageURL    string    `json:"image_url" example:"" extensions:"x-order=5"`
	}

	PhoneByIDResponse struct {
		PhoneID     int       `json:"phone_id" example:"1" extensions:"x-order=0"`
		Brand       string    `json:"brand" example:"Samsung" extensions:"x-order=1"`
		Model       string    `json:"model" example:"Galaxy" extensions:"x-order=2"`
		ReleaseDate time.Time `json:"release_date" example:"2023-11-11T00:00:00+07:00" extensions:"x-order=3"`
		Price       int       `json:"price" example:"10000" extensions:"x-order=4"`
		ImageURL    string    `json:"image_url" example:"" extensions:"x-order=5"`
		Reviews 		[]ReviewResponse	`json:"reviews" extensions:"x-order=6"`
	}

	ReviewResponse struct {
		ReviewID  	int      	`json:"review_id" example:"1"`
		PhoneID   	int      	`json:"phone_id" example:"1"`
		UserID    	int      	`json:"user_id" example:"1"`
		Rating    	int      	`json:"rating" example:"9"`
		ReviewText 	string    `json:"review_text" example:"product is nice"`
	}

	CreateReviewInput struct {
		Rating      int   `json:"rating" binding:"required" example:"1"`
    ReviewText  string `json:"review_text" binding:"required" example:"sample review text"`	
	}


	CreateReviewByPhoneID struct {
		Message string `json:"message" example:"review created successfully" extensions:"x-order=0"`
		Reviews []models.Review `json:"reviews" extensions:"x-order=1"`
	}
	
)

// @Summary List all phones
// @Description Get a list of Phones
// @Tags Phones
// @Produce json
// @Success 200 {object} []AllPhoneResponse
// @Router /phones [get]
func GetAllPhones(c *gin.Context) {
	var phones []models.Phone
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&phones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch phones"})
		return
	}

	var result []AllPhoneResponse

	for _, p := range phones {
		result = append(result, AllPhoneResponse{
			PhoneID: p.PhoneID,
			Brand: p.Brand,
			Model: p.Model,
			ReleaseDate: p.ReleaseDate,
			Price: p.Price,
			ImageURL: p.ImageURL,
		})
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Get phone details by ID
// @Description Get phone details by ID
// @Tags Phones
// @Produce json
// @Param phone_id path string true "PhoneID"
// @Success 200 {object} PhoneByIDResponse
// @Router /phones/{phone_id} [get]
func GetPhoneByID(c *gin.Context) {
	var phone models.Phone
	db := c.MustGet("db").(*gorm.DB)
	phoneID := c.Param("phone_id")

	if err := db.Preload("Reviews").First(&phone, phoneID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phone not found"})
		return
	}

  r := PhoneByIDResponse{
		PhoneID:     phone.PhoneID,
		Brand:       phone.Brand,
		Model:       phone.Model,
		ReleaseDate: phone.ReleaseDate,
		Price:       phone.Price,
		ImageURL:    phone.ImageURL,
	}

	reviews := make([]ReviewResponse, len(phone.Reviews))
	for i, review := range phone.Reviews {
		reviews[i] = ReviewResponse{
				ReviewID:   review.ReviewID,
				PhoneID:    review.PhoneID,
				UserID:     review.UserID,
				Rating:     review.Rating,
				ReviewText: review.ReviewText,
		}
	}
	r.Reviews = reviews

	c.JSON(http.StatusOK, r)
}

// @Summary Get all phone reviews by its ID
// @Description Get all phone reviews by ID
// @Tags Phones
// @Produce json
// @Param phone_id path string true "PhoneID"
// @Success 200 {object} models.Review
// @Router /phones/{phone_id}/reviews [get]
func GetReviewsForPhoneID(c *gin.Context) {
	var review models.Review
	db := c.MustGet("db").(*gorm.DB)
	phoneID := c.Param("phone_id")

	if err := db.Preload("Comments").Preload("Likes").Where("phone_id= ?", phoneID).Find(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// @Summary Create a phone review by its ID
// @Description Create phone review by ID
// @Tags Phones
// @Produce json
// @Param phone_id path string true "PhoneID"
// @Param Body body CreateReviewInput true "the body to create phone review"
// @Security Bearer
// @Success 200 {object} models.Review
// @Router /phones/{phone_id}/reviews [post]
func CreateReviewForPhone(c *gin.Context) {
	var input CreateReviewInput
	db := c.MustGet("db").(*gorm.DB)
	phoneID := c.Param("phone_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
	}
	
	parsedPhoneID, _ := strconv.Atoi(phoneID)
	review := models.Review{
		PhoneID		: parsedPhoneID,
		Rating		: input.Rating,
		UserID		: int(userID.(float64)),
		ReviewText: input.ReviewText,
		CreatedAt : time.Now(),
	}

	if err := db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review created successfully", "review": review})
}