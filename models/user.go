package models

import (
	"time"
	"golang.org/x/crypto/bcrypt"
  "gorm.io/gorm"
	"golang-phone-review-api/utils/token"
)

type (
	User struct {
		UserID    int      	`gorm:"primaryKey" json:"user_id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		Role			string		`json:"role"`
		CreatedAt time.Time `json:"created_at"`
		Reviews 	[]Review	`gorm:"foreignKey:user_id;references:user_id"`
		Comments 	[]Comment	`gorm:"foreignKey:user_id;references:user_id"`
		Likes 		[]Like		`gorm:"foreignKey:user_id;references:user_id"`
	}
)

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string, db *gorm.DB) (string, error) {
	var err error
	u := User{}
	
	err = db.Model(User{}).Where("username = ?", username).Take(&u).Error
	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.UserID)
	if err != nil {
		return "", err
	}
	
	return token, nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	//turn password into hash
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		return &User{}, errPassword
	}

	u.Password = string(hashedPassword)
	// remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	var err Error = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u,nil
}