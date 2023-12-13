package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"golang-phone-review-api/controllers"
	"golang-phone-review-api/middlewares"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/change-password", controllers.ChangePassword)

	/*Phone Endpoints*/
	r.GET("/phones", controllers.GetAllPhones)
	r.GET("/phones/:phone_id", controllers.GetPhoneByID)
	r.GET("/phones/:phone_id/reviews", controllers.GetReviewsForPhoneID)

	phonesMiddlewareRoute := r.Group("/phones")
	phonesMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
	phonesMiddlewareRoute.POST("/:phone_id/reviews", controllers.CreateReviewForPhone)
	phonesMiddlewareRoute.POST("/", controllers.CreatePhone)

	/*Review Endpoints*/
	r.GET("/reviews", controllers.GetAllReviews)

	reviewsMiddlewareRoute := r.Group("/reviews")
	reviewsMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
	reviewsMiddlewareRoute.PUT("/:review_id", controllers.EditReview)
	reviewsMiddlewareRoute.DELETE("/:review_id", controllers.DeleteReview)
	reviewsMiddlewareRoute.POST("/:review_id/comment", controllers.CreateReviewComment)
	// reviewsMiddlewareRoute.GET("/:review_id/comment", controllers.GetReviewComments)
	// reviewsMiddlewareRoute.POST("/:review_id/likes", controllers.CreateReviewLike)
	// reviewsMiddlewareRoute.GET("/:review_id/likes", controllers.GetReviewLikes)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
