package configs

import "github.com/gin-gonic/gin"

func SetupServer() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	return r
}