package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"golang-phone-review-api/controllers"
	// "golang-phone-review-api/middlewares"

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/change-password", controllers.ChangePassword)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
