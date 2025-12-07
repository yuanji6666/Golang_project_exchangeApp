package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{
	r := gin.Default()
	auth := r.Group("/api/auth") 
	auth.POST("/login", controllers.Login)
	auth.POST("/register", controllers.Register)

	api := r.Group("/api")
	api.GET("/exchangeRates", controllers.GetExchangeRates)
	api.Use(middlewares.AuthMiddleWare())
	{
	api.POST("/exchangeRates", controllers.CreateExchangeRate)
	}
	return r
	
}