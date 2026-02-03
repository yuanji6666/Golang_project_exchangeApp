package router

import (
	"exchangeapp/controllers"

	"github.com/gin-gonic/gin"
)

func registerExchangeRatesRouter(r *gin.RouterGroup){
	r.GET("/exchangeRates", controllers.GetExchangeRates)
	r.POST("/exchangeRates", controllers.CreateExchangeRate)
}