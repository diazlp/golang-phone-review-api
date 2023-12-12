package controllers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/bcrypt"
	// "github.com/golang-jwt/jwt"
	// "github.com/joho/godotenv"
	// "log"
	// "time"

	"golang-phone-review-api/models"
	"gorm.io/gorm"
	// "golang-phone-review-api/repositories"
)

type (
	LoginInput struct {
		Username string `json:"username" binding:"required" example:"admin"`
		Password string `json:"password" binding:"required" example:"admin"`
	}

	RegisterInput struct {
		Username 	string `json:"username" binding:"required" example:"admin"`
    Password 	string `json:"password" binding:"required" example:"admin"`
    Email    	string `json:"email" binding:"required" example:"admin@mail.com"`
    Role    	string `json:"role" binding:"required" example:"Admin"`
	}

	LoginResponse struct {
		Message 	string 	`json:"message" example:"login success"`
		User 			struct {
			Username string `json:"username" example:"John"`
			Email    string `json:"email" example:"john@example.com"`
			Role     string `json:"role" example:"user"`
		} `json:"user"`
		Token 		string 	`json:"token" example:"string"`
	}

	RegisterResponse struct {
		Message 	string 	`json:"message" example:"registration success"`
		User 			struct {
			Username string `json:"username" example:"John"`
			Email    string `json:"email" example:"john@example.com"`
			Role     string `json:"role" example:"user"`
		} `json:"user"`
		
	}
)

// @Summary Login as as user.
// @Description Logging in to get jwt token to access admin or user api by roles.
// @Tags Auth
// @Param Body body LoginInput true "the body to login a user"
// @Produce json
// @Success 200 {object} LoginResponse
// @Router /login [post]
func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password
	
	token, err := models.LoginCheck(u.Username, u.Password, db)
	_ = db.Model(models.User{}).Where("username = ?", input.Username).Take(&u).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	user := map[string]string {
		"username": u.Username,
		"email": u.Email,
	}

	c.JSON(http.StatusOK, gin.H{"message": "login success", "user": user, "token": token})
}

// @Summary Register a user.
// @Description registering a user from public access.
// @Tags Auth
// @Param Body body RegisterInput true "the body to register a user"
// @Produce json
// @Success 201 {object} RegisterResponse  "Register Success Response"
// @Router /register [post]
func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Email = input.Email
	u.Password = input.Password
	u.Role = input.Role

	_, err := u.SaveUser(db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := map[string]string {
		"username": input.Username,
		"email": input.Email,
		"role": input.Role,
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registration success", "user": user})
}



