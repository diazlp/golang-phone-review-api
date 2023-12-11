package main

import (
	"golang-phone-review-api/configs"
	// "github.com/gin-gonic/gin"
)

func main() {
	db, err := configs.SetupDatabase()
	if err != nil {
		panic(err.Error())
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// router := configs.SetupServer()
}