package configs

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupServer(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	return r
}