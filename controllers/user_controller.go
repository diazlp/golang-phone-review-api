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
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	RegisterInput struct {
		Username 	string `json:"username" binding:"required"`
    Password 	string `json:"password" binding:"required"`
    Email    	string `json:"email" binding:"required"`
    Role    	string `json:"role" binding:"required"`
	}
)

func LoginHandler(c *gin.Context) {
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



