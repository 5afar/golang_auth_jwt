package handler

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/token", getToken)
		auth.POST("/refresh", refreshToken)
	}
	return router
}
