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
		Username string `json:"username" binding:"required" example:"admin" extensions:"x-order=0"`
		Password string `json:"password" binding:"required" example:"admin" extensions:"x-order=1"`
	}

	RegisterInput struct {
		Username 	string `json:"username" binding:"required" example:"admin" extensions:"x-order=0"`
    Password 	string `json:"password" binding:"required" example:"admin" extensions:"x-order=1"`
    Email    	string `json:"email" binding:"required" example:"admin@mail.com" extensions:"x-order=2"`
    Role    	string `json:"role" binding:"required" example:"Admin" extensions:"x-order=3"`
	}

	ChangePasswordInput struct {
		Username 						string `json:"username" binding:"required" example:"admin" extensions:"x-order=0"`
    CurrentPassword 		string `json:"current_password" binding:"required" example:"admin" extensions:"x-order=1"`
    NewPassword 				string `json:"new_password" binding:"required" example:"admin1" extensions:"x-order=2"`
    ConfirmNewPassword 	string `json:"confirm_new_password" binding:"required" example:"admin1" extensions:"x-order=3"`
	}

	LoginResponse struct {
		Message 	string 	`json:"message" example:"login success" extensions:"x-order=0"`
		User 			struct {
			Username string `json:"username" example:"John"`
			Email    string `json:"email" example:"john@example.com"`
			Role     string `json:"role" example:"user"`
		} `json:"user" extensions:"x-order=1"`
		Token 		string 	`json:"token" example:"string" extensions:"x-order=2"`
	}

	RegisterResponse struct {
		Message 	string 	`json:"message" example:"registration success" extensions:"x-order=0"`
		User 			struct {
			Username string `json:"username" example:"John"`
			Email    string `json:"email" example:"john@example.com"`
			Role     string `json:"role" example:"Admin"`
		} `json:"user" extensions:"x-order=1"`
	}

	ChangePasswordResponse struct {
		Message		string 	`json:"message" example:"change password success"`
	}
)

// @Summary Login as user.
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

// @Summary Change user password.
// @Description Change user password by inputting the current password and the new password.
// @Tags Auth
// @Param Body body ChangePasswordInput true "the body to change user password"
// @Produce json
// @Success 200 {object} ChangePasswordResponse
// @Router /change-password [post]
func ChangePassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input ChangePasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.NewPassword != input.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new password should match confirm new password"})
		return
	}

	u := models.User{}
	u.Password = input.CurrentPassword

	_, err := models.LoginCheck(input.Username, input.CurrentPassword, db)
	_ = db.Model(models.User{}).Where("username = ?", input.Username).Take(&u).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	if _, err := u.UpdateUser(input.NewPassword, db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "change password success"})
}