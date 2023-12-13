package main

import (
	"golang-phone-review-api/configs"
	// "github.com/gin-gonic/gin"
	"golang-phone-review-api/routes"
	"golang-phone-review-api/docs"
)

// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Set up Swagger Info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Movie."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8070"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	db, err := configs.SetupDatabase()
	if err != nil {
		panic(err.Error())
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := configs.SetupServer(db)
	routes.SetupRoutes(r, db)
	r.Run("localhost:8070")
}